package routes

import (
	"net/http"
	"path/filepath"
	"strings"

	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/db/repositories"
	"infinite-experiment/politburo/internal/middleware"
	"infinite-experiment/politburo/internal/services"
	vizbuUI "infinite-experiment/politburo/vizburo/ui"

	"github.com/go-chi/chi/v5"
)

// RegisterUIRoutes registers all UI-related routes
func RegisterUIRoutes(
	r chi.Router,
	sessionSvc *common.SessionService,
	urlSigner *common.URLSignerService,
	userRepo *repositories.UserRepositoryGORM,
	vaRoleRepo *repositories.VAUserRoleRepository,
	vaRepo *repositories.VAGORMRepository,
	flightSvc *services.FlightsService,
	cache common.CacheInterface,
	liveAPI *common.LiveAPIService,
) {
	authHandler := vizbuUI.NewAuthHandler(sessionSvc, urlSigner, userRepo, vaRoleRepo, vaRepo)

	// Import middleware
	authMiddleware := middleware.AuthMiddleware(userRepo, nil, sessionSvc) // keysRepo is nil for UI routes

	// Static file serving (CSS, JS, images) with correct MIME types
	fileServer := http.FileServer(http.Dir("vizburo/ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", mimeTypeMiddleware(fileServer)))

	// Default route - redirect to login
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/auth/login", http.StatusMovedPermanently)
	})

	// Auth routes (public)
	r.Get("/auth/login", authHandler.TokenLoginHandler)
	r.Post("/auth/logout", authHandler.LogoutHandler)

	// Dashboard routes (require authentication)
	r.Route("/dashboard", func(dashboard chi.Router) {
		// Apply authentication middleware to all dashboard routes
		dashboard.Use(authMiddleware)

		// Main dashboard page
		dashboard.Get("/", vizbuUI.DashboardHandler)

		// Logbook page (role-checked in handler)
		dashboard.Get("/logbook", vizbuUI.LogbookHandler)

		// Logbook HTMX endpoints
		dashboard.Get("/logbook/flights", func(w http.ResponseWriter, r *http.Request) {
			vizbuUI.LogbookFlightsHandler(w, r, flightSvc)
		})
		dashboard.Get("/logbook/flight/{session_id}/{flight_id}/map", func(w http.ResponseWriter, r *http.Request) {
			vizbuUI.FlightMapHandler(w, r, cache, liveAPI, flightSvc)
		})
		dashboard.Get("/logbook/pilots/search", func(w http.ResponseWriter, r *http.Request) {
			vizbuUI.PilotSearchHandler(w, r, vaRoleRepo)
		})
		dashboard.Get("/logbook/map/reset", vizbuUI.MapResetHandler)

		// HTMX VA switch endpoint
		dashboard.Post("/switch-va", authHandler.SwitchVAHandler)
	})

	// UI API routes
	r.Route("/ui/api", func(uiApi chi.Router) {
		uiApi.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		})
	})
}

// mimeTypeMiddleware wraps a file server and sets correct MIME types for various file types
func mimeTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the file extension
		ext := filepath.Ext(r.URL.Path)

		// Set correct MIME type for .mjs files (ES modules)
		if strings.EqualFold(ext, ".mjs") {
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		}

		// Call the wrapped handler
		next.ServeHTTP(w, r)
	})
}
