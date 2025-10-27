package gorm

import "time"

// LiveryAirtableMapping represents the mapping between IF liveries and standardized Airtable field values
type LiveryAirtableMapping struct {
	ID          string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	VAID        string    `gorm:"column:va_id;type:uuid;not null;uniqueIndex:,composite:va_livery_field_unique"`
	LiveryID    string    `gorm:"column:livery_id;type:varchar(255);not null;uniqueIndex:,composite:va_livery_field_unique"`
	FieldType   string    `gorm:"column:field_type;type:varchar(50);not null;uniqueIndex:,composite:va_livery_field_unique"` // 'aircraft' or 'airline'
	SourceValue string    `gorm:"column:source_value;type:varchar(255);not null"`                                              // Raw from IF API
	TargetValue string    `gorm:"column:target_value;type:varchar(255);not null"`                                              // Standardized for Airtable
	IsActive    bool      `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	VA VA `gorm:"foreignKey:VAID"`
}

// TableName specifies the table name for GORM
func (LiveryAirtableMapping) TableName() string {
	return "livery_airtable_mappings"
}
