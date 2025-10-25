package api

import (
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/providers"
	"infinite-experiment/politburo/internal/services"
	"log"
	"os"
)

type Repositories struct {
	User            repositories.UserRepository
	UserGorm        *repositories.UserRepositoryGORM
	Keys            repositories.KeysRepo
	UserVASync      repositories.SyncRepository
	Va              repositories.VARepository
	VAGorm          *repositories.VAGormRepository
	DataProviderCfg *repositories.DataProviderConfigRepo
	VASyncHistory   *repositories.VASyncHistoryRepo
	PilotATSynced   *repositories.PilotATSyncedRepo
	RouteATSynced   *repositories.RouteATSyncedRepo
	PirepATSynced   *repositories.PirepATSyncedRepo
	AircraftLivery  *repositories.AircraftLiveryRepository
	LiveryAirtableMapping *repositories.LiveryAirtableMappingRepository
}

type Services struct {
	Cache              common.CacheInterface // Changed to interface to support Redis or in-memory
	LegacyCache        *common.CacheService  // For services that haven't been migrated to interface yet
	Live               common.LiveAPIService
	User               *services.UserService
	Reg                services.RegistrationService
	RegV2              *services.RegistrationServiceV2
	Conf               common.VAConfigService
	VaMgmt             services.VAManagementService
	AirtableApi        common.AirtableApiService
	AirtableSync       services.AtSyncService
	Flights            services.FlightsService
	PilotStats         *services.PilotStatsService
	DataProviderConfig *services.DataProviderConfigService
	AircraftLivery     *common.AircraftLiveryService
	RedisQueue         common.RedisQueueService
}
type Dependencies struct {
	Repo     *Repositories
	Services *Services
}

func InitDependencies() (*Dependencies, error) {

	repositories := &Repositories{
		User:       *repositories.NewUserRepository(db.DB),
		UserGorm:   repositories.NewUserRepositoryGORM(db.PgDB),
		Keys:       *repositories.NewApiKeysRepo(db.DB),
		Va:         *repositories.NewVARepository(db.DB),
		VAGorm:     repositories.NewVAGormRepository(db.PgDB),
		UserVASync: *repositories.NewSyncRepository(db.DB),
		DataProviderCfg: repositories.NewDataProviderConfigRepo(db.PgDB),
		VASyncHistory:   repositories.NewVASyncHistoryRepo(db.PgDB),
		PilotATSynced:   repositories.NewPilotATSyncedRepo(db.PgDB),
		RouteATSynced:   repositories.NewRouteATSyncedRepo(db.PgDB),
		PirepATSynced:   repositories.NewPirepATSyncedRepo(db.PgDB),
		AircraftLivery:  repositories.NewAircraftLiveryRepository(db.PgDB),
		LiveryAirtableMapping: repositories.NewLiveryAirtableMappingRepository(db.PgDB),
	}

	// Initialize cache service (Redis or in-memory based on USE_REDIS_CACHE env var)
	var cacheSvc common.CacheInterface
	useRedis := os.Getenv("USE_REDIS_CACHE") == "true"

	var redisQSvc common.RedisQueueService
	if useRedis {
		redisClient := common.NewRedisClient()
		redisCache, err := common.NewRedisCacheService(redisClient)
		if err != nil {
			log.Printf("Failed to initialize Redis cache, falling back to in-memory: %v", err)
			cacheSvc = common.NewCacheService(60000, 600)
		} else {
			log.Println("Using Redis cache")
			cacheSvc = redisCache
			redisQSvc = *common.NewRedisQueueService(redisClient)
		}
	} else {
		log.Println("Using in-memory cache")
		cacheSvc = common.NewCacheService(60000, 600)
	}

	// Create legacy cache wrapper for services that still need *CacheService
	var legacyCache *common.CacheService
	if cs, ok := cacheSvc.(*common.CacheService); ok {
		legacyCache = cs
	} else {
		// If using Redis, create a legacy in-memory cache for services that need it
		legacyCache = common.NewCacheService(60000, 600)
	}

	liveSvc := common.NewLiveAPIService()
	confSvc := common.NewVAConfigService(&repositories.Va, cacheSvc)

	// Initialize providers
	liveAPIProvider := providers.NewLiveAPIProvider()

	// Initialize pilot stats service first (needed by UserService)
	pilotStatsSvc := services.NewPilotStatsService(db.DB, db.PgDB, legacyCache, repositories.DataProviderCfg, &repositories.User, confSvc, repositories.PirepATSynced, repositories.RouteATSynced)

	// Initialize user service with both sqlx and GORM repositories and pilot stats service
	userSvc := services.NewUserService(&repositories.User, repositories.UserGorm, pilotStatsSvc)

	// Initialize data provider config service
	dataProviderConfigSvc := services.NewDataProviderConfigService(repositories.DataProviderCfg)

	// Initialize V2 registration service with GORM and LiveAPIProvider
	regServiceV2 := services.NewRegistrationServiceV2(db.PgDB, liveAPIProvider)

	// Initialize aircraft livery service
	aircraftLiverySvc := common.NewAircraftLiveryService(legacyCache, repositories.AircraftLivery)

	svc := &Services{
		User:               userSvc,
		Reg:                *services.NewRegistrationService(liveSvc, *legacyCache, repositories.User, repositories.Va),
		RegV2:              regServiceV2,
		Conf:               *confSvc,
		VaMgmt:             *services.NewVAManagementService(repositories.Va, repositories.User),
		AirtableApi:        *common.NewAirtableApiService(confSvc),
		AirtableSync:       *services.NewAtSyncService(legacyCache, &repositories.UserVASync),
		Flights:            *services.NewFlightsService(legacyCache, liveSvc, confSvc, aircraftLiverySvc),
		PilotStats:         pilotStatsSvc,
		DataProviderConfig: dataProviderConfigSvc,
		AircraftLivery:     aircraftLiverySvc,
		Cache:              cacheSvc,
		LegacyCache:        legacyCache,
		Live:               *liveSvc,
		RedisQueue:         redisQSvc,
	}

	return &Dependencies{
		Repo:     repositories,
		Services: svc,
	}, nil

}
