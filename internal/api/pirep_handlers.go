package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/models/dtos"
	gormModels "infinite-experiment/politburo/internal/models/gorm"
	"infinite-experiment/politburo/internal/services"
)

// GetPirepConfigHandler handles GET /api/v1/pireps/config
// Returns available flight modes and modal field configurations for the user's current flight
func (h *Handlers) GetPirepConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		// Get claims from context
		claims := auth.GetUserClaims(r.Context())
		if claims == nil {
			common.RespondError(w, initTime, nil, "Unauthorized: missing claims", http.StatusUnauthorized)
			return
		}

		vaDiscordServerID := claims.DiscordServerID()

		// Validate VA exists
		if vaDiscordServerID == "" {
			common.RespondError(w, initTime, fmt.Errorf("va not found"), "Virtual airline not found", http.StatusNotFound)
			return
		}

		// Get VA configuration with flight modes using Discord Server ID
		vaGorm, err := h.deps.Repo.VAGorm.GetByDiscordServerID(r.Context(), vaDiscordServerID)
		if err != nil {
			common.RespondError(w, initTime, err, "Failed to fetch VA configuration", http.StatusInternalServerError)
			return
		}

		if vaGorm == nil {
			common.RespondError(w, initTime, fmt.Errorf("va not found"), "Virtual airline not found", http.StatusNotFound)
			return
		}

		// Get the user by Discord ID and find their VA role/callsign
		discordID := claims.DiscordUserID()
		user, err := h.deps.Repo.UserGorm.GetUserWithVAAffiliations(r.Context(), discordID)
		if err != nil || user == nil {
			common.RespondError(w, initTime, fmt.Errorf("user not found"), "User not found", http.StatusNotFound)
			return
		}
		// Find user's callsign in their VA role
		userCallsign := ""
		for _, role := range user.UserVARoles {
			if role.VAID == vaGorm.ID {
				userCallsign = role.Callsign
				break
			}
		}

		if userCallsign == "" {
			common.RespondError(w, initTime, fmt.Errorf("user not member of va"), "User is not a member of this virtual airline", http.StatusForbidden)
			return
		}

		// Get VA config to retrieve prefix and suffix
		prefix, _ := h.deps.Services.Conf.GetConfigVal(r.Context(), vaGorm.ID, common.ConfigKeyCallsignPrefix)
		suffix, _ := h.deps.Services.Conf.GetConfigVal(r.Context(), vaGorm.ID, common.ConfigKeyCallsignSuffix)

		log.Printf("[GetPirepConfig] User callsign: %s, VA prefix: %s, VA suffix: %s", userCallsign, prefix, suffix)

		// Get VA live flights
		vaFlights, err := h.deps.Services.Flights.GetVALiveFlights(r.Context(), vaGorm.ID)
		if err != nil {
			log.Printf("[GetPirepConfig] Error fetching VA live flights: %v", err)
			// If we can't get live flights, continue with empty flight data
			// This allows the PIREP config to still be returned
			vaFlights = &[]dtos.LiveFlight{}
		}

		log.Printf("[GetPirepConfig] Fetched %d live flights for VA %s", len(*vaFlights), vaGorm.ID)

		// Find the user's current flight by matching the flight number
		// The user's callsign might be stored as just the number (e.g., "001")
		// and we need to construct what we're looking for: prefix + number + suffix
		expectedCallsignPattern := prefix + userCallsign + suffix
		log.Printf("[GetPirepConfig] Looking for flight matching pattern: %s (or just number: %s)", expectedCallsignPattern, userCallsign)

		flight := &common.FlightData{
			Callsign:    userCallsign,
			IFCUsername: user.IFCommunityID,
			Aircraft:    "",
			Livery:      "",
			LiveryID:    "",
			Route:       "",
		}

		// Find the user's current flight using unified method
		currentFlight, err := h.deps.Services.Flights.FindUserCurrentFlight(
			r.Context(),
			vaGorm.ID,
			userCallsign,
			prefix,
			suffix,
		)
		if err != nil {
			log.Printf("[GetPirepConfig] No matching flight found: %v", err)
			common.RespondError(w, initTime, fmt.Errorf("no live flight found"), "You are not currently flying. Please join a flight before filing a PIREP.", http.StatusNotFound)
			return
		}

		// Map the found flight to FlightData
		if currentFlight != nil {
			flight.Callsign = currentFlight.Callsign
			flight.Aircraft = currentFlight.Aircraft
			flight.Livery = currentFlight.Livery
			flight.LiveryID = currentFlight.LiveryId
			flight.Route = fmt.Sprintf("%s-%s", currentFlight.Origin, currentFlight.Destination)
			flight.Altitude = currentFlight.AltitudeFt
			flight.Speed = currentFlight.SpeedKts
		}

		// Build simplified response (without route details)
		response := h.buildSimplePirepConfigResponse(r.Context(), vaGorm, flight)
		common.RespondSuccess(w, initTime, "PIREP configuration fetched successfully", response)
	}
}

