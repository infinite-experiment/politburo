package jobs

import (
	"context"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db/repositories"
	"time"

	"gorm.io/gorm"
)

// InitializeJobs initializes and starts all background jobs
func InitializeJobs(
	ctx context.Context,
	db *gorm.DB,
	cache *common.CacheService,
	configRepo *repositories.DataProviderConfigRepo,
	syncHistoryRepo *repositories.VASyncHistoryRepo,
	pilotATSyncedRepo *repositories.PilotATSyncedRepo,
	vaConfigService *common.VAConfigService,
) *PilotSyncJob {
	// Initialize pilot sync job (syncs pilots from Airtable every hour)
	pilotSyncJob := NewPilotSyncJob(
		db,
		cache,
		configRepo,
		syncHistoryRepo,
		pilotATSyncedRepo,
		vaConfigService,
	)

	// Start scheduled sync job in background
	go pilotSyncJob.RunScheduled(ctx, 1*time.Hour)

	return pilotSyncJob
}
