package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

// Tenant ensures tenant_id is present in context (set by Auth middleware).
// Use after Auth middleware on all tenant-scoped routes so downstream handlers
// can safely use utils.GetTenantID(c) for data isolation (e.g. filter by tenant_id).
func Tenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := utils.GetTenantIDFromContext(c)
		if !ok || tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
			c.Abort()
			return
		}
		// tenant_id already attached by Auth middleware; continue to handlers
		c.Next()
	}
}