// SubmitPirepHandler handles POST /api/v1/pireps/submit
// Accepts PIREP submission data and processes it
func (h *Handlers) SubmitPirep() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		// Get claims from context
		claims := auth.GetUserClaims(r.Context())
		if claims == nil {
			common.RespondError(w, initTime, nil, "Unauthorized: missing claims", http.StatusUnauthorized)
			return
		}

		vaDiscordServerID := claims.DiscordServerID()

		// Validate VA exists
		if vaDiscordServerID == "" {
			common.RespondError(w, initTime, fmt.Errorf("va not found"), "Virtual airline not found", http.StatusNotFound)
			return
		}

		// Get VA configuration
		va, err := h.deps.Repo.VAGorm.GetByDiscordServerID(r.Context(), vaDiscordServerID)
		if err != nil {
			common.RespondError(w, initTime, err, "Failed to fetch VA configuration", http.StatusInternalServerError)
			return
		}

		if va == nil {
			common.RespondError(w, initTime, fmt.Errorf("va not found"), "Virtual airline not found", http.StatusNotFound)
			return
		}

		// Parse request body
		var submitRequest dtos.PirepSubmitRequest
		if err := json.NewDecoder(r.Body).Decode(&submitRequest); err != nil {
			common.RespondError(w, initTime, err, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get user and their current flight for livery mapping
		discordID := claims.DiscordUserID()
		user, err := h.deps.Repo.UserGorm.GetUserWithVAAffiliations(r.Context(), discordID)
		if err != nil || user == nil {
			common.RespondError(w, initTime, fmt.Errorf("user not found"), "User not found", http.StatusNotFound)
			return
		}

		// Get user's callsign for this VA
		userCallsign := ""
		for _, role := range user.UserVARoles {
			if role.VAID == va.ID {
				userCallsign = role.Callsign
				break
			}
		}

		if userCallsign == "" {
			common.RespondError(w, initTime, fmt.Errorf("user not member of va"), "User is not a member of this virtual airline", http.StatusForbidden)
			return
		}

		// Create submission service with all dependencies
		// Note: Service handles ALL flight data fetching internally (flight matching, livery resolution, aircraft/airline mapping)
		validator := services.NewFlightModeValidationService(&h.deps.Services.Live, h.deps.Services.Cache)
		submissionService := services.NewPirepSubmissionService(
			h.deps.Repo.UserGorm,
			h.deps.Repo.PilotATSynced,
			h.deps.Repo.RouteATSynced,
			h.deps.Repo.LiveryAirtableMapping,
			h.deps.Repo.DataProviderCfg,
			h.deps.Services.AirtableProvider,
			validator,
			h.deps.Services.Cache,
			&h.deps.Services.Flights,
			&h.deps.Services.Conf,
			h.deps.Services.DataProviderConfig,
		)

		// Submit PIREP (service handles all flight data fetching internally)
		response, err := submissionService.SubmitPirep(r.Context(), &submitRequest, va, claims)
		if err != nil {
			common.RespondError(w, initTime, err, "Failed to submit PIREP", http.StatusInternalServerError)
			return
		}

		// Return response (success or validation error)
		if response.Success {
			common.RespondSuccess(w, initTime, response.Message, response)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
		}
	}
}

