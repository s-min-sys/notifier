#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o notifier cmd/notifier/main.go
