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
	zpURL := fmt.Sprintf("https://api.zippopotam.us/us/%s", zip)
	resp, err := http.Get(zpURL)
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
// we can use the Overpass API quite similarly to how we used the Zippopotamus API using net/http
// expected error to be handled: Error: encoding error: Your input contains only whitespace." which just means "no query was given")
func LocateBusinesses(lat float64, lon float64, radius int) ([]Business, error) {
	fmt.Printf("Searching businesses around %.4f, %.4f within %dkm radius\n", lat, lon, radius)
	// TODO: Overpass logic
	// geoData should be the lat lon and radius from zippo + user input here I believe
	geoData := fmt.Sprintf("[out:json];node(around:%d,%.4f,%.4f)[amenity=all];out;", radius*1000, lat, lon)
	opURL := fmt.Sprintf("https://overpass-api.de/api/interpreter?data=" + geoData)
	// minimal example of how to use the overpass api
	resp, err := http.Get(opURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// return all the businesses found in the response
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

