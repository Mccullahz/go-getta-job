package server

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cliscraper/internal/backend/geo"
	"cliscraper/internal/utils"
	"cliscraper/internal/backend/web"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseHandlers struct {
	dbManager *utils.DatabaseManager
}

func NewDatabaseHandlers() (*DatabaseHandlers, error) {
	dbManager, err := utils.NewDatabaseManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	return &DatabaseHandlers{
		dbManager: dbManager,
	}, nil
}

func (h *DatabaseHandlers) Close() error {
	return h.dbManager.Close()
}

func (h *DatabaseHandlers) SearchHandlerDB(w http.ResponseWriter, r *http.Request) {
	zip := r.URL.Query().Get("zip")
	radiusStr := r.URL.Query().Get("radius")
	title := r.URL.Query().Get("title")

	radius, err := strconv.Atoi(radiusStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, Response{Status: "error", Message: "invalid radius"})
		return
	}

	userID := utils.GetDefaultUserID()

	// step 1: find businesses by zip
	businesses, err := geo.FindBusinessesByZip(zip, radius)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, Response{Status: "error", Message: err.Error()})
		return
	}

	// if no businesses found, return early with a helpful message
	if len(businesses) == 0 {
		writeJSON(w, http.StatusOK, Response{
			Status: "ok",
			Data: map[string]interface{}{
				"user_id":  userID.Hex(),
				"zip":      zip,
				"radius":   radius,
				"title":    title,
				"message":  "No businesses found in the specified area",
				"results":  []interface{}{},
			},
		})
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
	pool := web.NewWorkerPool(100, 300) // x workers, x s timeout
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
				Description:  "", // if available
			})
		}
	}

	// step 5: store results in MongoDB
	if err := h.dbManager.WriteResultsToDB(userID, title, jobResults); err != nil {
		writeJSON(w, http.StatusInternalServerError, Response{Status: "error", Message: err.Error()})
		return
	}

	// step 6: save geo result
	_, err = h.dbManager.WriteGeoResultsToDB(userID, zip, radius)
	if err != nil {
		fmt.Printf("Warning: failed to save geo result: %v\n", err)
		// don't fail the request for this
	}

	// convert worker pool results to JobPageResult format for API compatibility
	jobPageResults := make([]utils.JobPageResult, 0, len(results))
	for _, res := range results {
		if res.Error == nil && res.JobPage != "" {
			jobPageResults = append(jobPageResults, utils.JobPageResult{
				BusinessName: res.BusinessName,
				URL:          res.JobPage,
				Description:  "", // no description available from worker pool
			})
		}
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data: map[string]interface{}{
			"user_id": userID.Hex(),
			"zip":     zip,
			"radius":  radius,
			"title":   title,
			"results": jobPageResults,
		},
	})
}

// handle results requests with MongoDB retrieval
func (h *DatabaseHandlers) ResultsHandlerDB(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ResultsHandlerDB called\n")
	// this should come from authentication
	userID := utils.GetDefaultUserID()
	fmt.Printf("Using user ID: %s\n", userID.Hex())

	results, err := h.dbManager.LoadLatestResultsFromDB(userID)
	if err != nil {
		fmt.Printf("Failed to load results: %v\n", err)
		writeJSON(w, http.StatusNotFound, Response{Status: "error", Message: "results not found"})
		return
	}
	fmt.Printf("Loaded %d results\n", len(results))

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data: map[string]interface{}{
			"results": results,
		},
	})
}

func (h *DatabaseHandlers) StarredHandlerDB(w http.ResponseWriter, r *http.Request) {
	// TODO: implement starred jobs functionality with MongoDB
	starred := []utils.JobPageResult{
		{BusinessName: "Starred Co", URL: "https://starred.example.com/hiring"},
	}

	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data:   starred,
	})
}
