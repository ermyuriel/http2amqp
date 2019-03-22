#!/usr/bin/env bash
set -xe
go mod download
GOOS=linux GOARCH=amd64 go build -o bin/application -ldflags="-s -w"