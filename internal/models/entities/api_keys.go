package entities

import "time"

type ApiKey struct {
	ApiKey string `db:"id"`
	Status bool   `db:"status"`
}

type VAConfig struct {
	ID          string    `db:"id"`
	VAID        string    `db:"va_id"`
	ConfigKey   string    `db:"config_key"`
	ConfigValue string    `db:"config_value"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
