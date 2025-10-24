package gorm

import "time"

// AircraftLivery represents aircraft and livery metadata from Infinite Flight API
type AircraftLivery struct {
	ID           string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	LiveryID     string    `gorm:"column:livery_id;uniqueIndex;type:varchar(100)"`
	AircraftID   string    `gorm:"column:aircraft_id;index;type:varchar(100)"`
	AircraftName string    `gorm:"column:aircraft_name;type:text"`
	LiveryName   string    `gorm:"column:livery_name;type:text"`
	IsActive     bool      `gorm:"column:is_active;default:true;index"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
	LastSyncedAt time.Time `gorm:"column:last_synced_at;index"`
}

// TableName specifies the table name for GORM
func (AircraftLivery) TableName() string {
	return "aircraft_liveries"
}
