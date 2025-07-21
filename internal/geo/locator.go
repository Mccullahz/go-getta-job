// handles querying businesses near a zip
package geo

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
	"strconv"
)

type Business struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// structs to unmarshal the the json
type Places struct {
	PlaceName string `json:"place name"`
	Longitude string `json:"longitude"`
	Latitude  string `json:"latitude"`
	State      string `json:"state"`
	StateAbbr  string `json:"state abbreviation"`
}

type ZippoResponse struct {
	PostCode string `json:"post code"`
	Country  string `json:"country"`
	CountryAbbr string `json:"country abbreviation"`
	Places []Places `json:"places"`
}

// zippopotamus api allows us to extract coordinate data from a zip code. connect to the api via net/http, parse lat/lgn data from the response, and return it
func GetCoordinatesFromZip(zip string) (float64, float64, error) {
	url := fmt.Sprintf("https://api.zippopotam.us/us/%s", zip)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("reading response failed: %w", err)
	}

	var data ZippoResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return 0, 0, fmt.Errorf("JSON unmarshal failed: %w", err)
	}

	if len(data.Places) == 0 {
		return 0, 0, fmt.Errorf("no places found for zip %s", zip)
	}

	place := data.Places[0]
	lat, err1 := strconv.ParseFloat(place.Latitude, 64)
	lon, err2 := strconv.ParseFloat(place.Longitude, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid coordinates in API response")
	}

	return lat, lon, nil
}

// overpass api to locate businesses around x radius of a lat/lgn point, send a query to the overpass api, parse the response, and return a list of businesses to results.json
func LocateBusinesses(lat float64, lon float64, radius int) ([]Business, error) {
	fmt.Printf("Searching businesses around %.4f, %.4f within %dkm radius\n", lat, lon, radius)
	// TODO: Overpass logic


	// just returning example data for now
	return []Business{
		{Name: "Example Business", URL: "https://examplebuz.com/careers"},
	}, nil
}


func FindBusinessesByZip(zip string, radius int) ([]Business, error) {
	lat, lon, err := GetCoordinatesFromZip(zip)
	if err != nil {
		return nil, err
	}
	return LocateBusinesses(lat, lon, radius)
}

