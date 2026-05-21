#!/bin/bash

dest=./dest
rm -rf $dest

GOARCH=amd64 GOOS=linux go build -ldflags "-s -w" -o $dest/notifier cmd/notifier/main.go
