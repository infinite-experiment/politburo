package middleware

import (
	"infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"infinite-experiment/politburo/internal/constants"
	"net/http"
)

func IsStaffMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := auth.GetUserClaims(r.Context())

			if claims.Role() == constants.RoleAirlineManager.String() || claims.Role() == constants.RoleAdmin.String() || auth.IsGodMode(claims.DiscordUserID()) {
				next.ServeHTTP(w, r)
				return
			}
			common.RespondPermissionDenied(w, "staff (airline manager or admin)")
		})
	}
}
