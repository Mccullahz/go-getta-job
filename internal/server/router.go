// wiring handlers to routes for the api
package server

import (
	"fmt"
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

	// API routes -- should work for the things we have implemented so far 
	r.Get("/health", HealthHandler)
	r.Get("/search", SearchHandler)
	r.Get("/results", ResultsHandler)
	r.Get("/starred", StarredHandler)

	return r
}

// mongo handlers
func NewDatabaseRouter() (http.Handler, error) {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	fmt.Println("Creating database handlers...")
	dbHandlers, err := NewDatabaseHandlers()
	if err != nil {
		fmt.Printf("Failed to create database handlers: %v\n", err)
		return nil, fmt.Errorf("failed to create database handlers: %w", err)
	}
	fmt.Println("Database handlers created successfully")

	// cleanup middleware to close database connection
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			next.ServeHTTP(w, req)
		})
	})

	// API routes with database handlers
	r.Get("/health", HealthHandler)
	r.Get("/search", dbHandlers.SearchHandlerDB)
	r.Get("/results", dbHandlers.ResultsHandlerDB)
	r.Get("/starred", dbHandlers.StarredHandlerDB)

	return r, nil
}

