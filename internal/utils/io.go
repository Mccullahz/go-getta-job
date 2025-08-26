// io for managing results files
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"os"
	//"time"

)

// result struct to hold job page details
type JobPageResult struct {
	BusinessName string `json:"business_name"`
	URL	   string `json:"url"`
	Description string `json:"description"`
}

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

//lets make this better, ensure we are writing the files similar to the results file above
func WriteGeoResults(data []byte, outDir string) error {
	// Ensure output directory exists
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	
	}
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return fmt.Errorf("failed to pretty-print Geo results: %w", err)
	}

	filename := "geo_results.json"
	filePath := filepath.Join(outDir, filename)
	
	file, err := os.Create(filePath) // os.Create overwrites existing file
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

func DeleteOldestResults(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "results_*.json"))
	if err != nil || len(files) == 0 {
		return fmt.Errorf("no result files found")
	}

	// sort files by name in ascending order
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	// keep only the newest file
	if len(files) > 1 {

	}
		return nil
}

