package repositories

import (
	"context"
	"fmt"

	gormModels "infinite-experiment/politburo/internal/models/gorm"

	"gorm.io/gorm"
)

type UserRepositoryGORM struct {
	db *gorm.DB
}

// NewUserRepositoryGORM creates a new GORM-based user repository
func NewUserRepositoryGORM(db *gorm.DB) *UserRepositoryGORM {
	return &UserRepositoryGORM{db: db}
}

// GetUserWithVAAffiliations retrieves a user by Discord ID with all VA affiliations preloaded
func (r *UserRepositoryGORM) GetUserWithVAAffiliations(ctx context.Context, userDiscordID string) (*gormModels.User, error) {
	var user gormModels.User

	err := r.db.WithContext(ctx).
		Preload("UserVARoles.VA").
		Where("discord_id = ?", userDiscordID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found with discord_id: %s", userDiscordID)
		}
		return nil, fmt.Errorf("failed to fetch user with affiliations: %w", err)
	}

	return &user, nil
}

// GetUserByDiscordID retrieves a user by Discord ID without relationships
func (r *UserRepositoryGORM) GetUserByDiscordID(ctx context.Context, discordID string) (*gormModels.User, error) {
	var user gormModels.User

	err := r.db.WithContext(ctx).
		Where("discord_id = ?", discordID).
		First(&user).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return &user, nil
}
