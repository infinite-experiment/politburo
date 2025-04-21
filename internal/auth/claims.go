package auth

type UserClaims interface {
	UserID() string
	Role() string
	Source() string
	HasPermission(action string) bool
}

type JWTClaims struct {
	UserIDValue string
	RoleValue   string
}

func (c *JWTClaims) UserID() string { return c.UserIDValue }
func (c *JWTClaims) Role() string   { return c.RoleValue }
func (c *JWTClaims) Source() string { return "JWT" }
func (c *JWTClaims) HasPermission(p string) bool {
	return true
}

type APIKeyClaims struct {
	UserIDValue string
	RoleValue   string
}

func (c *APIKeyClaims) UserID() string { return c.UserIDValue }
func (c *APIKeyClaims) Role() string   { return c.RoleValue }
func (c *APIKeyClaims) Source() string { return "API_KEY" }
func (c *APIKeyClaims) HasPermission(p string) bool {
	return true
}
