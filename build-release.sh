#!/bin/zsh

# exit on error
set -e

echo "running all unit tests"
go test ./src/...
go test github.com/glxxyz/dedupe/param
go test github.com/glxxyz/dedupe/repo

# build on all the most common platforms
env GOOS=darwin GOARCH=amd64 go build -o bin/macos-amd64/dedupe ./src/...
env GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/dedupe ./src/...
env GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/dedupe.exe ./src/...

# run integration tests
./bb-test-integration.sh
