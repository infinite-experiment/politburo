package repositories

import (
	"context"
	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/entities"

	"github.com/jmoiron/sqlx"
)

type KeysRepo struct {
	db *sqlx.DB
}

func NewApiKeysRepo(db *sqlx.DB) *KeysRepo {
	return &KeysRepo{db}
}

func (r *KeysRepo) GetStatus(ctx context.Context, key string) (*entities.ApiKey, error) {
	var keyRes entities.ApiKey

	err := r.db.QueryRowxContext(ctx, constants.GetStatusByApiKey, key).StructScan(&keyRes)

	if err != nil {
		return nil, err
	}

	return &keyRes, nil
}
