// this files acts as an interface between the HTTP handlers and the backend logic.
package model

import (
	"cliscraper/internal/utils"
)

// define the interface the TUI depends on, should minimize refactoring where the backend logic changes
type Service interface {
	Health() error
	Search(zip, radius, title string) ([]utils.JobPageResult, error)
	Results() ([]utils.JobPageResult, error)
	Starred() ([]utils.JobPageResult, error)
	Register(username, email, password string) (*UserInfo, error)
	Login(username, password string) (*UserInfo, error)
}

// struct for user data from authentication, send to ui
type UserInfo struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

