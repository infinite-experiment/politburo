package middleware

import (
	"infinite-experiment/politburo/internal/context"
	"net/http"
)

func IsGodMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())

			if claims.DiscordUserID() == "988020008665882624" {
				http.Error(w, "Unauthorized. Need VA Admin perms", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)

		})
	}

}
