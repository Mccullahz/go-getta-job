// this files acts as an interface between the HTTP handlers and the backend logic.
package model

import (
	"cliscraper/internal/utils"
)

// define the interface the TUI depends on, should minimize refactoring where the backend logic changes
type Service interface {
	Health() error
	Search(zip, radius, title string) ([]utils.JobPageResult, error)
	Results(id string) ([]utils.JobPageResult, error)
	Starred() ([]utils.JobPageResult, error)
}

