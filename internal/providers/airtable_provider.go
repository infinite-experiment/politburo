package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/dtos"
	"io"
	"net/http"
	"time"
)

// AirtableProvider implements DataProvider for Airtable
type AirtableProvider struct {
	client *http.Client
	cache  *common.CacheService
}

// NewAirtableProvider creates a new Airtable provider
func NewAirtableProvider(cache *common.CacheService) *AirtableProvider {
	return &AirtableProvider{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: cache,
	}
}

// GetProviderType returns the provider type identifier
func (p *AirtableProvider) GetProviderType() string {
	return "airtable"
}

// FetchPilotRecord fetches a single pilot record by Airtable record ID
func (p *AirtableProvider) FetchPilotRecord(ctx context.Context, pilotID string, schema *dtos.EntitySchema) (*PilotRecord, error) {
	// Get config from context (should be set by service layer)
	config, ok := ctx.Value("provider_config").(*dtos.ProviderConfigData)
	if !ok {
		return nil, fmt.Errorf("provider config not found in context")
	}

	// Build Airtable API URL
	url := fmt.Sprintf("https://api.airtable.com/v0/%s/%s/%s",
		config.Credentials.BaseID,
		schema.TableName,
		pilotID,
	)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+config.Credentials.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Code:    constants.ErrCodeNetworkError,
			Message: constants.GetErrorMessage(constants.ErrCodeNetworkError),
			Err:     err,
		}
	}
	defer resp.Body.Close()

	// Handle error responses
	if err := p.handleHTTPError(resp); err != nil {
		return nil, err
	}

	// Parse response
	var airtableResp AirtableRecordResponse
	if err := json.NewDecoder(resp.Body).Decode(&airtableResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Transform to PilotRecord
	record := &PilotRecord{
		ProviderID: airtableResp.ID,
		RawFields:  airtableResp.Fields,
		Normalized: p.normalizeFields(airtableResp.Fields, schema),
	}

	return record, nil
}

// FetchRecords fetches multiple records with pagination
func (p *AirtableProvider) FetchRecords(ctx context.Context, schema *dtos.EntitySchema, filters *SyncFilters) (*RecordSet, error) {
	config, ok := ctx.Value("provider_config").(*dtos.ProviderConfigData)
	if !ok {
		return nil, fmt.Errorf("provider config not found in context")
	}

	// Build request payload
	payload := p.buildFetchPayload(schema, filters)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Build URL
	url := fmt.Sprintf("https://api.airtable.com/v0/%s/%s/listRecords",
		config.Credentials.BaseID,
		schema.TableName,
	)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+config.Credentials.APIKey)
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, &ProviderError{
			Code:    constants.ErrCodeNetworkError,
			Message: constants.GetErrorMessage(constants.ErrCodeNetworkError),
			Err:     err,
		}
	}
	defer resp.Body.Close()

	// Handle error responses
	if err := p.handleHTTPError(resp); err != nil {
		return nil, err
	}

	// Parse response
	var airtableResp AirtableListResponse
	if err := json.NewDecoder(resp.Body).Decode(&airtableResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Transform to RecordSet with IDs
	records := make([]RecordWithID, len(airtableResp.Records))
	for i, rec := range airtableResp.Records {
		records[i] = RecordWithID{
			ID:     rec.ID,
			Fields: rec.Fields,
		}
	}

	recordSet := &RecordSet{
		Records:      records,
		Offset:       airtableResp.Offset,
		HasMore:      airtableResp.Offset != "",
		TotalFetched: len(records),
	}

	return recordSet, nil
}

// ValidateConfig validates the Airtable configuration
func (p *AirtableProvider) ValidateConfig(ctx context.Context, config *dtos.ProviderConfigData) (*ValidationResult, error) {
	startTime := time.Now()
	result := &ValidationResult{
		IsValid:         true,
		PhasesCompleted: []string{},
		PhasesFailed:    []string{},
		Errors:          []dtos.ValidationError{},
		Warnings:        []dtos.ValidationError{},
	}

	// Phase 1: Credential Validation
	if err := p.validateCredentials(ctx, config); err != nil {
		result.IsValid = false
		result.PhasesFailed = append(result.PhasesFailed, "credential_validation")
		if provErr, ok := err.(*ProviderError); ok {
			result.Errors = append(result.Errors, dtos.ValidationError{
				Phase:     "credential_validation",
				Error:     provErr.Message,
				ErrorCode: provErr.Code,
				Timestamp: time.Now().Format(time.RFC3339),
			})
		}
		result.DurationMs = int(time.Since(startTime).Milliseconds())
		return result, nil
	}
	result.PhasesCompleted = append(result.PhasesCompleted, "credential_validation")

	// Phase 2: Table Validation
	// TODO: Implement in future
	result.PhasesCompleted = append(result.PhasesCompleted, "table_validation")

	// Phase 3: Field Validation
	// TODO: Implement in future
	result.PhasesCompleted = append(result.PhasesCompleted, "field_validation")

	result.DurationMs = int(time.Since(startTime).Milliseconds())
	return result, nil
}

