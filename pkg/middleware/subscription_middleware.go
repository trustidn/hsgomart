package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/internal/subscription"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

// SubscriptionChecker is implemented by subscription.Service for enforcement.
type SubscriptionChecker interface {
	CheckSubscription(tenantID string) (*subscription.SubscriptionWithPlan, error)
}

// Subscription ensures the tenant has a valid (active or trial) subscription.
// Use after Auth and Tenant middleware. Returns 402 with error message if subscription required/expired.
func Subscription(checker SubscriptionChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, ok := utils.GetTenantID(c)
		if !ok || tenantID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
			c.Abort()
			return
		}

		_, err := checker.CheckSubscription(tenantID)
		if err != nil {
			if err == subscription.ErrSubscriptionRequired {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "subscription required"})
				c.Abort()
				return
			}
			if err == subscription.ErrSubscriptionExpired {
				c.JSON(http.StatusPaymentRequired, gin.H{"error": "subscription expired"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "subscription check failed"})
			c.Abort()
			return
		}

		c.Next()
	}
}
