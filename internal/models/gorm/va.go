package gorm

import "time"

type VA struct {
	ID                  string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	DiscordID           string    `gorm:"column:discord_server_id;uniqueIndex"`
	Name                string    `gorm:"column:name"`
	Code                string    `gorm:"column:code;uniqueIndex"`
	IsActive            bool      `gorm:"column:is_active;default:true"`
	IsAirtableEnabled   bool      `gorm:"column:is_airtable_enabled;default:false"`
	FlightModesConfig   JSONB     `gorm:"column:flight_modes_config;type:jsonb;default:'{}'"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	UserVARoles []UserVARole `gorm:"foreignKey:VAID"`
	VAConfigs   []VAConfig   `gorm:"foreignKey:VAID"`
}

// TableName specifies the table name for GORM
func (VA) TableName() string {
	return "virtual_airlines"
}

type VAConfig struct {
	ID          string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	VAID        string    `gorm:"column:va_id;type:uuid"`
	ConfigKey   string    `gorm:"column:config_key"`
	ConfigValue string    `gorm:"column:config_value"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	VA VA `gorm:"foreignKey:VAID"`
}

// TableName specifies the table name for GORM
func (VAConfig) TableName() string {
	return "va_configs"
}
