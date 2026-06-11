# Services Health Check

Small Go service that runs scheduled HTTP checks, extracts values from JSON responses, evaluates expressions against the previous value, and opens PagerDuty incidents when a check fails.

## What It Does

Each fixture defines a namespace, a host, and one or more checks. For every check the service:

1. Sends an HTTP request to the configured host and relative path.
2. Fails the check when the request errors or returns a non-2xx status.
3. Extracts a JSON value using `gjson`.
4. Reads the previous value for the same namespace/check from Redis.
5. Stores the new value in Redis.
6. Evaluates the configured `expr` expression with `lastValue` and `newValue`.
7. Creates a PagerDuty incident when the expression returns `false`.

## Requirements

- Go 1.26+
- Redis
- PagerDuty REST API key, service ID, and escalation policy ID

## Quick Start

```sh
cp configs/fixtures.json /tmp/healthcheck-fixtures.json
export REDIS_URL=redis://localhost:6379/0
export PAGERDUTY_KEY=your-api-key
export PAGERDUTY_SERVICE=your-service-id
export PAGERDUTY_ESCALATION_POLICY=your-escalation-policy-id
go run . metrics --fixtures /tmp/healthcheck-fixtures.json
```

The service keeps running until it receives `SIGINT` or `SIGTERM`.

## Configuration

Environment variables:

| Name | Default | Description |
| --- | --- | --- |
| `REDIS_URL` | `redis://localhost:6379/0` | Redis URL used to store previous check values. |
| `PAGERDUTY_KEY` | empty | PagerDuty REST API key. Required at runtime. |
| `PAGERDUTY_SERVICE` | empty | PagerDuty service ID. Required at runtime. |
| `PAGERDUTY_ESCALATION_POLICY` | empty | PagerDuty escalation policy ID. Required at runtime. |
| `HEALTHCHECK_FIXTURES_FILE` | `configs/fixtures.json` | Fixture file path. |
| `HEALTHCHECK_HTTP_TIMEOUT` | `15s` | HTTP request timeout as a Go duration. |

CLI flags override environment values:

```sh
healthcheck metrics \
  --fixtures configs/fixtures.json \
  --redis-url redis://localhost:6379/0 \
  --http-timeout 10s
```

## Fixture Format

```json
[
  {
    "namespace": "bitcoin",
    "host": "https://btc1.trezor.io",
    "tests": [
      {
        "name": "block_height",
        "method": "GET",
        "url_path": "api",
        "json_path": "blockbook.bestHeight",
        "body": {},
        "expression": "lastValue <= newValue",
        "update_time": "30s"
      }
    ]
  }
]
```

Fields:

| Field | Required | Description |
| --- | --- | --- |
| `namespace` | yes | Logical group for checks and Redis storage. |
| `host` | yes | HTTP or HTTPS base URL. |
| `tests[].name` | yes | Check name, unique within a namespace. |
| `tests[].method` | yes | `GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `HEAD`, or `OPTIONS`. |
| `tests[].url_path` | no | Relative path appended to `host`; absolute URLs are rejected. |
| `tests[].json_path` | yes | `gjson` path used to extract the response value. |
| `tests[].body` | no | JSON request body. When present, `Content-Type: application/json` is sent. |
| `tests[].expression` | yes | `expr` expression that must return a boolean. |
| `tests[].update_time` | yes | Go duration used by the scheduler, for example `10s`, `1m`, or `5m`. |

Expressions receive:

| Variable | Description |
| --- | --- |
| `lastValue` | Value stored from the previous run. Defaults to `0` when Redis has no value yet. |
| `newValue` | Value extracted from the current response. |

## Development

```sh
make test
make vet
make build
make check
```

Useful docs:

- [Operations runbook](docs/OPERATIONS.md)
- [Audit and modernization plan](docs/AUDIT.md)
