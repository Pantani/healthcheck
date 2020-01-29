# TrustWallet Services Health Check

Light application to run http tests and evaluate expressions with results.

## Setup

### Quick start

Deploy it in less than 30 seconds!

### Prerequisite
* [GO](https://golang.org/doc/install) `1.13+`
* Locally running [Redis](https://redis.io/topics/quickstart) or url to remote instance (required for Observer only)

### From Source 

```shell
go get -u github.com/Pantani/healthcheck
cd $GOPATH/src/github.com/Pantani/healthcheck

// Make commands
- install     Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
- start       Start API in development mode.
- stop        Stop development mode.
- restart     Restart in development mode.
- compile     Compile the binary.
- exec        Run given command. e.g; make exec run="go test ./..."
- clean       Clean build files. Runs `go clean` internally.
- test        Run all unit tests.
- fmt         Run `go fmt` for all go files.
- goreleaser  Release the last tag version with GoReleaser.
- govet       Run go vet.
- golint      Run golint.
```
  
### Environment Variables

All environment variables for developing are set inside the .env file.

### Expressions

We are using `expr` to evaluate the expressions. You can find the expression language here:
https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md

### Tools

-   Setup Redis

```shell
brew install redis
```

-   Running in the IDE ( GoLand )

1.  Run
2.  Edit configuration
3.  New Go build configuration
4.  Select `directory` as configuration type
5.  Set `metrics` as program argument and `-i` as Go tools argument 

### Unit Tests

To run the unit tests: `make test`
