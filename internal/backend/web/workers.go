package web

import (
	"log"
	"sync"
	"time"
)

// a single scraping task.
type Job struct {
	BusinessName string
	URL    string
	Titles []string
}

type Result struct {
	BusinessName string
	URL          string
	JobPage      string
	Error        error
}

type WorkerPool struct {
	NumWorkers int
	Timeout    time.Duration
}

// defaults
func NewWorkerPool(numWorkers int, timeout time.Duration) *WorkerPool {
	return &WorkerPool{
		NumWorkers: numWorkers,
		Timeout:    timeout,
	}
}

// lauch scraper concurrently using ScrapeWebsite() from scraper.go.
func (wp *WorkerPool) Run(jobs []Job) []Result {
	log.Printf("Starting worker pool with %d workers for %d jobs", wp.NumWorkers, len(jobs))
	jobCh := make(chan Job, len(jobs))
	resultCh := make(chan Result, len(jobs))

	var wg sync.WaitGroup

	// launch workers -- SLOW: O(n^2) in worst case, 
	// TODO: we can optimize this by using a hashmap, but I have not approached this yet
	for i := 0; i < wp.NumWorkers; i++ {
		wg.Add(1)
		workerID := i + 1 // fix closure issue by capturing loop variable
		go func(id int) {
			defer wg.Done()
			for job := range jobCh {
				jobPage, err := ScrapeWebsite(job.URL, job.Titles)
				resultCh <- Result{
					BusinessName: job.BusinessName,
					URL:          job.URL,
					JobPage:      jobPage,
					Error:        err,
				}
			}
			log.Printf("Worker %d: Finished processing all jobs", id)
		}(workerID)
	}

	// pass jobs
	for _, j := range jobs {
		jobCh <- j
	}
	close(jobCh)

	go func() {
		wg.Wait()
		close(resultCh)
		log.Printf("All workers finished")
	}()

	results := make([]Result, 0, len(jobs))
	for res := range resultCh {
		results = append(results, res)
	}

	log.Printf("Worker pool completed! Processed %d jobs, got %d results", len(jobs), len(results))
	return results
}
