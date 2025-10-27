package common

import "context"

// FlightData represents current flight information for a user
type FlightData struct {
	Callsign     string `json:"callsign"`
	IFCUsername  string `json:"ifc_username"`
	Aircraft     string `json:"aircraft"`
	Livery       string `json:"livery"`
	LiveryID     string `json:"livery_id"`
	Route        string `json:"route"`
	Altitude     int    `json:"altitude"`   // Altitude in feet
	Speed        int    `json:"speed"`      // Speed in knots
	Multiplier   float64 `json:"multiplier"` // Mode multiplier
}

// GetUserFlight retrieves the current flight for a user by searching sessions
// The method looks for a flight matching the user's callsign across all active sessions
// Note: Aircraft, Livery, and Route information should be enriched by the caller
func (svc *LiveAPIService) GetUserFlight(ctx context.Context, callsign string) (*FlightData, error) {
	// Get all active sessions
	sessions, err := svc.GetSessions()
	if err != nil {
		return nil, err
	}

	if sessions == nil || len(sessions.Result) == 0 {
		return nil, nil
	}

	// Search for a flight matching the callsign
	for _, session := range sessions.Result {
		flights, _, err := svc.GetFlights(session.ID)
		if err != nil {
			continue
		}

		if flights == nil || len(flights.Flights) == 0 {
			continue
		}

		for _, flight := range flights.Flights {
			// Check if callsign matches
			if flight.Callsign == callsign {
				return &FlightData{
					Callsign:    flight.Callsign,
					IFCUsername: flight.Username,
					Aircraft:    "", // Will be enriched by handler or from aircraft liveries table
					Livery:      "", // Will be enriched by handler or from aircraft liveries table
					LiveryID:    flight.LiveryID,
					Route:       "", // Will be determined from flight plan or other sources
				}, nil
			}
		}
	}

	return nil, nil
}
