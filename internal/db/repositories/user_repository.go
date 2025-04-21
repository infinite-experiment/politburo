package repositories

import (
	"context"

	"infinite-experiment/infinite-experiment-backend/internal/models/entities"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) InsertUser(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO users (
			discord_id,
			if_community_id,
			is_active
		)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`

	return r.db.QueryRowxContext(ctx, query,
		user.DiscordID,
		user.IFCommunityID,
		user.IsActive,
	).StructScan(user)
}
