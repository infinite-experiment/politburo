package services

import (
	"context"
	"encoding/json"
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/models"
	"infinite-experiment/politburo/internal/models/dtos"
	"infinite-experiment/politburo/internal/providers"
	"time"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

type PilotStatsService struct {
	db               *sqlx.DB
	gormDB           *gorm.DB
	cache            *common.CacheService
	configRepo       *repositories.DataProviderConfigRepo
	userRepo         *repositories.UserRepository
	vaConfigService  *common.VAConfigService
	airtableProvider *providers.AirtableProvider
}

func NewPilotStatsService(
	db *sqlx.DB,
	gormDB *gorm.DB,
	cache *common.CacheService,
	configRepo *repositories.DataProviderConfigRepo,
	userRepo *repositories.UserRepository,
	vaConfigService *common.VAConfigService,
) *PilotStatsService {
	return &PilotStatsService{
		db:               db,
		gormDB:           gormDB,
		cache:            cache,
		configRepo:       configRepo,
		userRepo:         userRepo,
		vaConfigService:  vaConfigService,
		airtableProvider: providers.NewAirtableProvider(cache),
	}
}

// GetPilotStats fetches pilot stats from Airtable for the authenticated user
func (s *PilotStatsService) GetPilotStats(ctx context.Context, userDiscordID, vaID string) (*PilotStatsResponse, error) {
	// Step 1: Get user's VA membership to find airtable_pilot_id
	membership, err := s.getUserMembership(ctx, userDiscordID, vaID)
	if err != nil {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodePilotNotSynced,
			Message: constants.GetErrorMessage(constants.ErrCodePilotNotSynced),
			Err:     err,
		}
	}

	// Step 2: Validate airtable_pilot_id exists
	if membership.AirtablePilotID == nil || *membership.AirtablePilotID == "" {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodePilotAirtableIDMissing,
			Message: constants.GetErrorMessage(constants.ErrCodePilotAirtableIDMissing),
		}
	}

	airtablePilotID := *membership.AirtablePilotID

	// Step 3: Get active Airtable config for VA
	config, configData, err := s.getActiveAirtableConfig(ctx, vaID)
	if err != nil {
		return nil, err
	}

	// Step 4: Check cache first
	cacheKey := fmt.Sprintf("pilot_stats:%s:%s", vaID, airtablePilotID)
	if cachedData, found := s.cache.Get(cacheKey); found {
		if stats, ok := cachedData.(*PilotStatsResponse); ok {
			stats.Metadata.Cached = true
			return stats, nil
		}
	}

	// Step 5: Get pilot schema from config
	pilotSchema := configData.GetSchemaByType("pilot")
	if pilotSchema == nil {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodeConfigMalformed,
			Message: "Pilot schema not found in configuration",
		}
	}

	// Step 6: Fetch from Airtable
	ctx = context.WithValue(ctx, "provider_config", configData)
	pilotRecord, err := s.airtableProvider.FetchPilotRecord(ctx, airtablePilotID, pilotSchema)
	if err != nil {
		// Check if it's a provider error
		if provErr, ok := err.(*providers.ProviderError); ok {
			// Map to pilot-specific error if record not found
			if provErr.Code == constants.ErrCodeInvalidBaseID {
				return nil, &PilotStatsError{
					Code:    constants.ErrCodePilotNotFoundInAirtable,
					Message: constants.GetErrorMessage(constants.ErrCodePilotNotFoundInAirtable),
					Err:     err,
				}
			}
			return nil, &PilotStatsError{
				Code:    provErr.Code,
				Message: provErr.Message,
				Err:     err,
			}
		}
		return nil, &PilotStatsError{
			Code:    constants.ErrCodeNetworkError,
			Message: constants.GetErrorMessage(constants.ErrCodeNetworkError),
			Err:     err,
		}
	}

	// Step 7: Build response
	response := &PilotStatsResponse{
		AirtablePilotID: airtablePilotID,
		Raw:             pilotRecord.RawFields,
		Normalized:      pilotRecord.Normalized,
		Metadata: PilotStatsMetadata{
			SchemaVersion: configData.Version,
			LastFetched:   time.Now().Format(time.RFC3339),
			Cached:        false,
			VAName:        membership.VAName,
			ConfigActive:  config.IsActive,
		},
	}

	// Step 8: Cache the result (15 minutes TTL)
	s.cache.Set(cacheKey, response, 15*time.Minute)

	return response, nil
}

