# Healthcheck Audit and Modernization Plan

Date: 2026-06-10

## Scope

This audit reviewed the application as a small scheduled health-check worker:

- CLI and runtime configuration
- HTTP request execution
- Fixture parsing and validation
- Redis state storage
- PagerDuty alerting
- Dependency freshness
- Unit-test coverage
- Deployment and operator documentation
- Harness readiness for future maintenance

## Initial Findings

| Area | Finding | Risk | Status |
| --- | --- | --- | --- |
| Go toolchain | Module targeted Go 1.13. | Modern dependencies and tooling drift. | Updated to Go 1.26. |
| Dependency graph | `trustwallet/blockatlas` pulled a large legacy graph for logger/errors/Redis helpers. | Old transitive dependencies and harder maintenance. | Removed. |
| Redis | State storage was hidden behind an old project-specific wrapper. | Difficult to test and update. | Replaced with `go-redis/v9`. |
| HTTP client | Non-2xx responses were treated as successful checks. | Outages could be missed when an endpoint returned 500 with JSON. | Fixed. |
| HTTP body | JSON bodies were sent without `Content-Type`. | POST checks could be rejected by strict APIs. | Fixed. |
| Config | Startup logged `PAGERDUTY_KEY` directly. | Secret leakage in logs. | Fixed. |
| Config | Fixture path and HTTP timeout were hardcoded. | Harder deploy/runtime control. | Added env vars and CLI flags. |
| Fixtures | Loader relied on implicit working directory and did not validate schema. | Runtime failures inside scheduled jobs. | Added parser normalization and validation. |
| Tests | Existing fixture tests compared raw byte arrays and failed on the current checkout. | Low-signal tests and failing baseline. | Replaced with behavior tests. |
| Collector | Core logic had no focused tests. | High regression risk. | Added unit coverage with fakes. |
| Shutdown | Scheduler blocked forever on an untracked channel. | Poor shutdown behavior in workers. | Added context and signal-aware shutdown. |
| Deployment | Procfile used `web` even though the service does not listen on `$PORT`. | Heroku/router mismatch. | Changed to `worker`. |
| Example fixtures | Block-height expressions used equality or decreasing checks. | Normal height increases could create noisy incidents. | Changed examples to `lastValue <= newValue`. |

## Implemented Changes

### Runtime

- Kept the `metrics` command as the main entry point.
- Added `SIGINT`/`SIGTERM` handling through `context.Context`.
- Added runtime validation before the scheduler starts.
- Added CLI overrides:
  - `--fixtures`
  - `--redis-url`
  - `--http-timeout`

### Dependencies

- Updated direct dependencies to current stable releases:
  - `github.com/PagerDuty/go-pagerduty v1.8.0`
  - `github.com/expr-lang/expr v1.17.8`
  - `github.com/joho/godotenv v1.5.1`
  - `github.com/redis/go-redis/v9 v9.20.0`
  - `github.com/robfig/cron/v3 v3.0.1`
  - `github.com/spf13/cobra v1.10.2`
  - `github.com/tidwall/gjson v1.19.0`
- Removed:
  - `github.com/trustwallet/blockatlas`
  - `github.com/spf13/viper`
  - old Redis wrappers and legacy indirect modules from that graph
- Direct retained modules have no pending `go list -m -u` updates after the bump.
- Some modules still appear in `go list -m -u all` only through dependency/test graphs outside this module's direct imports. One example is `github.com/imdario/mergo`, whose newer release declares the module path `dario.cat/mergo`, so it cannot be applied as a plain transitive bump from this module.

### Tests

- Added focused tests for:
  - HTTP URL construction and execution
  - non-2xx handling
  - JSON body encoding
  - config loading and validation
  - fixture parsing and validation
  - Redis storage using `miniredis`
  - collector pass/fail/request-error/missing-path flows
  - PagerDuty constructor validation

## Current Behavior Contract

- Fixtures must use HTTP or HTTPS hosts.
- Check paths must be relative. Absolute URLs in `url_path` are rejected.
- Methods are normalized to uppercase and must be one of the supported HTTP verbs.
- `update_time` must be a positive Go duration.
- The first run of a check evaluates with `lastValue = 0` when Redis has no stored value.
- A request error or non-2xx status creates a PagerDuty incident.
- A missing `json_path` logs an error and skips storage/evaluation.
- A false expression creates a PagerDuty incident.

## Follow-Up Plan

1. Add an optional dry-run alert sender for local fixture testing without PagerDuty side effects.
2. Add duplicate fixture detection for repeated `namespace` + `name` pairs.
3. Add incident de-duplication or recovery behavior if PagerDuty noise becomes an issue.
4. Add structured JSON logging configuration if this is deployed into log aggregation.
5. Add CI with `go test ./...`, `go vet ./...`, and `go build ./...`.
6. Consider a lightweight `/healthz` HTTP server only if the deployment platform requires a web process.
