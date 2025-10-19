package routes

import (
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/api"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/workers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(upSince time.Time) http.Handler {

	// initialize Chi router
	r := chi.NewRouter()

	// global middleware
	//r.Use(middleware.Logging)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://localhost:8081"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// health check
	r.Get("/healthCheck", api.HealthCheckHandler(db.DB, upSince))

	// swagger UI
	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	// Initialize dependencies using DI pattern
	deps, err := api.InitDependencies()
	if err != nil {
		panic("Failed to initialize dependencies: " + err.Error())
	}

	// Initialize handlers with dependencies
	handlers := api.NewHandlers(deps)

	// Legacy: Keep individual references for old handlers that haven't been migrated yet
	userRepo := &deps.Repo.User
	keyRepo := &deps.Repo.Keys
	cacheSvc := &deps.Services.Cache
	liveSvc := &deps.Services.Live
	regSvc := &deps.Services.Reg
	cfgSvc := &deps.Services.Conf
	vaMgmtSvc := &deps.Services.VaMgmt
	atApiSvc := &deps.Services.AirtableApi
	syncSvc := &deps.Services.AirtableSync
	flightSvc := &deps.Services.Flights

	r.Get("/public/flight", api.UserFlightMapHandler(cacheSvc))
	r.Get("/public/flight/user", api.UserFlightsCacheHandler(cacheSvc))

	//Setup
	go workers.LogbookWorker(cacheSvc, liveSvc)
	go workers.StartCacheFiller(cacheSvc, liveSvc)
	// API v1 routes
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(middleware.AuthMiddleware(userRepo, keyRepo)) // global: all routes must be authenticated

		// Registered users group
		v1.Group(func(registered chi.Router) {
			// God-only group (admin + staff + member + registered)
			registered.Group(func(god chi.Router) {
				god.Use(middleware.IsGodMiddleware())
				god.Delete("/users/delete", handlers.DeleteAllUsers()) // MIGRATED to DI
			})
			registered.Use(middleware.IsRegisteredMiddleware())

			registered.Post("/server/init", api.InitRegisterServer(regSvc))

			// Member-only group (requires registered first)
			registered.Group(func(member chi.Router) {
				member.Use(middleware.IsMemberMiddleware())

				// NEW: User details endpoint with GORM
				member.Get("/user/details", handlers.GetUserDetails()) // MIGRATED to DI

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

					})
				})
			})
		})

		// Public
		v1.Post("/user/register/init", handlers.InitUserRegistration()) // MIGRATED to DI
	})

	return r
}
