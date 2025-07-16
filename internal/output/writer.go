// writes discovered job pages to a local results.json.
package output

import (
	"encoding/json"
	"os"
	"fmt"
	"time"
)
// result struct to hold job page details
type JobPageResult struct {
	BusinessName string `json:"business_name"`
	URL	   string `json:"url"`
	Description string `json:"description"`
}

func WriteResults(results []JobPageResult, outDir string) error {
	// output directory exist? if not create it then write results to a file
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		return fmt.Errorf("Failed to create output directory: %w", err)
	}

	filename := fmt.Sprintf("%s/results_%d.json", outDir, time.Now().Unix())
	file, err := os.Create(filename)
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

