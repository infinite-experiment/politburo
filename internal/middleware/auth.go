package middleware

import (
	"infinite-experiment/infinite-experiment-backend/internal/auth"
	"infinite-experiment/infinite-experiment-backend/internal/context"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		apiKey := r.Header.Get("X-API-Key")

		var claims auth.UserClaims

		switch {
		case strings.HasPrefix(authHeader, "Bearer "):
			// Parse JWT and validate
			claims = &auth.JWTClaims{
				UserIDValue: "user123",
				RoleValue:   "PILOT",
			}

		case apiKey != "":
			// TODO: Validate API Key from DB
			claims = &auth.APIKeyClaims{
				UserIDValue: "user123",
				RoleValue:   "MANAGER",
			}

		default:
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.SetUserClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
