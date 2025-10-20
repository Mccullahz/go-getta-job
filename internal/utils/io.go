// io for managing results files and database operations
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"os"
	"time"

	"cliscraper/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// result struct to hold job page details
type JobPageResult struct {
	BusinessName string `json:"business_name"`
	URL	   string `json:"url"`
	Description string `json:"description"`
}

type DatabaseManager struct {
	client     *database.Client
	jobRepo    *database.JobRepository
	businessRepo *database.BusinessRepository
	jobResultRepo *database.JobResultRepository
	geoResultRepo *database.GeoResultRepository
}

func NewDatabaseManager() (*DatabaseManager, error) {
	client, err := database.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create database client: %w", err)
	}

	repo := database.NewRepository(client)
	
	return &DatabaseManager{
		client:        client,
		jobRepo:       database.NewJobRepository(repo),
		businessRepo:  database.NewBusinessRepository(repo),
		jobResultRepo: database.NewJobResultRepository(repo),
		geoResultRepo: database.NewGeoResultRepository(repo),
	}, nil
}

func (dm *DatabaseManager) Close() error {
	return dm.client.Close()
}

// legacy, keeping for non-database use cases
func LoadLatestResults(dir string) ([]JobPageResult, error) {
	files, err := filepath.Glob(filepath.Join(dir, "results*.json"))
	if err != nil || len(files) == 0 {
		return nil,fmt.Errorf("no result files found")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	data, err := ioutil.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	var results []JobPageResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (dm *DatabaseManager) LoadLatestResultsFromDB(userID primitive.ObjectID) ([]JobPageResult, error) {
	// get the latest job result
	jobResult, err := dm.jobResultRepo.GetLatestJobResults(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest job results: %w", err)
	}

	// get the actual job details
	jobs, err := dm.jobRepo.GetJobsByIDs(jobResult.Jobs)
	if err != nil {
		return nil, fmt.Errorf("failed to get job details: %w", err)
	}

	businessIDs := make([]primitive.ObjectID, 0, len(jobs))
	for _, job := range jobs {
		businessIDs = append(businessIDs, job.BusinessID)
	}

	businesses, err := dm.businessRepo.GetBusinessesByIDs(businessIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get business details: %w", err)
	}

	businessMap := make(map[primitive.ObjectID]database.Business)
	for _, business := range businesses {
		businessMap[business.ID] = business
	}

	results := make([]JobPageResult, 0, len(jobs))
	for _, job := range jobs {
		business, exists := businessMap[job.BusinessID]
		if !exists {
			continue
		}

		results = append(results, JobPageResult{
			BusinessName: business.Name,
			URL:          job.URL,
			Description:  job.Description,
		})
	}

	return results, nil
}

// legacy function
func WriteResults(results []JobPageResult, outDir string) error {
	// output directory exist? if not create it then write results to a file
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create output directory: %w", err)
	}

	filename := "results.json"
	filepath := filepath.Join(outDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("Failed to create results file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(results); err != nil {
		return fmt.Errorf("Failed to write results to file: %w", err)
	}

	return nil
}

func (dm *DatabaseManager) WriteResultsToDB(userID primitive.ObjectID, queryTitle string, results []JobPageResult) error {
	fmt.Printf("WriteResultsToDB: Processing %d job results\n", len(results))
	
	businesses := make([]database.Business, 0, len(results))
	jobs := make([]database.Job, 0, len(results))

	businessMap := make(map[string]database.Business)
	businessCounter := 0

	for _, result := range results {
		businessKey := result.BusinessName + "|" + result.URL
		business, exists := businessMap[businessKey]
		if !exists {
			business = database.Business{
				Name:    result.BusinessName,
				URL:     result.URL,
				Address: "Address not available", // placeholder address
				Lat:     0,  // coordinates will be set from geo data if available
				Lon:     0,
			}
			businessMap[businessKey] = business
			businesses = append(businesses, business)
			businessCounter++
		}

		// Create job
		job := database.Job{
			Title:       queryTitle, // query title as job title
			Description: result.Description,
			URL:         result.URL,
			PostedAt:    &[]time.Time{time.Now()}[0],
		}
		jobs = append(jobs, job)
	}

	// save businesses to database
	fmt.Printf("WriteResultsToDB: Saving %d businesses\n", len(businesses))
	businessIDs, err := dm.businessRepo.SaveBusinesses(businesses)
	if err != nil {
		return fmt.Errorf("failed to save businesses: %w", err)
	}
	fmt.Printf("WriteResultsToDB: Saved %d businesses successfully\n", len(businessIDs))

	// update jobs with business IDs by matching business names
	jobIDs := make([]primitive.ObjectID, 0, len(jobs))
	for i, job := range jobs {
		found := false
		for j, business := range businesses {
			if business.Name == results[i].BusinessName {
				job.BusinessID = businessIDs[j]
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("Warning: No business found for job result: %s\n", results[i].BusinessName)
		}
		jobs[i] = job
	}

	fmt.Printf("WriteResultsToDB: Saving %d jobs\n", len(jobs))
	jobIDs, err = dm.jobRepo.SaveJobs(jobs)
	if err != nil {
		return fmt.Errorf("failed to save jobs: %w", err)
	}
	fmt.Printf("WriteResultsToDB: Saved %d jobs successfully\n", len(jobIDs))

	fmt.Printf("WriteResultsToDB: Saving job results for user %s\n", userID.Hex())
	_, err = dm.jobResultRepo.SaveJobResults(userID, queryTitle, jobIDs)
	if err != nil {
		return fmt.Errorf("failed to save job results: %w", err)
	}
	fmt.Printf("WriteResultsToDB: Saved job results successfully\n")

	return nil
} 

func WriteGeoResults(data []byte, outDir string) error {
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}
	
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return fmt.Errorf("failed to pretty-print Geo results: %w", err)
	}

	filename := "geo_results.json"
	filePath := filepath.Join(outDir, filename)
	
	file, err := os.Create(filePath) // overwrites existing file
	if err != nil {
		return fmt.Errorf("failed to create results file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := os.WriteFile(filePath, prettyJSON.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write Geo results to file: %w", err)
	}

	return nil
}

func (dm *DatabaseManager) WriteGeoResultsToDB(userID primitive.ObjectID, zip string, radius int) (*database.GeoResult, error) {
	return dm.geoResultRepo.SaveGeoResult(userID, zip, radius)
}


// legacy function
func DeleteOldestResults(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "results_*.json"))
	if err != nil || len(files) == 0 {
		return fmt.Errorf("no result files found")
	}

	// sort files by name in ascending order
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	
	// keep only the newest file, delete the rest
	if len(files) > 1 {
		for i := 0; i < len(files)-1; i++ {
			if err := os.Remove(files[i]); err != nil {
				return fmt.Errorf("failed to delete file %s: %w", files[i], err)
			}
		}
	}
	
	return nil
}

// returns a default user ID for testing/demo purposes
func GetDefaultUserID() primitive.ObjectID {
	// this should come from user authentication
	// Using a fixed ObjectID for demo purposes so we can find saved results
	userID, _ := primitive.ObjectIDFromHex("507f1f77bcf86cd799439011")
	return userID
}
