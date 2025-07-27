package routes

import (
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/api"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/services"
	"infinite-experiment/politburo/internal/workers"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(upSince time.Time) http.Handler {

	// initialize Chi router
	r := chi.NewRouter()

	// global middleware
	//r.Use(middleware.Logging)

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"https://*", "http://localhost:8081"}, // Allow all origins
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// 	ExposedHeaders:   []string{"Link"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// }))
	// health check
	r.Get("/healthCheck", api.HealthCheckHandler(db.DB, upSince))

	// swagger UI
	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	// services
	userRepo := repositories.NewUserRepository(db.DB)
	keyRepo := repositories.NewApiKeysRepo(db.DB)

	cacheSvc := common.NewCacheService(60000, 600)
	liveSvc := common.NewLiveAPIService()
	vaRepo := repositories.NewVARepository(db.DB)
	regSvc := services.NewRegistrationService(liveSvc, *cacheSvc, *userRepo, *vaRepo)
	cfgSvc := common.NewVAConfigService(vaRepo, cacheSvc)
	flightSvc := services.NewFlightsService(cacheSvc, liveSvc, cfgSvc)

	r.Get("/public/flight", api.UserFlightMapHandler(cacheSvc))

	//Setup
	go workers.LogbookWorker(cacheSvc, liveSvc)
	go workers.StartCacheFiller(cacheSvc, liveSvc)
	// API v1 routes
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Use(middleware.AuthMiddleware(userRepo, keyRepo))

		v1.Group(func(god chi.Router) {
			god.Use(middleware.IsGodMiddleware())
			god.Delete("/users/delete", api.DeleteAllUsers(userRepo))
		})

		// Admin routes
		v1.Group(func(admin chi.Router) {
			admin.Use(middleware.IsAdminMiddleware())
			admin.Post("/va/configs", api.SetConfigKeys(cfgSvc))
			admin.Get("/va/configs", api.GetVAConfigs(cfgSvc))
			admin.Get("/va/configs/keys", api.ListConfigKeys(cfgSvc))
		})

		// Pilot routes
		v1.Group(func(pilot chi.Router) {
			pilot.Use(middleware.IsMemberMiddleware())
			pilot.Get("/va/live", api.VaFlightsHandler(flightSvc))
			pilot.Get("/live/sessions", api.LiveServers(flightSvc))
		})

		// Staff routes
		v1.Group(func(staff chi.Router) {
			staff.Use(middleware.IsStaffMiddleware())
			staff.Get("/user/{user_id}/flights", api.UserFlightsHandler(flightSvc))
		})

		// Registered users routes
		v1.Group(func(mem chi.Router) {
			mem.Use(middleware.IsRegisteredMiddleware())
			mem.Post("/server/init", api.InitRegisterServer(regSvc))
		})

		// Publicly open routes
		v1.Post("/user/register/init", api.InitUserRegistrationHandler(regSvc))
	})

	return r
}
