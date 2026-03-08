package subscription

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetSubscription(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	sub, err := GetActiveSubscription(h.service.db, tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get subscription"})
		return
	}
	if sub == nil {
		c.JSON(http.StatusOK, gin.H{"subscription": nil})
		return
	}

	var trialDaysLeft *int
	var daysRemaining *int
	if sub.Subscription.EndDate != nil {
		days := int(math.Ceil(time.Until(*sub.Subscription.EndDate).Hours() / 24))
		if days < 0 {
			days = 0
		}
		daysRemaining = &days
		if sub.Subscription.Status == "trial" {
			trialDaysLeft = &days
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"subscription": gin.H{
			"id":         sub.Subscription.ID,
			"plan_name":  sub.Plan.Name,
			"plan_price": sub.Plan.Price,
			"status":     sub.Subscription.Status,
			"start_date": sub.Subscription.StartDate,
			"end_date":   sub.Subscription.EndDate,
		},
		"plan": gin.H{
			"id":           sub.Plan.ID,
			"name":         sub.Plan.Name,
			"price":        sub.Plan.Price,
			"max_users":    sub.Plan.MaxUsers,
			"max_products": sub.Plan.MaxProducts,
		},
		"trial_days_left": trialDaysLeft,
		"days_remaining":  daysRemaining,
	})
}

func (h *Handler) ListPlans(c *gin.Context) {
	var plans []Plan
	if err := h.service.db.Where("is_active = true OR is_active IS NULL").Order("price ASC").Find(&plans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list plans"})
		return
	}
	c.JSON(http.StatusOK, plans)
}

type ChangePlanInput struct {
	PlanID int `json:"plan_id" binding:"required"`
}

func (h *Handler) ChangePlan(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	var in ChangePlanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var plan Plan
	if err := h.service.db.Where("id = ?", in.PlanID).First(&plan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}

	res := h.service.db.Model(&Subscription{}).
		Where("tenant_id = ? AND status IN ?", tenantID, []string{"active", "trial"}).
		Updates(map[string]interface{}{"plan_id": in.PlanID, "status": "active"})
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "plan changed", "plan": plan.Name})
}
