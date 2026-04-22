# GoLDAP Server Makefile

APP_NAME := goldap-server
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-w -s -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

OUTPUT_DIR := bin

.PHONY: all build build-linux docker run clean test lint help

all: build

build:
	@mkdir -p $(OUTPUT_DIR)
	go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) .

build-linux:
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME) .
install: build-linux

build-all:
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-linux-amd64 .
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-linux-arm64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-amd64 .
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(OUTPUT_DIR)/$(APP_NAME)-darwin-arm64 .

docker: build-linux
	docker build -t $(APP_NAME):$(VERSION) -t $(APP_NAME):latest .

run:
	go run .

test:
	go test -v -race ./...

lint:
	golangci-lint run --timeout=5m ./...

clean:
	rm -rf $(OUTPUT_DIR)
	go clean -cache

help:
	@echo "GoLDAP Server Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make build       Build for current platform"
	@echo "  make build-linux Build for Linux amd64"
	@echo "  make build-all   Build for all platforms"
	@echo "  make docker      Build Docker image"
	@echo "  make run         Run in development mode"
	@echo "  make test        Run tests"
	@echo "  make lint        Run linter"
	@echo "  make clean       Clean build artifacts"
	@echo ""
	@echo "Version: $(VERSION)"
