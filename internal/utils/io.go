// io for managing results files
package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"

	"cliscraper/internal/output"
)

func LoadLatestResults(dir string) ([]output.JobPageResult, error) {
	files, err := filepath.Glob(filepath.Join(dir, "results_*.json"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("no result files found")
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i] > files[j]
	})

	data, err := ioutil.ReadFile(files[0])
	if err != nil {
		return nil, err
	}

	var results []output.JobPageResult
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func DeleteOldestResults(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "results_*.json"))
	if err != nil || len(files) == 0 {
		return fmt.Errorf("no result files found")
	}

	// Sort files by name in ascending order
	sort.Slice(files, func(i, j int) bool {
		return files[i] < files[j]
	})
	// Keep only the newest file
	if len(files) > 1 {

	}
		return nil
}

