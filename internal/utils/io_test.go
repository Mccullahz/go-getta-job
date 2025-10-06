package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteResults(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	// Test data
	results := []JobPageResult{
		{BusinessName: "Test Company 1", URL: "https://example1.com/careers", Description: "Software engineering positions"},
		{BusinessName: "Test Company 2", URL: "https://example2.com/jobs", Description: "Developer roles"},
	}

	// Test writing results
	err := WriteResults(results, tempDir)
	if err != nil {
		t.Fatalf("WriteResults failed: %v", err)
	}

	// Verify file was created
	resultsFile := filepath.Join(tempDir, "results.json")
	if _, err := os.Stat(resultsFile); os.IsNotExist(err) {
		t.Fatalf("Results file was not created")
	}

	// Verify file contents
	loadedResults, err := LoadLatestResults(tempDir)
	if err != nil {
		t.Fatalf("LoadLatestResults failed: %v", err)
	}

	if len(loadedResults) != len(results) {
		t.Fatalf("Expected %d results, got %d", len(results), len(loadedResults))
	}

	for i, result := range loadedResults {
		if result.BusinessName != results[i].BusinessName {
			t.Errorf("BusinessName mismatch: expected %s, got %s", results[i].BusinessName, result.BusinessName)
		}
		if result.URL != results[i].URL {
			t.Errorf("URL mismatch: expected %s, got %s", results[i].URL, result.URL)
		}
	}
}

func TestLoadLatestResults(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test with no files
	_, err := LoadLatestResults(tempDir)
	if err == nil {
		t.Error("Expected error when no result files exist")
	}

	// Create test results file
	results := []JobPageResult{
		{BusinessName: "Test Company", URL: "https://test.com/jobs", Description: "Test job"},
	}
	
	err = WriteResults(results, tempDir)
	if err != nil {
		t.Fatalf("WriteResults failed: %v", err)
	}

	// Test loading results
	loadedResults, err := LoadLatestResults(tempDir)
	if err != nil {
		t.Fatalf("LoadLatestResults failed: %v", err)
	}

	if len(loadedResults) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(loadedResults))
	}

	if loadedResults[0].BusinessName != "Test Company" {
		t.Errorf("Expected 'Test Company', got %s", loadedResults[0].BusinessName)
	}
}

func TestWriteGeoResults(t *testing.T) {
	tempDir := t.TempDir()
	
	// Test data
	geoData := []byte(`{"elements":[{"lat":40.7128,"lon":-74.0060,"tags":{"name":"Test Business"}}]}`)
	
	// Test writing geo results
	err := WriteGeoResults(geoData, tempDir)
	if err != nil {
		t.Fatalf("WriteGeoResults failed: %v", err)
	}

	// Verify file was created
	geoFile := filepath.Join(tempDir, "geo_results.json")
	if _, err := os.Stat(geoFile); os.IsNotExist(err) {
		t.Fatalf("Geo results file was not created")
	}

	// Verify file contents
	fileData, err := os.ReadFile(geoFile)
	if err != nil {
		t.Fatalf("Failed to read geo results file: %v", err)
	}

	// Should contain the original data (pretty-printed)
	if len(fileData) == 0 {
		t.Error("Geo results file is empty")
	}
}

func TestWriteResultsCreatesDirectory(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "nonexistent", "subdir")
	
	results := []JobPageResult{
		{BusinessName: "Test", URL: "https://test.com", Description: "Test"},
	}

	// Should create the directory structure
	err := WriteResults(results, subDir)
	if err != nil {
		t.Fatalf("WriteResults should create directory: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(subDir); os.IsNotExist(err) {
		t.Fatalf("Directory was not created: %s", subDir)
	}
}