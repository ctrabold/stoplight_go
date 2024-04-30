#!/bin/bash

set -eu

# Linter is disabled at the moment (shows linting issues but does not fail the build)
golangci-lint run ||true
go test -v -race . # Run all the tests with the race detector enabled
diff -u <(echo -n) <(gofmt -d .)
# vet is a tool for static analysis of Go programs.
go vet $(go list . | grep -v /vendor/)
