package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
        writeJSON(w, http.StatusBadRequest, Response{
            Status:  "error",
            Message: "invalid radius",
        })
        return
    }

    // step 1: find businesses by ZIP
    businesses, err := geo.FindBusinessesByZip(zip, radius)
    if err != nil {
        // "no input slice" case as no results, not failure
        if strings.Contains(err.Error(), "must provide at least one element in input slice") ||
            strings.Contains(strings.ToLower(err.Error()), "no businesses found") {
            writeJSON(w, http.StatusOK, Response{
                Status:  "ok",
                Message: "no businesses found in specified area",
                Data: map[string]interface{}{
                    "zip":     zip,
                    "radius":  radius,
                    "title":   title,
                    "results": []utils.JobPageResult{},
                },
            })
            return
        }

        // otherwise, it's a real server error
        writeJSON(w, http.StatusInternalServerError, Response{
            Status:  "error",
            Message: fmt.Sprintf("failed to locate businesses: %v", err),
        })
        return
    }

    // if no businesses exit cleanly
    if len(businesses) == 0 {
        writeJSON(w, http.StatusOK, Response{
            Status:  "ok",
            Message: "no businesses found in specified area",
            Data: map[string]interface{}{
                "zip":     zip,
                "radius":  radius,
                "title":   title,
                "results": []utils.JobPageResult{},
            },
        })
        return
    }

    // step 2: prepare jobs
    jobs := make([]web.Job, 0, len(businesses))
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

    // edge case where all businesses had no url
    if len(jobs) == 0 {
        writeJSON(w, http.StatusOK, Response{
            Status:  "ok",
            Message: "businesses found, but none have valid URLs",
            Data: map[string]interface{}{
                "zip":     zip,
                "radius":  radius,
                "title":   title,
                "results": []utils.JobPageResult{},
            },
        })
        return
    }

    // step 3: run worker pool
    pool := web.NewWorkerPool(100, 300)
    results := pool.Run(jobs)

    // step 4: collect successful results
    jobResults := make([]utils.JobPageResult, 0, len(results))
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

    // step 5: save only if there are valid results
    outDir := "./output"
    if len(jobResults) > 0 {
        if err := utils.WriteResults(jobResults, outDir); err != nil {
            writeJSON(w, http.StatusInternalServerError, Response{
                Status:  "error",
                Message: fmt.Sprintf("failed to save results: %v", err),
            })
            return
        }
    }

    id := uuid.New().String()

    // always return ok with structured data, even if results are empty
    writeJSON(w, http.StatusOK, Response{
        Status: "ok",
        Data: map[string]interface{}{
            "id":      id,
            "zip":     zip,
            "radius":  radius,
            "title":   title,
            "results": jobResults,
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
