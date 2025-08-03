package auth

import (
	"context"
)

type contextKey string

var userClaimsKey contextKey = "user_claims"

func SetUserClaims(ctx context.Context, claims UserClaims) context.Context {
	return context.WithValue(ctx, userClaimsKey, claims)
}

func GetUserClaims(ctx context.Context) UserClaims {
	val := ctx.Value(userClaimsKey)
	if claims, ok := val.(UserClaims); ok {
		return claims
	}
	return nil
}
