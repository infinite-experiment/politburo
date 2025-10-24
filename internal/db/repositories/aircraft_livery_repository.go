package repositories

import (
	"context"
	"fmt"
	"time"

	gormModels "infinite-experiment/politburo/internal/models/gorm"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AircraftLiveryRepository struct {
	db *gorm.DB
}

// NewAircraftLiveryRepository creates a new GORM-based aircraft livery repository
func NewAircraftLiveryRepository(db *gorm.DB) *AircraftLiveryRepository {
	return &AircraftLiveryRepository{db: db}
}

// GetByLiveryID fetches a single active livery by ID
func (r *AircraftLiveryRepository) GetByLiveryID(ctx context.Context, liveryID string) (*gormModels.AircraftLivery, error) {
	var livery gormModels.AircraftLivery

	err := r.db.WithContext(ctx).
		Where("livery_id = ? AND is_active = ?", liveryID, true).
		First(&livery).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("livery not found: %s", liveryID)
		}
		return nil, fmt.Errorf("failed to fetch livery: %w", err)
	}

	return &livery, nil
}

// GetAllActive fetches all active liveries
func (r *AircraftLiveryRepository) GetAllActive(ctx context.Context) ([]gormModels.AircraftLivery, error) {
	var liveries []gormModels.AircraftLivery

	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Find(&liveries).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch active liveries: %w", err)
	}

	return liveries, nil
}

// UpsertBatch performs bulk upsert with conflict resolution on livery_id
func (r *AircraftLiveryRepository) UpsertBatch(ctx context.Context, liveries []gormModels.AircraftLivery) error {
	if len(liveries) == 0 {
		return nil
	}

	// Set sync timestamp for all records
	now := time.Now()
	for i := range liveries {
		liveries[i].LastSyncedAt = now
		liveries[i].UpdatedAt = now
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "livery_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"aircraft_id",
				"aircraft_name",
				"livery_name",
				"is_active",
				"updated_at",
				"last_synced_at",
			}),
		}).
		Create(&liveries).Error

	if err != nil {
		return fmt.Errorf("failed to upsert liveries batch: %w", err)
	}

	return nil
}

// GetLastSyncTime returns the most recent sync timestamp
func (r *AircraftLiveryRepository) GetLastSyncTime(ctx context.Context) (time.Time, error) {
	var result struct {
		LastSynced time.Time
	}

	err := r.db.WithContext(ctx).
		Model(&gormModels.AircraftLivery{}).
		Select("MAX(last_synced_at) as last_synced").
		Scan(&result).Error

	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get last sync time: %w", err)
	}

	return result.LastSynced, nil
}

// MarkInactive marks liveries as inactive by their IDs
func (r *AircraftLiveryRepository) MarkInactive(ctx context.Context, liveryIDs []string) error {
	if len(liveryIDs) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Model(&gormModels.AircraftLivery{}).
		Where("livery_id IN ?", liveryIDs).
		Updates(map[string]interface{}{
			"is_active":      false,
			"updated_at":     time.Now(),
			"last_synced_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("failed to mark liveries inactive: %w", err)
	}

	return nil
}

// GetLiveryMap returns a map of liveryID -> AircraftLivery for fast lookups
func (r *AircraftLiveryRepository) GetLiveryMap(ctx context.Context) (map[string]gormModels.AircraftLivery, error) {
	liveries, err := r.GetAllActive(ctx)
	if err != nil {
		return nil, err
	}

	liveryMap := make(map[string]gormModels.AircraftLivery, len(liveries))
	for _, livery := range liveries {
		liveryMap[livery.LiveryID] = livery
	}

	return liveryMap, nil
}
