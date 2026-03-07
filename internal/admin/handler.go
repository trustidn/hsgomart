package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
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

// --- Subscription Order Management ---

type orderRow struct {
	ID              string  `json:"id"`
	TenantID        string  `json:"tenant_id"`
	TenantName      string  `json:"tenant_name"`
	PlanName        string  `json:"plan_name"`
	Amount          float64 `json:"amount"`
	Status          string  `json:"status"`
	InvoiceNumber   string  `json:"invoice_number"`
	PaymentProofURL string  `json:"payment_proof_url"`
	Notes           string  `json:"notes"`
	AdminNotes      string  `json:"admin_notes"`
	CreatedAt       string  `json:"created_at"`
	PaidAt          *string `json:"paid_at"`
	ReviewedAt      *string `json:"reviewed_at"`
}

func (h *Handler) ListOrders(c *gin.Context) {
	status := c.DefaultQuery("status", "")
	query := `
		SELECT o.id, o.tenant_id, t.name AS tenant_name, p.name AS plan_name,
		       o.amount, o.status, o.invoice_number, COALESCE(o.payment_proof_url,'') AS payment_proof_url,
		       COALESCE(o.notes,'') AS notes, COALESCE(o.admin_notes,'') AS admin_notes,
		       TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
		       TO_CHAR(o.paid_at, 'YYYY-MM-DD HH24:MI') AS paid_at,
		       TO_CHAR(o.reviewed_at, 'YYYY-MM-DD HH24:MI') AS reviewed_at
		FROM subscription_orders o
		JOIN tenants t ON t.id = o.tenant_id
		JOIN plans p ON p.id = o.plan_id
	`
	var rows []orderRow
	var err error
	if status != "" {
		err = h.db.Raw(query+" WHERE o.status = ? ORDER BY o.created_at DESC", status).Scan(&rows).Error
	} else {
		err = h.db.Raw(query + " ORDER BY o.created_at DESC").Scan(&rows).Error
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) GetOrderDetail(c *gin.Context) {
	id := c.Param("id")
	var row orderRow
	err := h.db.Raw(`
		SELECT o.id, o.tenant_id, t.name AS tenant_name, p.name AS plan_name,
		       o.amount, o.status, o.invoice_number, COALESCE(o.payment_proof_url,'') AS payment_proof_url,
		       COALESCE(o.notes,'') AS notes, COALESCE(o.admin_notes,'') AS admin_notes,
		       TO_CHAR(o.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
		       TO_CHAR(o.paid_at, 'YYYY-MM-DD HH24:MI') AS paid_at,
		       TO_CHAR(o.reviewed_at, 'YYYY-MM-DD HH24:MI') AS reviewed_at
		FROM subscription_orders o
		JOIN tenants t ON t.id = o.tenant_id
		JOIN plans p ON p.id = o.plan_id
		WHERE o.id = ?
	`, id).Scan(&row).Error
	if err != nil || row.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, row)
}

type ApproveOrderInput struct {
	AdminNotes string `json:"admin_notes"`
}

func (h *Handler) ApproveOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := utils.GetUserID(c)

	var in ApproveOrderInput
	c.ShouldBindJSON(&in)

	tx := h.db.Begin()

	var order struct {
		ID       string  `json:"id"`
		TenantID string  `json:"tenant_id"`
		PlanID   int     `json:"plan_id"`
		Status   string  `json:"status"`
		Amount   float64 `json:"amount"`
	}
	if err := tx.Raw("SELECT id, tenant_id, plan_id, status, amount FROM subscription_orders WHERE id = ?", id).Scan(&order).Error; err != nil || order.ID == "" {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.Status != "pending_review" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "order is not pending review"})
		return
	}

	now := time.Now()
	endDate := now.AddDate(0, 1, 0)
	if err := tx.Exec(`
		UPDATE subscription_orders SET status = 'approved', admin_notes = ?, reviewed_by = ?, reviewed_at = ? WHERE id = ?
	`, in.AdminNotes, userID, now, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve order"})
		return
	}

	// Activate or update subscription
	res := tx.Exec(`
		UPDATE subscriptions SET plan_id = ?, status = 'active', start_date = ?, end_date = ?
		WHERE tenant_id = ? AND status IN ('active', 'trial')
	`, order.PlanID, now, endDate, order.TenantID)
	if res.RowsAffected == 0 {
		tx.Exec(`
			INSERT INTO subscriptions (tenant_id, plan_id, status, start_date, end_date)
			VALUES (?, ?, 'active', ?, ?)
		`, order.TenantID, order.PlanID, now, endDate)
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "order approved, subscription activated"})
}

type RejectOrderInput struct {
	AdminNotes string `json:"admin_notes" binding:"required"`
}

func (h *Handler) RejectOrder(c *gin.Context) {
	id := c.Param("id")
	userID, _ := utils.GetUserID(c)

	var in RejectOrderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin_notes required"})
		return
	}

	now := time.Now()
	res := h.db.Exec(`
		UPDATE subscription_orders SET status = 'rejected', admin_notes = ?, reviewed_by = ?, reviewed_at = ?
		WHERE id = ? AND status = 'pending_review'
	`, in.AdminNotes, userID, now, id)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order not found or not pending review"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order rejected"})
}
