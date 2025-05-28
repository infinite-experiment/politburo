package routes

import (
	"net/http"
	"time"

	"infinite-experiment/politburo/internal/api"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/services"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(upSince time.Time) http.Handler {
	// initialize Chi router
	r := chi.NewRouter()

	// global middleware
	r.Use(middleware.Logging)
	r.Use(middleware.AuthMiddleware(repositories.NewUserRepository(db.DB)))

	// health check
	r.Get("/healthCheck", api.HealthCheckHandler(db.DB, upSince))

	// swagger UI
	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	// services
	userRepo := repositories.NewUserRepository(db.DB)
	cacheSvc := services.NewCacheService(60000, 600)
	liveSvc := services.NewLiveAPIService()
	flightSvc := services.NewFlightsService(*cacheSvc, liveSvc)
	regSvc := services.NewRegistrationService(liveSvc, *cacheSvc, *userRepo)
	api.SetUserService(services.NewUserService(userRepo))

	// API v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/user/register", api.RegisterUserHandler)
		r.Get("/user/{user_id}/flights", api.UserFlightsHandler(flightSvc))
		r.Get("/users/delete", api.DeleteAllUsers(userRepo))
		r.Post("/user/register/init", api.InitUserRegistrationHandler(regSvc))
	})

	return r
}
