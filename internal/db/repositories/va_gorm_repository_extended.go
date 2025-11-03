package repositories

import (
	"context"
	"fmt"

	models "infinite-experiment/politburo/internal/models/gorm"
	"gorm.io/gorm"
)

// VAGORMRepository provides GORM-based VA data access (extends va_gorm_repository)
type VAGORMRepository struct {
	db *gorm.DB
}

// NewVAGORMRepository creates a new VA GORM repository
func NewVAGORMRepository(db *gorm.DB) *VAGORMRepository {
	return &VAGORMRepository{db: db}
}

// GetByID retrieves a VA by UUID
func (r *VAGORMRepository) GetByID(ctx context.Context, vaID string) (*models.VA, error) {
	var va models.VA

	err := r.db.WithContext(ctx).
		Where("id = ?", vaID).
		First(&va).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("VA not found")
		}
		return nil, fmt.Errorf("failed to fetch VA: %w", err)
	}

	return &va, nil
}

// GetByDiscordServerID retrieves a VA by Discord server ID
func (r *VAGORMRepository) GetByDiscordServerID(ctx context.Context, discordServerID string) (*models.VA, error) {
	var va models.VA

	err := r.db.WithContext(ctx).
		Where("discord_server_id = ?", discordServerID).
		First(&va).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("VA not found with Discord server ID: %s", discordServerID)
		}
		return nil, fmt.Errorf("failed to fetch VA by Discord server ID: %w", err)
	}

	return &va, nil
}

// GetByCode retrieves a VA by its code (e.g., "SIA", "UAL")
func (r *VAGORMRepository) GetByCode(ctx context.Context, code string) (*models.VA, error) {
	var va models.VA

	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&va).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("VA not found with code: %s", code)
		}
		return nil, fmt.Errorf("failed to fetch VA by code: %w", err)
	}

	return &va, nil
}

// GetAllActive retrieves all active VAs
func (r *VAGORMRepository) GetAllActive(ctx context.Context) ([]models.VA, error) {
	var vas []models.VA

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&vas).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch active VAs: %w", err)
	}

	return vas, nil
}

// Create creates a new VA
func (r *VAGORMRepository) Create(ctx context.Context, va *models.VA) error {
	err := r.db.WithContext(ctx).Create(va).Error
	if err != nil {
		return fmt.Errorf("failed to create VA: %w", err)
	}
	return nil
}

// Update updates an existing VA
func (r *VAGORMRepository) Update(ctx context.Context, va *models.VA) error {
	err := r.db.WithContext(ctx).Save(va).Error
	if err != nil {
		return fmt.Errorf("failed to update VA: %w", err)
	}
	return nil
}
