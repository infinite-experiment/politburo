package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/routes"

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

	// Connect to DB with sqlx
	if err := db.InitPostgres(); err != nil {
		log.Fatalf("❌ Failed to connect to Postgres (sqlx): %v", err)
	}
	log.Println("✅ Connected to Postgres (sqlx)!")

	// Connect to DB with GORM
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	dbname := os.Getenv("PG_DB")
	password := os.Getenv("PG_PASSWORD")
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	if _, err := db.InitPostgresORM(dsn); err != nil {
		log.Fatalf("❌ Failed to connect to Postgres (GORM): %v", err)
	}
	log.Println("✅ Connected to Postgres (GORM)!")

	upSince := time.Now()

	// Initialize router with Chi
	router := routes.RegisterRoutes(upSince)

	// Attach logging middleware (already compatible)
	//loggedRouter := middleware.Logging(router)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
