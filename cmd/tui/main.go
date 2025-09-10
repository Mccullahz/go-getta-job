package main

import (
    "cliscraper/internal/api"
    "cliscraper/internal/ui"
    //tea "github.com/charmbracelet/bubbletea"
)

func main() {
    client := api.NewClient("http://localhost:8080")
    ui.Run(client)
}

