.PHONY: help build test clean dev dev-d dev-down dev-logs deps verify update-deps fmt vet lint security test-coverage

# Build variables
BINARY_NAME={{.ProjectName}}
MAIN_FILE=cmd/server/main.go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Container compose file
DC_FILE={{.ComposeFile}}

# Build flags
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# Help command (default)
help:
	@echo "{{.ProjectName}} - Development Commands"
	@echo ""
	@echo "BUILD:"
	@echo "  make build          - Build the application"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "TESTING:"
	@echo "  make test           - Run all tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "DEVELOPMENT:"
	@echo "  make dev            - Start dev environment (foreground)"
	@echo "  make dev-d          - Start dev environment (background)"
	@echo "  make dev-down       - Stop dev environment"
	@echo "  make dev-logs       - View dev environment logs"
	@echo ""
	@echo "CODE QUALITY:"
	@echo "  make fmt            - Format code (gofmt)"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make security       - Run gosec security checks"
	@echo ""
	@echo "DEPENDENCIES:"
	@echo "  make deps           - Download dependencies"
	@echo "  make verify         - Verify dependencies"
	@echo "  make update-deps    - Update/tidy dependencies"

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	CGO_ENABLED=0 $(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "✓ Build successful: ./$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "✓ Clean complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	@echo "✓ Dependencies downloaded"

# Verify dependencies
verify:
	@echo "Verifying dependencies..."
	$(GOMOD) verify
	@echo "✓ Dependencies verified"

# Update dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GOMOD) tidy
	@echo "✓ Dependencies updated"

# Start development environment
dev:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	@echo "Starting development environment..."
	{{.ContainerCmd}} compose --env-file .env -f $(DC_FILE) up --build

# Start development environment in background
dev-d:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy
	@echo "Starting development environment in background..."
	{{.ContainerCmd}} compose --env-file .env -f $(DC_FILE) up -d --build
	@echo "✓ Development environment started"

# Stop development environment
dev-down:
	@echo "Stopping development environment..."
	{{.ContainerCmd}} compose --env-file .env -f $(DC_FILE) down
	@echo "✓ Development environment stopped"

# View logs
dev-logs:
	{{.ContainerCmd}} compose --env-file .env -f $(DC_FILE) logs -f

# Lint code (requires golangci-lint)
lint:
	@echo "Running linter..."
	golangci-lint run ./...

# Security check (requires gosec)
security:
	@echo "Running security checks..."
	gosec -exclude-generated ./...

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...
	@echo "✓ Code formatted"

# Run vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...
	@echo "✓ Go vet passed"

.DEFAULT_GOAL := help
