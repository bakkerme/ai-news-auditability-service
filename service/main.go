package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bakkerme/ai-news-auditability-service/internal"
	"github.com/bakkerme/ai-news-auditability-service/internal/api"
	"github.com/bakkerme/ai-news-auditability-service/internal/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Initialize BadgerDB
	// Determine a base path for the database. For example, using the current working directory
	// or a specific path from an environment variable.
	basePath, err := os.Getwd() // Example: use current working directory
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	dbStoragePath := basePath // Or customize this path, e.g., "./data"

	// Load configuration
	spec, err := internal.GetConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := storage.InitDB(dbStoragePath, spec.RunDataTTLHours); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer storage.CloseDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: spec.CORSAllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	// Create API handler instance
	apiHandler := api.NewAPI(spec)

	// Routes
	api.RegisterRoutes(e, apiHandler)

	// Start server
	log.Println("Starting server on :8080")
	if err := e.Start(":8080"); err != nil {
		e.Logger.Fatal(err)
	}
}
