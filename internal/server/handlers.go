package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	//"cliscraper/internal/backend/web"
	"cliscraper/internal/utils"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// health check
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, Response{Status: "ok"})
}

// scrape/search trigger
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	// expecting zip, radius, title in query params
	zip := r.URL.Query().Get("zip")
	radius := r.URL.Query().Get("radius")
	title := r.URL.Query().Get("title")

	// TODO: actually call into backend/web scraper logic, for now just mock
	results := []utils.JobPageResult{
		{BusinessName: "Example Co", URL: "https://example.com/careers"},
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data: map[string]interface{}{
			"zip":     zip,
			"radius":  radius,
			"title":   title,
			"results": results,
		},
	})
}

// fetch search results by ID
func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// TODO: load from storage
	results := []utils.JobPageResult{
		{BusinessName: "Cached Co", URL: "https://cached.example.com/jobs"},
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data: map[string]interface{}{
			"id":      id,
			"results": results,
		},
	})
}

// fetch starred jobs
func StarredHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: pull starred items, will need to actually store them somewhere out of memory
	starred := []utils.JobPageResult{
		{BusinessName: "Starred Co", URL: "https://starred.example.com/hiring"},
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data:   starred,
	})
}

