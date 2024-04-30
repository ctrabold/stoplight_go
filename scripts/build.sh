#!/bin/bash

set -eu

golangci-lint run
go test -v -race . # Run all the tests with the race detector enabled
diff -u <(echo -n) <(gofmt -d .)
go vet $(go list . | grep -v /vendor/)

rice embed-go
go build -o bin/stoplight .
