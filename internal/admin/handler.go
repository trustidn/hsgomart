package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type tenantRow struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	UserCount int    `json:"user_count"`
	PlanName  string `json:"plan_name"`
}

func (h *Handler) ListTenants(c *gin.Context) {
	var rows []tenantRow
	err := h.db.Raw(`
		SELECT t.id, t.name, t.email, t.status,
		       (SELECT COUNT(*) FROM users u WHERE u.tenant_id = t.id) AS user_count,
		       COALESCE((SELECT p.name FROM subscriptions s JOIN plans p ON p.id = s.plan_id WHERE s.tenant_id = t.id AND s.status IN ('active','trial') ORDER BY s.end_date DESC NULLS LAST LIMIT 1), 'None') AS plan_name
		FROM tenants t
		ORDER BY t.created_at DESC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tenants"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	var row tenantRow
	err := h.db.Raw(`
		SELECT t.id, t.name, t.email, t.status,
		       (SELECT COUNT(*) FROM users u WHERE u.tenant_id = t.id) AS user_count,
		       COALESCE((SELECT p.name FROM subscriptions s JOIN plans p ON p.id = s.plan_id WHERE s.tenant_id = t.id AND s.status IN ('active','trial') ORDER BY s.end_date DESC NULLS LAST LIMIT 1), 'None') AS plan_name
		FROM tenants t WHERE t.id = ?
	`, id).Scan(&row).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	c.JSON(http.StatusOK, row)
}

type UpdateTenantInput struct {
	Status string `json:"status" binding:"required"`
}

func (h *Handler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	var in UpdateTenantInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res := h.db.Table("tenants").Where("id = ?", id).Update("status", in.Status)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) ListSubscriptions(c *gin.Context) {
	var rows []struct {
		ID         string  `json:"id"`
		TenantID   string  `json:"tenant_id"`
		TenantName string  `json:"tenant_name"`
		PlanName   string  `json:"plan_name"`
		Status     string  `json:"status"`
		EndDate    *string `json:"end_date"`
	}
	err := h.db.Raw(`
		SELECT s.id, s.tenant_id, t.name AS tenant_name, p.name AS plan_name, s.status,
		       TO_CHAR(s.end_date, 'YYYY-MM-DD') AS end_date
		FROM subscriptions s
		JOIN tenants t ON t.id = s.tenant_id
		JOIN plans p ON p.id = s.plan_id
		ORDER BY s.end_date DESC NULLS LAST
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

type UpdateSubscriptionInput struct {
	PlanID *int    `json:"plan_id"`
	Status *string `json:"status"`
}

func (h *Handler) UpdateSubscription(c *gin.Context) {
	id := c.Param("id")
	var in UpdateSubscriptionInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if in.PlanID != nil {
		updates["plan_id"] = *in.PlanID
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	res := h.db.Table("subscriptions").Where("id = ?", id).Updates(updates)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) Stats(c *gin.Context) {
	var totalTenants int64
	h.db.Table("tenants").Count(&totalTenants)

	var totalTransactions int64
	h.db.Table("transactions").Count(&totalTransactions)

	var totalRevenue float64
	h.db.Table("transactions").Where("status = 'completed'").Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	c.JSON(http.StatusOK, gin.H{
		"total_tenants":      totalTenants,
		"total_transactions": totalTransactions,
		"total_revenue":      totalRevenue,
	})
}
