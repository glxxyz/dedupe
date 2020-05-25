#!/bin/bash

# exit on error
set -e

# run all unit tests
go test ./...
go test github.com/glxxyz/dedupe/param
go test github.com/glxxyz/dedupe/repo

env GOOS=darwin GOARCH=amd64 go build -o bin/macos-amd64/dedupe *.go
env GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/dedupe *.go
env GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/dedupe.exe *.go
