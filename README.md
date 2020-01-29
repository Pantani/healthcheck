# Services Health Check

Light application to run HTTP tests and evaluate expressions with results.

## How It Works?

The application parse all tests inside the test file. After that, it executes the test, saves the value inside the `Redis` database and evaluates the expression, if the expression is false or the request fails, the application going to create an incident inside the PargerDuty. 

## Setup

### Quick start

Deploy it in less than 30 seconds!

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/Pantani/healthcheck)

### Prerequisite
* [GO](https://golang.org/doc/install) `1.13+`
* Locally running [Redis](https://redis.io/topics/quickstart) or URL to remote instance (required for Observer only).

### From Source 

```shell
go get -u github.com/Pantani/healthcheck
cd $GOPATH/src/github.com/Pantani/healthcheck
```

### Make commands

```makefile
- install     Install missing dependencies. Runs `go get` internally. e.g.; make install get=github.com/foo/bar
- start       Start API in development mode.
- stop        Stop development mode.
- restart     Restart in development mode.
- compile     Compile the binary.
- exec        Run given command. e.g.; make exec run="go test ./..."
- clean       Clean build files. Runs `go clean` internally.
- test        Run all unit tests.
- fmt         Run `go fmt` for all go files.
- goreleaser  Release the last tag version with GoReleaser.
- govet       Run go vet.
- golint      Run golint.
```
  
### Environment Variables

All environment variables for developing are set inside the .env file.

```dotenv
REDIS_URL=Redis Database URL
PAGERDUTY_KEY=PargerDuty Access Key
PAGERDUTY_ESCALATION_POLICY=PagerDuty Escalation Policy ID
PAGERDUTY_SERVICE=PagerDuty Service ID
```

### Create Tests

To create a new test, do you need to edit the file `config/fixtures.json`. 

#### Test Structure
```
[
  {
    "namespace": string, # test namespace.
    "host": string, # the host to be reached.
    "tests": [ # Array of tests to be executed in the current namespace.
      {
        "name": string, # test name.
        "method": string, # HTTP method type.
        "url_path": string, # URL path to the test.
        "json_path": string, # path inside JSON to get the information to be saved.
        "body": any, # body for POST tests.
        "expression": string, # expression used to evaluate the test.
        "update_time": string # test update time.
      }
      ...
    ]
  },
  ...
]
```

#### Test Example

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
        "expression": "lastValue <= newValue",
        "update_time": "5s"
      },
      {
        "name": "host",
        "method": "GET",
        "url_path": "api",
        "json_path": "blockbook.host",
        "expression": "len(newValue) > 0",
        "update_time": "10s"
      }
    ]
  },
  {
    "namespace": "ethereum",
    "host": "https://infura.io",
    "tests": [
      {
        "name": "eth_getBlockByNumber",
        "method": "POST",
        "json_path": "result.hash",
        "body": {"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1b4", true],"id":1},
        "expression": "newValue matches lastValue",
        "update_time": "10s"
      }
    ]
  }
]
```

#### Expressions

We are using `expr` to evaluate the expressions. You can find the expression language [here](https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md)

#### Default Variables

The Default variables give you a convenient way to evaluate expressions in your tests. 
These variables are automatically set by the application:

| Variable       | Description                                   |
| :------------: | :-------------------------------------------- |
| lastValue      | The result value from the last test executed  | 
| newValue       | The result value from the current tests       | 


### Tools

-   Setup Redis:

```shell
brew install redis
```

-   Running in the IDE ( GoLand ):

1.  Run;
2.  Edit configuration;
3.  New Go build configuration;
4.  Select `directory` as configuration type;
5.  Set `metrics` as program argument and `-i` as Go tools argument; 

### Unit Tests

To run the unit tests: `make test`.
