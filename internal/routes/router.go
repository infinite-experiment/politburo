package routes

import (
	"context"
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/api"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/jobs"
	"infinite-experiment/politburo/internal/middleware"
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

	r.Get("/public/flight", api.UserFlightMapHandler(legacyCacheSvc))
	r.Get("/public/flight/user", api.UserFlightsCacheHandler(legacyCacheSvc))

	// Setup workers

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
	// API v1 routes
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(middleware.AuthMiddleware(userRepoGorm, keyRepo)) // global: all routes must be authenticated (using GORM)
		v1.Get("/user/details", handlers.GetUserDetails())       // MIGRATED to DI
		v1.Get("/admin/verify-god", handlers.VerifyGodMode())    // God mode verification

		// Registered users group
		v1.Group(func(registered chi.Router) {
			// God-only group (admin + staff + member + registered)
			registered.Group(func(god chi.Router) {
				god.Use(middleware.IsGodMiddleware())
				god.Delete("/users/delete", handlers.DeleteAllUsers()) // MIGRATED to DI
			})
			registered.Use(middleware.IsRegisteredMiddleware())

			registered.Post("/server/init", handlers.InitServerRegistrationV2()) // V2: Uses GORM

			// Member-only group (requires registered first)
			registered.Group(func(member chi.Router) {
				member.Use(middleware.IsMemberMiddleware())

				// Pilot stats endpoint - comprehensive stats including game stats (future) and provider data
				member.Get("/pilot/stats", handlers.GetPilotStats())

				member.Get("/va/live", api.VaFlightsHandler(flightSvc))
				member.Get("/live/sessions", api.LiveServers(flightSvc))

				// Staff-only group (requires member + registered)
				member.Group(func(staff chi.Router) {
					staff.Use(middleware.IsStaffMiddleware())
					staff.Get("/user/{user_id}/flights", api.UserFlightsHandler(flightSvc, cfgSvc))
					staff.Post("/va/userSync", api.SyncUser(vaMgmtSvc))

					// Admin-only group (staff + member + registered)
					staff.Group(func(admin chi.Router) {
						admin.Use(middleware.IsAdminMiddleware())

						admin.Post("/va/setRole", api.SetRole(vaMgmtSvc))
						admin.Post("/va/configs", api.SetConfigKeys(cfgSvc))
						admin.Get("/va/configs", api.GetVAConfigs(cfgSvc))
						admin.Get("/va/configs/keys", api.ListConfigKeys(cfgSvc))
						admin.Get("/debug", api.DebugHandler(*atApiSvc, *syncSvc))

						// Data provider configuration management
						admin.Post("/admin/data-provider/config", api.SaveDataProviderConfigHandler(deps))

						// Background jobs management
						admin.Post("/admin/jobs/sync-pilots", jobsHandler.TriggerPilotSync())
						admin.Get("/admin/jobs/status", jobsHandler.GetJobStatus())

					})
				})
			})
		})

		// Public
		v1.Post("/user/register/init", handlers.InitUserRegistrationV2()) // V2: Uses GORM + LiveAPIProvider
		v1.Post("/user/register/link", handlers.LinkUserToVA())           // Link existing user to VA
	})

	return r
}
