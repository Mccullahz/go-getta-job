package main

import (
	"fmt"
	"cliscraper/internal/api"
)

func main() {
	client := api.NewClient("http://localhost:8080")
	
	fmt.Println("Testing search...")
	results, err := client.Search("45150", "2", "job")
	if err != nil {
		fmt.Printf("Search error: %v\n", err)
		return
	}
	
	fmt.Printf("Search returned %d results\n", len(results))
	for i, result := range results {
		if i < 3 { // show first 3 results
			fmt.Printf("Result %d: %s - %s\n", i+1, result.BusinessName, result.URL)
		}
	}
	
	fmt.Println("\nTesting results...")
	results2, err := client.Results()
	if err != nil {
		fmt.Printf("Results error: %v\n", err)
		return
	}
	
	fmt.Printf("Results returned %d items\n", len(results2))
	for i, result := range results2 {
		if i < 3 { // show first 3 results
			fmt.Printf("Result %d: %s - %s\n", i+1, result.BusinessName, result.URL)
		}
	}
}