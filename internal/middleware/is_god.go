package middleware

import (
	context "infinite-experiment/politburo/internal/auth"
	"log"
	"net/http"
)

func IsGodMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())
			log.Printf("Discurd User ID: %s", claims.DiscordUserID())

			if claims.DiscordUserID() != "988020008665882624" {
				http.Error(w, "Unauthorized. Need VA Admin perms", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)

		})
	}

}
