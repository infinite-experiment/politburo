package middleware

import (
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"log"
	"net/http"
)

func IsGodMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := auth.GetUserClaims(r.Context())
			log.Printf("Discord User ID: %s", claims.DiscordUserID())

			if auth.IsGodMode(claims.DiscordUserID()) {
				next.ServeHTTP(w, r)
				return
			}
			common.RespondPermissionDenied(w, "god mode (system administrator)")

		})
	}

}
