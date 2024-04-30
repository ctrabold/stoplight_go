#!/bin/bash

set -eu

rice embed-go
go build -o bin/stoplight .
