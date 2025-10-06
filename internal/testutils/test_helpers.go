package testutils

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cliscraper/internal/backend/geo"
	"cliscraper/internal/utils"
)

// MockHTTPServer creates a test HTTP server with predefined responses
func MockHTTPServer(t *testing.T, responses map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		response, exists := responses[r.URL.Path]
		if !exists {
			http.NotFound(w, r)
			return
		}
		
		json.NewEncoder(w).Encode(response)
	}))
}

// LoadTestData loads JSON test data from testdata directory
func LoadTestData(t *testing.T, filename string, v interface{}) {
	path := filepath.Join("testdata", filename)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to load test data %s: %v", filename, err)
	}
	
	if err := json.Unmarshal(data, v); err != nil {
		t.Fatalf("Failed to unmarshal test data %s: %v", filename, err)
	}
}

// CreateTempDir creates a temporary directory for testing
func CreateTempDir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "go-getta-job-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir
}

// CleanupTempDir removes a temporary directory
func CleanupTempDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir %s: %v", dir, err)
	}
}

// MockBusinesses returns test business data
func MockBusinesses() []geo.Business {
	return []geo.Business{
		{Name: "Test Company 1", URL: "https://example1.com", Lat: 40.7128, Lon: -74.0060},
		{Name: "Test Company 2", URL: "https://example2.com", Lat: 40.7589, Lon: -73.9851},
		{Name: "Test Company 3", URL: "https://example3.com", Lat: 40.7505, Lon: -73.9934},
	}
}

// MockJobResults returns test job results
func MockJobResults() []utils.JobPageResult {
	return []utils.JobPageResult{
		{BusinessName: "Test Company 1", URL: "https://example1.com/careers", Description: "Software engineering positions"},
		{BusinessName: "Test Company 2", URL: "https://example2.com/jobs", Description: "Developer roles"},
	}
}

// MockHTTPResponse creates a mock HTTP response
func MockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}