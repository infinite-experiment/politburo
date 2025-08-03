package entities

type PilotATSynced struct {
	ID         string `db:"id"`         // UUID
	ATID       string `db:"at_id"`      // varchar(20)
	Callsign   string `db:"callsign"`   // varchar(20)
	Registered bool   `db:"registered"` // boolean
	ServerID   string `db:"server_id"`
}
