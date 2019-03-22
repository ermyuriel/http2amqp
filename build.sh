#!/usr/bin/env bash
set -xe
go get -u -v "github.com/eucj/amqphelper"
GOOS=linux GOARCH=amd64 go build -o bin/application -ldflags="-s -w"