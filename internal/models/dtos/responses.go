package dtos

import "time"

// dtos/user_stats.go
type UserStatsResponse struct {
	ErrorCode int         `json:"errorCode"`
	Result    []UserStats `json:"result"`
}

type UserStats struct {
	OnlineFlights         int                   `json:"onlineFlights"`
	Violations            int                   `json:"violations"`
	XP                    int                   `json:"xp"`
	LandingCount          int                   `json:"landingCount"`
	FlightTime            int                   `json:"flightTime"`
	ATCOperations         int                   `json:"atcOperations"`
	ATCRank               *int                  `json:"atcRank"` // nullable
	Grade                 int                   `json:"grade"`
	Hash                  string                `json:"hash"`
	ViolationCountByLevel ViolationCountByLevel `json:"violationCountByLevel"`
	Roles                 []int                 `json:"roles"`
	UserID                string                `json:"userId"`
	VirtualOrganization   *string               `json:"virtualOrganization"` // nullable
	DiscourseUsername     *string               `json:"discourseUsername"`   // nullable
	Groups                []string              `json:"groups"`
	ErrorCode             int                   `json:"errorCode"`
}

type ViolationCountByLevel struct {
	Level1 int `json:"level1"`
	Level2 int `json:"level2"`
	Level3 int `json:"level3"`
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

// ---- AIRCRAFT LIVERIES ----
type AircraftLiveriesResponse struct {
	Liveries  []AircraftLivery `json:"result"`
	ErrorCode int              `json:"errorCode"`
}
type AircraftLivery struct {
	LiveryId     string `json:"id"`
	AircraftID   string `json:"aircraftID"`
	LiveryName   string `json:"liveryName"`
	AircraftName string `json:"aircraftName"`
}

// ---- USER FLIGHTS ----

type UserFlightEntry struct {
	ID                 string    `json:"id"`
	Created            time.Time `json:"created"`
	UserID             string    `json:"userId"`
	AircraftID         string    `json:"aircraftId"`
	LiveryID           string    `json:"liveryId"`
	Callsign           string    `json:"callsign"`
	Server             string    `json:"server"`
	DayTime            float32   `json:"dayTime"`
	NightTime          float32   `json:"nightTime"`
	TotalTime          float32   `json:"totalTime"`
	LandingCount       int       `json:"landingCount"`
	OriginAirport      string    `json:"originAirport"`
	DestinationAirport string    `json:"destinationAirport"`
	XP                 int       `json:"xp"`
	WorldType          int       `json:"worldType"`
	Violations         []any     `json:"violations"` // If you know violation type, replace 'any'
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
	IfcId   string             `json:"ifc_id"`
	Status  bool               `json:"status"`
	Message string             `json:"message"`
	Steps   []RegistrationStep `json:"steps"`
}

type RegistrationStep struct {
	Name    string `json:"name"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type InitServerResponse struct {
	VACode  string             `json:"va_code"`
	Status  bool               `json:"status"`
	Message string             `json:"message,omitempty"`
	Steps   []RegistrationStep `json:"steps"`
}

type UserFlightsResponse struct {
	PageIndex   int               `json:"pageIndex"`
	TotalPages  int               `json:"totalPages"`
	TotalCount  int               `json:"totalCount"`
	HasPrevious bool              `json:"hasPreviousPage"`
	HasNext     bool              `json:"hasNextPage"`
	Flights     []UserFlightEntry `json:"data"`
}

type UserFlightsRawResponse struct {
	ErrorCode int                 `json:"errorCode"`
	Result    UserFlightsResponse `json:"result"`
}

type HistoryRecord struct {
	Origin     string    `json:"origin"`
	Dest       string    `json:"dest"`
	TimeStamp  time.Time `json:"timestamp"`
	EndTime    time.Time `json:"endtime"`
	Landings   int       `json:"landings"`
	Server     string    `json:"server"`
	Aircraft   string    `json:"aircraft"`
	Livery     string    `json:"livery"`
	MapUrl     string    `json:"mapUrl"`
	Callsign   string    `json:"callsign"`
	Violations int       `json:"violations"`
	Equipment  string    `json:"equipment"`
	Duration   string    `json:"duration"`
}

type FlightHistoryDto struct {
	PageNo  int             `json:"page"`
	Records []HistoryRecord `json:"records"`
	Error   string          `json:"error"`
}

type SessionsResponse struct {
	ErrorCode int       `json:"errorCode"`
	Result    []Session `json:"result"`
}

type Session struct {
	MaxUsers          int     `json:"maxUsers"`
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	UserCount         int     `json:"userCount"`
	Type              int     `json:"type"`
	WorldType         int     `json:"worldType"`
	MinimumGradeLevel int     `json:"minimumGradeLevel"`
	MinimumAppVersion string  `json:"minimumAppVersion"`
	MaximumAppVersion *string `json:"maximumAppVersion"` // nullable
}

type FlightRouteResponse struct {
	ErrorCode int              `json:"errorCode"`
	Result    []FlightPosition `json:"result"`
}

type FlightPosition struct {
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Altitude    float64   `json:"altitude"`
	Track       float64   `json:"track"`
	GroundSpeed float64   `json:"groundSpeed"`
	Date        time.Time `json:"date"`
	ID          *string   `json:"id"`
	FID         string    `json:"fid"`
}
