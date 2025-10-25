package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/lib/pq"
)

// ValidationStatus represents the validation state of a config
type ValidationStatus string

const (
	ValidationStatusPending    ValidationStatus = "pending"
	ValidationStatusValidating ValidationStatus = "validating"
	ValidationStatusValid      ValidationStatus = "valid"
	ValidationStatusInvalid    ValidationStatus = "invalid"
)

// Scan implements the sql.Scanner interface for ValidationStatus
func (vs *ValidationStatus) Scan(value interface{}) error {
	if value == nil {
		*vs = ValidationStatusPending
		return nil
	}
	*vs = ValidationStatus(value.(string))
	return nil
}

// Value implements the driver.Valuer interface for ValidationStatus
func (vs ValidationStatus) Value() (driver.Value, error) {
	return string(vs), nil
}

// DataProviderConfig represents the GORM model for va_data_provider_configs
type DataProviderConfig struct {
	ID               string           `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	VAID             string           `gorm:"column:va_id;type:uuid;not null"`
	ProviderType     string           `gorm:"column:provider_type;type:varchar(50);not null"`
	ConfigData       JSONB            `gorm:"column:config_data;type:jsonb;not null"`
	ConfigVersion    int              `gorm:"column:config_version;default:1"`
	IsActive         bool             `gorm:"column:is_active;default:false"`
	ValidationStatus ValidationStatus `gorm:"column:validation_status;type:validation_status;default:'pending'"`
	FeaturesEnabled  pq.StringArray   `gorm:"column:features_enabled;type:text[]"`
	LastValidatedAt  *time.Time       `gorm:"column:last_validated_at"`
	ValidationErrors JSONB            `gorm:"column:validation_errors;type:jsonb"`
	CreatedAt        time.Time        `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt        time.Time        `gorm:"column:updated_at;autoUpdateTime"`
	CreatedBy        *string          `gorm:"column:created_by;type:uuid"`
	UpdatedBy        *string          `gorm:"column:updated_by;type:uuid"`

	// Relationships
	VA VA `gorm:"foreignKey:VAID"`
}

func (DataProviderConfig) TableName() string {
	return "va_data_provider_configs"
}

// ProviderValidationHistory represents the GORM model for va_provider_validation_history
type ProviderValidationHistory struct {
	ID               string           `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	ConfigID         string           `gorm:"column:config_id;type:uuid;not null"`
	ValidationStatus ValidationStatus `gorm:"column:validation_status;type:validation_status;not null"`
	ValidationErrors JSONB            `gorm:"column:validation_errors;type:jsonb"`
	PhasesCompleted  pq.StringArray   `gorm:"column:phases_completed;type:text[]"`
	PhasesFailed     pq.StringArray   `gorm:"column:phases_failed;type:text[]"`
	DurationMs       *int             `gorm:"column:duration_ms"`
	ValidatedAt      time.Time        `gorm:"column:validated_at;autoCreateTime"`
	TriggeredBy      *string          `gorm:"column:triggered_by;type:varchar(50)"`

	// Relationships
	Config DataProviderConfig `gorm:"foreignKey:ConfigID"`
}

func (ProviderValidationHistory) TableName() string {
	return "va_provider_validation_history"
}

// JSONB is a custom type for JSONB fields
type JSONB map[string]interface{}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result
	return nil
}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}
