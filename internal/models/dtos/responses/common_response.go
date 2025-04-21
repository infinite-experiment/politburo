package responses

import "time"

type APIResponse[T any] struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"error,omitempty"`
	Data      *T        `json:"data,omitempty"`
}
