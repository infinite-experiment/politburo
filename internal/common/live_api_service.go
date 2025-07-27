package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
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

	// log.Printf("%v", req)
	// log.Printf("%v", svc.BaseURL+endpoint)

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return resp.StatusCode, errors.New("resource not found")
	}
	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	return resp.StatusCode, json.NewDecoder(resp.Body).Decode(result)
}

func (svc *LiveAPIService) doPost(
	endpoint string,
	payload interface{},
	result interface{},
) (int, error) {
	// 1) serialize body
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		return 0, err
	}

	// 2) build request
	req, err := http.NewRequest("POST", svc.BaseURL+endpoint, buf)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+svc.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// 3) log full request
	// if dumpReq, err := httputil.DumpRequestOut(req, true); err == nil {
	// 	log.Printf("→ HTTP Request:\n%s\n", dumpReq)
	// } else {
	// 	log.Printf("→ Request dump error: %v\n", err)
	// }

	// 4) do request
	resp, err := svc.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// 5) log headers (no body)
	// if dumpRespHeader, err := httputil.DumpResponse(resp, false); err == nil {
	// 	log.Printf("← HTTP Response Headers:\n%s\n", dumpRespHeader)
	// }

	// 6) read & log body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}
	// log.Printf("← HTTP Response Body:\n%s\n", string(bodyBytes))

	// 7) restore Body for JSON decode
	resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	// 8) status check + unmarshal
	switch resp.StatusCode {
	case http.StatusOK:
		return resp.StatusCode, json.NewDecoder(resp.Body).Decode(result)
	case http.StatusNotFound:
		return resp.StatusCode, errors.New("resource not found")
	default:
		return resp.StatusCode, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
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

func (svc *LiveAPIService) GetUserByIfcId(ifcId string) (*dtos.UserStatsResponse, int, error) {
	var (
		r dtos.UserStatsResponse
	)
	reqBody := dtos.LiveApiUserStatsReq{
		DiscourseNames: []string{ifcId},
	}
	status, err := svc.doPost("/users", reqBody, &r)
	if err != nil {
		return nil, status, err
	}
	return &r, status, nil
}

// Get Sessions
func (svc *LiveAPIService) GetSessions() (*dtos.SessionsResponse, error) {
	var r dtos.SessionsResponse
	_, err := svc.doGET("/sessions", &r)

	if err != nil {
		return nil, err
	}
	return &r, nil
}

// Flight Route
func (svc *LiveAPIService) GetFlightRoute(flightID string, sessionId string) (*dtos.FlightRouteResponse, int, error) {
	var r dtos.FlightRouteResponse
	status, err := svc.doGET("/sessions/"+sessionId+"/flights/"+flightID+"/route", &r)
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
func (svc *LiveAPIService) GetFlights(sId string) (*dtos.FlightsResponse, int, error) {
	var r dtos.FlightsResponse
	status, err := svc.doGET("/sessions/"+sId+"/flights", &r)
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

/**
GET https://api.infiniteflight.com/public/v2/users/813ef838-f55f-40ba-99a1-594c4c28c86f/flights?page=1
*/
// User Flights
func (svc *LiveAPIService) GetUserFlights(userID string, page int) (*dtos.UserFlightsResponse, int, error) {
	var r dtos.UserFlightsRawResponse
	status, err := svc.doGET("/users/"+userID+"/flights?page="+fmt.Sprint(page), &r)
	if err != nil {
		log.Print("---------Flights log-------------")
		log.Printf("Params: %s -> %d \n%v \n %v", userID, page, status, err)
		log.Print("---------Flights log-------------")
		return nil, status, err
	}
	return &r.Result, status, nil
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

func (svc *LiveAPIService) GetFlightPlan(sessionID, flightID string) (*dtos.FlightPlanResponse, int, error) {

	var wrap dtos.FlightPlanWrapper
	endpoint := "/sessions/" + sessionID + "/flights/" + flightID + "/flightplan"

	status, err := svc.doGET(endpoint, &wrap)
	if err != nil {
		return nil, status, err
	}
	if wrap.ErrorCode != 0 {
		return nil, status,
			fmt.Errorf("live-api returned errorCode %d", wrap.ErrorCode)
	}
	return &wrap.Result, status, nil
}

// func dumpJSONBody(r io.Reader) (pretty string, reread io.Reader, err error) {
// 	raw, err := io.ReadAll(r)
// 	if err != nil {
// 		return "", nil, err
// 	}

// 	var buf bytes.Buffer
// 	if err := json.Indent(&buf, raw, "", "  "); err != nil {
// 		// not JSON; fall back to raw
// 		pretty = string(raw)
// 	} else {
// 		pretty = buf.String()
// 	}
// 	// give the caller a new reader so they can decode again
// 	return pretty, bytes.NewReader(raw), nil
// }
