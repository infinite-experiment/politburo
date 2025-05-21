package api

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/models/dtos/responses"
	"net/http"
	"time"
)

func respondWithSuccess[T any](w http.ResponseWriter, statusCode int, data *T) {
	resp := responses.APIResponse[T]{
		Status:    "success",
		Timestamp: time.Now().UTC(),
		Data:      data,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(resp)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	resp := responses.APIResponse[any]{
		Status:    "error",
		Timestamp: time.Now().UTC(),
		Error:     message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(resp)
}
