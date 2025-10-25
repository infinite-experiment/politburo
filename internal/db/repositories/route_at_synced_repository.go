package repositories

import (
	"context"

	"infinite-experiment/politburo/internal/models/gorm"

	gormlib "gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RouteATSyncedRepo handles route_at_synced table operations
type RouteATSyncedRepo struct {
	db *gormlib.DB
}

// NewRouteATSyncedRepo creates a new route at synced repository
func NewRouteATSyncedRepo(db *gormlib.DB) *RouteATSyncedRepo {
	return &RouteATSyncedRepo{db: db}
}

// Upsert inserts or updates a route record from Airtable
// ON CONFLICT (server_id, at_id) DO UPDATE
func (r *RouteATSyncedRepo) Upsert(ctx context.Context, route *gorm.RouteATSynced) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "server_id"},
				{Name: "at_id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"origin", "destination", "route", "updated_at"}),
		}).
		Create(route).Error
}

// FindByATID finds a route by VA ID and Airtable ID
func (r *RouteATSyncedRepo) FindByATID(ctx context.Context, vaID string, atID string) (*gorm.RouteATSynced, error) {
	var route gorm.RouteATSynced

	err := r.db.WithContext(ctx).
		Where("server_id = ? AND at_id = ?", vaID, atID).
		First(&route).Error

	if err != nil {
		if err == gormlib.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &route, nil
}

// GetAllByVA returns all routes for a specific VA
func (r *RouteATSyncedRepo) GetAllByVA(ctx context.Context, vaID string) ([]gorm.RouteATSynced, error) {
	var routes []gorm.RouteATSynced

	err := r.db.WithContext(ctx).
		Where("server_id = ?", vaID).
		Order("origin ASC, destination ASC").
		Find(&routes).Error

	if err != nil {
		return nil, err
	}

	return routes, nil
}

// CountByVA returns the total number of routes for a specific VA
func (r *RouteATSyncedRepo) CountByVA(ctx context.Context, vaID string) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&gorm.RouteATSynced{}).
		Where("server_id = ?", vaID).
		Count(&count).Error

	return count, err
}
