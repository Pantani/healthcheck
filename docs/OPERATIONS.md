# Operations Runbook

## Local Run

Start Redis:

```sh
redis-server
```

Run the worker:

```sh
export REDIS_URL=redis://localhost:6379/0
export PAGERDUTY_KEY=your-api-key
export PAGERDUTY_SERVICE=your-service-id
export PAGERDUTY_ESCALATION_POLICY=your-escalation-policy-id
go run . metrics --fixtures configs/fixtures.json
```

Stop it with `Ctrl-C`.

## Runtime Settings

| Setting | How to set | Notes |
| --- | --- | --- |
| Fixture file | `HEALTHCHECK_FIXTURES_FILE` or `--fixtures` | CLI flag wins over env. |
| Redis URL | `REDIS_URL` or `--redis-url` | Supports full Redis URLs or `host:port`. |
| HTTP timeout | `HEALTHCHECK_HTTP_TIMEOUT` or `--http-timeout` | Use Go duration syntax such as `5s` or `1m`. |
| PagerDuty key | `PAGERDUTY_KEY` | Required. Logged only as `<set>` or `<unset>`. |
| PagerDuty service | `PAGERDUTY_SERVICE` | Required. |
| PagerDuty escalation policy | `PAGERDUTY_ESCALATION_POLICY` | Required. |

## Deployment Notes

This application is a worker process, not an HTTP server. On Heroku-style platforms it should run as:

```Procfile
worker: bin/healthcheck metrics
```

If the platform requires a web process that binds `$PORT`, add a separate lightweight health endpoint rather than running the scheduler as a fake web dyno.

## Check Triage

When an incident is created:

1. Confirm the target endpoint is reachable outside the worker.
2. Check whether the endpoint returned a non-2xx status.
3. Verify the configured `json_path` still exists in the response.
4. Compare `lastValue` in Redis with the new response value.
5. Re-run the expression manually with those two values.

## Redis State

Values are stored in Redis hashes:

- Hash key: `namespace`
- Field: check `name`
- Value: JSON-encoded extracted value

Example:

```sh
redis-cli HGET bitcoin block_height
```

## Development Verification

Run the full local gate:

```sh
make check
```

This formats code, tidies modules, runs unit tests, runs `go vet`, and builds the binary.
