GOOS=darwin
GOARCH=amd64
go build -o bin/macos-amd64/dedupe *.go

GOOS=linux
GOARCH=amd64
go build -o bin/linux-amd64/dedupe *.go

GOOS=windows
GOARCH=amd64
go build -o bin/windows-amd64/dedupe.exe *.go