// validateCredentials checks if the API key and base ID are valid
func (p *AirtableProvider) validateCredentials(ctx context.Context, config *dtos.ProviderConfigData) error {
	url := fmt.Sprintf("https://api.airtable.com/v0/meta/bases/%s/tables", config.Credentials.BaseID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+config.Credentials.APIKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return &ProviderError{
			Code:    constants.ErrCodeNetworkError,
			Message: constants.GetErrorMessage(constants.ErrCodeNetworkError),
			Err:     err,
		}
	}
	defer resp.Body.Close()

	return p.handleHTTPError(resp)
}

// handleHTTPError converts HTTP errors to ProviderError
func (p *AirtableProvider) handleHTTPError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)

	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return &ProviderError{
			Code:    constants.ErrCodeInvalidAPIKey,
			Message: constants.GetErrorMessage(constants.ErrCodeInvalidAPIKey),
			Details: string(body),
		}
	case http.StatusNotFound:
		return &ProviderError{
			Code:    constants.ErrCodeInvalidBaseID,
			Message: constants.GetErrorMessage(constants.ErrCodeInvalidBaseID),
			Details: string(body),
		}
	case http.StatusTooManyRequests:
		return &ProviderError{
			Code:    constants.ErrCodeRateLimited,
			Message: constants.GetErrorMessage(constants.ErrCodeRateLimited),
			Details: string(body),
		}
	default:
		return &ProviderError{
			Code:    constants.ErrCodeNetworkError,
			Message: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)),
			Details: string(body),
		}
	}
}

// normalizeFields maps raw Airtable fields to internal field names
func (p *AirtableProvider) normalizeFields(rawFields map[string]interface{}, schema *dtos.EntitySchema) map[string]interface{} {
	normalized := make(map[string]interface{})

	for _, fieldMapping := range schema.Fields {
		if value, exists := rawFields[fieldMapping.AirtableName]; exists {
			normalized[fieldMapping.InternalName] = value
		} else if fieldMapping.DefaultValue != nil {
			normalized[fieldMapping.InternalName] = *fieldMapping.DefaultValue
		}
	}

	return normalized
}

// buildFetchPayload builds the request payload for fetching records
func (p *AirtableProvider) buildFetchPayload(schema *dtos.EntitySchema, filters *SyncFilters) map[string]interface{} {
	payload := make(map[string]interface{})

	// Add fields to fetch
	fieldNames := schema.GetAirtableFieldNames()
	if len(fieldNames) > 0 {
		payload["fields"] = fieldNames
	}

	// Add filter - prioritize custom filter formula over modified since
	if filters != nil {
		if filters.FilterFormula != "" {
			// Use custom filter formula if provided
			payload["filterByFormula"] = filters.FilterFormula
		} else if filters.ModifiedSince != nil && schema.LastModifiedField != "" {
			// Fall back to modified since filter
			formula := fmt.Sprintf("IS_AFTER({%s}, '%s')", schema.LastModifiedField, *filters.ModifiedSince)
			payload["filterByFormula"] = formula
		}
	}

	// Add pagination
	if filters != nil && filters.Offset != "" {
		payload["offset"] = filters.Offset
	}

	// Add limit
	if filters != nil && filters.Limit > 0 {
		payload["pageSize"] = filters.Limit
	}

	return payload
}

// Airtable API response structures

type AirtableRecordResponse struct {
	ID     string                 `json:"id"`
	Fields map[string]interface{} `json:"fields"`
}

type AirtableListResponse struct {
	Records []AirtableRecordResponse `json:"records"`
	Offset  string                   `json:"offset,omitempty"`
}

// ProviderError represents a provider-specific error
type ProviderError struct {
	Code    string
	Message string
	Details string
	Err     error
}

func (e *ProviderError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ProviderError) Unwrap() error {
	return e.Err
}
