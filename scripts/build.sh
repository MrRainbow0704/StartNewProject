#!/usr/bin/env sh

go mod tidy
go build -o ./bin/start-new-project.exe ./cmd