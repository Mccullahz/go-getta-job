// handles querying businesses near a zip
package geo

import (
	"encoding/json"
	"net/http"
	"fmt"
	"io"
	"strconv"
	"io/ioutil"
	//"cliscraper/internal/utils"
)

type Business struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
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

// overpass api to locate businesses around x radius of a lat/lgn point, send a query to the overpass api, parse the response, and return a list of businesses to geo-results.json
// geoData should be the lat lon from zippo + radius from user input. OverPass is going to read this distance in meters, so we need to convert to miles. -- 1 mile = 1609.34 meters, so we can multiply the radius by 1609.34 to get the distance in meters.
// expected error to be handled: Error: encoding error: Your input contains only whitespace." which just means "no query was given")
func LocateBusinesses(lat, lon float64, radius int) ([]Business, error) {
	// miles → meters
	radiusMeters := radius * 1609

	query := fmt.Sprintf(`
		[out:json];
		(
			node["amenity"]["name"](around:%d,%f,%f);
			node["shop"]["name"](around:%d,%f,%f);
			node["office"]["name"](around:%d,%f,%f);
			node["craft"]["name"](around:%d,%f,%f);
			node["tourism"]["name"](around:%d,%f,%f);
		);
		out;
	`, radiusMeters, lat, lon,
		radiusMeters, lat, lon,
		radiusMeters, lat, lon,
		radiusMeters, lat, lon,
		radiusMeters, lat, lon,
	)

	url := "https://overpass-api.de/api/interpreter?data=" + query

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query Overpass: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("bad Overpass response: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Elements []struct {
			Tags struct {
				Name string `json:"name"`
				URL  string `json:"website"`
			} `json:"tags"`
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"elements"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Overpass response: %w", err)
	}

	var businesses []Business
	for _, el := range result.Elements {
		if el.Tags.Name == "" {
			continue
		}
		businesses = append(businesses, Business{
			Name: el.Tags.Name,
			URL:  el.Tags.URL,
			Lat:  el.Lat,
			Lon:  el.Lon,
		})
	}

	return businesses, nil
}


func FindBusinessesByZip(zip string, radius int) ([]Business, error) {
	lat, lon, err := GetCoordinatesFromZip(zip)
	if err != nil {
		return nil, err
	}
	return LocateBusinesses(lat, lon, radius)
}

