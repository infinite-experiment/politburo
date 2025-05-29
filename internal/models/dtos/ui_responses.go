package dtos

import "time"

type FlightMeta struct {
	Aircraft   string    `json:"aircraft"`
	Livery     string    `json:"livery"`
	MaxSpeed   float32   `json:"maxSpeed"`
	Violations int       `json:"violations"`
	Landings   int       `json:"landings"`
	Duration   int       `json:"durationSeconds"`
	StartedAt  time.Time `json:"startedAt"`
}

type RouteWaypoint struct {
	Lat       string    `json:"lat"`
	Long      string    `json:"long"`
	Timestamp time.Time `json:"timestamp"`
	Altitude  int       `json:"alt"`
}

type RouteNode struct {
	Name string `json:"name"`
	Icao string `json:"icao"`
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

type FlightInfo struct {
	Meta   FlightMeta      `json:"meta"`
	Route  []RouteWaypoint `json:"path"`
	Origin RouteNode       `json:"origin"`
	Dest   RouteNode       `json:"dest"`
}
