package repositories

import (
	"context"

	"infinite-experiment/politburo/internal/models/gorm"

	gormlib "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PilotATSyncedRepo handles pilot_at_synced table operations
type PilotATSyncedRepo struct {
	db *gormlib.DB
}

// UnlinkedUser represents a user that needs to be linked to Airtable
type UnlinkedUser struct {
	ID       string `gorm:"column:id"`
	Callsign string `gorm:"column:callsign"`
}

// NewPilotATSyncedRepo creates a new pilot at synced repository
func NewPilotATSyncedRepo(db *gormlib.DB) *PilotATSyncedRepo {
	return &PilotATSyncedRepo{db: db}
}

// Upsert inserts or updates a pilot record from Airtable
// ON CONFLICT (server_id, at_id) DO UPDATE
func (r *PilotATSyncedRepo) Upsert(ctx context.Context, pilot *gorm.PilotATSynced) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "server_id"},
				{Name: "at_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"callsign", "registered"}),
		}).
		Create(pilot).Error
}

// FindByCallsign finds a pilot by VA ID and callsign (case-insensitive)
func (r *PilotATSyncedRepo) FindByCallsign(ctx context.Context, vaID string, callsign string) (*gorm.PilotATSynced, error) {
	var pilot gorm.PilotATSynced

	err := r.db.WithContext(ctx).
		Where("server_id = ? AND LOWER(callsign) = LOWER(?)", vaID, callsign).
		First(&pilot).Error

	if err != nil {
		if err == gormlib.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &pilot, nil
}

// GetUnlinkedUsers returns all users in va_user_roles who don't have airtable_pilot_id set
func (r *PilotATSyncedRepo) GetUnlinkedUsers(ctx context.Context, vaID string) ([]UnlinkedUser, error) {
	var users []UnlinkedUser

	err := r.db.WithContext(ctx).
		Table("va_user_roles").
		Where("va_id = ? AND callsign IS NOT NULL AND callsign != '' AND (airtable_pilot_id IS NULL OR airtable_pilot_id = '')", vaID).
		Select("id, callsign").
		Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}

// UpdateUserAirtableID updates the airtable_pilot_id for a user in va_user_roles
func (r *PilotATSyncedRepo) UpdateUserAirtableID(ctx context.Context, userRoleID string, airtableID string) error {
	return r.db.WithContext(ctx).
		Table("va_user_roles").
		Where("id = ?", userRoleID).
		Update("airtable_pilot_id", airtableID).Error
}
