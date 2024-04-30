#!/bin/bash

set -eu

# Run linter if available
if [ -x "$(command -v golangci-lint)" ]; then
    golangci-lint run .
else
    echo "INFO golangci-lint not available. Skipping..."
fi
go test -v -race . # Run all the tests with the race detector enabled
go fmt .
# vet is a tool for static analysis of Go programs.
go vet $(go list . | grep -v /vendor/)
