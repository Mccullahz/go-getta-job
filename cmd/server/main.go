// starting http server using routes from server package (/internal/server/router.go)
package main

import (
    "log"
    "net/http"
    "os"

    "cliscraper/internal/server"
)

func main() {
    // check if should use database mode
    useDB := os.Getenv("USE_DATABASE")
    
    var r http.Handler
    var err error
    
    if useDB == "true" {
        log.Println("Starting server with MongoDB support...")
        r, err = server.NewDatabaseRouter()
        if err != nil {
            log.Fatal("Failed to create database router:", err)
        }
    } else {
        log.Println("Starting server with file-based storage...")
        r = server.NewRouter()
    }
    
    log.Println("Server running on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}

