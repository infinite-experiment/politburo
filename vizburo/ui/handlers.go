package ui

import (
	"encoding/json"
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/services"
	"log"
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
	// Use pagination metadata from the Live API response
	data := map[string]interface{}{
		"Flights":     flightHistory.Records,
		"PageNo":      page,
		"HasNext":     flightHistory.HasNext,
		"HasPrevious": flightHistory.HasPrevious,
		"TotalPages":  flightHistory.TotalPages,
		"TotalCount":  flightHistory.TotalCount,
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
	flightSvc *services.FlightsService,
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

	// Get page from query params for metadata lookup (optional, defaults to 1)
	pageStr := r.URL.Query().Get("page")
	page := 1
	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Get IFC ID (username) from query params to look up actual IF user ID
	ifcID := r.URL.Query().Get("user_id")
	if ifcID == "" {
		log.Printf("[FlightMapHandler] Warning: user_id (IFC ID) not provided, cannot fetch metadata")
	}

	// Debug logging
	log.Printf("[FlightMapHandler] Processing flight: session_id=%s, flight_id=%s, ifc_id=%s, page=%d", sessionID, flightID, ifcID, page)

	// Look up actual IF user ID from IFC ID
	var actualUserID string
	if ifcID != "" {
		userStats, _, err := liveAPI.GetUserByIfcId(ifcID)
		if err != nil {
			log.Printf("[FlightMapHandler] Failed to lookup IF user ID for IFC ID %s: %v", ifcID, err)
		} else if len(userStats.Result) > 0 {
			actualUserID = userStats.Result[0].UserID
			log.Printf("[FlightMapHandler] Resolved IFC ID %s to IF user ID %s", ifcID, actualUserID)
		} else {
			log.Printf("[FlightMapHandler] No results when looking up IFC ID %s", ifcID)
		}
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
				log.Printf("[FlightMapHandler] Route data cache hit: %s_%s, route points=%d", sessionID, flightID, len(flight.Route))
			} else {
				log.Printf("[FlightMapHandler] Route data unmarshal failed: %v", err)
			}
		} else {
			log.Printf("[FlightMapHandler] Route data marshal failed: %v", err)
		}
	} else {
		log.Printf("[FlightMapHandler] Route data cache miss: %s_%s", sessionID, flightID)
	}

	// Fetch metadata using FlightsService (handles caching internally)
	// This ensures we use the same cache instance and get properly typed data
	var metaRecord *dtos.HistoryRecord
	if actualUserID != "" && page > 0 && ifcID != "" {
		historyDto, err := flightSvc.GetUserFlights(ifcID, page, sessionID)
		if err != nil {
			log.Printf("[FlightMapHandler] Failed to fetch user flights: %v", err)
		} else if historyDto != nil && len(historyDto.Records) > 0 {
			// Find the matching flight in the history
			for _, record := range historyDto.Records {
				if record.FlightID == flightID {
					metaRecord = &record
					log.Printf("[FlightMapHandler] Metadata found: %s, aircraft=%s, callsign=%s", flightID, record.Aircraft, record.Callsign)
					break
				}
			}
			if metaRecord == nil {
				log.Printf("[FlightMapHandler] Flight %s not found in history page %d (found %d records)", flightID, page, len(historyDto.Records))
			}
		} else {
			log.Printf("[FlightMapHandler] No flight history returned for user %s, page %d", actualUserID, page)
		}
	} else {
		log.Printf("[FlightMapHandler] Skipping metadata fetch: user_id=%s, page=%d, ifc_id=%s", actualUserID, page, ifcID)
	}

	apiKey := os.Getenv("IF_API_KEY")

	// Determine which partial to render based on data availability
	var partialPath string
	var data map[string]interface{}

	// State 1: Flight with route available
	if flightInfo != nil {
		log.Printf("[FlightMapHandler] Rendering with-route state (route available)")
		partialPath = "partials/flight-map/with-route.html"

		route := flightInfo.Route
		if len(route) > 500 {
			route = downsampleRoute(route)
		}

		data = map[string]interface{}{
			"FlightID":  flightID,
			"SessionID": sessionID,
			"APIKey":    apiKey,
			"Meta":      metaRecord, // Metadata from flight history cache
			"Path":      route,
			"Origin":    flightInfo.Origin,
			"Dest":      flightInfo.Dest,
			"RouteMeta": flightInfo.Meta, // Max speed/alt from route
		}
	} else if metaRecord != nil {
		// State 2: Metadata available but no route
		log.Printf("[FlightMapHandler] Rendering metadata-only state (route unavailable)")
		partialPath = "partials/flight-map/metadata-only.html"

		data = map[string]interface{}{
			"FlightID":  flightID,
			"SessionID": sessionID,
			"APIKey":    apiKey,
			"Meta":      metaRecord,
		}
	} else {
		// State 3: No flight data available
		log.Printf("[FlightMapHandler] Rendering empty state (no flight data)")
		partialPath = "partials/flight-map/empty.html"

		data = map[string]interface{}{
			"FlightID":  flightID,
			"SessionID": sessionID,
			"APIKey":    apiKey,
		}
	}

	// Render appropriate partial
	if err := RenderPartial(w, partialPath, data); err != nil {
		log.Printf("[FlightMapHandler] Error rendering partial %s: %v", partialPath, err)
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
