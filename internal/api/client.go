// this file is essentially a wrapper for GET and POST requests to the server, replacing direct function calls to the web and geo packages

/*
Moving from a monolithic server to a client-server architecture, we need to replace the direct calls that we made to the web and geo packages (think UpdateTitle() and UpdateSearching() with HTTP requests to server endpoints. Instead of calling the functions directly, we will make a GET or POST request to the server, which will then call the functions and return the results. 
*/
package api

import (
	//"cliscraper/internal/backend/web"
	//"cliscraper/internal/backend/geo"
	"cliscraper/internal/utils"
	"net/http"
	"time"
	"fmt"
	"encoding/json"
)

type Client struct {
	BaseURL string
	HTTPClient *http.Client
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
		Timeout: 10 * time.Second,
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
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp Response
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}

	if apiResp.Status != "ok" {
		return nil, fmt.Errorf("search failed: %s", apiResp.Message)
	}

	// decoding Data.results into []utils.JobPageResult
	var payload struct {
		Zip    string               `json:"zip"`
		Radius string               `json:"radius"`
		Title  string               `json:"title"`
		Results []utils.JobPageResult `json:"results"`
	}

	if err := json.Unmarshal(apiResp.Data, &payload); err !=nil {
		return nil, err
	}
	return payload.Results, nil
}

func (c *Client) Results(id string) ([]utils.JobPageResult, error) {
	url := fmt.Sprintf("%s/results/%s", c.BaseURL, id)
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
		return nil, fmt.Errorf("results failed: %s", apiResp.Message)
	}

	var payload struct {
		ID      string               `json:"id"`
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
		return nil, fmt.Errorf("starred failed: %s", apiResp.Message)
	}

	var starred []utils.JobPageResult
	if err := json.Unmarshal(apiResp.Data, &starred); err != nil {
		return nil, err
	}

	return starred, nil
}
