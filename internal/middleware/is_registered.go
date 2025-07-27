package middleware

import (
	"infinite-experiment/politburo/internal/context"
	"net/http"
)

func IsRegisteredMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())

			if claims.UserID() == "" {
				http.Error(w, "Unauthorized. Not a registered user", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
