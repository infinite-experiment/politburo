package routes

import (
	"infinite-experiment/infinite-experiment-backend/internal/api"
	"infinite-experiment/infinite-experiment-backend/internal/db"
	"infinite-experiment/infinite-experiment-backend/internal/db/repositories"
	"infinite-experiment/infinite-experiment-backend/internal/services"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes(upSince time.Time) http.Handler {
	r := mux.NewRouter()

	// TODO: Add Middlewares
	// Global Middleware
	// r.Use(middleware.RateLimitMiddleware)

	userRepo := repositories.NewUserRepository(db.DB)
	// userService := services.NewUserService(userRepo)
	api.SetUserService(services.NewUserService(userRepo))

	r.HandleFunc("/healthCheck", api.HealthCheckHandler(db.DB, upSince)).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/user/register", api.RegisterUserHandler).Methods("POST")

	return r

}
