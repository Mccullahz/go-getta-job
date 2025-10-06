# Go Getta Job - Testing Makefile

.PHONY: test test-unit test-integration test-coverage test-race test-verbose clean test-setup

# Default test target
test: test-no-external

# Run all unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v ./internal/...

# Run integration tests (requires external APIs)
test-integration:
	@echo "Running integration tests..."
	go test -v -run TestFullWorkflowIntegration ./integration_test.go

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./internal/...

# Run tests in verbose mode
test-verbose:
	@echo "Running tests in verbose mode..."
	go test -v -count=1 ./internal/...

# Run specific package tests
test-utils:
	go test -v ./internal/utils/...

test-web:
	go test -v ./internal/backend/web/...

test-geo:
	go test -v ./internal/backend/geo/...

test-api:
	go test -v ./internal/api/...

test-server:
	go test -v ./internal/server/...

test-ui:
	go test -v ./internal/ui/...

# Run tests without external API calls
test-no-external:
	@echo "Running tests without external API calls..."
	SKIP_EXTERNAL_APIS=true go test -v ./internal/...

# Run tests with minimal external API usage (small radius)
test-minimal-external:
	@echo "Running tests with minimal external API usage..."
	go test -v ./internal/...

# Setup test environment
test-setup:
	@echo "Setting up test environment..."
	mkdir -p testdata
	@echo "Test data directory created"

# Clean test artifacts
clean:
	@echo "Cleaning test artifacts..."
	rm -f coverage.out coverage.html
	rm -rf testdata/output
	@echo "Test artifacts cleaned"

# Run all tests (unit + integration)
test-all: test-unit test-integration

# Run tests in CI mode (no external APIs, with coverage)
test-ci: test-setup test-no-external test-coverage

# Help target
help:
	@echo "Available targets:"
	@echo "  test              - Run unit tests (default)"
	@echo "  test-unit         - Run unit tests"
	@echo "  test-integration  - Run integration tests"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-race         - Run tests with race detection"
	@echo "  test-verbose      - Run tests in verbose mode"
	@echo "  test-utils        - Run utils package tests"
	@echo "  test-web          - Run web package tests"
	@echo "  test-geo          - Run geo package tests"
	@echo "  test-api          - Run API package tests"
	@echo "  test-server       - Run server package tests"
	@echo "  test-ui           - Run UI package tests"
	@echo "  test-no-external  - Run tests without external API calls"
	@echo "  test-setup        - Setup test environment"
	@echo "  test-all          - Run all tests"
	@echo "  test-ci           - Run tests in CI mode"
	@echo "  clean             - Clean test artifacts"
	@echo "  help              - Show this help"
