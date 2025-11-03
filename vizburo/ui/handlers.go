package ui

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/services"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// This package contains authentication and UI handlers for Vizburo.
// See auth.go for authentication handlers.

// DashboardHandler serves the main dashboard page for authenticated users
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context (injected by auth middleware)
	claims := auth.GetUserClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Get session data from context
	sessionDataInterface := auth.GetSessionData(r.Context())
	if sessionDataInterface == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Cast to SessionData
	sessionData, ok := sessionDataInterface.(*common.SessionData)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	// Get active VA
	activeVA := sessionData.GetActiveVA()

	// Prepare template data
	data := map[string]interface{}{
		"ActiveVA":        activeVA,
		"VirtualAirlines": sessionData.VirtualAirlines,
		"Username":        sessionData.Username,
		"UserID":          sessionData.UserID,
		"ActiveVAID":      sessionData.ActiveVAID,
		"PageTitle":       "Dashboard",
	}

	// Render template
	if err := RenderTemplate(w, "pages/dashboard.html", data); err != nil {
		http.Error(w, "Error rendering dashboard", http.StatusInternalServerError)
		return
	}
}

// LogbookHandler serves the logbook page for staff and admin users
func LogbookHandler(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims := auth.GetUserClaims(r.Context())
	if claims == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Get session data from context
	sessionDataInterface := auth.GetSessionData(r.Context())
	if sessionDataInterface == nil {
		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
		return
	}

	// Cast to SessionData
	sessionData, ok := sessionDataInterface.(*common.SessionData)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	// Get active VA
	activeVA := sessionData.GetActiveVA()
	if activeVA == nil {
		http.Error(w, "No active VA found", http.StatusInternalServerError)
		return
	}

	// Role check: only staff and admin can access logbook
	if activeVA.Role != "staff" && activeVA.Role != "admin" {
		http.Error(w, "You do not have permission to access the logbook. Staff or Admin privileges required.", http.StatusForbidden)
		return
	}

	// Prepare template data
	data := map[string]interface{}{
		"ActiveVA":        activeVA,
		"VirtualAirlines": sessionData.VirtualAirlines,
		"Username":        sessionData.Username,
		"UserID":          sessionData.UserID,
		"ActiveVAID":      sessionData.ActiveVAID,
		"PageTitle":       "Logbook",
	}

	// Render template
	if err := RenderTemplate(w, "pages/logbook.html", data); err != nil {
		http.Error(w, "Error rendering logbook", http.StatusInternalServerError)
		return
	}
}

