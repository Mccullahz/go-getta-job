// handles querying businesses near a zip using public dataset.
package geo

import (
)

type Business struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func LocateBusinesses(zip string, radius int) ([]Business, error) {
	// Could use Google Places API, Yelp, or OpenCage Geocoder

	return []Business{
		{Name: "Example Business", URL: "https://example.com/careers"},
	}, nil

}

