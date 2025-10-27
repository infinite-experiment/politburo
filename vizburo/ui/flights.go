package ui

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Flight represents a single flight record from the API
type Flight struct {
	ID           string  `json:"id"`
	CallSign     string  `json:"call_sign"`
	Aircraft     string  `json:"aircraft"`
	Livery       string  `json:"livery"`
	Route        string  `json:"route"`
	Distance     float64 `json:"distance"`
	MaxSpeed     float64 `json:"max_speed"`
	MaxAlt       float64 `json:"max_alt"`
	Duration     int     `json:"duration_seconds"`
	Violations   int     `json:"violations"`
	Landings     int     `json:"landings"`
	StartedAt    string  `json:"started_at"`
}

// FlightPath represents a waypoint on the flight path
type FlightPath struct {
	Lat  string `json:"lat"`
	Long string `json:"long"`
}

// FlightDetails represents detailed flight information with path
type FlightDetails struct {
	ID       string        `json:"id"`
	Meta     Flight        `json:"meta"`
	Path     []FlightPath  `json:"path"`
}

// FlightsListHandler returns the main flights page with sidebar and map
func (h *UIHandler) FlightsListHandler(w http.ResponseWriter, r *http.Request) {
	flightsListContent := `
<div class="flex flex-col h-full w-full gap-4">
    <!-- Map Container (Top 3/4) -->
    <div id="map-container" class="h-3/4 w-full bg-gray-100 dark:bg-gray-800 rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700 shadow-md">
        <div id="map" class="w-full h-full" style="position: relative;">
            <div class="w-full h-full flex items-center justify-center text-gray-500 dark:text-gray-400">
                <div class="text-center">
                    <p class="mb-4 text-5xl">📍</p>
                    <p class="text-xl font-semibold text-gray-700 dark:text-gray-300">Select Flight</p>
                    <p class="text-sm mt-2 text-gray-600 dark:text-gray-400">Choose a flight from the list below to view its route</p>
                </div>
            </div>
        </div>
    </div>

    <!-- Bottom Container (Bottom 1/4) - Split into two columns -->
    <div class="h-1/4 w-full flex gap-4 overflow-hidden">
        <!-- Left Panel: Flight List (Scrollable) -->
        <div id="flight-list-panel" class="flex-1 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col shadow-md">
            <div class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
                <h2 class="text-sm font-semibold text-gray-900 dark:text-white uppercase tracking-wide">Recent Flights</h2>
            </div>
            <div class="flex-1 overflow-y-auto p-2">
                <!-- Flight list will be loaded via HTMX -->
                <div id="flight-list" hx-get="/flights/api/list" hx-trigger="load" hx-swap="innerHTML">
                    <div class="text-center py-8 text-gray-500 dark:text-gray-400">
                        <p class="text-sm">Loading flights...</p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Right Panel: Flight Details (Card-based) -->
        <div id="flight-details-panel" class="flex-1 bg-gray-50 dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col shadow-md">
            <div class="px-4 py-3 border-b border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
                <h2 class="text-sm font-semibold text-gray-900 dark:text-white uppercase tracking-wide">Flight Details</h2>
            </div>
            <div class="flex-1 overflow-y-auto p-4">
                <p class="text-gray-600 dark:text-gray-400 text-center py-8">Select a flight from the list to see details</p>
            </div>
        </div>
    </div>
</div>

<script type="importmap">
{
  "imports": {
    "gleo/": "https://unpkg.com/gleo/src/"
  }
}
</script>
<script src="/static/js/flight-map.js" type="module"></script>

<script>
// Handle flight selection from list via event delegation
// Wait for DOM to be ready and flight-map.js to load
document.addEventListener('DOMContentLoaded', function() {
    document.addEventListener('click', function(event) {
        const flightElement = event.target.closest('[data-flight-id]');
        if (flightElement) {
            const flightId = flightElement.getAttribute('data-flight-id');
            console.log('Flight clicked:', flightId);
            if (typeof window.loadFlightDetails === 'function') {
                window.loadFlightDetails(flightId);
            } else {
                console.error('loadFlightDetails not available');
            }
        }
    });
});
</script>
`

	data := map[string]interface{}{
		"Title":   "Flight Visualization",
		"Content": flightsListContent,
		"Theme":   getThemeFromRequest(r),
	}
	RenderTemplate(w, "layouts/sidebar.html", data)
}

// FlightsListAPIHandler returns a list of flights as HTML (for HTMX)
func (h *UIHandler) FlightsListAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Mock flight data - in production, this would fetch from the API
	flights := []Flight{
		{
			ID:        "1",
			CallSign:  "DEMO001",
			Aircraft:  "Boeing 777-300ER",
			Livery:    "Air France",
			Route:     "LFPG-EGLL",
			Distance:  344,
			MaxSpeed:  450,
			MaxAlt:    43000,
			Duration:  8100,
			Violations: 0,
			Landings:  1,
			StartedAt: "2024-10-26T10:30:00Z",
		},
		{
			ID:        "2",
			CallSign:  "DEMO002",
			Aircraft:  "Airbus A380",
			Livery:    "Emirates",
			Route:     "KJFK-LFPG",
			Distance:  3626,
			MaxSpeed:  488,
			MaxAlt:    43000,
			Duration:  30600,
			Violations: 1,
			Landings:  1,
			StartedAt: "2024-10-25T14:15:00Z",
		},
	}

	html := `<div class="space-y-2">`
	for _, flight := range flights {
		violationBadge := ""
		if flight.Violations > 0 {
			violationBadge = fmt.Sprintf(`<span class="text-xs bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400 px-2 py-1 rounded-full">%d</span>`, flight.Violations)
		}

		html += fmt.Sprintf(`
<div data-flight-id="%s" class="p-3 bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-flight-primary dark:hover:border-flight-accent hover:shadow-lg cursor-pointer transition-all duration-200 hover:scale-[1.02]">
    <div class="flex items-start justify-between">
        <div class="flex-1 min-w-0">
            <div class="font-bold text-sm text-gray-900 dark:text-white truncate">%s</div>
            <div class="text-xs text-gray-600 dark:text-gray-400 mt-0.5 truncate">%s</div>
        </div>
        %s
    </div>
    <div class="flex items-center justify-between mt-2">
        <div class="text-xs text-gray-500 dark:text-gray-400">%s</div>
        <div class="text-xs font-semibold text-flight-primary">%.0f nm</div>
    </div>
</div>
`, flight.ID, flight.CallSign, flight.Aircraft, violationBadge, flight.Route, flight.Distance)
	}
	html += `</div>`

	w.Write([]byte(html))
}

// FlightsDetailsAPIHandler returns flight details with path as JSON
func (h *UIHandler) FlightsDetailsAPIHandler(w http.ResponseWriter, r *http.Request) {
	flightID := r.PathValue("id")

	// Mock flight data with path - in production, this would fetch from the API
	details := FlightDetails{
		ID: flightID,
		Meta: Flight{
			ID:        flightID,
			CallSign:  "DEMO001",
			Aircraft:  "Boeing 777-300ER",
			Livery:    "Air France",
			Route:     "LFPG-EGLL",
			Distance:  344,
			MaxSpeed:  450,
			MaxAlt:    43000,
			Duration:  8100,
			Violations: 0,
			Landings:  1,
			StartedAt: "2024-10-26T10:30:00Z",
		},
		Path: []FlightPath{
			{Lat: "48.9627", Long: "2.5500"},   // CDG
			{Lat: "48.8500", Long: "2.3000"},   // Mid-point
			{Lat: "51.4700", Long: "-0.0000"},  // LHR
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}
