package dtos

import "time"

type FlightMeta struct {
	Aircraft          string    `json:"aircraft"`
	Livery            string    `json:"livery"`
	MaxSpeed          int       `json:"maxSpeed"`
	MaxAlt            int       `json:"maxAlt"`
	Violations        int       `json:"violations"`
	Landings          int       `json:"landings"`
	Duration          int       `json:"durationSeconds"`
	DurationFormatted string    `json:"durationFormatted"` // HH:MM format for display
	StartedAt         time.Time `json:"startedAt"`
}

type RouteWaypoint struct {
	Lat         string `json:"lat"`
	Long        string `json:"long"`
	Timestamp   time.Time `json:"timestamp"`
	Altitude    int    `json:"alt"`
	GroundSpeed int    `json:"gs"`
	Color       string `json:"color"` // Altitude-based gradient color
}

type RouteNode struct {
	Name string `json:"name"`
	Icao string `json:"icao"`
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

type FlightInfo struct {
	Meta      FlightMeta      `json:"meta"`
	Route     []RouteWaypoint `json:"path"`
	Origin    RouteNode       `json:"origin"`
	Dest      RouteNode       `json:"dest"`
	SessionID string          `json:"sessionId"` // Session ID for the flight (for debug panel)
}

type UserFlights struct {
	Flights []FlightSummary `json:"flights"`
}
