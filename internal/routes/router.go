package routes

import (
	"context"
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/api"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/jobs"
	"infinite-experiment/politburo/internal/workers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func RegisterRoutes(upSince time.Time) http.Handler {

	// initialize Chi router
	r := chi.NewRouter()

	// global middleware
	//r.Use(middleware.Logging)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:8081"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-API-Key", "X-Server-Id", "X-Discord-Id"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// health check
	r.Get("/healthCheck", api.HealthCheckHandler(db.DB, upSince))

	// Initialize dependencies using DI pattern
	deps, err := api.InitDependencies()
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}

	// Initialize handlers with dependencies
	handlers := api.NewHandlers(deps)

	// Legacy: Keep individual references for old handlers that haven't been migrated yet
	userRepoGorm := deps.Repo.UserGorm
	keyRepo := &deps.Repo.Keys
	legacyCacheSvc := deps.Services.LegacyCache
	cfgSvc := &deps.Services.Conf
	vaMgmtSvc := &deps.Services.VaMgmt
	atApiSvc := &deps.Services.AirtableApi
	syncSvc := &deps.Services.AirtableSync
	flightSvc := &deps.Services.Flights

	// Register UI routes (separate from API)
	RegisterUIRoutes(r)

	// Setup workers and jobs first
	// Setup scheduled jobs (both pilot and route sync run every hour)
	jobsContainer := jobs.InitializeJobs(
		context.Background(),
		db.PgDB,
		deps.Services.Cache, // Use CacheInterface (supports Redis or in-memory)
		deps.Repo.DataProviderCfg,
		deps.Repo.VASyncHistory,
		deps.Repo.PilotATSynced,
		deps.Repo.RouteATSynced,
		deps.Repo.PirepATSynced,
		cfgSvc,
		&deps.Services.RedisQueue,
	)

	workers.InitWorkers(
		db.PgDB,
		&deps.Services.Cache,
		&deps.Services.Live,
		deps.Services.AircraftLivery,
		&deps.Services.RedisQueue,
		deps.Repo.AircraftLivery,
		deps.Repo.DataProviderCfg,
		deps.Repo.PirepATSynced,
		deps.Repo.VASyncHistory,
	)

	// Initialize jobs handler for manual triggering
	jobsHandler := api.NewJobsHandler(jobsContainer.PilotSync)

	// Register API routes (after jobsHandler is initialized)
	RegisterAPIRoutes(r, userRepoGorm, keyRepo, handlers, legacyCacheSvc, cfgSvc, vaMgmtSvc, atApiSvc, syncSvc, flightSvc, jobsHandler, deps)

	return r
}
