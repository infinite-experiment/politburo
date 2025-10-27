package routes

import (
	"net/http"

	vizbuUI "infinite-experiment/politburo/vizburo/ui"
	"github.com/go-chi/chi/v5"
)

// RegisterUIRoutes registers all UI-related routes
func RegisterUIRoutes(r chi.Router) {
	uiHandler := vizbuUI.NewUIHandler()

	// Static file serving (CSS, JS, images)
	fileServer := http.FileServer(http.Dir("vizburo/ui/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	// Default route - redirect to dashboard
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
	})

	// Dashboard routes (no authentication required for Phase 1)
	r.Get("/dashboard", uiHandler.DashboardHandler)
	r.Get("/dashboard/minimal", uiHandler.MinimalDashboardHandler)

	// Component Showcase routes (for testing and design)
	r.Get("/dashboard/showcase", uiHandler.ShowcaseHandler)
	r.Get("/dashboard/showcase/page", uiHandler.ShowcasePageHandler)

	// Component endpoints (HTMX)
	r.Route("/dashboard/showcase/component", func(componentApi chi.Router) {
		componentApi.Get("/buttons", uiHandler.ComponentButtonsHandler)
		componentApi.Get("/forms", uiHandler.ComponentFormsHandler)
		componentApi.Get("/typography", uiHandler.ComponentTypographyHandler)
		componentApi.Get("/cards", uiHandler.ComponentCardsHandler)
		componentApi.Get("/badges", uiHandler.ComponentBadgesHandler)
		componentApi.Get("/alerts", uiHandler.ComponentAlertsHandler)
	})

	// Flight visualization routes
	r.Get("/flights", uiHandler.FlightsListHandler)
	r.Route("/flights/api", func(flightsApi chi.Router) {
		flightsApi.Get("/list", uiHandler.FlightsListAPIHandler)
		flightsApi.Get("/details/{id}", uiHandler.FlightsDetailsAPIHandler)
	})

	// UI API routes
	r.Route("/ui/api", func(uiApi chi.Router) {
		uiApi.Post("/theme", uiHandler.SetThemeHandler)
		uiApi.Get("/health", uiHandler.HealthCheckHandler)
	})
}
