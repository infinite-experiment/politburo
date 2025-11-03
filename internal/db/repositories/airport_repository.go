package repositories

import (
	"context"

	"infinite-experiment/politburo/internal/models/gorm"

	gormlib "gorm.io/gorm"
)

// AirportRepository handles airport table operations
type AirportRepository struct {
	db *gormlib.DB
}

// NewAirportRepository creates a new airport repository
func NewAirportRepository(db *gormlib.DB) *AirportRepository {
	return &AirportRepository{db: db}
}

// FindByICAO finds an airport by ICAO code (case-insensitive)
func (r *AirportRepository) FindByICAO(ctx context.Context, icao string) (*gorm.Airport, error) {
	var airport gorm.Airport

	err := r.db.WithContext(ctx).
		Where("UPPER(icao) = UPPER(?)", icao).
		First(&airport).Error

	if err != nil {
		if err == gormlib.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &airport, nil
}

// FindByIATA finds an airport by IATA code (case-insensitive)
func (r *AirportRepository) FindByIATA(ctx context.Context, iata string) (*gorm.Airport, error) {
	var airport gorm.Airport

	err := r.db.WithContext(ctx).
		Where("UPPER(iata) = UPPER(?)", iata).
		First(&airport).Error

	if err != nil {
		if err == gormlib.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &airport, nil
}

// BatchInsert inserts multiple airports
func (r *AirportRepository) BatchInsert(ctx context.Context, airports []gorm.Airport) error {
	return r.db.WithContext(ctx).
		CreateInBatches(airports, 100).Error
}

// DeleteAll deletes all airports (useful for re-importing)
func (r *AirportRepository) DeleteAll(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("1 = 1").
		Delete(&gorm.Airport{}).Error
}

// Count returns total number of airports
func (r *AirportRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&gorm.Airport{}).Count(&count).Error
	return count, err
}
