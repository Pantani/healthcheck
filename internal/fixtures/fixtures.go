package fixtures

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DefaultPath = "configs/fixtures.json"
)

func GeFixtures() (testFixtures Fixtures, err error) {
	return GetFixtures(DefaultPath)
}

func GetFixtures(path string) (Fixtures, error) {
	var testFixtures Fixtures
	if err := loadFixtures(path, &testFixtures); err != nil {
		return nil, err
	}
	testFixtures.Normalize()
	if err := testFixtures.Validate(); err != nil {
		return nil, err
	}
	return testFixtures, nil
}

func loadFixtures(path string, result interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read fixtures %q: %w", path, err)
	}
	if err := json.Unmarshal(b, result); err != nil {
		return fmt.Errorf("parse fixtures %q: %w", path, err)
	}
	return nil
}

func (fxs Fixtures) Normalize() {
	for i := range fxs {
		fxs[i].Namespace = strings.TrimSpace(fxs[i].Namespace)
		fxs[i].Host = strings.TrimRight(strings.TrimSpace(fxs[i].Host), "/")
		for j := range fxs[i].Tests {
			fxs[i].Tests[j].Name = strings.TrimSpace(fxs[i].Tests[j].Name)
			fxs[i].Tests[j].Method = strings.ToUpper(strings.TrimSpace(fxs[i].Tests[j].Method))
			fxs[i].Tests[j].URLPath = strings.TrimLeft(strings.TrimSpace(fxs[i].Tests[j].URLPath), "/")
			fxs[i].Tests[j].JSONPath = strings.TrimSpace(fxs[i].Tests[j].JSONPath)
			fxs[i].Tests[j].Expression = strings.TrimSpace(fxs[i].Tests[j].Expression)
			fxs[i].Tests[j].UpdateTime = strings.TrimSpace(fxs[i].Tests[j].UpdateTime)
		}
	}
}

func (fxs Fixtures) Validate() error {
	if len(fxs) == 0 {
		return fmt.Errorf("fixtures file does not contain any namespace")
	}
	for i, fixture := range fxs {
		if fixture.Namespace == "" {
			return fmt.Errorf("fixture[%d].namespace is required", i)
		}
		if err := validateHost(fixture.Host); err != nil {
			return fmt.Errorf("fixture[%d].host: %w", i, err)
		}
		if len(fixture.Tests) == 0 {
			return fmt.Errorf("fixture[%d].tests must contain at least one test", i)
		}
		for j, test := range fixture.Tests {
			if err := test.Validate(); err != nil {
				return fmt.Errorf("fixture[%d].tests[%d]: %w", i, j, err)
			}
		}
	}
	return nil
}

func (fxs Fixtures) CheckCount() int {
	total := 0
	for _, fixture := range fxs {
		total += len(fixture.Tests)
	}
	return total
}

func (t Test) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("name is required")
	}
	if !validMethod(t.Method) {
		return fmt.Errorf("unsupported method %q", t.Method)
	}
	if t.JSONPath == "" {
		return fmt.Errorf("json_path is required")
	}
	if t.Expression == "" {
		return fmt.Errorf("expression is required")
	}
	duration, err := time.ParseDuration(t.UpdateTime)
	if err != nil {
		return fmt.Errorf("invalid update_time %q: %w", t.UpdateTime, err)
	}
	if duration <= 0 {
		return fmt.Errorf("update_time must be greater than zero")
	}
	return nil
}

func validateHost(raw string) error {
	if raw == "" {
		return fmt.Errorf("is required")
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return err
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return fmt.Errorf("must use http or https")
	}
	if parsed.Host == "" {
		return fmt.Errorf("must include a host")
	}
	return nil
}

func validMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodDelete, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}
