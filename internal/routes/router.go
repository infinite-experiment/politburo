package routes

import (
	"infinite-experiment/politburo/internal/api"
	"infinite-experiment/politburo/internal/db"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/services"
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
	cacheService := services.NewCacheService(60000, 600)
	liveApiService := services.NewLiveAPIService()

	userRegistationService := services.NewRegistrationService(liveApiService, *cacheService, *userRepo)
	// userService := services.NewUserService(userRepo)
	api.SetUserService(services.NewUserService(userRepo))

	r.Use(middleware.AuthMiddleware(userRepo))

	r.HandleFunc("/healthCheck", api.HealthCheckHandler(db.DB, upSince)).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/user/register", api.RegisterUserHandler).Methods("POST")
	apiV1.HandleFunc("/users/delete", api.DeleteAllUsers(userRepo)).Methods("GET")
	apiV1.HandleFunc("/user/register/init", api.InitUserRegistrationHandler(userRegistationService)).Methods("POST")

	return r

}
