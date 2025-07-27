package common

import (
	stdCtx "context"
	"errors"
	"fmt"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/context"
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
)

var AllowedVAConfigKeys = []string{
	ConfigKeyIFServerID,
	ConfigKeyTest,
	ConfigKeyCallsignPrefix,
	ConfigKeyCallsignSuffix,
}

// String slice ready for JSON response
func ListAllowedVAConfigKeys() []string { return AllowedVAConfigKeys }

// O(n) validator (fine for small list)
func IsValidVAConfigKey(k string) bool {
	for _, allowed := range AllowedVAConfigKeys {
		if allowed == k {
			return true
		}
	}
	return false
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
	key string,
	value string,
) (*map[string]string, error) {

	if !IsValidVAConfigKey(key) {
		return nil, fmt.Errorf("%q is not a valid key", key)
	}

	claims := context.GetUserClaims(ctx)
	va_id := claims.ServerID()

	// upsert
	if err := s.repo.UpsertVAConfig(ctx, va_id, key, value); err != nil {
		return nil, fmt.Errorf("failed to set config: %w", err)
	}
	cKey := configCacheKey(va_id)
	fmt.Printf("Evicting: %s", cKey)

	s.cache.Delete(cKey)

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
) (string, error) {

	if !IsValidVAConfigKey(key) {
		return "", fmt.Errorf("%q is not a valid key", key)
	}

	cfgs, err := s.GetAllConfigValues(ctx, vaID)
	if err != nil {
		return "", err
	}
	return cfgs[key], nil
}
