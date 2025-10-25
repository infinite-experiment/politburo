package gorm

import "time"

// VASyncHistory tracks sync operations for virtual airlines
type VASyncHistory struct {
	ID         string     `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	VAID       string     `gorm:"column:va_id;type:uuid;not null"`
	Event      string     `gorm:"column:event;type:varchar(50);not null"`
	CreatedAt  time.Time  `gorm:"column:created_at;autoCreateTime"`
	LastSyncAt *time.Time `gorm:"column:last_sync_at"`

	// Relationships
	VA VA `gorm:"foreignKey:VAID"`
}

// TableName specifies the table name for GORM
func (VASyncHistory) TableName() string {
	return "va_sync_history"
}
