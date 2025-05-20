package entities

import "time"

type ServiceStatus struct {
	Status  string `json:"status"`
	Details string `json:"details"`
}

type HealthCheckResponse struct {
	Status   string                   `json:"status"`
	Services map[string]ServiceStatus `json:"services"`
	UpSince  time.Time                `json:"up_since"`
	Uptime   string                   `json:"uptime"`
}
