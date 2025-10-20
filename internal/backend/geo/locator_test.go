package geo

import (
	"os"
	"testing"
)

func TestBusinessStruct(t *testing.T) {
	business := Business{
		Name:   "Test Company",
		URL:    "https://example.com",
		Titles: []string{"job"},
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
		PlaceName:     "Loveland",
		Longitude:     "-84.2638",
		Latitude:      "39.2689",
		State:         "Ohio",
		StateAbbr:     "OH",
	}
	
	if place.PlaceName != "Loveland" {
		t.Errorf("Expected place name 'Loveland', got %s", place.PlaceName)
	}
	
	if place.Longitude != "-84.2638" {
		t.Errorf("Expected longitude '-84.2638', got %s", place.Longitude)
	}
	
	if place.Latitude != "39.2689" {
		t.Errorf("Expected latitude '39.2689', got %s", place.Latitude)
	}
	
	if place.State != "Ohio" {
		t.Errorf("Expected state 'Ohio', got %s", place.State)
	}
	
	if place.StateAbbr != "OH" {
		t.Errorf("Expected state abbreviation 'OH', got %s", place.StateAbbr)
	}
}

func TestZippoResponseStruct(t *testing.T) {
	response := ZippoResponse{
		PostCode:      "45140",
		Country:       "United States",
		CountryAbbr:   "US",
		Places:        []Places{
			{
				PlaceName: "Loveland",
				Longitude: "-84.2638",
				Latitude:  "39.2689",
				State:     "Ohio",
				StateAbbr: "OH",
			},
		},
	}
	
	if response.PostCode != "45140" {
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
				Lat:  39.2689,
				Lon:  -84.2638,
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
	if element.Lat != 39.2689 {
		t.Errorf("Expected lat 39.2689, got %f", element.Lat)
	}
	
	if element.Lon != -84.2638 {
		t.Errorf("Expected lon -84.2638, got %f", element.Lon)
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
	
	// Test with a known zip code
	zip := "45140" // Known zip code in Ohio, not too many businesses to scrape
	radius := 2    // 2 mile radius (very small for testing)
	
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

