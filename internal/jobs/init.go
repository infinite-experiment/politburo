package jobs

import (
	"context"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/workers"
	"time"

	"gorm.io/gorm"
)

// JobsContainer holds all initialized jobs
type JobsContainer struct {
	PilotSync     *PilotSyncJob
	RouteSync     *RouteSyncJob
	PirepSync     *PirepSyncJob
	PIREPBackfill *workers.PIREPBackfill
}

// InitializeJobs initializes and starts all background jobs
func InitializeJobs(
	ctx context.Context,
	db *gorm.DB,
	cache common.CacheInterface,
	configRepo *repositories.DataProviderConfigRepo,
	syncHistoryRepo *repositories.VASyncHistoryRepo,
	pilotATSyncedRepo *repositories.PilotATSyncedRepo,
	routeATSyncedRepo *repositories.RouteATSyncedRepo,
	pirepATSyncedRepo *repositories.PirepATSyncedRepo,
	vaConfigService *common.VAConfigService,
	redisQueue *common.RedisQueueService,
) *JobsContainer {
	// Initialize pilot sync job (syncs pilots from Airtable every hour)
	pilotSyncJob := NewPilotSyncJob(
		db,
		cache,
		configRepo,
		syncHistoryRepo,
		pilotATSyncedRepo,
		vaConfigService,
	)

	// Initialize route sync job (syncs routes from Airtable every hour)
	routeSyncJob := NewRouteSyncJob(
		db,
		cache,
		configRepo,
		syncHistoryRepo,
		routeATSyncedRepo,
	)

	// Initialize PIREP sync job (syncs PIREPs from Airtable every hour)
	pirepSyncJob := NewPirepSyncJob(
		db,
		cache,
		configRepo,
		syncHistoryRepo,
		pirepATSyncedRepo,
		redisQueue,
	)

	// Initialize PIREP backfill job (backfills missing pilot/route data every 15 minutes)
	pirepBackfillJob := workers.NewPIREPBackfill(
		db,
		cache,
		*routeATSyncedRepo,
		*pilotATSyncedRepo,
	)

	// Start scheduled sync jobs in background (all run every hour)
	go pilotSyncJob.RunScheduled(ctx, 10*time.Minute)
	go routeSyncJob.RunScheduled(ctx, 10*time.Minute)
	go pirepSyncJob.RunScheduled(ctx, 10*time.Minute)
	go pirepBackfillJob.RunScheduled(ctx, 10*time.Minute)

	return &JobsContainer{
		PilotSync:     pilotSyncJob,
		RouteSync:     routeSyncJob,
		PirepSync:     pirepSyncJob,
		PIREPBackfill: pirepBackfillJob,
	}
}
