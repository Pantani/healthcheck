package config

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

const (
	DefaultFixturesPath = "configs/fixtures.json"
	DefaultRedisURL     = "redis://localhost:6379/0"
	DefaultHTTPTimeout  = 15 * time.Second
)

type Redis struct {
	URL string
}

type PagerDuty struct {
	Key              string
	Service          string
	EscalationPolicy string
}

type Fixtures struct {
	Path string
}

type HTTP struct {
	Timeout time.Duration
}

type Configuration struct {
	Redis     Redis
	PagerDuty PagerDuty
	Fixtures  Fixtures
	HTTP      HTTP
}

var Active = Default()

func Default() Configuration {
	return Configuration{
		Redis: Redis{
			URL: DefaultRedisURL,
		},
		Fixtures: Fixtures{
			Path: DefaultFixturesPath,
		},
		HTTP: HTTP{
			Timeout: DefaultHTTPTimeout,
		},
	}
}

func Load() (Configuration, error) {
	cfg := Default()
	cfg.Redis.URL = env("REDIS_URL", cfg.Redis.URL)
	cfg.PagerDuty.Key = env("PAGERDUTY_KEY", cfg.PagerDuty.Key)
	cfg.PagerDuty.Service = env("PAGERDUTY_SERVICE", cfg.PagerDuty.Service)
	cfg.PagerDuty.EscalationPolicy = env("PAGERDUTY_ESCALATION_POLICY", cfg.PagerDuty.EscalationPolicy)
	cfg.Fixtures.Path = env("HEALTHCHECK_FIXTURES_FILE", cfg.Fixtures.Path)

	timeout := env("HEALTHCHECK_HTTP_TIMEOUT", "")
	if timeout != "" {
		parsed, err := time.ParseDuration(timeout)
		if err != nil {
			return Configuration{}, fmt.Errorf("invalid HEALTHCHECK_HTTP_TIMEOUT %q: %w", timeout, err)
		}
		cfg.HTTP.Timeout = parsed
	}
	return cfg, nil
}

func (cfg Configuration) ValidateRuntime() error {
	if cfg.Redis.URL == "" {
		return fmt.Errorf("REDIS_URL is required")
	}
	if cfg.PagerDuty.Key == "" {
		return fmt.Errorf("PAGERDUTY_KEY is required")
	}
	if cfg.PagerDuty.Service == "" {
		return fmt.Errorf("PAGERDUTY_SERVICE is required")
	}
	if cfg.PagerDuty.EscalationPolicy == "" {
		return fmt.Errorf("PAGERDUTY_ESCALATION_POLICY is required")
	}
	if cfg.Fixtures.Path == "" {
		return fmt.Errorf("HEALTHCHECK_FIXTURES_FILE cannot be empty")
	}
	if cfg.HTTP.Timeout <= 0 {
		return fmt.Errorf("HEALTHCHECK_HTTP_TIMEOUT must be greater than zero")
	}
	return nil
}

func Apply(cfg Configuration) {
	Active = cfg
}

func InitConfig() {
	cfg, err := Load()
	if err != nil {
		log.Printf("configuration error: %v", err)
	}
	Apply(cfg)
	log.Printf("REDIS_URL: %s", Active.Redis.URL)
	log.Printf("PAGERDUTY_KEY: %s", secretState(Active.PagerDuty.Key))
	log.Printf("PAGERDUTY_SERVICE: %s", Active.PagerDuty.Service)
	log.Printf("PAGERDUTY_ESCALATION_POLICY: %s", Active.PagerDuty.EscalationPolicy)
	log.Printf("HEALTHCHECK_FIXTURES_FILE: %s", Active.Fixtures.Path)
	log.Printf("HEALTHCHECK_HTTP_TIMEOUT: %s", Active.HTTP.Timeout)
}

func env(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func secretState(value string) string {
	if value == "" {
		return "<unset>"
	}
	return "<set>"
}
