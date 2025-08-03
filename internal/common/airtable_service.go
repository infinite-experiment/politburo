package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/models/dtos"
	"net/http"
	"net/url"
	"time"
)

type AirtableConfig struct {
	ApiKey         string
	TableName      string
	BaseId         string
	LastModifiedAt string
	Keys           *[]string
}

type AirtableApiService struct {
	conf *VAConfigService
}

const (
	ATTypePIREP = "pirep"
	ATTypeRoute = "route"
	ATTypePilot = "pilot"
)

func NewAirtableApiService(c *VAConfigService) *AirtableApiService {
	return &AirtableApiService{
		conf: c,
	}
}

func (s *AirtableApiService) LoadAirtableCfgs(ctx context.Context, sId string, event string) (*AirtableConfig, bool) {
	var airtableConfigKeys = []string{
		ConfigKeyAirtableAPIKey,
		ConfigKeyAirtableVABase,
		ConfigKeyATTablePilots,
		ConfigKeyATTableRoutes,
		ConfigKeyATTablePIREPs,
		ConfigKeyATFieldPilotsCallsign,
		ConfigKeyATFieldRoutesOrigin,
		ConfigKeyATFieldRoutesDest,
		ConfigKeyATFieldPIREPsCallsign,
		ConfigKeyATFieldPIREPsRoute,
		ConfigKeyATFieldPIREPsFlightTime,
		ConfigKeyATFieldLastModified,
		ConfigKeyATFieldRoutesRoute,
	}

	cnf, ok := s.conf.GetConfigValues(ctx, sId, airtableConfigKeys)

	data := &AirtableConfig{
		ApiKey:         cnf[ConfigKeyAirtableAPIKey],
		BaseId:         cnf[ConfigKeyAirtableVABase],
		LastModifiedAt: cnf[ConfigKeyATFieldLastModified],
	}

	keys := []string{}
	switch event {
	case ATTypePIREP:
		data.TableName = cnf[ConfigKeyATTablePIREPs]
		keys = append(keys,
			cnf[ConfigKeyATFieldPIREPsCallsign],
			cnf[ConfigKeyATFieldPIREPsFlightTime],
			cnf[ConfigKeyATFieldPIREPsRoute])
	case ATTypeRoute:
		data.TableName = cnf[ConfigKeyATTableRoutes]
		keys = append(keys,
			cnf[ConfigKeyATFieldRoutesOrigin],
			cnf[ConfigKeyATFieldRoutesDest],
			cnf[ConfigKeyATFieldRoutesRoute])
	case ATTypePilot:
		data.TableName = cnf[ConfigKeyATTablePilots]
		keys = append(keys, cnf[ConfigKeyATFieldPilotsCallsign])
	}

	if len(keys) > 0 {
		data.Keys = &keys
	}

	return data, ok
}

func (s *AirtableApiService) FetchRecords(
	ctx context.Context,
	event string,
	filterByModified bool,
	since *time.Time,
	offset string,
) (*dtos.AirtableListRecords, error) {
	claims := auth.GetUserClaims(ctx)

	cfg, ok := s.LoadAirtableCfgs(ctx, claims.ServerID(), event)
	if !ok || cfg == nil || cfg.Keys == nil {
		return nil, fmt.Errorf("airtable config missing or incomplete")
	}

	fields := append(*cfg.Keys, cfg.LastModifiedAt)

	payload := map[string]interface{}{
		"fields": fields,
	}

	if filterByModified && since != nil {
		filterFormula := fmt.Sprintf(
			`IS_AFTER({%s}, "%s")`,
			cfg.LastModifiedAt,
			since.UTC().Format(time.RFC3339),
		)
		payload["filterByFormula"] = filterFormula
	}

	if offset != "" {
		payload["offset"] = offset
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	encodedTable := url.PathEscape(cfg.TableName)
	url := fmt.Sprintf("https://api.airtable.com/v0/%s/%s/listRecords", cfg.BaseId, encodedTable)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	// LogHTTPRequest(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var raw struct {
		Records []map[string]interface{} `json:"records"`
		Offset  string                   `json:"offset"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}

	return &dtos.AirtableListRecords{
		Results: raw.Records,
		Offset:  raw.Offset,
	}, nil
}

func (s *AirtableApiService) ExtractFieldsFromRecord(
	ctx context.Context,
	serverID string,
	event string,
	rec map[string]interface{},
) (map[string]string, string, error) {
	fields, ok := rec["fields"].(map[string]interface{})
	if !ok {
		return nil, "", fmt.Errorf("record missing fields block: %v", rec)
	}

	id, ok := rec["id"].(string)
	if !ok {
		return nil, "", fmt.Errorf("record missing Airtable ID")
	}

	cfg, ok := s.LoadAirtableCfgs(ctx, serverID, event)
	if !ok || cfg == nil || cfg.Keys == nil {
		return nil, "", fmt.Errorf("airtable config not available")
	}

	result := map[string]string{}
	for _, key := range *cfg.Keys {
		if val, exists := fields[key]; exists {
			result[key] = fmt.Sprintf("%v", val)
		}
	}

	return result, id, nil
}
