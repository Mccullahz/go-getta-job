package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"

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

    // Step 1: find businesses by zip
    businesses, err := geo.FindBusinessesByZip(zip, radius)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, Response{Status: "error", Message: err.Error()})
        return
    }

    // Step 2: scrape job pages concurrently
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := []utils.JobPageResult{}

    for _, b := range businesses {
        if b.URL == "" {
            continue
        }
        wg.Add(1)
        go func(b geo.Business) {
            defer wg.Done()
            jobURL, err := web.ScrapeWebsite(b.URL, []string{title})
            if err != nil {
                fmt.Printf("scrape error for %s: %v\n", b.URL, err)
                return
            }
            if jobURL != "" {
                mu.Lock()
                results = append(results, utils.JobPageResult{
                    BusinessName: b.Name,
                    URL:          jobURL,
                })
                mu.Unlock()
            }
        }(b)
    }
    wg.Wait()

    // Step 3: store results (to ./output/results.json)
    outDir := "./output"
    if err := utils.WriteResults(results, outDir); err != nil {
        writeJSON(w, http.StatusInternalServerError, Response{Status: "error", Message: err.Error()})
        return
    }

    id := uuid.New().String() // still generate ID for client compatibility

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

// fetch starred jobs (stub)
func StarredHandler(w http.ResponseWriter, r *http.Request) {
	starred := []utils.JobPageResult{
		{BusinessName: "Starred Co", URL: "https://starred.example.com/hiring"},
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data:   starred,
	})
}

