#/usr/bin/env sh sh

go mod download && go mod verify
go build -o ./bin/ ./cmd/...