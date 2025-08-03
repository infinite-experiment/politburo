package middleware

import (
	context "infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/constants"
	"net/http"
)

func IsStaffMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())

			if claims.Role() == constants.RoleAirlineManager.String() || claims.Role() == constants.RoleAdmin.String() {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, "Unauthorized. Need staff perms", http.StatusUnauthorized)
			return
		})
	}
}
