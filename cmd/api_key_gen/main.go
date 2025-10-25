package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://ieuser:iepass@localhost:5432/infinite?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	var id int64 // adapt to your column type
	if err := db.QueryRow(`INSERT INTO api_keys (status) VALUES (true) RETURNING id`).Scan(&id); err != nil {
		log.Fatalf("insert api key: %v", err)
	}

	fmt.Println("New API Key:", id)
}
