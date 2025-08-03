package repositories

import (
	"context"
	"infinite-experiment/politburo/internal/models/entities"

	"github.com/jmoiron/sqlx"
)

type SyncRepository struct {
	db *sqlx.DB
}

func NewSyncRepository(db *sqlx.DB) *SyncRepository {
	return &SyncRepository{
		db: db,
	}
}

func (svc SyncRepository) UpsertPilot(
	ctx context.Context,
	pilot *entities.PilotATSynced) error {
	const query = `
		INSERT INTO pilot_at_synced ( at_id, callsign, registered, server_id)
		VALUES (:at_id, :callsign, :registered, :server_id)
		ON CONFLICT (server_id, at_id) DO UPDATE
		SET callsign = EXCLUDED.callsign,
		    registered = EXCLUDED.registered
	`

	_, err := svc.db.NamedExecContext(ctx, query, pilot)
	return err

}

// func (svc SyncRepository) FindPilotByCallsign(
// 	ctx context.Context,
// 	callsign string,
// ) (*entities.PilotATSynced, error) {
// 	// Find pilot by callsign
// }

func (svc SyncRepository) UpsertRoute(
	ctx context.Context,
	route *entities.RouteATSynced) error {
	const query = `
		INSERT INTO route_at_synced (at_id, origin, destination, server_id, route)
		VALUES (:at_id, :origin, :destination, :server_id, :route)
		ON CONFLICT (server_id, at_id) DO UPDATE
		SET origin = EXCLUDED.origin,
			destination = EXCLUDED.destination,
			route = EXCLUDED.route;
	`

	_, err := svc.db.NamedExecContext(ctx, query, route)
	return err
}
