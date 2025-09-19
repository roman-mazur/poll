#!/usr/bin/env bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o pollsvc-x86_64 ./cmd/pollsvc
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o pollsvc-arm64 ./cmd/pollsvc
