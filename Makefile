VERSION=0.1.3
GITCOMMIT?=$(shell git describe --dirty --always)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"
all: mackerel-plugin-linux-usage

.PHONY: mackerel-plugin-linux-usage

mackerel-plugin-linux-usage: main.go
	go build $(LDFLAGS) -o mackerel-plugin-linux-usage

linux: main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o mackerel-plugin-linux-usage

check:
	go test -v ./...
	go test -race ./...

lint:
	golangci-lint run --timeout 5m ./...