package auth

import (
	"context"
)

type contextKey string

var userClaimsKey contextKey = "user_claims"
var sessionDataKey contextKey = "session_data"

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

// SetSessionData stores session data in context for use by handlers (e.g., VA switcher)
func SetSessionData(ctx context.Context, sessionData interface{}) context.Context {
	return context.WithValue(ctx, sessionDataKey, sessionData)
}

// GetSessionData retrieves session data from context
func GetSessionData(ctx context.Context) interface{} {
	return ctx.Value(sessionDataKey)
}
