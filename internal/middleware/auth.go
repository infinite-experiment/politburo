package middleware

import (
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/context"
	"infinite-experiment/politburo/internal/db/repositories"
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(userRepo *repositories.UserRepository, keysRepo *repositories.KeysRepo) func(http.Handler) http.Handler {
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
				http.Error(w, "Unauthorized. Missing API Key", http.StatusUnauthorized)
				return

			case apiKey != "":
				// TODO: Validate API Key from DB
				serverId := r.Header.Get("X-Server-Id")
				userId := r.Header.Get("X-Discord-Id")

				keyRes, err := keysRepo.GetStatus(r.Context(), apiKey)
				if err != nil {
					http.Error(w, "Unauthorized. Invalid API Key", http.StatusUnauthorized)
					return
				}

				if !keyRes.Status {
					http.Error(w, "Unauthorized. Inactive API Key", http.StatusUnauthorized)
					return
				}

				claims = auth.MakeClaimsFromApi(r.Context(), userRepo, serverId, userId)
				log.Printf("Claims: %v", claims)

			default:
				http.Error(w, "Unauthorized. Unknown Error", http.StatusUnauthorized)
				return
			}

			ctx := context.SetUserClaims(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
