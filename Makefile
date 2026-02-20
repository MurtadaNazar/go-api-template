.PHONY: help build run test lint clean install-deps release version

# Variables
BINARY_NAME=go-platform
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_FLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-deps: ## Install development dependencies
	go mod download
	go mod verify
	@command -v golangci-lint >/dev/null 2>&1 || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: ## Build the binary
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	CGO_ENABLED=0 go build $(GO_FLAGS) -o $(BINARY_NAME)
	@echo "✓ Binary created: ./$(BINARY_NAME)"

run: build ## Build and run the application
	./$(BINARY_NAME)

dev: ## Run in development mode (with hot reload)
	@command -v air >/dev/null 2>&1 || go install github.com/cosmtrek/air@latest
	air

test: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

test-quick: ## Run tests without coverage
	go test -v -race ./...

lint: ## Run linter
	@command -v golangci-lint >/dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run --timeout=5m

fmt: ## Format code
	go fmt ./...
	go mod tidy

clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean -cache -testcache

# Release targets (requires git tags)
release-patch: ## Create patch release (e.g., v1.0.0 -> v1.0.1)
	@bash scripts/release.sh patch

release-minor: ## Create minor release (e.g., v1.0.0 -> v1.1.0)
	@bash scripts/release.sh minor

release-major: ## Create major release (e.g., v1.0.0 -> v2.0.0)
	@bash scripts/release.sh major

version: ## Show current version
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Binary: $(BINARY_NAME)"

# Cross-platform builds
build-linux-amd64: ## Build for Linux amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(GO_FLAGS) -o $(BINARY_NAME)-linux-amd64

build-linux-arm64: ## Build for Linux arm64
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(GO_FLAGS) -o $(BINARY_NAME)-linux-arm64

build-darwin-amd64: ## Build for macOS amd64
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(GO_FLAGS) -o $(BINARY_NAME)-darwin-amd64

build-darwin-arm64: ## Build for macOS arm64 (Apple Silicon)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(GO_FLAGS) -o $(BINARY_NAME)-darwin-arm64

build-windows-amd64: ## Build for Windows amd64
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(GO_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe

build-all: build-linux-amd64 build-linux-arm64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64 ## Build for all platforms
	@echo "✓ All binaries built"
	@ls -lh $(BINARY_NAME)-*

# Code quality
check: fmt lint test ## Run all checks (format, lint, test)
	@echo "✓ All checks passed"

.DEFAULT_GOAL := help
