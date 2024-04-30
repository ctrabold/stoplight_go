#!/bin/bash

set -eu

# Linter is disabled at the moment (shows linting issues but does not fail the build)
golangci-lint run .
go test -v -race . # Run all the tests with the race detector enabled
go fmt .
# vet is a tool for static analysis of Go programs.
go vet $(go list . | grep -v /vendor/)
