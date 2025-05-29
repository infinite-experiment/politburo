package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/routes"
	"infinite-experiment/politburo/internal/workers"

	// Swagger docs
	_ "infinite-experiment/politburo/docs"
)

// @title Infinite Experiment API
// @version 1.0
// @description Backend for Infinite Experiment bot and web client.
// @contact.name Sanket Pandia
// @contact.email sanket@example.com
// @host localhost:8080
// @BasePath /
func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Connect to DB
	if err := db.InitPostgres(); err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}

	upSince := time.Now()
	log.Println("✅ Connected to Postgres!")

	// Initialize router with Chi
	router := routes.RegisterRoutes(upSince)

	// Attach logging middleware (already compatible)
	loggedRouter := middleware.Logging(router)

	go workers.LogbookWorker()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
