package utils

import "github.com/gin-gonic/gin"

// Context keys set by JWT auth middleware. All business modules must use
// tenant_id from context to isolate data (multi-tenant SaaS).
const (
	KeyUserID   = "user_id"
	KeyTenantID = "tenant_id"
	KeyRole     = "role"
)

// GetUserID returns the authenticated user's ID from context (set by JWT middleware).
// Second return is false if not set (e.g. route not protected).
func GetUserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(KeyUserID)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetTenantID returns the current tenant ID from context (set by JWT middleware).
// Use this in all tenant-scoped handlers to isolate data: e.g. WHERE tenant_id = ?.
func GetTenantID(c *gin.Context) (string, bool) {
	v, ok := c.Get(KeyTenantID)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// GetTenantIDFromContext returns the tenant_id from request context.
// Alias for GetTenantID for clarity at call sites.
func GetTenantIDFromContext(c *gin.Context) (string, bool) {
	return GetTenantID(c)
}

// GetUserRole returns the authenticated user's role from context (set by JWT middleware).
func GetUserRole(c *gin.Context) (string, bool) {
	v, ok := c.Get(KeyRole)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}
