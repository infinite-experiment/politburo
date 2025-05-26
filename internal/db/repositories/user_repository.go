package repositories

import (
	"context"
	"log"

	"infinite-experiment/politburo/internal/constants"
	"infinite-experiment/politburo/internal/models/entities"

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

func (r *UserRepository) FindUserByDiscordId(ctx context.Context, discordId string) (*entities.User, error) {

	var user entities.User

	err := r.db.QueryRowxContext(ctx, constants.GetUserByDiscordId, discordId).StructScan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) DeleteAllUsers(ctx context.Context) error {
	err := r.db.QueryRowxContext(ctx, constants.DeleteAllUsers)
	log.Printf("Query output: %v", err)
	return nil
}
