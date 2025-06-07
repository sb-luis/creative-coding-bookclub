package main

import (
	"log"
	"net/http"
	"os"

	"github.com/sb-luis/creative-coding-bookclub/internal/routes"
	"github.com/sb-luis/creative-coding-bookclub/internal/services"
	"github.com/sb-luis/creative-coding-bookclub/internal/utils"
)

var globalServices *services.Services

func main() {
	// Configure logger to write to stdout
	log.SetOutput(os.Stdout)

	// Load .env file during development only
	// In production, use environment variables that are already set
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "production" {
		if err := utils.LoadEnvFile(); err != nil {
			log.Printf("Note: Could not load .env file (this is normal in production): %v", err)
		}
	}

	// Initialize i18n
	utils.I18nInit()

	// Initialize database
	if err := utils.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	globalServices = services.NewServices(utils.GetDB())

	// Ensure database is properly closed on shutdown
	defer func() {
		if err := utils.CloseDatabase(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Create a new custom router
	router := utils.NewRouter()

	routes.RegisterRoutes(router, globalServices)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" // Default to port 8000 if PORT is not set
	}

	log.Printf("Starting server on :%s\n", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
