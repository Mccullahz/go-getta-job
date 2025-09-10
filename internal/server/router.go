package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// set up all routes for the API server
func NewRouter() http.Handler {
	r := chi.NewRouter()

	// middleware probablt want logging, recovery, etc, can adjust later 
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	// API routes
	r.Get("/health", HealthHandler)
	r.Post("/search", SearchHandler)
	r.Get("/results/{id}", ResultsHandler) // fetch results by search id
	r.Get("/starred", StarredHandler)      // fetch all starred jobs

	return r
}

