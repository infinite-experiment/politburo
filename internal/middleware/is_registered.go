package middleware

import (
	context "infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"log"
	"net/http"
)

func IsRegisteredMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())

			log.Printf("User ID: %s, God ID: %s", claims.UserID(), claims.DiscordUserID())
			if claims.UserID() == "" && !context.IsGodMode(claims.DiscordUserID()) {
				common.RespondPermissionDenied(w, "registered user")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
