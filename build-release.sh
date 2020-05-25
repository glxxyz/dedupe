#!/bin/zsh

# exit on error
set -e

# run all unit tests
go test ./...
go test github.com/glxxyz/dedupe/param
go test github.com/glxxyz/dedupe/repo

# build on all the most common platforms
env GOOS=darwin GOARCH=amd64 go build -o bin/macos-amd64/dedupe ./...
env GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/dedupe ./...
env GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/dedupe.exe ./...

# run integration tests
./bb-test-integration.sh
