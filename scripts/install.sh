#!/bin/bash

set -eu

go get github.com/GeertJohan/go.rice
go get github.com/GeertJohan/go.rice/rice
go get github.com/mitchellh/gox
go get -t -v ./...