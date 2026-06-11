BINARY := healthcheck
GOFILES := $(shell find . -name '*.go' -not -path './vendor/*')

.PHONY: build clean fmt help run test tidy vet check

build:
	go build -o bin/$(BINARY) .

clean:
	rm -rf bin
	go clean

fmt:
	gofmt -w $(GOFILES)

run:
	go run . metrics

test:
	go test ./...

tidy:
	go mod tidy

vet:
	go vet ./...

check: fmt tidy test vet build

help:
	@echo "Available targets:"
	@echo "  build  Build bin/$(BINARY)"
	@echo "  clean  Remove build output and clean Go cache metadata"
	@echo "  fmt    Format Go files"
	@echo "  run    Run the metrics scheduler locally"
	@echo "  test   Run unit tests"
	@echo "  tidy   Tidy go.mod and go.sum"
	@echo "  vet    Run go vet"
	@echo "  check  Run fmt, tidy, test, vet, and build"
