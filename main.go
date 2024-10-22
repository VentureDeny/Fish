package main

import (
	"fish/config"
	"fish/db"
	"fish/router"
	"log"
	"net/http"
	"time"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize the database
	db.InitDB(cfg)
	defer db.CloseDB()

	// Initialize routes
	r := router.InitializeRoutes()

	// Create the HTTP server
	server := &http.Server{
		Addr:         cfg.ServerAddr,
		Handler:      r, // Use the router as the handler
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Server is running on %s", cfg.ServerAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
