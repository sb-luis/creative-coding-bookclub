package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sb-luis/creative-coding-bookclub/internal/routes"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

func main() {
	// Configure logger to write to stdout
	log.SetOutput(os.Stdout)

	// Initialize i18n
	utils.I18nInit()

	// Create a new custom router
	router := utils.NewRouter()

	routes.RegisterRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default to port 8000 if PORT is not set
	}

	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
