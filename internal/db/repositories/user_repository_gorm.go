package repositories

import (
	"context"
	"fmt"

	"infinite-experiment/politburo/internal/models/entities"
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

// FindUserMembership retrieves membership information for a user in a VA by Discord IDs
// Returns a Membership DTO with nullable fields (UserID, VAID, Role can all be nil)
// This replicates the original CTE-based SQL query behavior
func (r *UserRepositoryGORM) FindUserMembership(ctx context.Context, discordServerID string, userDiscordID string) (*entities.Membership, error) {
	// Step 1: Find user by discord_id
	var user gormModels.User
	userErr := r.db.WithContext(ctx).
		Where("discord_id = ?", userDiscordID).
		First(&user).Error

	// Step 2: Find VA by discord_server_id
	var va gormModels.VA
	vaErr := r.db.WithContext(ctx).
		Where("discord_server_id = ?", discordServerID).
		First(&va).Error

	// Initialize result with nullable fields
	result := &entities.Membership{
		UserID: nil,
		VAID:   nil,
		Role:   nil,
	}

	// If user found, set UserID
	if userErr == nil {
		result.UserID = &user.ID
	} else if userErr != gorm.ErrRecordNotFound {
		// Real error (not just "not found")
		return nil, fmt.Errorf("failed to fetch user: %w", userErr)
	}

	// If VA found, set VAID
	if vaErr == nil {
		result.VAID = &va.ID
	} else if vaErr != gorm.ErrRecordNotFound {
		// Real error (not just "not found")
		return nil, fmt.Errorf("failed to fetch VA: %w", vaErr)
	}

	// Step 3: If both user and VA exist, try to find the role relationship
	if result.UserID != nil && result.VAID != nil {
		var userVARole gormModels.UserVARole
		roleErr := r.db.WithContext(ctx).
			Where("user_id = ? AND va_id = ?", user.ID, va.ID).
			First(&userVARole).Error

		if roleErr == nil {
			result.Role = &userVARole.Role
		} else if roleErr != gorm.ErrRecordNotFound {
			// Real error (not just "not found")
			return nil, fmt.Errorf("failed to fetch role: %w", roleErr)
		}
		// If role not found, it stays nil (which is fine)
	}

	return result, nil
}