// getUserMembership fetches the user's membership in the VA
func (s *PilotStatsService) getUserMembership(ctx context.Context, userDiscordID, vaID string) (*MembershipWithAirtable, error) {
	query := `
		SELECT
			u.id as user_id,
			u.discord_id,
			u.if_community_id,
			vur.airtable_pilot_id,
			vur.callsign,
			vur.role,
			va.name as va_name
		FROM users u
		JOIN va_user_roles vur ON u.id = vur.user_id
		JOIN virtual_airlines va ON vur.va_id = va.id
		WHERE u.discord_id = $1 AND vur.va_id = $2 AND vur.is_active = true
		LIMIT 1
	`

	var membership MembershipWithAirtable
	err := s.db.GetContext(ctx, &membership, query, userDiscordID, vaID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user membership: %w", err)
	}

	return &membership, nil
}

// getActiveAirtableConfig fetches and parses the active Airtable config for a VA
func (s *PilotStatsService) getActiveAirtableConfig(ctx context.Context, vaID string) (*models.DataProviderConfig, *dtos.ProviderConfigData, error) {
	// Get config entity from database
	config, err := s.configRepo.GetActiveConfig(ctx, vaID, "airtable")
	if err != nil {
		return nil, nil, &PilotStatsError{
			Code:    constants.ErrCodeConfigNotFound,
			Message: constants.GetErrorMessage(constants.ErrCodeConfigNotFound),
			Err:     err,
		}
	}

	if config == nil {
		return nil, nil, &PilotStatsError{
			Code:    constants.ErrCodeVAAirtableNotEnabled,
			Message: constants.GetErrorMessage(constants.ErrCodeVAAirtableNotEnabled),
		}
	}

	// Parse JSONB config_data
	configData, err := repositories.ParseConfigData(config.ConfigData)
	if err != nil {
		return nil, nil, &PilotStatsError{
			Code:    constants.ErrCodeConfigMalformed,
			Message: constants.GetErrorMessage(constants.ErrCodeConfigMalformed),
			Err:     err,
		}
	}

	return config, configData, nil
}

// PilotStatsResponse represents the API response for pilot stats
type PilotStatsResponse struct {
	AirtablePilotID string                 `json:"airtable_pilot_id"`
	Raw             map[string]interface{} `json:"raw"`
	Normalized      map[string]interface{} `json:"normalized"`
	Metadata        PilotStatsMetadata     `json:"metadata"`
}

// PilotStatsMetadata contains metadata about the response
type PilotStatsMetadata struct {
	SchemaVersion string `json:"schema_version"`
	LastFetched   string `json:"last_fetched"`
	Cached        bool   `json:"cached"`
	VAName        string `json:"va_name"`
	ConfigActive  bool   `json:"config_active"`
}

// MembershipWithAirtable extends membership data with Airtable ID
type MembershipWithAirtable struct {
	UserID          string  `db:"user_id"`
	DiscordID       string  `db:"discord_id"`
	IFCommunityID   string  `db:"if_community_id"`
	AirtablePilotID *string `db:"airtable_pilot_id"`
	Callsign        string  `db:"callsign"`
	Role            string  `db:"role"`
	VAName          string  `db:"va_name"`
}

