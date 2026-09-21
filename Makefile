VERSION=0.1.2
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION}"
all: mackerel-plugin-linux-usage

.PHONY: mackerel-plugin-linux-usage linux check lint

mackerel-plugin-linux-usage: *.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-usage

linux: *.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-usage

check:
	go test -v ./...
	go test -race ./...

lint:
	golangci-lint run --timeout 5m ./...