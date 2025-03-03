#!/usr/bin/env sh

go mod tidy
go build -o ./bin/main ./cmd