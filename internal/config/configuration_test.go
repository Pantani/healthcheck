package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	clearEnv(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Redis.URL != DefaultRedisURL {
		t.Errorf("Redis.URL = %q, want %q", cfg.Redis.URL, DefaultRedisURL)
	}
	if cfg.Fixtures.Path != DefaultFixturesPath {
		t.Errorf("Fixtures.Path = %q, want %q", cfg.Fixtures.Path, DefaultFixturesPath)
	}
	if cfg.HTTP.Timeout != DefaultHTTPTimeout {
		t.Errorf("HTTP.Timeout = %s, want %s", cfg.HTTP.Timeout, DefaultHTTPTimeout)
	}
}

func TestLoadFromEnvironment(t *testing.T) {
	clearEnv(t)
	t.Setenv("REDIS_URL", "redis://redis:6379/1")
	t.Setenv("PAGERDUTY_KEY", "key")
	t.Setenv("PAGERDUTY_SERVICE", "service")
	t.Setenv("PAGERDUTY_ESCALATION_POLICY", "policy")
	t.Setenv("HEALTHCHECK_FIXTURES_FILE", "fixtures.json")
	t.Setenv("HEALTHCHECK_HTTP_TIMEOUT", "3s")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Redis.URL != "redis://redis:6379/1" {
		t.Errorf("Redis.URL = %q", cfg.Redis.URL)
	}
	if cfg.PagerDuty.Key != "key" {
		t.Errorf("PagerDuty.Key = %q", cfg.PagerDuty.Key)
	}
	if cfg.HTTP.Timeout != 3*time.Second {
		t.Errorf("HTTP.Timeout = %s", cfg.HTTP.Timeout)
	}
}

func TestLoadRejectsInvalidDuration(t *testing.T) {
	clearEnv(t)
	t.Setenv("HEALTHCHECK_HTTP_TIMEOUT", "not-a-duration")

	if _, err := Load(); err == nil {
		t.Fatal("Load() error = nil, want invalid duration")
	}
}

func TestValidateRuntime(t *testing.T) {
	cfg := Default()
	if err := cfg.ValidateRuntime(); err == nil {
		t.Fatal("ValidateRuntime() error = nil, want missing PagerDuty config")
	}

	cfg.PagerDuty.Key = "key"
	cfg.PagerDuty.Service = "service"
	cfg.PagerDuty.EscalationPolicy = "policy"
	if err := cfg.ValidateRuntime(); err != nil {
		t.Fatalf("ValidateRuntime() error = %v", err)
	}
}

func clearEnv(t *testing.T) {
	t.Helper()
	t.Setenv("REDIS_URL", "")
	t.Setenv("PAGERDUTY_KEY", "")
	t.Setenv("PAGERDUTY_SERVICE", "")
	t.Setenv("PAGERDUTY_ESCALATION_POLICY", "")
	t.Setenv("HEALTHCHECK_FIXTURES_FILE", "")
	t.Setenv("HEALTHCHECK_HTTP_TIMEOUT", "")
}
