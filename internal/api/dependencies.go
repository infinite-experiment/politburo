package api

import (
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/services"
)

type Repositories struct {
	User       repositories.UserRepository
	UserGorm   *repositories.UserRepositoryGORM
	Keys       repositories.KeysRepo
	UserVASync repositories.SyncRepository
	Va         repositories.VARepository
}

type Services struct {
	Cache        common.CacheService
	Live         common.LiveAPIService
	User         *services.UserService
	Reg          services.RegistrationService
	Conf         common.VAConfigService
	VaMgmt       services.VAManagementService
	AirtableApi  common.AirtableApiService
	AirtableSync services.AtSyncService
	Flights      services.FlightsService
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
		UserVASync: *repositories.NewSyncRepository(db.DB),
	}

	cacheSvc := common.NewCacheService(60000, 600)
	liveSvc := common.NewLiveAPIService()
	confSvc := common.NewVAConfigService(&repositories.Va, cacheSvc)

	// Initialize user service with both sqlx and GORM repositories
	userSvc := services.NewUserService(&repositories.User, repositories.UserGorm)

	services := &Services{
		User:         userSvc,
		Reg:          *services.NewRegistrationService(liveSvc, *cacheSvc, repositories.User, repositories.Va),
		Conf:         *confSvc,
		VaMgmt:       *services.NewVAManagementService(repositories.Va, repositories.User),
		AirtableApi:  *common.NewAirtableApiService(confSvc),
		AirtableSync: *services.NewAtSyncService(cacheSvc, &repositories.UserVASync),
		Flights:      *services.NewFlightsService(cacheSvc, liveSvc, confSvc),
		Cache:        *cacheSvc,
		Live:         *liveSvc,
	}

	return &Dependencies{
		Repo:     repositories,
		Services: services,
	}, nil

}
