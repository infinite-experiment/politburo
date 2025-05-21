package main

import (
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/routes"
	"log"
	"net/http"
	"os"
	"time"

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

	var upSince = time.Now()
	log.Println("✅ Connected to Postgres!")

	router := routes.RegisterRoutes(upSince)
	loggedRouter := middleware.Logging(router)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", loggedRouter))
}
