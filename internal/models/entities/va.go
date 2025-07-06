package entities

import "time"

type VA struct {
	ID             string    `db:"id"`
	DiscordID      string    `db:"discord_server_id"`
	Name           string    `db:"name"`
	Code           string    `db:"code"`
	IsActive       bool      `db:"is_active"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	CallsignPrefix string    `db:"callsign_prefix"`
	CallsignSuffix string    `db:"callsign_suffix"`
}
