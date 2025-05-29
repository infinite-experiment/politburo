package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	// dsn := os.Getenv("POSTGRES_DSN")
	dsn := "postgres://ieuser:iepass@localhost:5432/infinite?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var apiKey string
	err = db.QueryRow(
		`INSERT INTO api_keys (status) VALUES (true) RETURNING id`,
	).Scan(&apiKey)
	if err != nil {
		panic(err)
	}

	fmt.Println("New API Key:", apiKey)
}
