package context

import (
	"context"
	"infinite-experiment/politburo/internal/auth"
)

type contextKey string

var userClaimsKey contextKey = "user_claims"

func SetUserClaims(ctx context.Context, claims auth.UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

func GetUserClaims(ctx context.Context) auth.UserClaims {
	val := ctx.Value(userClaimsKey)
	if claims, ok := val.(auth.UserClaims); ok {
		return claims
	}
	return nil
}