// GetPilotStatusByCallsign fetches pilot data from Airtable by searching for the callsign
// This method constructs the full callsign using the configured prefix
func (s *PilotStatsService) GetPilotStatusByCallsign(ctx context.Context, userDiscordID, vaID string) (*PilotStatusResponse, error) {
	// Step 1: Get user's VA membership to check role and get callsign
	membership, err := s.getUserMembership(ctx, userDiscordID, vaID)
	if err != nil {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodePilotNotSynced,
			Message: constants.GetErrorMessage(constants.ErrCodePilotNotSynced),
			Err:     err,
		}
	}

	// Step 2: Check if user has a role (is a member)
	if membership.Role == "" {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodePilotNotSynced,
			Message: "User is not a member of this VA",
		}
	}

	// Step 3: Get callsign prefix from VA config
	callsignPrefix, ok := s.vaConfigService.GetConfigVal(ctx, vaID, common.ConfigKeyAirtableCallsignColumnPrefix)
	if !ok {
		callsignPrefix = "" // Default to no prefix if not configured
	}

	// Step 4: Construct full callsign
	fullCallsign := callsignPrefix + membership.Callsign
	fmt.Printf("Searching for pilot with callsign: %s (prefix: %s, base: %s)\n", fullCallsign, callsignPrefix, membership.Callsign)

	// Step 5: Get active Airtable config for VA
	config, configData, err := s.getActiveAirtableConfig(ctx, vaID)
	if err != nil {
		return nil, err
	}

	// Step 6: Get pilot schema from config
	pilotSchema := configData.GetSchemaByType("pilot")
	if pilotSchema == nil {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodeConfigMalformed,
			Message: "Pilot schema not found in configuration",
		}
	}

	// Step 7: Get the callsign field name from schema
	var callsignFieldName string
	for _, field := range pilotSchema.Fields {
		if field.InternalName == "callsign" {
			callsignFieldName = field.AirtableName
			break
		}
	}

	if callsignFieldName == "" {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodeConfigMalformed,
			Message: "Callsign field not found in pilot schema",
		}
	}

	// Step 8: Build Airtable filter formula
	// Airtable formula: {Callsign} = 'TEST012'
	filterFormula := fmt.Sprintf("{%s} = '%s'", callsignFieldName, fullCallsign)
	fmt.Printf("Airtable filter formula: %s\n", filterFormula)

	// Step 9: Fetch from Airtable using filter
	ctx = context.WithValue(ctx, "provider_config", configData)
	filters := &providers.SyncFilters{
		FilterFormula: filterFormula,
		Limit:         1, // We only expect one record
	}

	recordSet, err := s.airtableProvider.FetchRecords(ctx, pilotSchema, filters)
	if err != nil {
		if provErr, ok := err.(*providers.ProviderError); ok {
			return nil, &PilotStatsError{
				Code:    provErr.Code,
				Message: provErr.Message,
				Err:     err,
			}
		}
		return nil, &PilotStatsError{
			Code:    constants.ErrCodeNetworkError,
			Message: constants.GetErrorMessage(constants.ErrCodeNetworkError),
			Err:     err,
		}
	}

	// Step 10: Check if we found a record
	if len(recordSet.Records) == 0 {
		return nil, &PilotStatsError{
			Code:    constants.ErrCodePilotNotFoundInAirtable,
			Message: fmt.Sprintf("No pilot found with callsign: %s", fullCallsign),
		}
	}

	// Step 11: Log the raw response data
	record := recordSet.Records[0]
	fmt.Printf("\n=== PILOT STATUS RESPONSE FROM AIRTABLE ===\n")
	fmt.Printf("Record ID: %s\n", record.ID)

	// Pretty print the fields as JSON
	fieldsJSON, err := json.MarshalIndent(record.Fields, "", "  ")
	if err != nil {
		fmt.Printf("Raw Fields (error formatting): %+v\n", record.Fields)
	} else {
		fmt.Printf("Raw Fields (JSON):\n%s\n", string(fieldsJSON))
	}
	fmt.Printf("===========================================\n\n")

	// Step 12: Build response
	response := &PilotStatusResponse{
		AirtablePilotID: record.ID,
		Callsign:        membership.Callsign,
		FullCallsign:    fullCallsign,
		Role:            membership.Role,
		RawFields:       record.Fields,
		Metadata: PilotStatusMetadata{
			SchemaVersion: configData.Version,
			FetchedAt:     time.Now().Format(time.RFC3339),
			VAName:        membership.VAName,
			ConfigActive:  config.IsActive,
		},
	}

	return response, nil
}

// PilotStatusResponse represents the API response for pilot status search by callsign
type PilotStatusResponse struct {
	AirtablePilotID string                 `json:"airtable_pilot_id"`
	Callsign        string                 `json:"callsign"`
	FullCallsign    string                 `json:"full_callsign"`
	Role            string                 `json:"role"`
	RawFields       map[string]interface{} `json:"raw_fields"`
	Metadata        PilotStatusMetadata    `json:"metadata"`
}

// PilotStatusMetadata contains metadata about the response
type PilotStatusMetadata struct {
	SchemaVersion string `json:"schema_version"`
	FetchedAt     string `json:"fetched_at"`
	VAName        string `json:"va_name"`
	ConfigActive  bool   `json:"config_active"`
}

// PilotStatsError represents a pilot stats specific error
type PilotStatsError struct {
	Code    string
	Message string
	Err     error
}

func (e *PilotStatsError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *PilotStatsError) Unwrap() error {
	return e.Err
}
