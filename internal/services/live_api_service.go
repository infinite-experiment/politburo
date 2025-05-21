package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"infinite-experiment/politburo/internal/models/dtos"
)

type LiveAPIService struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// NewLiveAPIService creates a new instance, reading config from environment variables
func NewLiveAPIService() *LiveAPIService {
	baseURL := os.Getenv("IF_API_BASE_URL")
	if baseURL == "" {
		baseURL = "https://api.infiniteflight.com/public/v2" // Default
	}
	apiKey := os.Getenv("IF_API_KEY")
	client := &http.Client{Timeout: 10 * time.Second}
	return &LiveAPIService{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client:  client,
	}
}

// helper: does GET with auth header, parses json into result, returns status code
func (svc *LiveAPIService) doGET(endpoint string, result interface{}) (int, error) {
	req, err := http.NewRequest("GET", svc.BaseURL+endpoint, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+svc.APIKey)

	resp, err := svc.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, errors.New("resource not found")
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(result)
}

// User Grade
func (svc *LiveAPIService) GetUserGrade(userID string) (*dtos.UserGradeResponse, int, error) {
	var r dtos.UserGradeResponse
	status, err := svc.doGET("/user/grade/"+userID, &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// ATC
func (svc *LiveAPIService) GetATC() (*dtos.ATCResponse, int, error) {
	var r dtos.ATCResponse
	status, err := svc.doGET("/atc", &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// Flights
func (svc *LiveAPIService) GetFlights() (*dtos.FlightsResponse, int, error) {
	var r dtos.FlightsResponse
	status, err := svc.doGET("/flights", &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// Flight Route
func (svc *LiveAPIService) GetFlightRoute(flightID string) (*dtos.FlightRouteResponse, int, error) {
	var r dtos.FlightRouteResponse
	status, err := svc.doGET("/flight/route/"+flightID, &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// Aircraft Liveries
func (svc *LiveAPIService) GetAircraftLiveries() (*dtos.AircraftLiveriesResponse, int, error) {
	var r dtos.AircraftLiveriesResponse
	status, err := svc.doGET("/aircraft/liveries", &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// User Flights
func (svc *LiveAPIService) GetUserFlights(userID string) (*dtos.UserFlightsResponse, int, error) {
	var r dtos.UserFlightsResponse
	status, err := svc.doGET("/user/flights/"+userID, &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// World Status
func (svc *LiveAPIService) GetWorldStatus() (*dtos.WorldStatusResponse, int, error) {
	var r dtos.WorldStatusResponse
	status, err := svc.doGET("/world/status", &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// ATIS
func (svc *LiveAPIService) GetATIS() (*dtos.ATISResponse, int, error) {
	var r dtos.ATISResponse
	status, err := svc.doGET("/atis", &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}
