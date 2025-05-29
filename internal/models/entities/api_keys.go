package entities

type ApiKey struct {
	ApiKey string `db:"id"`
	Status bool   `db:"status"`
}
