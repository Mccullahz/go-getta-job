// this file is essentially a wrapper for GET and POST requests to the server, replacing direct function calls to the web and geo packages

/*
Moving from a monolithic server to a client-server architecture, we need to replace the direct calls that we made to the web and geo packages (think UpdateTitle() and UpdateSearching() with HTTP requests to server endpoints. Instead of calling the functions directly, we will make a GET or POST request to the server, which will then call the functions and return the results. 
*/
package api

import (
	"cliscraper/internal/utils"
	"net/http"
	"net/url"
	"time"
	"fmt"
	"io"
	"bytes"
	"strings"
	"encoding/json"

)

type Client struct {
	BaseURL string
	HTTPClient *http.Client
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
		Timeout: 1800 * time.Second, // outrageous timeout for scraping, just dont want to deal with issues with tomeouts right now
	},
	}
}

// backend health endpoint.
func (c *Client) Health() error {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/health")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed: %s", resp.Status)
	}
	return nil
}

func (c *Client) Search(zip, radius, title string) ([]utils.JobPageResult, error) {
	params := url.Values{}
	params.Set("zip", zip)
	params.Set("radius", radius)
	params.Set("title", title)

	url := fmt.Sprintf("%s/search?%s", c.BaseURL, params.Encode())

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// handle Overpass HTML responses early
	if bytes.HasPrefix(bytes.TrimSpace(body), []byte("<")) {
		return nil, fmt.Errorf("backend returned HTML instead of JSON (likely Overpass error or rate limit)\nURL: %s\nBody: %s", url, string(body[:120]))
	}

	// decode the unified response first to extract user-friendly messages
	var apiResp Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		// If we can't decode JSON and status is not OK, return a cleaner error
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("search failed: server returned an error. Please try again later.")
		}
		return nil, fmt.Errorf("invalid JSON from API: %w\nBody:\n%s", err, string(body[:200]))
	}

	// handle non-200 statuses - extract message from response if available
	if resp.StatusCode != http.StatusOK {
		msg := apiResp.Message
		if msg == "" {
			msg = "Search failed. Please try again later."
		}
		// Check if this is a "no results" case that should be treated as success
		msgLower := strings.ToLower(msg)
		if strings.Contains(msgLower, "no matching jobs") ||
			strings.Contains(msgLower, "no businesses found") ||
			strings.Contains(msgLower, "must provide at least one element") ||
			strings.Contains(msgLower, "no results") {
			return []utils.JobPageResult{}, nil
		}
		return nil, fmt.Errorf("%s", msg)
	}

	// handle no results cleanly
	if apiResp.Status != "ok" {
		msg := apiResp.Message
		if strings.Contains(strings.ToLower(msg), "must provide at least one element") ||
			strings.Contains(strings.ToLower(msg), "no businesses found") ||
			strings.Contains(strings.ToLower(msg), "no matching jobs") ||
			strings.Contains(strings.ToLower(msg), "no results") {
			// treat as valid empty result
			return []utils.JobPageResult{}, nil
		}
		if msg == "" {
			msg = "Search failed. Please try again later."
		}
		return nil, fmt.Errorf("%s", msg)
	}

	// gracefully handle missing or empty data
	if len(apiResp.Data) == 0 {
		return []utils.JobPageResult{}, nil
	}

	var payload struct {
		Zip     string               `json:"zip"`
		Radius  int                  `json:"radius"`
		Title   string               `json:"title"`
		Results []utils.JobPageResult `json:"results"`
	}

	if err := json.Unmarshal(apiResp.Data, &payload); err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w\nData:\n%s", err, string(apiResp.Data))
	}

	// return empty but valid slice if nothing found
	if len(payload.Results) == 0 {
		return []utils.JobPageResult{}, nil
	}

	return payload.Results, nil
}

func (c *Client) Results() ([]utils.JobPageResult, error) {
	url := fmt.Sprintf("%s/results", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if apiResp.Status != "ok" {
		msg := apiResp.Message
		if msg == "" {
			msg = "Unable to retrieve results. Please try again later."
		}
		return nil, fmt.Errorf("%s", msg)
	}

	var payload struct {
		Results []utils.JobPageResult `json:"results"`
	}
	if err := json.Unmarshal(apiResp.Data, &payload); err != nil {
		return nil, err
	}

	return payload.Results, nil
}

func (c *Client) Starred() ([]utils.JobPageResult, error) {
	url := fmt.Sprintf("%s/starred", c.BaseURL)
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if apiResp.Status != "ok" {
		msg := apiResp.Message
		if msg == "" {
			msg = "Unable to retrieve starred jobs. Please try again later."
		}
		return nil, fmt.Errorf("%s", msg)
	}

	var starred []utils.JobPageResult
	if err := json.Unmarshal(apiResp.Data, &starred); err != nil {
		return nil, err
	}

	return starred, nil
}
