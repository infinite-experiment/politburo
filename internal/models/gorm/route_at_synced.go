package gorm

import "time"

// RouteATSynced represents a route record synced from Airtable
type RouteATSynced struct {
	ID          string    `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	ATID        string    `gorm:"column:at_id;type:varchar(20);not null"`
	ServerID    string    `gorm:"column:server_id;type:uuid;not null"`
	Origin      string    `gorm:"column:origin;type:varchar(10)"`
	Destination string    `gorm:"column:destination;type:varchar(10)"`
	Route       string    `gorm:"column:route;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;default:now()"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:now()"`
}

// TableName specifies the table name for GORM
func (RouteATSynced) TableName() string {
	return "route_at_synced"
}
