package main

import (
	"infinite-experiment/infinite-experiment-backend/internal/routes"
	"log"
	"net/http"
)

func main() {

	router := routes.RegisterRoutes()

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
