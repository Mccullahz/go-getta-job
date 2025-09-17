// starting http server using routes from server package (/internal/server/router.go)
package main

import (
    "log"
    "net/http"

    "cliscraper/internal/server"
)

func main() {
    r := server.NewRouter()
    log.Println("Server running on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
}

