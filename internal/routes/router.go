package routes

import (
	"infinite-experiment/infinite-experiment-backend/internal/api"
	"net/http"

	"github.com/gorilla/mux"
	// "infinite-experiment-backend/internal/api"
	// "infinite-experiment-backend/internal/middleware"
)

func RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// TODO: Add Middlewares
	// Global Middleware
	// r.Use(middleware.RateLimitMiddleware)

	r.HandleFunc("/healthCheck", api.HealthCheckHandler).Methods("GET")

	// apiV1 := r.PathPrefix("/api/v1").Subrouter()

	// apiV1.HandleFunc

	return r

}