// LogbookFlightsHandler returns a paginated list of flights for a given user (HTMX partial)
func LogbookFlightsHandler(
	w http.ResponseWriter,
	r *http.Request,
	flightSvc *services.FlightsService,
) {
	// Get user claims from context
	claims := auth.GetUserClaims(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session data from context
	sessionDataInterface := auth.GetSessionData(r.Context())
	if sessionDataInterface == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get query parameters
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "user_id parameter required", http.StatusBadRequest)
		return
	}

	pageStr := r.URL.Query().Get("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// For now, assume we can use a placeholder session ID
	// In production, this should come from the request context
	sessionID := ""

	// Fetch flights from service
	flightHistory, err := flightSvc.GetUserFlights(userID, page, sessionID)
	if err != nil {
		http.Error(w, "Failed to fetch flights: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if flightHistory == nil || flightHistory.Records == nil {
		flightHistory = &dtos.FlightHistoryDto{
			PageNo:  page,
			Records: []dtos.HistoryRecord{},
		}
	}

	// Prepare template data with pagination
	// The flightHistory.Records contains the flights for the current page
	// We need to determine if there are more pages based on the record count
	const recordsPerPage = 20
	hasMore := len(flightHistory.Records) == recordsPerPage // If we got a full page, there might be more

	data := map[string]interface{}{
		"Flights":     flightHistory.Records,
		"PageNo":      page,
		"HasNext":     hasMore,
		"HasPrevious": page > 1,
		"NextPage":    page + 1,
		"PrevPage":    page - 1,
		"UserID":      userID,
	}

	// Render partial (no base layout)
	if err := RenderPartial(w, "partials/flight-list.html", data); err != nil {
		http.Error(w, "Error rendering flight list", http.StatusInternalServerError)
		return
	}
}

// FlightMapHandler returns flight route map data with cached route (HTMX partial)
func FlightMapHandler(
	w http.ResponseWriter,
	r *http.Request,
	cache common.CacheInterface,
	liveAPI *common.LiveAPIService,
) {
	// Get user claims
	claims := auth.GetUserClaims(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session data
	sessionDataInterface := auth.GetSessionData(r.Context())
	if sessionDataInterface == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get flight_id and session_id from URL params
	flightID := chi.URLParam(r, "flight_id")
	sessionID := chi.URLParam(r, "session_id")
	if flightID == "" {
		http.Error(w, "flight_id parameter required", http.StatusBadRequest)
		return
	}
	if sessionID == "" {
		http.Error(w, "session_id parameter required", http.StatusBadRequest)
		return
	}

	// Try to get from cache first - use combo key with session_id and flight_id
	cacheKey := string(constants.CachePrefixFlightHistory) + sessionID + "_" + flightID
	cachedVal, found := cache.Get(cacheKey)
	var flightInfo *dtos.FlightInfo

	if found && cachedVal != nil {
		// Redis returns the value as a generic interface{} (map[string]interface{} from JSON unmarshal)
		// We need to marshal it back to JSON then unmarshal to FlightInfo struct
		jsonData, err := json.Marshal(cachedVal)
		if err == nil {
			var flight dtos.FlightInfo
			if err := json.Unmarshal(jsonData, &flight); err == nil {
				flightInfo = &flight
			}
		}
	}

	// If not in cache, render the empty state with debug panel
	// This allows the user to test the Live API endpoint directly
	if flightInfo == nil {
		apiKey := os.Getenv("IF_API_KEY")

		data := map[string]interface{}{
			"FlightID":  flightID,
			"APIKey":    apiKey,
			"SessionID": sessionID, // Now we have it from the URL param
		}
		if err := RenderPartial(w, "partials/flight-map-empty.html", data); err != nil {
			http.Error(w, "Error rendering map", http.StatusInternalServerError)
		}
		return
	}

	// Downsample route if it has too many waypoints
	route := flightInfo.Route
	if len(route) > 500 {
		route = downsampleRoute(route)
	}

	// Prepare template data
	data := map[string]interface{}{
		"Path":      route,
		"Origin":    flightInfo.Origin,
		"Dest":      flightInfo.Dest,
		"Meta":      flightInfo.Meta,
		"FlightID":  flightID,
		"SessionID": flightInfo.SessionID, // Include session ID for debug panel if needed
	}

	// Render partial
	if err := RenderPartial(w, "partials/flight-map.html", data); err != nil {
		http.Error(w, "Error rendering flight map", http.StatusInternalServerError)
		return
	}
}

// PilotSearchHandler returns pilot search results for autocomplete (HTMX partial)
func PilotSearchHandler(
	w http.ResponseWriter,
	r *http.Request,
	vaRoleRepo *repositories.VAUserRoleRepository,
) {
	// Get user claims
	claims := auth.GetUserClaims(r.Context())
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get session data
	sessionDataInterface := auth.GetSessionData(r.Context())
	if sessionDataInterface == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionData, ok := sessionDataInterface.(*common.SessionData)
	if !ok {
		http.Error(w, "Invalid session data", http.StatusInternalServerError)
		return
	}

	// Get active VA
	activeVA := sessionData.GetActiveVA()
	if activeVA == nil {
		http.Error(w, "No active VA found", http.StatusInternalServerError)
		return
	}

	// Role check
	if activeVA.Role != "staff" && activeVA.Role != "admin" {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}

	// Get search query
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		// Return empty results partial
		if err := RenderPartial(w, "partials/search-results.html", map[string]interface{}{
			"Results": []interface{}{},
		}); err != nil {
			http.Error(w, "Error rendering search results", http.StatusInternalServerError)
		}
		return
	}

	// Get all VA members
	vaRoles, err := vaRoleRepo.GetAllByVAID(r.Context(), activeVA.VAID)
	if err != nil {
		http.Error(w, "Failed to search pilots", http.StatusInternalServerError)
		return
	}

	// Filter by IFCommunityID match (case-insensitive)
	queryLower := strings.ToLower(query)
	var results []map[string]interface{}
	for _, vaRole := range vaRoles {
		ifcID := strings.ToLower(vaRole.User.IFCommunityID)
		if strings.Contains(ifcID, queryLower) {
			results = append(results, map[string]interface{}{
				"Username": vaRole.User.IFCommunityID,
				"Role":     vaRole.Role,
				"VAName":   activeVA.VAName,
			})
			if len(results) >= 10 {
				break
			}
		}
	}

	// Prepare template data
	data := map[string]interface{}{
		"Results": results,
	}

	// Render partial
	if err := RenderPartial(w, "partials/search-results.html", data); err != nil {
		http.Error(w, "Error rendering search results", http.StatusInternalServerError)
		return
	}
}

// downsampleRoute reduces the number of waypoints in a flight route
// For 500-1500 points: keeps every 3rd point
// For 1500+ points: keeps every 5th point
// Always keeps first and last points
func downsampleRoute(route []dtos.RouteWaypoint) []dtos.RouteWaypoint {
	if len(route) <= 500 {
		return route
	}

	var step int
	if len(route) <= 1500 {
		step = 3
	} else {
		step = 5
	}

	result := []dtos.RouteWaypoint{route[0]} // Always include first point

	for i := step; i < len(route)-1; i += step {
		result = append(result, route[i])
	}

	result = append(result, route[len(route)-1]) // Always include last point

	return result
}

// MapResetHandler returns empty map state (HTMX partial)
func MapResetHandler(w http.ResponseWriter, r *http.Request) {
	// Get session data
	sessionDataInterface := auth.GetSessionData(r.Context())
	if sessionDataInterface == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	data := map[string]interface{}{
		"Path": []interface{}{},
	}

	// Render empty map
	if err := RenderPartial(w, "partials/flight-map-empty.html", data); err != nil {
		http.Error(w, "Error rendering map", http.StatusInternalServerError)
		return
	}
}
