// handles querying businesses near a zip
package geo

import (
	"encoding/json"
	"net/http"
	"net/url"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Business struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Titles []string `json:"titles,omitempty"`
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

type OverpassResponse struct {
	Elements []struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Tags map[string]string `json:"tags"`
	} `json:"elements"`
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
func LocateBusinesses(lat float64, lon float64, radius int) ([]Business, error) {

	rawQuery := fmt.Sprintf(`
[out:json];
(
  node["shop"]["name"](around:%d,%f,%f);
  node["amenity"]["name"](around:%d,%f,%f);
  node["office"]["name"](around:%d,%f,%f);
  node["craft"]["name"](around:%d,%f,%f);
  node["tourism"]["name"](around:%d,%f,%f);
);
out;`,
		radius*1609, lat, lon,
		radius*1609, lat, lon,
		radius*1609, lat, lon,
		radius*1609, lat, lon,
		radius*1609, lat, lon,
	)

	// collapse whitespace so it doesn’t break encoding
	compressed := strings.Join(strings.Fields(rawQuery), " ")

	// making the request and error handling
	baseURL := "https://overpass-api.de/api/interpreter"
	params := url.Values{}
	params.Set("data", compressed)
	opURL := baseURL + "?" + params.Encode()

	resp, err := http.Get(opURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}

	// NOTE: geo results are now stored in MongoDB instead of on go server as json

	var opResp OverpassResponse
	if err := json.Unmarshal(body, &opResp); err != nil {
		return nil, fmt.Errorf("Error parsing Overpass JSON: %w", err)
	}

	// []Business is intended return
	var businesses []Business
	for _, el := range opResp.Elements {
		name := el.Tags["name"]
		if name == "" {
			continue
		}

		// urls usually populate under the website tag, but sometimes under contact:website or contact:url
		url := ""
		if v, ok := el.Tags["website"]; ok {
			url = v
		} else if v, ok := el.Tags["contact:website"]; ok {
			url = v
		} else if v, ok := el.Tags["contact:url"]; ok {
			url = v
		}

		businesses = append(businesses, Business{
			Name: name,
			URL:  url,
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

