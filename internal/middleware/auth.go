package middleware

import (
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/context"
	"infinite-experiment/politburo/internal/db/repositories"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(userRepo *repositories.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			apiKey := r.Header.Get("X-API-Key")

			log.Printf("API_KEY: %q, Authorization: %q", apiKey, authHeader)

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
				serverId := r.Header.Get("X-Server-Id")
				userId := r.Header.Get("X-Discord-Id")

				if claims = auth.MakeClaimsFromApi(r.Context(), userRepo, serverId, userId); claims == nil {
					log.Printf("No claims Found")
				}

			default:
				//http.Error(w, "Unauthorized", http.StatusUnauthorized)
				//return
			}

			ctx := context.SetUserClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
