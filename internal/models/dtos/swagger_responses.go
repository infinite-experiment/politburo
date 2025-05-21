package dtos

import (
	"infinite-experiment/politburo/internal/models/dtos/responses"
)

// UserRegisterSwaggerResponse wraps the actual response inside your API response format.
// This struct is ONLY used for Swagger documentation.
type UserRegisterSwaggerResponse struct {
	Status    string                         `json:"status"`
	Timestamp string                         `json:"timestamp"`
	Error     string                         `json:"error,omitempty"`
	Data      responses.UserRegisterResponse `json:"data"`
}
