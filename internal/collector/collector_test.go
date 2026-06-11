package collector

import (
	"context"
	"errors"
	"testing"

	"github.com/Pantani/healthcheck/internal/fixtures"
)

func TestRunnerCollectPassesAndStoresValue(t *testing.T) {
	store := newMemoryStore()
	store.data["api.status"] = float64(1)
	alerts := 0
	runner := Runner{
		Store: store,
		Alert: AlertFunc(func(_, _, _ string) error {
			alerts++
			return nil
		}),
	}

	runner.Collect(context.Background(), "api", fixtures.Test{
		Name:       "status",
		Method:     "GET",
		URLPath:    "health",
		JSONPath:   "version",
		Expression: "lastValue < newValue",
	}, requesterFunc(func(context.Context, string, string, interface{}) (string, error) {
		return `{"version":2}`, nil
	}))

	if alerts != 0 {
		t.Fatalf("alerts = %d, want 0", alerts)
	}
	if got := store.data["api.status"]; got != float64(2) {
		t.Fatalf("stored value = %v (%T), want 2", got, got)
	}
}

func TestRunnerCollectAlertsWhenExpressionFails(t *testing.T) {
	store := newMemoryStore()
	store.data["api.status"] = float64(3)
	alerts := 0
	runner := Runner{
		Store: store,
		Alert: AlertFunc(func(namespace, name, path string) error {
			alerts++
			if namespace != "api" || name != "status" || path != "health" {
				t.Fatalf("alert = %s %s %s", namespace, name, path)
			}
			return nil
		}),
	}

	runner.Collect(context.Background(), "api", fixtures.Test{
		Name:       "status",
		Method:     "GET",
		URLPath:    "health",
		JSONPath:   "version",
		Expression: "lastValue < newValue",
	}, requesterFunc(func(context.Context, string, string, interface{}) (string, error) {
		return `{"version":2}`, nil
	}))

	if alerts != 1 {
		t.Fatalf("alerts = %d, want 1", alerts)
	}
}

func TestRunnerCollectAlertsOnRequestError(t *testing.T) {
	store := newMemoryStore()
	alerts := 0
	runner := Runner{
		Store: store,
		Alert: AlertFunc(func(_, _, _ string) error {
			alerts++
			return nil
		}),
	}

	runner.Collect(context.Background(), "api", fixtures.Test{
		Name:       "status",
		Method:     "GET",
		URLPath:    "health",
		JSONPath:   "version",
		Expression: "newValue == true",
	}, requesterFunc(func(context.Context, string, string, interface{}) (string, error) {
		return "", errors.New("timeout")
	}))

	if alerts != 1 {
		t.Fatalf("alerts = %d, want 1", alerts)
	}
	if len(store.data) != 0 {
		t.Fatalf("store writes = %v, want none", store.data)
	}
}

func TestRunnerCollectSkipsMissingJSONPath(t *testing.T) {
	store := newMemoryStore()
	alerts := 0
	runner := Runner{
		Store: store,
		Alert: AlertFunc(func(_, _, _ string) error {
			alerts++
			return nil
		}),
	}

	runner.Collect(context.Background(), "api", fixtures.Test{
		Name:       "status",
		Method:     "GET",
		URLPath:    "health",
		JSONPath:   "missing",
		Expression: "newValue == true",
	}, requesterFunc(func(context.Context, string, string, interface{}) (string, error) {
		return `{"ok":true}`, nil
	}))

	if alerts != 0 {
		t.Fatalf("alerts = %d, want 0", alerts)
	}
	if len(store.data) != 0 {
		t.Fatalf("store writes = %v, want none", store.data)
	}
}

type requesterFunc func(context.Context, string, string, interface{}) (string, error)

func (f requesterFunc) Execute(ctx context.Context, method string, path string, body interface{}) (string, error) {
	return f(ctx, method, path, body)
}

type memoryStore struct {
	data map[string]interface{}
}

func newMemoryStore() *memoryStore {
	return &memoryStore{data: map[string]interface{}{}}
}

func (s *memoryStore) SaveData(_ context.Context, entity, key string, value interface{}) error {
	s.data[entity+"."+key] = value
	return nil
}

func (s *memoryStore) GetData(_ context.Context, entity, key string, result interface{}) error {
	value, ok := s.data[entity+"."+key]
	if !ok {
		return errors.New("missing")
	}
	ptr, ok := result.(*interface{})
	if !ok {
		return errors.New("result must be *interface{}")
	}
	*ptr = value
	return nil
}
