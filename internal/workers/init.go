package workers

import (
	"context"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db/repositories"
	"time"

	"gorm.io/gorm"
)

type WorkersContainer struct {
	CacheFiller MetaCacheWorker
}

func InitWorkers(
	db *gorm.DB,
	c *common.CacheInterface,
	api *common.LiveAPIService,
	liverySvc *common.AircraftLiveryService,
	redQ *common.RedisQueueService,
	liveryRepo *repositories.AircraftLiveryRepository,
	dataProvCfg *repositories.DataProviderConfigRepo,
	pirepSyncedRepo *repositories.PirepATSyncedRepo,
	vaSyncHRepo *repositories.VASyncHistoryRepo,
) *WorkersContainer {
	mcf := NewMetaCacheFiller(c, api, liveryRepo, liverySvc)

	// Start the logbook worker to cache flight routes on-demand
	go LogbookWorker(*c, api, liverySvc)

	qWorker := NewPirepQueueWorker("pirep_queue", db, redQ, dataProvCfg, pirepSyncedRepo, vaSyncHRepo)
	monitor := NewPirepQueueMonitor(db, redQ)

	go qWorker.Start(context.Background(), 5)
	go monitor.Start(context.Background(), 30*time.Second)

	// Start workers
	go mcf.Start()

	return &WorkersContainer{
		CacheFiller: *mcf,
	}
}
