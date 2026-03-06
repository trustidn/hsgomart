package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/internal/auth"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type TokenValidator interface {
	ValidateToken(tokenString string) (*auth.JWTClaims, error)
}

// Auth validates the JWT from Authorization: Bearer <token> and sets user_id,
// tenant_id, and role on the request context for tenant isolation in downstream handlers.
func Auth(validator TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		claims, err := validator.ValidateToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set(utils.KeyUserID, claims.UserID)
		c.Set(utils.KeyTenantID, claims.TenantID)
		c.Set(utils.KeyRole, claims.Role)
		c.Next()
	}
}
