package gorm

// PilotATSynced represents a pilot record synced from Airtable
type PilotATSynced struct {
	ID         string `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	ATID       string `gorm:"column:at_id;type:varchar(20);not null"`
	Callsign   string `gorm:"column:callsign;type:varchar(20)"`
	Registered bool   `gorm:"column:registered;default:false"`
	ServerID   string `gorm:"column:server_id;type:uuid"`
}

// TableName specifies the table name for GORM
func (PilotATSynced) TableName() string {
	return "pilot_at_synced"
}
