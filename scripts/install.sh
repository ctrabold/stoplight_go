#!/bin/bash

set -eu

go install github.com/GeertJohan/go.rice/rice@latest
go install -v .