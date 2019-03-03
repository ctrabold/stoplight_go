#!/bin/bash

set -eu

go test -v -race ./...  # Run all the tests with the race detector enabled
diff -u <(echo -n) <(gofmt -d .)
go vet $(go list ./... | grep -v /vendor/)

rice embed-go
gox -os="linux darwin windows" \
    -arch="amd64" \
    -ldflags "-X main.Rev=`git rev-parse --short HEAD`" \
    -verbose \
    -output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}" \
    ./...
