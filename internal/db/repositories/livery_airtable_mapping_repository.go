package repositories

import (
	"context"
	"fmt"

	gormModels "infinite-experiment/politburo/internal/models/gorm"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LiveryAirtableMappingRepository struct {
	db *gorm.DB
}

// NewLiveryAirtableMappingRepository creates a new GORM-based livery airtable mapping repository
func NewLiveryAirtableMappingRepository(db *gorm.DB) *LiveryAirtableMappingRepository {
	return &LiveryAirtableMappingRepository{db: db}
}

// GetMapping retrieves a single mapping by VA ID, livery ID, and field type
func (r *LiveryAirtableMappingRepository) GetMapping(ctx context.Context, vaID, liveryID, fieldType string) (*gormModels.LiveryAirtableMapping, error) {
	var mapping gormModels.LiveryAirtableMapping

	err := r.db.WithContext(ctx).
		Where("va_id = ? AND livery_id = ? AND field_type = ? AND is_active = ?",
			vaID, liveryID, fieldType, true).
		First(&mapping).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("mapping not found: va_id=%s, livery_id=%s, field_type=%s", vaID, liveryID, fieldType)
		}
		return nil, fmt.Errorf("failed to fetch mapping: %w", err)
	}

	return &mapping, nil
}

// GetMappingsByLivery retrieves both aircraft and airline mappings for a livery
// Returns a map with "aircraft" and "airline" keys containing their target values
func (r *LiveryAirtableMappingRepository) GetMappingsByLivery(ctx context.Context, vaID, liveryID string) (map[string]string, error) {
	var mappings []gormModels.LiveryAirtableMapping

	err := r.db.WithContext(ctx).
		Where("va_id = ? AND livery_id = ? AND is_active = ?",
			vaID, liveryID, true).
		Find(&mappings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch livery mappings: %w", err)
	}

	result := make(map[string]string)
	for _, m := range mappings {
		result[m.FieldType] = m.TargetValue
	}

	return result, nil
}

// UpsertMappings performs an upsert operation for multiple mappings
// Uses the composite unique index (va_id, livery_id, field_type) for conflict resolution
func (r *LiveryAirtableMappingRepository) UpsertMappings(ctx context.Context, mappings []gormModels.LiveryAirtableMapping) error {
	if len(mappings) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "va_id"},
				{Name: "livery_id"},
				{Name: "field_type"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"source_value",
				"target_value",
				"is_active",
				"updated_at",
			}),
		}).
		Create(&mappings).Error

	if err != nil {
		return fmt.Errorf("failed to upsert livery mappings: %w", err)
	}

	return nil
}

// GetMappingsByVA retrieves all active mappings for a specific VA
func (r *LiveryAirtableMappingRepository) GetMappingsByVA(ctx context.Context, vaID string) ([]gormModels.LiveryAirtableMapping, error) {
	var mappings []gormModels.LiveryAirtableMapping

	err := r.db.WithContext(ctx).
		Where("va_id = ? AND is_active = ?", vaID, true).
		Find(&mappings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch VA mappings: %w", err)
	}

	return mappings, nil
}

// GetMappingsByLiveryIDs retrieves mappings for multiple livery IDs in a VA
func (r *LiveryAirtableMappingRepository) GetMappingsByLiveryIDs(ctx context.Context, vaID string, liveryIDs []string) (map[string]map[string]string, error) {
	var mappings []gormModels.LiveryAirtableMapping

	err := r.db.WithContext(ctx).
		Where("va_id = ? AND livery_id IN ? AND is_active = ?",
			vaID, liveryIDs, true).
		Find(&mappings).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch livery mappings: %w", err)
	}

	// Build nested map: livery_id -> field_type -> target_value
	result := make(map[string]map[string]string)
	for _, m := range mappings {
		if result[m.LiveryID] == nil {
			result[m.LiveryID] = make(map[string]string)
		}
		result[m.LiveryID][m.FieldType] = m.TargetValue
	}

	return result, nil
}

// DeleteByLiveryID deletes all mappings for a livery (soft delete via is_active flag)
func (r *LiveryAirtableMappingRepository) DeleteByLiveryID(ctx context.Context, vaID, liveryID string) error {
	err := r.db.WithContext(ctx).
		Model(&gormModels.LiveryAirtableMapping{}).
		Where("va_id = ? AND livery_id = ?", vaID, liveryID).
		Update("is_active", false).Error

	if err != nil {
		return fmt.Errorf("failed to delete livery mappings: %w", err)
	}

	return nil
}
