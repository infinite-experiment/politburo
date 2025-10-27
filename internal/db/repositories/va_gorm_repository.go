package repositories

import (
	"context"
	"fmt"

	gormModels "infinite-experiment/politburo/internal/models/gorm"

	"gorm.io/gorm"
)

// VAGormRepository handles VA table operations using GORM
type VAGormRepository struct {
	db *gorm.DB
}

// NewVAGormRepository creates a new GORM-based VA repository
func NewVAGormRepository(db *gorm.DB) *VAGormRepository {
	return &VAGormRepository{db: db}
}

// GetByID retrieves a VA by its ID
func (r *VAGormRepository) GetByID(ctx context.Context, vaID string) (*gormModels.VA, error) {
	var va gormModels.VA

	err := r.db.WithContext(ctx).
		Where("id = ?", vaID).
		First(&va).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch VA: %w", err)
	}

	return &va, nil
}

// GetByDiscordServerID retrieves a VA by Discord server ID
func (r *VAGormRepository) GetByDiscordServerID(ctx context.Context, discordServerID string) (*gormModels.VA, error) {
	var va gormModels.VA

	err := r.db.WithContext(ctx).
		Where("discord_server_id = ?", discordServerID).
		First(&va).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch VA: %w", err)
	}

	return &va, nil
}

// GetByCode retrieves a VA by its code
func (r *VAGormRepository) GetByCode(ctx context.Context, code string) (*gormModels.VA, error) {
	var va gormModels.VA

	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&va).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch VA: %w", err)
	}

	return &va, nil
}

// GetAll retrieves all active VAs
func (r *VAGormRepository) GetAll(ctx context.Context) ([]gormModels.VA, error) {
	var vas []gormModels.VA

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&vas).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch VAs: %w", err)
	}

	return vas, nil
}

// UpdateFlightModesConfig updates the flight modes configuration for a VA
func (r *VAGormRepository) UpdateFlightModesConfig(ctx context.Context, vaID string, config gormModels.JSONB) error {
	result := r.db.WithContext(ctx).
		Model(&gormModels.VA{}).
		Where("id = ?", vaID).
		Update("flight_modes_config", config)

	if result.Error != nil {
		return fmt.Errorf("failed to update flight modes config: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("VA not found with ID: %s", vaID)
	}

	return nil
}
