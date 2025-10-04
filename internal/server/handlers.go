package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"

	"cliscraper/internal/backend/geo"
	"cliscraper/internal/utils"
	"cliscraper/internal/backend/web"
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
	zip := r.URL.Query().Get("zip")
	radiusStr := r.URL.Query().Get("radius")
	title := r.URL.Query().Get("title")

	radius, err := strconv.Atoi(radiusStr)
    	if err != nil {
        	writeJSON(w, http.StatusBadRequest, Response{Status: "error", Message: "invalid radius"})
        	return
    	}

    	// step 1: find businesses by zip
    	businesses, err := geo.FindBusinessesByZip(zip, radius)
    	if err != nil {
    	    writeJSON(w, http.StatusInternalServerError, Response{Status: "error", Message: err.Error()})
    	    return
    	}

    	// step 2: create workers and prepare jobs for pooling
	jobs := []web.Job{}
	for _, b := range businesses {
		if b.URL == "" {
			continue
		}
		jobs = append(jobs, web.Job{
			BusinessName: b.Name,
			URL:          b.URL,
			Titles:       []string{title},
		})
	}

	// step 3: run worker pool
	pool := web.NewWorkerPool(10, 300) // x workers, x s timeout
	results := pool.Run(jobs)

	// step 4: collect results
	jobResults := []utils.JobPageResult{}
	for _, res := range results {
		if res.Error != nil {
			fmt.Printf("Error scraping %s: %v\n", res.URL, res.Error)
			continue
		}
		if res.JobPage != "" {
			jobResults = append(jobResults, utils.JobPageResult{
				BusinessName: res.BusinessName, 
				URL:          res.JobPage,
			})
		}
	}

	// step 5: store results
    	outDir := "./output"
    	if err := utils.WriteResults(jobResults, outDir); err != nil {
    	    writeJSON(w, http.StatusInternalServerError, Response{Status: "error", Message: err.Error()})
    	    return
    	}

    	id := uuid.New().String() // still generating ID, not using now, but will need to when moving to DB

    	writeJSON(w, http.StatusOK, Response{
        	Status: "ok",
        	Data: map[string]interface{}{
        	    "id":      id,
        	    "zip":     zip,
        	    "radius":  radius,
        	    "title":   title,
        	    "results": results,
        	},
     	})
}

// fetch search results by latest file
func ResultsHandler(w http.ResponseWriter, r *http.Request) {
    // NOTE: ignoring {id}, just load the latest results.json
    outDir := "./output"
    results, err := utils.LoadLatestResults(outDir)
    if err != nil {
        writeJSON(w, http.StatusNotFound, Response{Status: "error", Message: "results not found"})
        return
    }
	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data:   map[string]interface{}{
		"results": results,
		},
	})
}

// fetch starred jobs (stub), Starred works on client side ONLY for now. When DB is added, this function will fetch from there/
func StarredHandler(w http.ResponseWriter, r *http.Request) {
	starred := []utils.JobPageResult{
		{BusinessName: "Starred Co", URL: "https://starred.example.com/hiring"},
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data:   starred,
	})
}

