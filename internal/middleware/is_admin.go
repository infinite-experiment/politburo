package middleware

import (
	context "infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/constants"
	"net/http"
)

func IsAdminMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())

			if claims.Role() != constants.RoleAdmin.String() {
				http.Error(w, "Unauthorized. Need VA Admin perms", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
