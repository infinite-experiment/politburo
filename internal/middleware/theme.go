package middleware

import (
	"context"
	"net/http"
)

// ThemeMiddleware injects the user's theme preference into the request context
func ThemeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get theme from cookie, default to "light"
		theme := "light"
		cookie, err := r.Cookie("theme_preference")
		if err == nil && cookie.Value != "" {
			theme = cookie.Value
		}

		// Validate theme value (prevent injection)
		validThemes := map[string]bool{
			"light":         true,
			"dark":          true,
			"high-contrast": true,
		}

		if !validThemes[theme] {
			theme = "light"
		}

		// Store in context for template access
		ctx := context.WithValue(r.Context(), "theme", theme)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
