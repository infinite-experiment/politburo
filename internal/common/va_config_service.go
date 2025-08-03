package common

import (
	stdCtx "context"
	"errors"
	"fmt"
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/db/repositories"
	"time"
)

///////////////////////////////////////////////////////////////////////////////
// Public “enum” — just string constants
///////////////////////////////////////////////////////////////////////////////

const (
	ConfigKeyIFServerID     = "if_server_id"
	ConfigKeyTest           = "test"
	ConfigKeyCallsignPrefix = "callsign_prefix"
	ConfigKeyCallsignSuffix = "callsign_suffix"
	ConfigKeyAirtableAPIKey = "airtable_api_key"
	ConfigKeyAirtableVABase = "airtable_va_base"

	// New table keys
	ConfigKeyATTablePilots = "at_table_pilots"
	ConfigKeyATTableRoutes = "at_table_routes"
	ConfigKeyATTablePIREPs = "at_table_pireps"

	// Field mapping keys
	ConfigKeyATFieldPilotsCallsign   = "at_field_pilots_callsign"
	ConfigKeyATFieldRoutesOrigin     = "at_field_routes_origin"
	ConfigKeyATFieldRoutesDest       = "at_field_routes_dest"
	ConfigKeyATFieldRoutesRoute      = "at_field_routes_route"
	ConfigKeyATFieldPIREPsCallsign   = "at_field_pireps_callsign"
	ConfigKeyATFieldPIREPsRoute      = "at_field_pireps_route"
	ConfigKeyATFieldPIREPsFlightTime = "at_field_pireps_ft"

	ConfigKeyATFieldLastModified = "at_field_last_modified"
)

var AllowedVAConfigKeys = map[string]struct{}{
	ConfigKeyIFServerID:              {},
	ConfigKeyTest:                    {},
	ConfigKeyCallsignPrefix:          {},
	ConfigKeyCallsignSuffix:          {},
	ConfigKeyAirtableAPIKey:          {},
	ConfigKeyAirtableVABase:          {},
	ConfigKeyATTablePilots:           {},
	ConfigKeyATTableRoutes:           {},
	ConfigKeyATTablePIREPs:           {},
	ConfigKeyATFieldPilotsCallsign:   {},
	ConfigKeyATFieldRoutesOrigin:     {},
	ConfigKeyATFieldRoutesDest:       {},
	ConfigKeyATFieldPIREPsCallsign:   {},
	ConfigKeyATFieldPIREPsRoute:      {},
	ConfigKeyATFieldPIREPsFlightTime: {},
	ConfigKeyATFieldLastModified:     {},
	ConfigKeyATFieldRoutesRoute:      {},
}

func ListAllowedVAConfigKeys() []string { return GetKeysStructMap(AllowedVAConfigKeys) }

func IsValidVAConfigKey(k string) bool {
	_, ok := AllowedVAConfigKeys[k]
	return ok
}

///////////////////////////////////////////////////////////////////////////////
// Service
///////////////////////////////////////////////////////////////////////////////

type VAConfigService struct {
	repo  *repositories.VARepository
	cache *CacheService
}

func NewVAConfigService(r *repositories.VARepository, c *CacheService) *VAConfigService {
	return &VAConfigService{repo: r, cache: c}
}

func configCacheKey(vaID string) string {
	return string(constants.CachePrefixVAConfig) + vaID
}

// Expose constants to API callers
func (s *VAConfigService) ListPossibleKeys() []string { return ListAllowedVAConfigKeys() }

// ---------------------------------------------------------------------------
// Set VA config and return updated map
// ---------------------------------------------------------------------------
func (s *VAConfigService) SetVaConfig(
	ctx stdCtx.Context,
	cfgs map[string]string,
) (*map[string]string, error) {

	claims := auth.GetUserClaims(ctx)
	fmt.Printf("Request Map: \n %v", cfgs)
	for key, value := range cfgs {

		if !IsValidVAConfigKey(key) {
			return nil, fmt.Errorf("%q is not a valid key", key)
		}

		va_id := claims.ServerID()

		// upsert
		if err := s.repo.UpsertVAConfig(ctx, va_id, key, value); err != nil {
			return nil, fmt.Errorf("failed to set config: %w", err)
		}
		cKey := configCacheKey(va_id)
		fmt.Printf("Evicting: %s", cKey)

		s.cache.Delete(cKey)
	}

	cfgs, err := s.GetAllConfigValues(ctx, claims.ServerID())
	if err != nil {
		return nil, err
	}
	return &cfgs, nil
}

// ---------------------------------------------------------------------------
// Get *all* values (cached)             map[string]string
// ---------------------------------------------------------------------------
func (s *VAConfigService) GetAllConfigValues(
	ctx stdCtx.Context,
	vaID string,
) (map[string]string, error) {

	ttl := 10 * time.Minute
	cacheKey := configCacheKey(vaID)

	val, err := s.cache.GetOrSet(cacheKey, ttl, func() (any, error) {
		rows, err := s.repo.GetVAConfigs(ctx, vaID)
		if err != nil {
			return nil, err
		}
		m := make(map[string]string, len(*rows))
		for _, r := range *rows {
			m[r.ConfigKey] = r.ConfigValue
		}

		return m, nil
	})
	if err != nil {
		return nil, err
	}

	cfgs, ok := val.(map[string]string)
	if !ok {
		return nil, errors.New("cache type assertion to map[string]string failed")
	}
	return cfgs, nil
}

// ---------------------------------------------------------------------------
// Get single value
// ---------------------------------------------------------------------------
func (s *VAConfigService) GetConfigVal(
	ctx stdCtx.Context,
	vaID string,
	key string, // callers import ConfigKeyIFServerID etc.
) (string, bool) {

	if !IsValidVAConfigKey(key) {
		return "", false
	}

	cfgs, err := s.GetAllConfigValues(ctx, vaID)
	if err != nil {
		return "", false
	}
	return cfgs[key], true
}

func (s *VAConfigService) GetConfigValues(
	ctx stdCtx.Context,
	vaID string,
	keys []string, // callers import ConfigKeyIFServerID etc.
) (map[string]string, bool) {

	conf := make(map[string]string, len(keys))
	cfgs, err := s.GetAllConfigValues(ctx, vaID)

	if err != nil {
		return conf, false
	}

	for _, key := range keys {
		if !IsValidVAConfigKey(key) {
			return conf, false
		}
		val, ok := cfgs[key]
		if ok {
			conf[key] = val
		} else {
			conf[key] = ""
		}
	}

	return conf, true
}
