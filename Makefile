.PHONY: all build test lint clean install release

BINARY := cursor-tool
MODULE := github.com/WarnetBes/cursor-tool
CMD_PATH := ./cmd/cursor-tool
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X $(MODULE)/cmd/cursor-tool/commands.Version=$(VERSION) -w -s"

all: build

build:
	go build $(LDFLAGS) -o bin/$(BINARY) $(CMD_PATH)

build-all:
	GOOS=linux   GOARCH=amd64  go build $(LDFLAGS) -o bin/$(BINARY)_linux_amd64   $(CMD_PATH)
	GOOS=linux   GOARCH=arm64  go build $(LDFLAGS) -o bin/$(BINARY)_linux_arm64   $(CMD_PATH)
	GOOS=darwin  GOARCH=amd64  go build $(LDFLAGS) -o bin/$(BINARY)_darwin_amd64  $(CMD_PATH)
	GOOS=darwin  GOARCH=arm64  go build $(LDFLAGS) -o bin/$(BINARY)_darwin_arm64  $(CMD_PATH)
	GOOS=windows GOARCH=amd64  go build $(LDFLAGS) -o bin/$(BINARY)_windows_amd64.exe $(CMD_PATH)

test:
	go test ./... -v -race -cover

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

install: build
	cp bin/$(BINARY) /usr/local/bin/$(BINARY)

run: build
	./bin/$(BINARY)

tidy:
	go mod tidy

vet:
	go vet ./...
