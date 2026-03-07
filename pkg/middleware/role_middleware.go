package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

// Role allows only requests whose context role is in the allowed list.
// Use after Auth (and Tenant). Reads role from context (set by Auth middleware).
func Role(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := utils.GetUserRole(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		c.Abort()
	}
}
