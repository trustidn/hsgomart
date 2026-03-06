package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

// Owner ensures the authenticated user has role "owner".
// Use after Auth (and Tenant) for owner-only routes (e.g. purchase management).
func Owner() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := utils.GetUserRole(c)
		if !ok || role != "owner" {
			c.JSON(http.StatusForbidden, gin.H{"error": "owner access required"})
			c.Abort()
			return
		}
		c.Next()
	}
}
