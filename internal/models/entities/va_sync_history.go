package entities

import "time"

type VASyncHistory struct {
	ID         string    `db:"id"`           // UUID
	VAID       string    `db:"va_id"`        // UUID
	Event      string    `db:"event"`        // varchar(10)
	CreatedAt  time.Time `db:"created_at"`   // timestamp
	LastSyncAt time.Time `db:"last_sync_at"` // nullable timestamp
}
