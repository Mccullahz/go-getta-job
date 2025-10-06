package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cliscraper/internal/testutils"
)

func TestNewClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewClient(baseURL)
	
	if client.BaseURL != baseURL {
		t.Errorf("Expected base URL %s, got %s", baseURL, client.BaseURL)
	}
	
	if client.HTTPClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}
	
	if client.HTTPClient.Timeout != 1800*time.Second {
		t.Errorf("Expected timeout 1800s, got %v", client.HTTPClient.Timeout)
	}
}

func TestClientHealth(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("Expected path /health, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	err := client.Health()
	if err != nil {
		t.Errorf("Health check failed: %v", err)
	}
}

func TestClientHealthFailure(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error"}`))
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	err := client.Health()
	if err == nil {
		t.Error("Expected health check to fail")
	}
}

func TestClientSearch(t *testing.T) {
	// Create mock server
	expectedResults := testutils.MockJobResults()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Errorf("Expected path /search, got %s", r.URL.Path)
		}
		
		// Check query parameters
		zip := r.URL.Query().Get("zip")
		radius := r.URL.Query().Get("radius")
		title := r.URL.Query().Get("title")
		
		if zip != "10001" {
			t.Errorf("Expected zip 10001, got %s", zip)
		}
		if radius != "5" {
			t.Errorf("Expected radius 5, got %s", radius)
		}
		if title != "engineer" {
			t.Errorf("Expected title engineer, got %s", title)
		}
		
		response := Response{
			Status: "ok",
			Data: json.RawMessage(`{
				"zip": "10001",
				"radius": 5,
				"title": "engineer",
				"results": [
					{"business_name": "Test Company 1", "url": "https://example1.com/careers", "description": "Software engineering positions"},
					{"business_name": "Test Company 2", "url": "https://example2.com/jobs", "description": "Developer roles"}
				]
			}`),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	results, err := client.Search("10001", "5", "engineer")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	
	if len(results) != len(expectedResults) {
		t.Fatalf("Expected %d results, got %d", len(expectedResults), len(results))
	}
	
	for i, result := range results {
		if result.BusinessName != expectedResults[i].BusinessName {
			t.Errorf("Result %d: expected business name %s, got %s", i, expectedResults[i].BusinessName, result.BusinessName)
		}
		if result.URL != expectedResults[i].URL {
			t.Errorf("Result %d: expected URL %s, got %s", i, expectedResults[i].URL, result.URL)
		}
	}
}

func TestClientSearchError(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Status:  "error",
			Message: "Invalid parameters",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	_, err := client.Search("invalid", "invalid", "invalid")
	if err == nil {
		t.Error("Expected search to fail")
	}
}

func TestClientResults(t *testing.T) {
	// Create mock server
	expectedResults := testutils.MockJobResults()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/results" {
			t.Errorf("Expected path /results, got %s", r.URL.Path)
		}
		
		response := Response{
			Status: "ok",
			Data: json.RawMessage(`{
				"results": [
					{"business_name": "Test Company 1", "url": "https://example1.com/careers", "description": "Software engineering positions"},
					{"business_name": "Test Company 2", "url": "https://example2.com/jobs", "description": "Developer roles"}
				]
			}`),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	results, err := client.Results()
	if err != nil {
		t.Fatalf("Results failed: %v", err)
	}
	
	if len(results) != len(expectedResults) {
		t.Fatalf("Expected %d results, got %d", len(expectedResults), len(results))
	}
}

func TestClientStarred(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/starred" {
			t.Errorf("Expected path /starred, got %s", r.URL.Path)
		}
		
		response := Response{
			Status: "ok",
			Data:   json.RawMessage(`[{"business_name": "Starred Co", "url": "https://starred.example.com/hiring", "description": "Starred job"}]`),
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	starred, err := client.Starred()
	if err != nil {
		t.Fatalf("Starred failed: %v", err)
	}
	
	if len(starred) != 1 {
		t.Fatalf("Expected 1 starred result, got %d", len(starred))
	}
	
	if starred[0].BusinessName != "Starred Co" {
		t.Errorf("Expected business name 'Starred Co', got %s", starred[0].BusinessName)
	}
}

func TestResponseStruct(t *testing.T) {
	response := Response{
		Status:  "ok",
		Message: "Success",
		Data:    json.RawMessage(`{"key": "value"}`),
	}
	
	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got %s", response.Status)
	}
	
	if response.Message != "Success" {
		t.Errorf("Expected message 'Success', got %s", response.Message)
	}
	
	if string(response.Data) != `{"key": "value"}` {
		t.Errorf("Expected data '{\"key\": \"value\"}', got %s", string(response.Data))
	}
}