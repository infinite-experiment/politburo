package entities

import "time"

// ValidationStatus represents the validation state of a config
type ValidationStatus string

const (
	ValidationStatusPending    ValidationStatus = "pending"
	ValidationStatusValidating ValidationStatus = "validating"
	ValidationStatusValid      ValidationStatus = "valid"
	ValidationStatusInvalid    ValidationStatus = "invalid"
)

// DataProviderConfig represents the database entity for va_data_provider_configs
type DataProviderConfig struct {
	ID               string           `db:"id"`
	VAID             string           `db:"va_id"`
	ProviderType     string           `db:"provider_type"`
	ConfigData       []byte           `db:"config_data"` // JSONB stored as bytes
	ConfigVersion    int              `db:"config_version"`
	IsActive         bool             `db:"is_active"`
	ValidationStatus ValidationStatus `db:"validation_status"`
	FeaturesEnabled  []string         `db:"features_enabled"` // PostgreSQL array
	LastValidatedAt  *time.Time       `db:"last_validated_at"`
	ValidationErrors []byte           `db:"validation_errors"` // JSONB stored as bytes
	CreatedAt        time.Time        `db:"created_at"`
	UpdatedAt        time.Time        `db:"updated_at"`
	CreatedBy        *string          `db:"created_by"`
	UpdatedBy        *string          `db:"updated_by"`
}

// ProviderValidationHistory represents validation history records
type ProviderValidationHistory struct {
	ID                string           `db:"id"`
	ConfigID          string           `db:"config_id"`
	ValidationStatus  ValidationStatus `db:"validation_status"`
	ValidationErrors  []byte           `db:"validation_errors"` // JSONB
	PhasesCompleted   []string         `db:"phases_completed"`  // PostgreSQL array
	PhasesFailed      []string         `db:"phases_failed"`     // PostgreSQL array
	DurationMs        *int             `db:"duration_ms"`
	ValidatedAt       time.Time        `db:"validated_at"`
	TriggeredBy       *string          `db:"triggered_by"`
}
