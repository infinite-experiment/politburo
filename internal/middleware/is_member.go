package middleware

import (
	context "infinite-experiment/politburo/internal/auth"
	"infinite-experiment/politburo/internal/common"
	"net/http"
)

func IsMemberMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			claims := context.GetUserClaims(r.Context())

			if claims.Role() == "" {
				common.RespondPermissionDenied(w, "member (pilot)")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
