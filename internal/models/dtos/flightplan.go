package dtos

import "time"

// -------- helper types -----------------------------------------------------

type APITime struct{ time.Time }

const apiLayout = "2006-01-02 15:04:05Z07:00"

func (t *APITime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1] // strip quotes
	if s == "" || s == "null" {
		return nil
	}
	tt, err := time.Parse(apiLayout, s)
	if err != nil {
		return err
	}
	t.Time = tt
	return nil
}

// -------- main DTOs --------------------------------------------------------

type FlightPlanWrapper struct {
	ErrorCode int                `json:"errorCode"`
	Result    FlightPlanResponse `json:"result"`
}

type FlightPlanResponse struct {
	FlightPlanID    string           `json:"flightPlanId"`
	FlightID        string           `json:"flightId"`
	Waypoints       []string         `json:"waypoints"`
	LastUpdate      APITime          `json:"lastUpdate"`
	FlightPlanItems []FlightPlanItem `json:"flightPlanItems"`
}

type FlightPlanItem struct {
	Name       string           `json:"name"`
	Type       int              `json:"type"`
	Children   []FlightPlanItem `json:"children"`   // nil when no children
	Identifier *string          `json:"identifier"` // nullable
	Altitude   float64          `json:"altitude"`
	Location   Location         `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Altitude  float64 `json:"altitude"`
}

type FlightSummary struct {
	FlightID    string
	Origin      string
	Destination string
	Aircraft    string
	Livery      string
}
