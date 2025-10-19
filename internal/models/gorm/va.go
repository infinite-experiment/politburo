package gorm

import "time"

type VA struct {
	ID                string    `gorm:"column:id;primaryKey;type:uuid"`
	DiscordID         string    `gorm:"column:discord_server_id;uniqueIndex"`
	Name              string    `gorm:"column:name"`
	Code              string    `gorm:"column:code;uniqueIndex"`
	IsActive          bool      `gorm:"column:is_active;default:true"`
	IsAirtableEnabled bool      `gorm:"column:is_airtable_enabled;default:false"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoUpdateTime"`

	// Relationships
	UserVARoles []UserVARole `gorm:"foreignKey:VAID"`
}

// TableName specifies the table name for GORM
func (VA) TableName() string {
	return "virtual_airlines"
}