// SetFlightModesConfig handles POST /api/v1/va/flight-modes/config
// Stores or updates flight mode configuration for a VA (admin-only)
func (h *Handlers) SetFlightModesConfig() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		initTime := time.Now()

		// Get claims from context
		claims := auth.GetUserClaims(r.Context())
		if claims == nil {
			common.RespondError(w, initTime, nil, "Unauthorized: missing claims", http.StatusUnauthorized)
			return
		}

		vaDiscordServerID := claims.DiscordServerID()

		// Validate VA exists
		if vaDiscordServerID == "" {
			common.RespondError(w, initTime, fmt.Errorf("va not found"), "Virtual airline not found", http.StatusNotFound)
			return
		}

		// Get VA by Discord Server ID using GORM repository
		vaGorm, err := h.deps.Repo.VAGorm.GetByDiscordServerID(r.Context(), vaDiscordServerID)
		if err != nil {
			common.RespondError(w, initTime, err, "Failed to fetch VA configuration", http.StatusInternalServerError)
			return
		}

		if vaGorm == nil {
			common.RespondError(w, initTime, fmt.Errorf("va not found"), "Virtual airline not found", http.StatusNotFound)
			return
		}

		// Parse request body
		var configPayload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&configPayload); err != nil {
			common.RespondError(w, initTime, err, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Use service to validate and save configuration
		configSvc := services.NewFlightModesConfigService(h.deps.Repo.VAGorm)
		if err := configSvc.ValidateAndSaveConfig(r.Context(), vaGorm.ID, configPayload); err != nil {
			common.RespondError(w, initTime, err, "Invalid configuration", http.StatusBadRequest)
			return
		}

		// Get the number of modes for response
		flightModes := configPayload["flight_modes"].(map[string]interface{})

		response := map[string]interface{}{
			"success": true,
			"message": "Flight modes configuration saved successfully",
			"va_id":   vaGorm.ID,
			"modes":   len(flightModes),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}

// buildPirepConfigResponse constructs the ConfigResponse with available modes and routes
func (h *Handlers) buildPirepConfigResponse(
	ctx context.Context,
	va *gormModels.VA,
	flight *common.FlightData,
	userDiscordID string,
) *dtos.ConfigResponse {
	response := &dtos.ConfigResponse{
		UserInfo: dtos.UserInfo{
			Callsign:            flight.Callsign,
			IFCUsername:         flight.IFCUsername,
			CurrentAircraft:     flight.Aircraft,
			CurrentLivery:       flight.Livery,
			CurrentRoute:        flight.Route,
			CurrentFlightStatus: "in_flight",
		},
		AvailableModes: []dtos.ModeResponse{},
	}

	// Extract flight modes from config
	if va.FlightModesConfig == nil || len(va.FlightModesConfig) == 0 {
		return response
	}

	flightModes, ok := va.FlightModesConfig["flight_modes"].(map[string]interface{})
	if !ok {
		return response
	}

	// Create validator
	validator := services.NewFlightModeValidationService(&h.deps.Services.Live, h.deps.Services.Cache)

	// Get all routes for this VA for route selection modes
	allRoutes, err := h.deps.Repo.RouteATSynced.GetAllByVA(ctx, va.ID)
	if err != nil {
		allRoutes = []gormModels.RouteATSynced{}
	}

	// Process each configured mode
	for modeID, modeData := range flightModes {
		modeConfig, ok := modeData.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if enabled
		enabled, _ := modeConfig["enabled"].(bool)
		if !enabled {
			continue
		}

		// Convert to FlightModeConfig struct
		modeConfigJSON, _ := json.Marshal(modeConfig)
		var flightModeConfig dtos.FlightModeConfig
		if err := json.Unmarshal(modeConfigJSON, &flightModeConfig); err != nil {
			continue
		}

		// Validate mode
		validationResult := validator.ValidateFlightForMode(ctx, flight.Route, &flightModeConfig.Validations)

		modeResponse := dtos.ModeResponse{
			ModeID:                 modeID,
			DisplayName:            flightModeConfig.DisplayName,
			RequiresRouteSelection: flightModeConfig.RequiresRouteSelection,
			Fields:                 flightModeConfig.Fields,
		}

		if validationResult.Valid {
			modeResponse.Status = "valid"

			// Add available routes if mode requires selection
			if flightModeConfig.RequiresRouteSelection {
				modeResponse.AvailableRoutes = h.buildAvailableRoutes(allRoutes, &flightModeConfig)
			} else if flightModeConfig.AutoRoute != nil {
				// Add auto-route information
				autoRoute := h.findAutoRoute(allRoutes, flightModeConfig.AutoRoute.RouteName)
				if autoRoute != nil {
					modeResponse.AutoRoute = autoRoute
				}
			}
		} else {
			modeResponse.Status = "invalid"
			modeResponse.ErrorReason = validationResult.ErrorMsg
		}

		response.AvailableModes = append(response.AvailableModes, modeResponse)
	}

	return response
}

// buildAvailableRoutes creates route options for a mode
func (h *Handlers) buildAvailableRoutes(allRoutes []gormModels.RouteATSynced, modeConfig *dtos.FlightModeConfig) []dtos.RouteOption {
	var routes []dtos.RouteOption

	// If mode has validation rules, filter routes
	if modeConfig.Validations.AllowedRoutes != nil && len(modeConfig.Validations.AllowedRoutes) > 0 {
		for _, route := range allRoutes {
			for _, allowed := range modeConfig.Validations.AllowedRoutes {
				if route.Route == allowed {
					routes = append(routes, dtos.RouteOption{
						RouteID:    route.ATID,
						Name:       route.Route,
						Multiplier: h.getRouteMultiplier(route, modeConfig),
					})
					break
				}
			}
		}
	} else {
		// Return all routes
		for _, route := range allRoutes {
			routes = append(routes, dtos.RouteOption{
				RouteID:    route.ATID,
				Name:       route.Route,
				Multiplier: h.getRouteMultiplier(route, modeConfig),
			})
		}
	}

	return routes
}

// findAutoRoute finds an auto-route by name
func (h *Handlers) findAutoRoute(allRoutes []gormModels.RouteATSynced, routeName string) *dtos.RouteOption {
	for _, route := range allRoutes {
		if route.Route == routeName {
			return &dtos.RouteOption{
				RouteID:    route.ATID,
				Name:       route.Route,
				Multiplier: 1.0, // Will be overridden by mode config
			}
		}
	}
	return nil
}

// getRouteMultiplier retrieves the multiplier for a route in a mode
// TODO: Implement route-specific multiplier lookup when route multiplier system is available
func (h *Handlers) getRouteMultiplier(route gormModels.RouteATSynced, modeConfig *dtos.FlightModeConfig) float64 {
	// For now, return mode default multiplier
	// In future, this can look up route-specific multipliers
	if modeConfig.Metadata != nil {
		if multiplier, ok := modeConfig.Metadata["multiplier"].(float64); ok {
			return multiplier
		}
	}
	if modeConfig.AutoRoute != nil {
		return modeConfig.AutoRoute.Multiplier
	}
	return 1.0
}

// buildSimplePirepConfigResponse constructs a minimal SimpleConfigResponse with just available modes and user info
// This is used for the GET /api/v1/pireps/config endpoint to provide a lightweight response
func (h *Handlers) buildSimplePirepConfigResponse(
	ctx context.Context,
	va *gormModels.VA,
	flight *common.FlightData,
) *dtos.SimpleConfigResponse {
	response := &dtos.SimpleConfigResponse{
		UserInfo: dtos.UserInfo{
			Callsign:            flight.Callsign,
			IFCUsername:         flight.IFCUsername,
			CurrentAircraft:     flight.Aircraft,
			CurrentLivery:       flight.Livery,
			CurrentRoute:        flight.Route,
			CurrentFlightStatus: "in_flight",
			CurrentAltitude:     flight.Altitude,
			CurrentSpeed:        flight.Speed,
		},
		AvailableModes: []dtos.SimpleModeResponse{},
	}

	// Extract flight modes from config
	if va.FlightModesConfig == nil || len(va.FlightModesConfig) == 0 {
		return response
	}

	flightModes, ok := va.FlightModesConfig["flight_modes"].(map[string]interface{})
	if !ok {
		return response
	}

	// Create validator
	validator := services.NewFlightModeValidationService(&h.deps.Services.Live, h.deps.Services.Cache)

	// Process each configured mode
	for modeID, modeData := range flightModes {
		modeConfig, ok := modeData.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if enabled
		enabled, _ := modeConfig["enabled"].(bool)
		if !enabled {
			continue
		}

		// Get display name
		displayName, _ := modeConfig["display_name"].(string)
		if displayName == "" {
			displayName = modeID
		}

		// Get requires_route_selection
		requiresRouteSelection, _ := modeConfig["requires_route_selection"].(bool)

		// Convert to FlightModeConfig struct for validation
		modeConfigJSON, _ := json.Marshal(modeConfig)
		var flightModeConfig dtos.FlightModeConfig
		if err := json.Unmarshal(modeConfigJSON, &flightModeConfig); err != nil {
			continue
		}

		// Validate mode
		validationResult := validator.ValidateFlightForMode(ctx, flight.Route, &flightModeConfig.Validations)

		modeResponse := dtos.SimpleModeResponse{
			ModeID:                 modeID,
			DisplayName:            displayName,
			RequiresRouteSelection: requiresRouteSelection,
			AutofillRoute:          flightModeConfig.AutofillRoute,
			Fields:                 flightModeConfig.Fields,
		}

		if validationResult.Valid {
			modeResponse.Status = "valid"
		} else {
			modeResponse.Status = "invalid"
			modeResponse.ErrorReason = validationResult.ErrorMsg
		}

		response.AvailableModes = append(response.AvailableModes, modeResponse)
	}

	return response
}
