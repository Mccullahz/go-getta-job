package model

import (
	"testing"

	"cliscraper/internal/testutils"
	"cliscraper/internal/utils"
)

func TestStateConstants(t *testing.T) {
	// Test that state constants are properly defined
	states := []state{
		StateHome,
		StateZipInput,
		StateRadiusInput,
		StateTitleInput,
		StateSearching,
		StateStarred,
		StateDone,
	}
	
	expectedCount := 7
	if len(states) != expectedCount {
		t.Errorf("Expected %d states, got %d", expectedCount, len(states))
	}
}

func TestPreviousState(t *testing.T) {
	tests := []struct {
		name     string
		current  state
		expected state
	}{
		{
			name:     "ZipInput to Home",
			current:  StateZipInput,
			expected: StateHome,
		},
		{
			name:     "RadiusInput to ZipInput",
			current:  StateRadiusInput,
			expected: StateZipInput,
		},
		{
			name:     "TitleInput to RadiusInput",
			current:  StateTitleInput,
			expected: StateRadiusInput,
		},
		{
			name:     "Searching to TitleInput",
			current:  StateSearching,
			expected: StateTitleInput,
		},
		{
			name:     "Starred to Home",
			current:  StateStarred,
			expected: StateHome,
		},
		{
			name:     "Done to Home",
			current:  StateDone,
			expected: StateHome,
		},
		{
			name:     "Home to Home (default)",
			current:  StateHome,
			expected: StateHome,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := PreviousState(tt.current)
			if result != tt.expected {
				t.Errorf("PreviousState(%v) = %v, want %v", tt.current, result, tt.expected)
			}
		})
	}
}

func TestInitialModel(t *testing.T) {
	// Create a mock service
	mockService := &mockService{}
	
	model := InitialModel(mockService)
	
	if model.CurrentState != StateHome {
		t.Errorf("Expected initial state to be StateHome, got %v", model.CurrentState)
	}
	
	if model.service != mockService {
		t.Error("Expected service to be set")
	}
	
	if model.Starred == nil {
		t.Error("Expected Starred slice to be initialized")
	}
	
	if len(model.Starred) != 0 {
		t.Errorf("Expected empty Starred slice, got %d items", len(model.Starred))
	}
	
	if model.TopCursor != 0 {
		t.Errorf("Expected TopCursor to be 0, got %d", model.TopCursor)
	}
	
	if model.InnerCursor != 0 {
		t.Errorf("Expected InnerCursor to be 0, got %d", model.InnerCursor)
	}
}

func TestModelService(t *testing.T) {
	mockService := &mockService{}
	model := Model{service: mockService}
	
	service := model.Service()
	if service != mockService {
		t.Error("Expected Service() to return the set service")
	}
}

func TestModelFields(t *testing.T) {
	model := Model{
		Zip:          "10001",
		Radius:       "5",
		Title:        "engineer",
		Err:          "test error",
		Businesses:   testutils.MockBusinesses(),
		Results:      testutils.MockJobResults(),
		ShowResults:  true,
		Width:        80,
		Height:       24,
	}
	
	if model.Zip != "10001" {
		t.Errorf("Expected Zip '10001', got %s", model.Zip)
	}
	
	if model.Radius != "5" {
		t.Errorf("Expected Radius '5', got %s", model.Radius)
	}
	
	if model.Title != "engineer" {
		t.Errorf("Expected Title 'engineer', got %s", model.Title)
	}
	
	if model.Err != "test error" {
		t.Errorf("Expected Err 'test error', got %s", model.Err)
	}
	
	if len(model.Businesses) != 3 {
		t.Errorf("Expected 3 businesses, got %d", len(model.Businesses))
	}
	
	if len(model.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(model.Results))
	}
	
	if !model.ShowResults {
		t.Error("Expected ShowResults to be true")
	}
	
	if model.Width != 80 {
		t.Errorf("Expected Width 80, got %d", model.Width)
	}
	
	if model.Height != 24 {
		t.Errorf("Expected Height 24, got %d", model.Height)
	}
}

// Mock service for testing
type mockService struct{}

func (m *mockService) Health() error {
	return nil
}

func (m *mockService) Search(zip, radius, title string) ([]utils.JobPageResult, error) {
	return testutils.MockJobResults(), nil
}

func (m *mockService) Results() ([]utils.JobPageResult, error) {
	return testutils.MockJobResults(), nil
}

func (m *mockService) Starred() ([]utils.JobPageResult, error) {
	return testutils.MockJobResults(), nil
}