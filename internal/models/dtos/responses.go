package dtos

import "time"

type UserStatsResponse struct {
	UserID             string  `json:"UserID"`
	XP                 int     `json:"XP"`
	FlightTime         float64 `json:"FlightTime"`
	Landings           int     `json:"Landings"`
	Violations         int     `json:"Violations"`
	Reports            int     `json:"Reports"`
	ATCOperations      int     `json:"ATCOperations"`
	VirtualAirline     string  `json:"VirtualAirline"`
	Grade              int     `json:"Grade"`
	LastFlight         string  `json:"LastFlight"` // Or time.Time if you want
	OnlineFlights      int     `json:"OnlineFlights"`
	LandingCount90Days int     `json:"LandingCount90Days"`
	// Add other fields as needed from the API response
}

// ---- USER GRADE ----
type UserGradeResponse struct {
	UserID string `json:"UserID"`
	Grade  int    `json:"Grade"`
}

// ---- ATC ----
type ATCResponse struct {
	ATC []ATCEntry `json:"ATC"`
}
type ATCEntry struct {
	ID        string  `json:"Id"`
	Type      int     `json:"Type"`
	Frequency string  `json:"Frequency"`
	Facility  int     `json:"Facility"`
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Altitude  int     `json:"Altitude"`
	Airport   string  `json:"Airport"`
	Active    bool    `json:"Active"`
	Username  string  `json:"Username"`
	UserID    string  `json:"UserID"`
}

// ---- FLIGHTS ----
type FlightsResponse struct {
	Flights []FlightEntry `json:"Flights"`
}
type FlightEntry struct {
	ID            string    `json:"Id"`
	UserID        string    `json:"UserID"`
	Username      string    `json:"Username"`
	FlightPlanID  string    `json:"FlightPlanID"`
	Server        int       `json:"Server"`
	StartTime     time.Time `json:"StartTime"`
	EndTime       time.Time `json:"EndTime"`
	AircraftID    int       `json:"AircraftID"`
	LiveryID      int       `json:"LiveryID"`
	Latitude      float64   `json:"Latitude"`
	Longitude     float64   `json:"Longitude"`
	Altitude      int       `json:"Altitude"`
	Heading       float64   `json:"Heading"`
	Speed         float64   `json:"Speed"`
	VerticalSpeed float64   `json:"VerticalSpeed"`
	Status        int       `json:"Status"`
	Origin        string    `json:"Origin"`
	Destination   string    `json:"Destination"`
	FlightPlan    string    `json:"FlightPlan"`
	// Add other fields as needed from docs
}

// ---- FLIGHT ROUTE ----
type FlightRouteResponse struct {
	Route []FlightRouteEntry `json:"Route"`
}
type FlightRouteEntry struct {
	Latitude  float64 `json:"Latitude"`
	Longitude float64 `json:"Longitude"`
	Altitude  int     `json:"Altitude"`
	Time      string  `json:"Time"` // ISO8601
}

// ---- AIRCRAFT LIVERIES ----
type AircraftLiveriesResponse struct {
	Liveries []AircraftLivery `json:"Liveries"`
}
type AircraftLivery struct {
	ID         int    `json:"Id"`
	AircraftID int    `json:"AircraftId"`
	Name       string `json:"Name"`
	ShortName  string `json:"ShortName"`
	IsDefault  bool   `json:"IsDefault"`
}

// ---- USER FLIGHTS ----
type UserFlightsResponse struct {
	Flights []UserFlightEntry `json:"Flights"`
}
type UserFlightEntry struct {
	ID           string    `json:"Id"`
	StartTime    time.Time `json:"StartTime"`
	EndTime      time.Time `json:"EndTime"`
	Origin       string    `json:"Origin"`
	Destination  string    `json:"Destination"`
	AircraftID   int       `json:"AircraftID"`
	LiveryID     int       `json:"LiveryID"`
	FlightPlanID string    `json:"FlightPlanID"`
	// Add additional fields if present
}

// ---- WORLD STATUS ----
type WorldStatusResponse struct {
	Status           string         `json:"Status"`
	Servers          []WorldServer  `json:"Servers"`
	OnlineFlights    int            `json:"OnlineFlights"`
	ActiveATC        int            `json:"ActiveATC"`
	OnlineUsers      int            `json:"OnlineUsers"`
	RecentViolations map[string]int `json:"RecentViolations"`
}

type WorldServer struct {
	ID     int    `json:"Id"`
	Name   string `json:"Name"`
	Status int    `json:"Status"`
}

// ---- ATIS ----
type ATISResponse struct {
	ATIS []ATISEntry `json:"ATIS"`
}
type ATISEntry struct {
	Airport   string `json:"Airport"`
	Frequency string `json:"Frequency"`
	Text      string `json:"Text"`
	Updated   string `json:"Updated"`
}

// --- Controller endpoints ----

type APIResponse struct {
	Status       string `json:"status"`
	Message      string `json:"message"`
	ResponseTime string `json:"response_time"`
	Data         any    `json:"data,omitempty"`
}

type InitApiResponse struct {
	IfcId                   string `json:"ifc_id"`
	IsVerificationInitiated bool   `json:"is_verification_initiated"`
	Message                 string `json:"message"`
	LastFlight              string `json:"last_flight"`
}
