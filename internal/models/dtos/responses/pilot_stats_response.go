package responses

// PilotStatsResponse is the response for GET /api/v1/pilot/stats
// It combines game statistics and provider data into a unified response
type PilotStatsResponse struct {
	// Infinite Flight game statistics (from Live API) - Future implementation
	GameStats *IFGameStats `json:"game_stats,omitempty"`

	// Provider data (Airtable, Google Sheets, etc.)
	ProviderData *ProviderPilotData `json:"provider_data,omitempty"`

	// Career Mode data (if configured)
	CareerModeData *CareerModeData `json:"career_mode_data,omitempty"`

	// Recent PIREPs (flight logs) from synced data
	RecentPIREPs []RecentPIREP `json:"recent_pireps,omitempty"`

	// Metadata about the response
	Metadata PilotStatsMetadata `json:"metadata"`
}

// IFGameStats represents Infinite Flight Live API statistics
// This will be populated in a future implementation
type IFGameStats struct {
	FlightTime    int `json:"flight_time,omitempty"`
	OnlineFlights int `json:"online_flights,omitempty"`
	LandingCount  int `json:"landing_count,omitempty"`
	XP            int `json:"xp,omitempty"`
	Grade         int `json:"grade,omitempty"`
	Violations    int `json:"violations,omitempty"`
}

// ProviderPilotData contains standardized + custom fields from data provider
// Only fields marked as is_user_visible=true in the config will be included
type ProviderPilotData struct {
	// Standardized fields (all optional - only present if configured and available)
	FlightHours  *interface{} `json:"flight_hours,omitempty"`  // Can be int or float depending on provider
	Rank         *string      `json:"rank,omitempty"`          // Pilot rank/category
	JoinDate     *string      `json:"join_date,omitempty"`     // When pilot joined the VA
	LastActivity *string      `json:"last_activity,omitempty"` // Last activity date
	LastFlight   *string      `json:"last_flight,omitempty"`   // Last flight date
	Region       *string      `json:"region,omitempty"`        // Geographic region
	TotalFlights *int         `json:"total_flights,omitempty"` // Number of flights
	Status       *string      `json:"status,omitempty"`        // Active/inactive status

	// All other fields that don't map to standard names
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`
}

// CareerModeData contains career mode specific data from the provider
type CareerModeData struct {
	// Standardized career mode fields (all optional)
	TotalCMHours              *interface{} `json:"total_cm_hours,omitempty"`               // Career mode hours completed
	RequiredHoursToNext       *interface{} `json:"required_hours_to_next,omitempty"`       // Hours needed for next aircraft
	LastActivityCM            *string      `json:"last_activity_cm,omitempty"`             // Last career mode activity
	AssignedRoutes            *interface{} `json:"assigned_routes,omitempty"`              // Assigned flight routes (can be array)
	Aircraft                  *string      `json:"aircraft,omitempty"`                     // Current aircraft
	Airline                   *string      `json:"airline,omitempty"`                      // Current airline
	LastFlownRoute            *string      `json:"last_flown_route,omitempty"`             // Last PIREP route
	LastCareerModePIREP       *interface{} `json:"last_career_mode_pirep,omitempty"`       // Last PIREP log reference (Airtable IDs)
	LastCareerModeFlight      *string      `json:"last_career_mode_flight,omitempty"`      // Last career mode flight route (enriched from route_at_synced)

	// All other career mode fields that don't map to standard names
	AdditionalFields map[string]interface{} `json:"additional_fields,omitempty"`
}

// PilotStatsMetadata provides context about the data source and freshness
type PilotStatsMetadata struct {
	ProviderType       string `json:"provider_type,omitempty"`   // e.g., "airtable", "google_sheets"
	ProviderConfigured bool   `json:"provider_configured"`       // Whether a provider is configured for this VA
	SchemaVersion      string `json:"schema_version,omitempty"`  // Config schema version
	LastFetched        string `json:"last_fetched"`              // ISO 8601 timestamp
	Cached             bool   `json:"cached"`                    // Whether data came from cache
	VAName             string `json:"va_name,omitempty"`         // Name of the virtual airline
}

// RecentPIREP represents a recent PIREP (flight log) record
type RecentPIREP struct {
	ATID          string   `json:"at_id"`                    // Airtable record ID
	Route         string   `json:"route"`                    // Flight route (e.g., "KLAX-KSFO")
	FlightMode    string   `json:"flight_mode,omitempty"`    // Flight mode (e.g., "Casual", "Expert")
	FlightTime    *float64 `json:"flight_time,omitempty"`    // Flight duration in hours
	PilotCallsign string   `json:"pilot_callsign,omitempty"` // Pilot callsign
	Aircraft      string   `json:"aircraft,omitempty"`       // Aircraft type (e.g., "B738")
	Livery        string   `json:"livery,omitempty"`         // Aircraft livery/airline
	ATCreatedTime *string  `json:"at_created_time,omitempty"` // Airtable creation timestamp
}
