package main

import (
	"infinite-experiment/infinite-experiment-backend/internal/db"
	"infinite-experiment/infinite-experiment-backend/internal/routes"
	"infinite-experiment/infinite-experiment-backend/pkg/queries"
	"log"
	"net/http"

	_ "infinite-experiment/infinite-experiment-backend/docs"
)

// @title Infinite Experiment API
// @version 1.0
// @description Backend for Infinite Experiment bot and web client.
// @contact.name Sanket Pandia
// @contact.email sanket@example.com
// @host localhost:8080
// @BasePath /
func main() {

	// Connect to DB
	if err := db.InitPostgres(); err != nil {
		log.Fatalf("❌ Failed to connect to Postgres: %v", err)
	}

	log.Println("✅ Connected to Postgres!")

	// Load Queries
	err := queries.LoadAll()
	if err != nil {
		log.Fatalf("Failed to load queries: %v", err)
	}
	log.Println("Queries loaded!")

	router := routes.RegisterRoutes()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
