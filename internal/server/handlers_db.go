package server

import (
	//"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cliscraper/internal/api"
	"cliscraper/internal/backend/geo"
	"cliscraper/internal/utils"
	"cliscraper/internal/backend/web"
	//"go.mongodb.org/mongo-driver/bson/primitive"
)

type DatabaseHandlers struct {
	dbManager  *utils.DatabaseManager
	userHandler *api.UserHandler
}

func NewDatabaseHandlers() (*DatabaseHandlers, error) {
	dbManager, err := utils.NewDatabaseManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create database manager: %w", err)
	}

	userHandler := api.NewUserHandler(dbManager.GetUserRepo())

	return &DatabaseHandlers{
		dbManager:  dbManager,
		userHandler: userHandler,
	}, nil
}

func (h *DatabaseHandlers) Close() error {
	return h.dbManager.Close()
}

func (h *DatabaseHandlers) GetUserHandler() *api.UserHandler {
	return h.userHandler
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
		// Check if this is a "no results" error (not a real failure)
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "must provide at least one element") ||
			strings.Contains(errMsg, "no businesses found") ||
			strings.Contains(errMsg, "no results") {
			writeJSON(w, http.StatusOK, Response{
				Status: "ok",
				Message: "No matching jobs found. Please try adjusting your search criteria (zip code, radius, or job title).",
				Data: map[string]interface{}{
					"user_id": userID.Hex(),
					"zip":     zip,
					"radius":  radius,
					"title":   title,
					"results": []interface{}{},
				},
			})
			return
		}
		// Real server error - return user-friendly message
		writeJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Unable to search for businesses. Please try again later.",
		})
		return
	}

	// if no businesses found, return early with a helpful message
	if len(businesses) == 0 {
		writeJSON(w, http.StatusOK, Response{
			Status: "ok",
			Message: "No businesses found in the specified area. Please try adjusting your zip code or search radius.",
			Data: map[string]interface{}{
				"user_id": userID.Hex(),
				"zip":     zip,
				"radius":  radius,
				"title":   title,
				"results": []interface{}{},
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
	// If no job results found, return a user-friendly message instead of error
	if len(jobResults) == 0 {
		writeJSON(w, http.StatusOK, Response{
			Status: "ok",
			Message: fmt.Sprintf("No matching jobs found for '%s' in the specified area. Please try adjusting your search criteria.", title),
			Data: map[string]interface{}{
				"user_id": userID.Hex(),
				"zip":     zip,
				"radius":  radius,
				"title":   title,
				"results": []interface{}{},
			},
		})
		return
	}

	if err := h.dbManager.WriteResultsToDB(userID, title, jobResults); err != nil {
		// Check if this is a "no data" error that we should handle gracefully
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "must provide at least one element") ||
			strings.Contains(errMsg, "no data") ||
			strings.Contains(errMsg, "empty") {
			writeJSON(w, http.StatusOK, Response{
				Status: "ok",
				Message: fmt.Sprintf("No matching jobs found for '%s' in the specified area. Please try adjusting your search criteria.", title),
				Data: map[string]interface{}{
					"user_id": userID.Hex(),
					"zip":     zip,
					"radius":  radius,
					"title":   title,
					"results": []interface{}{},
				},
			})
			return
		}
		// Real database error
		writeJSON(w, http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Unable to save search results. Please try again later.",
		})
		return
	}

	// step 6: save geo result
	_, err = h.dbManager.WriteGeoResultsToDB(userID, zip, radius)
	if err != nil {
		fmt.Printf("Warning: failed to save geo result: %v\n", err)
		// don't fail the request for this
	}

	// Return successful response with job results (already checked for empty above)
	writeJSON(w, http.StatusOK, Response{
		Status: "ok",
		Data: map[string]interface{}{
			"user_id": userID.Hex(),
			"zip":     zip,
			"radius":  radius,
			"title":   title,
			"results": jobResults,
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
		writeJSON(w, http.StatusNotFound, Response{
			Status:  "error",
			Message: "No search results found. Please perform a search first.",
		})
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
