.PHONY: build test clean run deps verify update-deps dev dev-down help

# Build variables
BINARY_NAME=go-platform-template
MAIN_FILE=cmd/server/main.go
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GORUN=$(GOCMD) run

# Docker compose file
DC_FILE=docker-compose/docker-compose.yml

# Build flags
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

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

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	$(GORUN) $(MAIN_FILE)

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

# Start development environment with Docker
dev:
	@echo "Starting development environment..."
	docker compose -f $(DC_FILE) up --build

# Start development environment in background
dev-d:
	@echo "Starting development environment in background..."
	docker compose -f $(DC_FILE) up -d --build
	@echo "✓ Development environment started"

# Stop development environment
dev-down:
	@echo "Stopping development environment..."
	docker compose -f $(DC_FILE) down
	@echo "✓ Development environment stopped"

# View logs
dev-logs:
	docker compose -f $(DC_FILE) logs -f

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

# Help command
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build targets:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application locally"
	@echo "  clean          - Clean build artifacts"
	@echo ""
	@echo "Testing targets:"
	@echo "  test           - Run all tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo ""
	@echo "Dependency targets:"
	@echo "  deps           - Download dependencies"
	@echo "  verify         - Verify dependencies"
	@echo "  update-deps    - Update dependencies (tidy)"
	@echo ""
	@echo "Development targets:"
	@echo "  dev            - Start development environment (foreground)"
	@echo "  dev-d          - Start development environment (background)"
	@echo "  dev-down       - Stop development environment"
	@echo "  dev-logs       - View development environment logs"
	@echo ""
	@echo "Code quality targets:"
	@echo "  fmt            - Format code"
	@echo "  vet            - Run go vet"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  security       - Run security checks (requires gosec)"
	@echo ""
	@echo "Other:"
	@echo "  help           - Show this help message"

.DEFAULT_GOAL := help
