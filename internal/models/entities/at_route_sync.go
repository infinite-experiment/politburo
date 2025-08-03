package entities

import "time"

type RouteATSynced struct {
	ID          string    `db:"id"`
	ATID        string    `db:"at_id"`
	ServerID    string    `db:"server_id"`
	Origin      string    `db:"origin"`
	Destination string    `db:"destination"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Route       string    `db:"route"`
}
