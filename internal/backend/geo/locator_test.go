package geo

import (
	"os"
	"testing"
)

func TestBusinessStruct(t *testing.T) {
	business := Business{
		Name:   "Test Company",
		URL:    "https://example.com",
		Titles: []string{"engineer"},
		Lat:    40.7128,
		Lon:    -74.0060,
	}
	
	if business.Name != "Test Company" {
		t.Errorf("Expected name 'Test Company', got %s", business.Name)
	}
	
	if business.URL != "https://example.com" {
		t.Errorf("Expected URL 'https://example.com', got %s", business.URL)
	}
	
	if len(business.Titles) != 1 {
		t.Errorf("Expected 1 title, got %d", len(business.Titles))
	}
	
	if business.Lat != 40.7128 {
		t.Errorf("Expected lat 40.7128, got %f", business.Lat)
	}
	
	if business.Lon != -74.0060 {
		t.Errorf("Expected lon -74.0060, got %f", business.Lon)
	}
}

func TestPlacesStruct(t *testing.T) {
	place := Places{
		PlaceName:     "New York",
		Longitude:     "-74.0060",
		Latitude:      "40.7128",
		State:         "New York",
		StateAbbr:     "NY",
	}
	
	if place.PlaceName != "New York" {
		t.Errorf("Expected place name 'New York', got %s", place.PlaceName)
	}
	
	if place.Longitude != "-74.0060" {
		t.Errorf("Expected longitude '-74.0060', got %s", place.Longitude)
	}
	
	if place.Latitude != "40.7128" {
		t.Errorf("Expected latitude '40.7128', got %s", place.Latitude)
	}
	
	if place.State != "New York" {
		t.Errorf("Expected state 'New York', got %s", place.State)
	}
	
	if place.StateAbbr != "NY" {
		t.Errorf("Expected state abbreviation 'NY', got %s", place.StateAbbr)
	}
}

func TestZippoResponseStruct(t *testing.T) {
	response := ZippoResponse{
		PostCode:      "10001",
		Country:       "United States",
		CountryAbbr:   "US",
		Places:        []Places{
			{
				PlaceName: "New York",
				Longitude: "-74.0060",
				Latitude:  "40.7128",
				State:     "New York",
				StateAbbr: "NY",
			},
		},
	}
	
	if response.PostCode != "10001" {
		t.Errorf("Expected post code '10001', got %s", response.PostCode)
	}
	
	if response.Country != "United States" {
		t.Errorf("Expected country 'United States', got %s", response.Country)
	}
	
	if response.CountryAbbr != "US" {
		t.Errorf("Expected country abbreviation 'US', got %s", response.CountryAbbr)
	}
	
	if len(response.Places) != 1 {
		t.Errorf("Expected 1 place, got %d", len(response.Places))
	}
}

func TestOverpassResponseStruct(t *testing.T) {
	response := OverpassResponse{
		Elements: []struct {
			Lat  float64           `json:"lat"`
			Lon  float64           `json:"lon"`
			Tags map[string]string `json:"tags"`
		}{
			{
				Lat:  40.7128,
				Lon:  -74.0060,
				Tags: map[string]string{
					"name":    "Test Business",
					"website": "https://example.com",
				},
			},
		},
	}
	
	if len(response.Elements) != 1 {
		t.Errorf("Expected 1 element, got %d", len(response.Elements))
	}
	
	element := response.Elements[0]
	if element.Lat != 40.7128 {
		t.Errorf("Expected lat 40.7128, got %f", element.Lat)
	}
	
	if element.Lon != -74.0060 {
		t.Errorf("Expected lon -74.0060, got %f", element.Lon)
	}
	
	if element.Tags["name"] != "Test Business" {
		t.Errorf("Expected name 'Test Business', got %s", element.Tags["name"])
	}
	
	if element.Tags["website"] != "https://example.com" {
		t.Errorf("Expected website 'https://example.com', got %s", element.Tags["website"])
	}
}

// Note: The actual API functions (GetCoordinatesFromZip, LocateBusinesses, FindBusinessesByZip)
// make real HTTP requests to external APIs, so they need to be tested with Mock HTTP servers established

func TestFindBusinessesByZipIntegration(t *testing.T) {
	// This is an integration test that makes real API calls
	// Skip in short mode to avoid external dependencies. We are testing our code, not the APIs code.
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	// Skip if external APIs are disabled
	if os.Getenv("SKIP_EXTERNAL_APIS") == "true" {
		t.Skip("Skipping integration test when external APIs are disabled")
	}
	
	zip := "45140" // Known zip code in Ohio, not too many businesses to scrape, this is working 
	radius := 1    // 1 mile radius (very small for testing)
	
	businesses, err := FindBusinessesByZip(zip, radius)
	if err != nil {
		t.Fatalf("FindBusinessesByZip failed: %v", err)
	}
	
	// We should get some businesses back
	if len(businesses) == 0 {
		t.Error("Expected to find some businesses, got none")
	}
	
	// Verify business structure
	for i, business := range businesses {
		if business.Name == "" {
			t.Errorf("Business %d has empty name", i)
		}
		
		// Lat/Lon should be reasonable for the Loveland, OH area
		if business.Lat < 35 || business.Lat > 50 {
			t.Errorf("Business %d has unreasonable latitude: %f", i, business.Lat)
		}
		
		if business.Lon < -90.0 || business.Lon > -75.0 {
			t.Errorf("Business %d has unreasonable longitude: %f", i, business.Lon)
		}
	}
}
