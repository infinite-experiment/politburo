package repositories

import (
	"context"

	"infinite-experiment/politburo/internal/models/gorm"

	gormlib "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PirepATSyncedRepo handles pirep_at_synced table operations
type PirepATSyncedRepo struct {
	db *gormlib.DB
}

// NewPirepATSyncedRepo creates a new PIREP at synced repository
func NewPirepATSyncedRepo(db *gormlib.DB) *PirepATSyncedRepo {
	return &PirepATSyncedRepo{db: db}
}

// Upsert inserts or updates a PIREP record from Airtable
// ON CONFLICT (server_id, at_id) DO UPDATE
func (r *PirepATSyncedRepo) Upsert(ctx context.Context, pirep *gorm.PirepATSynced) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "server_id"},
				{Name: "at_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"route", "flight_mode", "flight_time", "pilot_callsign",
				"aircraft", "livery", "route_at_id", "pilot_at_id", "at_created_time", "updated_at",
			}),
		}).
		Create(pirep).Error
}

// FindByATID finds a PIREP by VA ID and Airtable ID
func (r *PirepATSyncedRepo) FindByATID(ctx context.Context, vaID string, atID string) (*gorm.PirepATSynced, error) {
	var pirep gorm.PirepATSynced

	err := r.db.WithContext(ctx).
		Where("server_id = ? AND at_id = ?", vaID, atID).
		First(&pirep).Error

	if err != nil {
		if err == gormlib.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &pirep, nil
}

// GetByPilot returns all PIREPs for a specific pilot callsign, ordered by creation time descending
func (r *PirepATSyncedRepo) GetByPilot(ctx context.Context, vaID string, pilotCallsign string, limit int) ([]gorm.PirepATSynced, error) {
	var pireps []gorm.PirepATSynced

	query := r.db.WithContext(ctx).
		Where("server_id = ? AND pilot_callsign = ?", vaID, pilotCallsign).
		Order("at_created_time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&pireps).Error
	if err != nil {
		return nil, err
	}

	return pireps, nil
}

// GetByVA returns all PIREPs for a specific VA, ordered by creation time descending
func (r *PirepATSyncedRepo) GetByVA(ctx context.Context, vaID string, limit int) ([]gorm.PirepATSynced, error) {
	var pireps []gorm.PirepATSynced

	query := r.db.WithContext(ctx).
		Where("server_id = ?", vaID).
		Order("at_created_time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&pireps).Error
	if err != nil {
		return nil, err
	}

	return pireps, nil
}

// CountByVA returns the total number of PIREPs for a specific VA
func (r *PirepATSyncedRepo) CountByVA(ctx context.Context, vaID string) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&gorm.PirepATSynced{}).
		Where("server_id = ?", vaID).
		Count(&count).Error

	return count, err
}

// CountByPilot returns the total number of PIREPs for a specific pilot
func (r *PirepATSyncedRepo) CountByPilot(ctx context.Context, vaID string, pilotCallsign string) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&gorm.PirepATSynced{}).
		Where("server_id = ? AND pilot_callsign = ?", vaID, pilotCallsign).
		Count(&count).Error

	return count, err
}

// FindByATIDs finds PIREPs by VA ID and a list of Airtable IDs, ordered by creation time descending
func (r *PirepATSyncedRepo) FindByATIDs(ctx context.Context, vaID string, atIDs []string, limit int) ([]gorm.PirepATSynced, error) {
	var pireps []gorm.PirepATSynced

	if len(atIDs) == 0 {
		return pireps, nil
	}

	query := r.db.WithContext(ctx).
		Where("server_id = ? AND at_id IN ?", vaID, atIDs).
		Order("at_created_time DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&pireps).Error
	if err != nil {
		return nil, err
	}

	return pireps, nil
}
