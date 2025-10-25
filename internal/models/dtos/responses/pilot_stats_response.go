package responses

// PilotStatsResponse is the response for GET /api/v1/pilot/stats
// It combines game statistics and provider data into a unified response
type PilotStatsResponse struct {
	// Infinite Flight game statistics (from Live API) - Future implementation
	GameStats *IFGameStats `json:"game_stats,omitempty"`

	// Provider data (Airtable, Google Sheets, etc.)
	ProviderData *ProviderPilotData `json:"provider_data,omitempty"`

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

// PilotStatsMetadata provides context about the data source and freshness
type PilotStatsMetadata struct {
	ProviderType       string `json:"provider_type,omitempty"`   // e.g., "airtable", "google_sheets"
	ProviderConfigured bool   `json:"provider_configured"`       // Whether a provider is configured for this VA
	SchemaVersion      string `json:"schema_version,omitempty"`  // Config schema version
	LastFetched        string `json:"last_fetched"`              // ISO 8601 timestamp
	Cached             bool   `json:"cached"`                    // Whether data came from cache
	VAName             string `json:"va_name,omitempty"`         // Name of the virtual airline
}
