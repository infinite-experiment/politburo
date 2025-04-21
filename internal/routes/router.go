package routes

import (
	"infinite-experiment/infinite-experiment-backend/internal/api"
	"infinite-experiment/infinite-experiment-backend/internal/db"
	"infinite-experiment/infinite-experiment-backend/internal/db/repositories"
	"infinite-experiment/infinite-experiment-backend/internal/services"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// TODO: Add Middlewares
	// Global Middleware
	// r.Use(middleware.RateLimitMiddleware)

	userRepo := repositories.NewUserRepository(db.DB)
	api.SetUserService(services.NewUserService(userRepo))

	r.HandleFunc("/healthCheck", api.HealthCheckHandler).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.HandleFunc("/user/register", api.RegisterUserHandler).Methods("POST")

	return r

}
