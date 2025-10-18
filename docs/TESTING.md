# Testing Guide for Go Getta Job

This document provides a comprehensive guide to testing the Go Getta Job application. Currently, this only includes unit and integration tests, but future plans include end to end tests and UI / Mobile tests.

## Overview

The testing suite is organized into several layers:

1. **Unit Tests** - Test individual functions and components in isolation
2. **Integration Tests** - Test the interaction between components
3. **End-to-End Tests** - Test complete workflows

## Test Structure

```
├── internal/
│   ├── testutils/           # Test utilities and helpers
│   ├── utils/               # Unit tests for utility functions
│   ├── backend/
│   │   ├── geo/            # Unit tests for geo-location functionality
│   │   └── web/            # Unit tests for web scraping
│   ├── api/                # Unit tests for API client
│   ├── server/             # Unit tests for HTTP handlers
│   └── ui/                 # Unit tests for UI components
├── testdata/               # Test data files
├── integration_test.go     # Integration tests
└── Makefile               # Test runner commands
```

## Running Tests

### Quick Start

```bash
# Run all unit tests
make test

# Run tests with coverage
make test-coverage

# Run tests without external API calls
make test-no-external
```

### Available Make Commands

| Command | Description |
|---------|-------------|
| `make test` | Run unit tests (default) |
| `make test-unit` | Run unit tests |
| `make test-integration` | Run integration tests |
| `make test-coverage` | Run tests with coverage report |
| `make test-race` | Run tests with race detection |
| `make test-verbose` | Run tests in verbose mode |
| `make test-utils` | Run utils package tests |
| `make test-web` | Run web package tests |
| `make test-geo` | Run geo package tests |
| `make test-api` | Run API package tests |
| `make test-server` | Run server package tests |
| `make test-ui` | Run UI package tests |
| `make test-no-external` | Run tests without external API calls |
| `make test-all` | Run all tests |
| `make test-ci` | Run tests in CI mode |
| `make clean` | Clean test artifacts |

## Test Categories

### Unit Tests

Unit tests focus on testing individual functions and methods in isolation.

#### Utils Package (`internal/utils/`)
- **helper_test.go**: Tests for URL normalization, validation functions
- **io_test.go**: Tests for file I/O operations, JSON serialization

#### Web Package (`internal/backend/web/`)
- **detector_test.go**: Tests for job page detection, title matching
- **workers_test.go**: Tests for worker pool functionality

#### Geo Package (`internal/backend/geo/`)
- **locator_test.go**: Tests for coordinate lookup, business finding

#### API Package (`internal/api/`)
- **client_test.go**: Tests for HTTP client functionality

#### Server Package (`internal/server/`)
- **handlers_test.go**: Tests for HTTP handlers

#### UI Package (`internal/ui/`)
- **model_test.go**: Tests for state management

### Integration Tests

Integration tests verify that different components work together correctly.

#### Full Workflow Tests
- Complete search workflow from API request to results
- File operations integration
- Client server communication

#### Mock Server Tests
- API client with mock HTTP server
- Error handling scenarios
- Response parsing validation

## Test Data

Test data is stored in the `testdata/` directory:

- `test_businesses.json`: Sample business data for testing
- `test_job_pages.html`: Sample HTML for job page detection
- `test_results.json`: Sample job results for testing

## Mocking and Test Utilities

### Test Utilities (`internal/testutils/`)

The `testutils` package provides:

- **MockHTTPServer**: Creates test HTTP servers with predefined responses
- **LoadTestData**: Loads JSON test data from files
- **CreateTempDir/CleanupTempDir**: Manages temporary directories
- **MockBusinesses/MockJobResults**: Provides sample data

### Example Usage

```go
// Create a mock server
server := testutils.MockHTTPServer(t, map[string]interface{}{
    "/health": map[string]string{"status": "ok"},
    "/search": map[string]interface{}{
        "status": "ok",
        "data": map[string]interface{}{
            "results": testutils.MockJobResults(),
        },
    },
})
defer server.Close()

// Create client with mock server
client := api.NewClient(server.URL)

// Test the client
err := client.Health()
if err != nil {
    t.Fatalf("Health check failed: %v", err)
}
```

## Coverage

The test suite aims for comprehensive coverage of:

- **Core functionality**: All business logic functions
- **Error handling**: Edge cases and error conditions
- **Data validation**: Input validation and sanitization
- **API interactions**: HTTP client and server handlers
- **File operations**: Reading and writing data files

### Coverage Reports

```bash
# Generate coverage report
make test-coverage

# View coverage in browser
open coverage.html
```

## External Dependencies

Some tests make real HTTP requests to external APIs:

- **Zippopotamus API**: For zip code to coordinates conversion
- **Overpass API**: For business location data
- **Target websites**: For web scraping functionality

### Skipping External API Calls

To run tests without external API calls:

```bash
# Set environment variable
export SKIP_EXTERNAL_APIS=true
make test

# Or use the convenience command
make test-no-external
```

## Continuous Integration

The project includes GitHub Actions workflows for:

- **Unit Tests**: Run on every push and PR
- **Integration Tests**: Run on PRs and main branch
- **Linting**: Code quality checks
- **Security Scanning**: Vulnerability detection

## Best Practices

### Writing Tests

1. **Test naming**: Use descriptive test names that explain the scenario
2. **Arrange-Act-Assert**: Structure tests clearly
3. **Test isolation**: Each test should be independent
4. **Mock external dependencies**: Use mocks for external services
5. **Test edge cases**: Include boundary conditions and error cases

### Example Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := FunctionName(input)
    
    // Assert
    if result != expected {
        t.Errorf("FunctionName(%q) = %q, want %q", input, result, expected)
    }
}
```

### Test Data Management

1. **Use testdata directory**: Store test files in `testdata/`
2. **Create temporary directories**: Use `t.TempDir()` for file tests
3. **Clean up resources**: Ensure proper cleanup in tests
4. **Use constants**: Define test data as constants when possible

## Debugging Tests

### Verbose Output

```bash
# Run tests with verbose output
make test-verbose

# Run specific test with verbose output
go test -v -run TestSpecificFunction ./internal/utils/
```

### Race Detection

```bash
# Run tests with race detection
make test-race
```

### Debugging Individual Tests

```bash
# Run specific test
go test -v -run TestFunctionName ./internal/package/

# Run tests in specific package
go test -v ./internal/package/
```

## Performance Testing

While not included in the current test suite, consider adding:

- **Benchmark tests**: For performance-critical functions
- **Load testing**: For API endpoints
- **Memory profiling**: For memory usage optimization

## Future Improvements

1. **UI Testing**: Add tests for TUI components
2. **End-to-End Testing**: Complete user workflow tests
3. **Performance Testing**: Benchmark and load tests
4. **Property-Based Testing**: Using tools like `quick`
5. **Mutation Testing**: To verify test quality

## Troubleshooting

### Common Issues

1. **Import cycles**: Ensure test files don't create import cycles
2. **External API failures**: Use mocks or skip external tests
3. **File permissions**: Ensure test directories are writable
4. **Race conditions**: Use race detection to identify issues

### Getting Help

- Check the test output for specific error messages
- Use `go test -v` for detailed output
- Review the test utilities in `internal/testutils/`
- Check the CI logs for integration test failures
