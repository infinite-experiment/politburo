package common

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"log"
	"net/http"
	"time"
)

// RespondSuccess sends a standardized JSON success response.
func RespondSuccess(w http.ResponseWriter, initTime time.Time, message string, data any, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	response := dtos.APIResponse{
		Status:       string(constants.APIStatusOk),
		Message:      message,
		ResponseTime: GetResponseTime(initTime),
		Data:         data,
	}

	writeJSON(w, code, response)
}

// RespondError sends a standardized JSON error response.
func RespondError(w http.ResponseWriter, initTime time.Time, err error, message string, statusCode ...int) {
	code := http.StatusInternalServerError
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	msg := message
	if err != nil && err.Error() != "" {
		msg = err.Error()
	}

	response := dtos.APIResponse{
		Status:       string(constants.APIStatusError),
		Message:      msg,
		ResponseTime: GetResponseTime(initTime),
	}

	writeJSON(w, code, response)
}

// writeJSON marshals data and writes it to the HTTP response.
func writeJSON(w http.ResponseWriter, code int, body dtos.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("JSON encode failed: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
