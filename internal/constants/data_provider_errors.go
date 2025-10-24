package constants

// Data Provider Error Codes
// These constants define specific error scenarios for external data providers

// Credential-related errors
const (
	ErrCodeInvalidAPIKey        = "INVALID_API_KEY"
	ErrCodeInvalidBaseID        = "INVALID_BASE_ID"
	ErrCodeRateLimited          = "RATE_LIMITED"
	ErrCodeNetworkError         = "NETWORK_ERROR"
	ErrCodeAuthenticationFailed = "AUTHENTICATION_FAILED"
)

// Table-related errors
const (
	ErrCodeTableNotFound       = "TABLE_NOT_FOUND"
	ErrCodeTableAccessDenied   = "TABLE_ACCESS_DENIED"
	ErrCodeTableEmpty          = "TABLE_EMPTY"
)

// Field-related errors
const (
	ErrCodeFieldNotFound         = "FIELD_NOT_FOUND"
	ErrCodeFieldRenamed          = "FIELD_RENAMED"
	ErrCodeFieldTypeMismatch     = "FIELD_TYPE_MISMATCH"
)

// Data validation errors
const (
	ErrCodeTypeConversionError      = "TYPE_CONVERSION_ERROR"
	ErrCodeRequiredFieldEmpty       = "REQUIRED_FIELD_EMPTY"
	ErrCodeRequiredFieldMostlyEmpty = "REQUIRED_FIELD_MOSTLY_EMPTY"
	ErrCodeInvalidDataFormat        = "INVALID_DATA_FORMAT"
	ErrCodeDataOutOfRange           = "DATA_OUT_OF_RANGE"
)

// Configuration errors
const (
	ErrCodeConfigMalformed           = "CONFIG_MALFORMED"
	ErrCodeSchemaVersionUnsupported  = "SCHEMA_VERSION_UNSUPPORTED"
	ErrCodeValidationTimeout         = "VALIDATION_TIMEOUT"
	ErrCodeConfigNotFound            = "CONFIG_NOT_FOUND"
	ErrCodeConfigNotActive           = "CONFIG_NOT_ACTIVE"
	ErrCodeConfigNotValidated        = "CONFIG_NOT_VALIDATED"
)

// User/Pilot-specific errors
const (
	ErrCodePilotNotSynced           = "PILOT_NOT_SYNCED"
	ErrCodePilotAirtableIDMissing   = "PILOT_AIRTABLE_ID_MISSING"
	ErrCodePilotNotFoundInAirtable  = "PILOT_NOT_FOUND_IN_AIRTABLE"
	ErrCodeVAAirtableNotEnabled     = "VA_AIRTABLE_NOT_ENABLED"
)

// Error Messages
// Human-readable messages corresponding to error codes

var DataProviderErrorMessages = map[string]string{
	// Credentials
	ErrCodeInvalidAPIKey:        "The Airtable API key is invalid or has been revoked",
	ErrCodeInvalidBaseID:        "The Airtable Base ID is invalid or you don't have access to it",
	ErrCodeRateLimited:          "Rate limit exceeded. Please try again later",
	ErrCodeNetworkError:         "Unable to connect to Airtable. Please check your internet connection",
	ErrCodeAuthenticationFailed: "Authentication with Airtable failed",

	// Tables
	ErrCodeTableNotFound:     "The specified table was not found in the Airtable base",
	ErrCodeTableAccessDenied: "You don't have permission to access this table",
	ErrCodeTableEmpty:        "The table exists but contains no records",

	// Fields
	ErrCodeFieldNotFound:     "The specified field was not found in the table",
	ErrCodeFieldRenamed:      "The field appears to have been renamed in Airtable",
	ErrCodeFieldTypeMismatch: "The field type in Airtable doesn't match the expected type",

	// Data validation
	ErrCodeTypeConversionError:      "Unable to convert the field value to the expected type",
	ErrCodeRequiredFieldEmpty:       "A required field is empty in the Airtable record",
	ErrCodeRequiredFieldMostlyEmpty: "A required field is empty in more than 10% of records",
	ErrCodeInvalidDataFormat:        "The data format is invalid",
	ErrCodeDataOutOfRange:           "The data value is outside the acceptable range",

	// Configuration
	ErrCodeConfigMalformed:          "The configuration structure is invalid",
	ErrCodeSchemaVersionUnsupported: "The configuration uses an unsupported schema version",
	ErrCodeValidationTimeout:        "Configuration validation timed out",
	ErrCodeConfigNotFound:           "No data provider configuration found for this VA",
	ErrCodeConfigNotActive:          "The data provider configuration is not active",
	ErrCodeConfigNotValidated:       "The configuration has not been validated yet",

	// User/Pilot
	ErrCodePilotNotSynced:          "This pilot has not been synced with Airtable yet",
	ErrCodePilotAirtableIDMissing:  "The pilot's Airtable ID is not set in the system",
	ErrCodePilotNotFoundInAirtable: "The pilot's record was not found in Airtable. They may have been removed",
	ErrCodeVAAirtableNotEnabled:    "Airtable integration is not enabled for this Virtual Airline",
}

// GetErrorMessage returns the human-readable message for an error code
func GetErrorMessage(code string) string {
	if msg, exists := DataProviderErrorMessages[code]; exists {
		return msg
	}
	return "An unknown error occurred"
}

// Feature Requirements
// Defines which fields are required for specific features

type FeatureRequirement struct {
	FeatureName    string
	RequiredFields map[string][]string // entity_type -> required internal_names
}

var FeatureRequirements = map[string]FeatureRequirement{
	"sync_pilots": {
		FeatureName: "sync_pilots",
		RequiredFields: map[string][]string{
			"pilot": {"callsign"},
		},
	},
	"sync_routes": {
		FeatureName: "sync_routes",
		RequiredFields: map[string][]string{
			"route": {"origin", "destination"},
		},
	},
	"pilot_stats": {
		FeatureName: "pilot_stats",
		RequiredFields: map[string][]string{
			"pilot": {"callsign"}, // At minimum need callsign
		},
	},
	"flight_tracking": {
		FeatureName: "flight_tracking",
		RequiredFields: map[string][]string{
			"pilot": {"callsign"},
			"route": {"origin", "destination"},
		},
	},
}
