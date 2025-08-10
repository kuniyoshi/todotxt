# todo.txt - Simple and powerful task management
# Makefile for building, testing, and installing

# Variables
BINARY_NAME=todo
BINARY_DIR=bin
DIST_DIR=dist
VERSION?=v1.0.0
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Platform-specific settings
ifdef WINDOWS
    BINARY_NAME=todo.exe
endif

# Default target
.PHONY: all
all: build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME) .
	@echo "Build complete: $(BINARY_DIR)/$(BINARY_NAME)"

# Build for current platform (simpler version)
.PHONY: compile
compile:
	$(GOBUILD) -o $(BINARY_NAME) .

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

# Run go mod tidy
.PHONY: tidy
tidy:
	@echo "Running go mod tidy..."
	$(GOMOD) tidy

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Running linter..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Static analysis (requires staticcheck)
.PHONY: staticcheck
staticcheck:
	@echo "Running staticcheck..."
	@if command -v staticcheck >/dev/null 2>&1; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed. Install with: go install honnef.co/go/tools/cmd/staticcheck@latest"; \
	fi

# Vet code
.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

# Run all quality checks
.PHONY: check
check: fmt vet staticcheck lint test

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(BINARY_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html


# Cross-compilation targets
.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-linux-amd64 .

.PHONY: build-windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-windows-amd64.exe .

.PHONY: build-darwin
build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-amd64 .

.PHONY: build-darwin-arm64
build-darwin-arm64:
	@echo "Building for macOS ARM64..."
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-darwin-arm64 .

# Build for all platforms
.PHONY: build-all
build-all: build-linux build-windows build-darwin build-darwin-arm64
	@echo "Cross-compilation complete. Binaries in $(DIST_DIR)/"

# Create release archives
.PHONY: release
release: build-all
	@echo "Creating release archives..."
	@cd $(DIST_DIR) && \
	tar -czf $(BINARY_NAME)-linux-amd64.tar.gz $(BINARY_NAME)-linux-amd64 && \
	zip $(BINARY_NAME)-windows-amd64.zip $(BINARY_NAME)-windows-amd64.exe && \
	tar -czf $(BINARY_NAME)-darwin-amd64.tar.gz $(BINARY_NAME)-darwin-amd64 && \
	tar -czf $(BINARY_NAME)-darwin-arm64.tar.gz $(BINARY_NAME)-darwin-arm64
	@echo "Release archives created in $(DIST_DIR)/"

# Development server (for testing)
.PHONY: dev
dev: build
	@echo "Running development version..."
	./$(BINARY_DIR)/$(BINARY_NAME) --help

# Show help
.PHONY: help
help:
	@echo "todo.txt Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build          Build the binary"
	@echo "  compile        Simple build (current platform)"
	@echo "  test           Run tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  bench          Run benchmarks"
	@echo "  fmt            Format code"
	@echo "  tidy           Run go mod tidy"
	@echo "  lint           Run linter (requires golangci-lint)"
	@echo "  staticcheck    Run static analysis (requires staticcheck)"
	@echo "  vet            Run go vet"
	@echo "  check          Run all quality checks"
	@echo "  clean          Clean build artifacts"
	@echo "  build-all      Cross-compile for all platforms"
	@echo "  release        Create release archives"
	@echo "  dev            Build and run development version"
	@echo "  help           Show this help message"
	@echo ""
	@echo "Cross-compilation:"
	@echo "  build-linux    Build for Linux"
	@echo "  build-windows  Build for Windows"
	@echo "  build-darwin   Build for macOS (Intel)"
	@echo "  build-darwin-arm64  Build for macOS (ARM64)"
	@echo ""
	@echo "Environment variables:"
	@echo "  VERSION        Set version for build (default: v1.0.0)"