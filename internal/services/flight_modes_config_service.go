package services

import (
	"context"
	"fmt"

	"infinite-experiment/politburo/internal/db/repositories"
	gormModels "infinite-experiment/politburo/internal/models/gorm"
)

// FlightModesConfigService handles flight modes configuration management
type FlightModesConfigService struct {
	vaGormRepo *repositories.VAGormRepository
}

// NewFlightModesConfigService creates a new flight modes config service
func NewFlightModesConfigService(vaGormRepo *repositories.VAGormRepository) *FlightModesConfigService {
	return &FlightModesConfigService{
		vaGormRepo: vaGormRepo,
	}
}

// ValidateAndSaveConfig validates the flight modes configuration and saves it to the database
// Validates against the complete schema from PIREP_LOGGING_IMPLEMENTATION_PLAN.md
func (s *FlightModesConfigService) ValidateAndSaveConfig(ctx context.Context, vaID string, configPayload map[string]interface{}) error {
	// Validate basic structure - must have flight_modes key
	if _, ok := configPayload["flight_modes"]; !ok {
		return fmt.Errorf("configuration must contain 'flight_modes' key")
	}

	flightModes, ok := configPayload["flight_modes"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("'flight_modes' must be an object/map")
	}

	// Validate each mode has required keys and proper structure
	for modeID, modeData := range flightModes {
		modeConfig, ok := modeData.(map[string]interface{})
		if !ok {
			return fmt.Errorf("mode '%s' must be an object/map", modeID)
		}

		// Required keys for a mode
		if _, hasEnabled := modeConfig["enabled"]; !hasEnabled {
			return fmt.Errorf("mode '%s' must have 'enabled' field", modeID)
		}

		if _, hasDisplayName := modeConfig["display_name"]; !hasDisplayName {
			return fmt.Errorf("mode '%s' must have 'display_name' field", modeID)
		}

		// Validate requires_route_selection
		if _, hasRouteSelection := modeConfig["requires_route_selection"]; !hasRouteSelection {
			return fmt.Errorf("mode '%s' must have 'requires_route_selection' field", modeID)
		}

		// Validate fields array
		fields, ok := modeConfig["fields"].([]interface{})
		if !ok {
			return fmt.Errorf("mode '%s': 'fields' must be an array", modeID)
		}

		for idx, fieldData := range fields {
			field, ok := fieldData.(map[string]interface{})
			if !ok {
				return fmt.Errorf("mode '%s': field[%d] must be an object", modeID, idx)
			}

			// Validate field properties
			if _, hasName := field["name"]; !hasName {
				return fmt.Errorf("mode '%s': field[%d] must have 'name'", modeID, idx)
			}

			if _, hasType := field["type"]; !hasType {
				return fmt.Errorf("mode '%s': field[%d] must have 'type' (text, textarea, number)", modeID, idx)
			}

			if _, hasLabel := field["label"]; !hasLabel {
				return fmt.Errorf("mode '%s': field[%d] must have 'label'", modeID, idx)
			}

			if _, hasRequired := field["required"]; !hasRequired {
				return fmt.Errorf("mode '%s': field[%d] must have 'required'", modeID, idx)
			}
		}

		// Validate validations object
		validations, ok := modeConfig["validations"].(map[string]interface{})
		if !ok {
			return fmt.Errorf("mode '%s': 'validations' must be an object", modeID)
		}

		if _, hasAllowAny := validations["allow_any_current_route"]; !hasAllowAny {
			return fmt.Errorf("mode '%s': validations must have 'allow_any_current_route'", modeID)
		}

		if _, hasValidationMode := validations["validation_mode"]; !hasValidationMode {
			return fmt.Errorf("mode '%s': validations must have 'validation_mode' (any, exact_match)", modeID)
		}

		// Validate metadata exists (can be empty but must exist)
		if _, hasMetadata := modeConfig["metadata"]; !hasMetadata {
			return fmt.Errorf("mode '%s' must have 'metadata' object", modeID)
		}

		// auto_route is optional, but if present must be valid
		if autoRoute, hasAutoRoute := modeConfig["auto_route"]; hasAutoRoute && autoRoute != nil {
			if autoRouteObj, ok := autoRoute.(map[string]interface{}); ok {
				if _, hasRouteName := autoRouteObj["route_name"]; !hasRouteName {
					return fmt.Errorf("mode '%s': auto_route must have 'route_name'", modeID)
				}

				if _, hasMultiplier := autoRouteObj["multiplier"]; !hasMultiplier {
					return fmt.Errorf("mode '%s': auto_route must have 'multiplier'", modeID)
				}
			}
		}
	}

	// If validation passes, save to database
	// Convert map[string]interface{} to gormModels.JSONB
	jsonbConfig := gormModels.JSONB(configPayload)

	if err := s.vaGormRepo.UpdateFlightModesConfig(ctx, vaID, jsonbConfig); err != nil {
		return fmt.Errorf("failed to save flight modes configuration: %w", err)
	}

	return nil
}

// GetConfig retrieves the flight modes configuration for a VA
func (s *FlightModesConfigService) GetConfig(ctx context.Context, vaID string) (map[string]interface{}, error) {
	va, err := s.vaGormRepo.GetByID(ctx, vaID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch VA: %w", err)
	}

	if va == nil {
		return nil, fmt.Errorf("VA not found with ID: %s", vaID)
	}

	if va.FlightModesConfig == nil {
		return map[string]interface{}{}, nil
	}

	return va.FlightModesConfig, nil
}
