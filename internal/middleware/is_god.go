package middleware

import (
	context "infinite-experiment/politburo/internal/auth"
	"log"
	"net/http"
	"os"
)

func IsGodMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())
			log.Printf("Discurd User ID: %s", claims.DiscordUserID())

			god_key := os.Getenv("GOD_MODE")

			if god_key != "" && claims.DiscordUserID() == god_key {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized. Need VA Admin perms", http.StatusUnauthorized)

		})
	}

}
