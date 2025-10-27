package ui

import (
	"net/http"
)

// UIHandler manages all UI routes
type UIHandler struct{}

// NewUIHandler creates a new UI handler
func NewUIHandler() *UIHandler {
	return &UIHandler{}
}

// DashboardHandler renders the main dashboard with sidebar layout
func (h *UIHandler) DashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Vizburo Dashboard",
		"Content": "Dashboard content will go here",
		"Theme":   getThemeFromRequest(r),
	}
	RenderTemplate(w, "layouts/sidebar.html", data)
}

// MinimalDashboardHandler renders a minimal layout without sidebar
func (h *UIHandler) MinimalDashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"Title":   "Vizburo - Minimal View",
		"Content": "Minimal layout content will go here",
		"Theme":   getThemeFromRequest(r),
	}
	RenderTemplate(w, "layouts/minimal.html", data)
}

// SetThemeHandler handles theme changes via POST request
func (h *UIHandler) SetThemeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	theme := r.FormValue("theme")
	if theme == "" {
		theme = "light"
	}

	// Set theme cookie (HTTP-only, Secure in production, expires in 1 year)
	http.SetCookie(w, &http.Cookie{
		Name:     "theme_preference",
		Value:    theme,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 year
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": true, "theme": "` + theme + `"}`))
}

// HealthCheckHandler is a simple health check for the UI service
func (h *UIHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`))
}
