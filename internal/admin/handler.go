package admin

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

// ─── Tenant Management ──────────────────────────────────────────────────────

type tenantRow struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Email          string  `json:"email"`
	Phone          string  `json:"phone"`
	Status         string  `json:"status"`
	UserCount      int     `json:"user_count"`
	PlanName       string  `json:"plan_name"`
	PlanID         *int    `json:"plan_id"`
	SubStatus      string  `json:"sub_status"`
	EndDate        *string `json:"end_date"`
	DaysRemaining  *int    `json:"days_remaining"`
}

func (h *Handler) ListTenants(c *gin.Context) {
	statusFilter := c.DefaultQuery("status", "")
	var rows []tenantRow
	query := `
		SELECT t.id, t.name, t.email, COALESCE(t.phone, '') AS phone, t.status,
		       (SELECT COUNT(*) FROM users u WHERE u.tenant_id = t.id AND u.deleted_at IS NULL) AS user_count,
		       COALESCE(s.plan_name, 'None') AS plan_name,
		       s.plan_id,
		       COALESCE(s.sub_status, '') AS sub_status,
		       s.end_date,
		       CASE WHEN s.end_date IS NOT NULL
		            THEN GREATEST(0, EXTRACT(DAY FROM (s.end_date::timestamp - NOW())))::int
		            ELSE NULL END AS days_remaining
		FROM tenants t
		LEFT JOIN LATERAL (
			SELECT p.name AS plan_name, p.id AS plan_id, sub.status AS sub_status,
			       TO_CHAR(sub.end_date, 'YYYY-MM-DD') AS end_date,
			       sub.end_date AS end_date_raw
			FROM subscriptions sub
			JOIN plans p ON p.id = sub.plan_id
			WHERE sub.tenant_id = t.id AND sub.status IN ('active','trial')
			ORDER BY sub.end_date DESC NULLS LAST
			LIMIT 1
		) s ON true
		WHERE t.status != 'deleted'
	`
	var err error
	if statusFilter != "" {
		err = h.db.Raw(query+" AND t.status = ? ORDER BY t.created_at DESC", statusFilter).Scan(&rows).Error
	} else {
		err = h.db.Raw(query + " ORDER BY t.created_at DESC").Scan(&rows).Error
	}
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
		SELECT t.id, t.name, t.email, COALESCE(t.phone, '') AS phone, t.status,
		       (SELECT COUNT(*) FROM users u WHERE u.tenant_id = t.id AND u.deleted_at IS NULL) AS user_count,
		       COALESCE(s.plan_name, 'None') AS plan_name,
		       s.plan_id,
		       COALESCE(s.sub_status, '') AS sub_status,
		       s.end_date,
		       CASE WHEN s.end_date IS NOT NULL
		            THEN GREATEST(0, EXTRACT(DAY FROM (s.end_date::timestamp - NOW())))::int
		            ELSE NULL END AS days_remaining
		FROM tenants t
		LEFT JOIN LATERAL (
			SELECT p.name AS plan_name, p.id AS plan_id, sub.status AS sub_status,
			       TO_CHAR(sub.end_date, 'YYYY-MM-DD') AS end_date
			FROM subscriptions sub
			JOIN plans p ON p.id = sub.plan_id
			WHERE sub.tenant_id = t.id AND sub.status IN ('active','trial')
			ORDER BY sub.end_date DESC NULLS LAST
			LIMIT 1
		) s ON true
		WHERE t.id = ?
	`, id).Scan(&row).Error
	if err != nil || row.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
		return
	}
	c.JSON(http.StatusOK, row)
}

type CreateTenantInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	PlanID   *int   `json:"plan_id"`
}

func (h *Handler) CreateTenant(c *gin.Context) {
	var in CreateTenantInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := utils.ValidatePasswordStrength(in.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing int64
	h.db.Table("users").Where("email = ?", in.Email).Count(&existing)
	if existing > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}

	tx := h.db.Begin()

	var tenantID string
	if err := tx.Raw(`INSERT INTO tenants (name, email, status) VALUES (?, ?, 'active') RETURNING id`, in.Name, in.Email).Scan(&tenantID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	if err := tx.Exec(`INSERT INTO users (tenant_id, name, email, password_hash, role, status) VALUES (?, ?, ?, ?, 'owner', 'active')`, tenantID, in.Name, in.Email, string(hash)).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	planID := 0
	subStatus := "trial"
	if in.PlanID != nil && *in.PlanID > 0 {
		var plan struct {
			ID    int
			Price float64
		}
		if err := tx.Raw("SELECT id, price FROM plans WHERE id = ?", *in.PlanID).Scan(&plan).Error; err != nil || plan.ID == 0 {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan"})
			return
		}
		planID = plan.ID
		if plan.Price > 0 {
			subStatus = "active"
		}
	} else {
		var trialPlan struct{ ID int }
		if err := tx.Raw("SELECT id FROM plans WHERE name = 'Trial' LIMIT 1").Scan(&trialPlan).Error; err != nil || trialPlan.ID == 0 {
			tx.Exec("INSERT INTO plans (name, price, max_users, max_products) VALUES ('Trial', 0, 5, 100)")
			tx.Raw("SELECT id FROM plans WHERE name = 'Trial' LIMIT 1").Scan(&trialPlan)
		}
		planID = trialPlan.ID
	}

	now := time.Now()
	endDate := now.AddDate(0, 1, 0)
	if err := tx.Exec(`INSERT INTO subscriptions (tenant_id, plan_id, status, start_date, end_date) VALUES (?, ?, ?, ?, ?)`, tenantID, planID, subStatus, now, endDate).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create subscription"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "tenant created", "tenant_id": tenantID})
}

type UpdateTenantInput struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Phone  *string `json:"phone"`
	Status *string `json:"status"`
	PlanID *int    `json:"plan_id"`
}

func (h *Handler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	var in UpdateTenantInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := h.db.Begin()

	tenantUpdates := map[string]interface{}{}
	if in.Name != nil {
		tenantUpdates["name"] = *in.Name
	}
	if in.Email != nil {
		tenantUpdates["email"] = *in.Email
	}
	if in.Phone != nil {
		tenantUpdates["phone"] = *in.Phone
	}
	if in.Status != nil {
		tenantUpdates["status"] = *in.Status
	}
	if len(tenantUpdates) > 0 {
		if err := tx.Table("tenants").Where("id = ?", id).Updates(tenantUpdates).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tenant"})
			return
		}
	}

	if in.PlanID != nil {
		now := time.Now()
		endDate := now.AddDate(0, 1, 0)
		res := tx.Exec(`UPDATE subscriptions SET plan_id = ?, status = 'active', start_date = ?, end_date = ? WHERE tenant_id = ? AND status IN ('active','trial')`, *in.PlanID, now, endDate, id)
		if res.RowsAffected == 0 {
			tx.Exec(`INSERT INTO subscriptions (tenant_id, plan_id, status, start_date, end_date) VALUES (?, ?, 'active', ?, ?)`, id, *in.PlanID, now, endDate)
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	tx := h.db.Begin()
	tx.Exec("UPDATE tenants SET status = 'deleted' WHERE id = ?", id)
	tx.Exec("UPDATE subscriptions SET status = 'expired' WHERE tenant_id = ?", id)
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted"})
}

// ─── Subscriptions ──────────────────────────────────────────────────────────

func (h *Handler) ListSubscriptions(c *gin.Context) {
	var rows []struct {
		ID            string  `json:"id"`
		TenantID      string  `json:"tenant_id"`
		TenantName    string  `json:"tenant_name"`
		PlanName      string  `json:"plan_name"`
		Status        string  `json:"status"`
		EndDate       *string `json:"end_date"`
		DaysRemaining *int    `json:"days_remaining"`
	}
	err := h.db.Raw(`
		SELECT s.id, s.tenant_id, t.name AS tenant_name, p.name AS plan_name, s.status,
		       TO_CHAR(s.end_date, 'YYYY-MM-DD') AS end_date,
		       CASE WHEN s.end_date IS NOT NULL
		            THEN GREATEST(0, EXTRACT(DAY FROM (s.end_date::timestamp - NOW())))::int
		            ELSE NULL END AS days_remaining
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

// ─── Plan CRUD ──────────────────────────────────────────────────────────────

type planRow struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	MaxUsers    int     `json:"max_users"`
	MaxProducts int     `json:"max_products"`
	Description string  `json:"description"`
	IsActive    bool    `json:"is_active"`
	TenantCount int     `json:"tenant_count"`
}

func (h *Handler) ListPlans(c *gin.Context) {
	var rows []planRow
	err := h.db.Raw(`
		SELECT p.id, p.name, p.price, p.max_users, p.max_products,
		       COALESCE(p.description, '') AS description,
		       COALESCE(p.is_active, true) AS is_active,
		       (SELECT COUNT(*) FROM subscriptions s WHERE s.plan_id = p.id AND s.status IN ('active','trial')) AS tenant_count
		FROM plans p
		ORDER BY p.price ASC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list plans"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

type CreatePlanInput struct {
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price"`
	MaxUsers    int     `json:"max_users" binding:"required"`
	MaxProducts int     `json:"max_products" binding:"required"`
	Description string  `json:"description"`
}

func (h *Handler) CreatePlan(c *gin.Context) {
	var in CreatePlanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id int
	err := h.db.Raw(`INSERT INTO plans (name, price, max_users, max_products, description, is_active) VALUES (?, ?, ?, ?, ?, true) RETURNING id`,
		in.Name, in.Price, in.MaxUsers, in.MaxProducts, in.Description).Scan(&id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create plan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "plan created", "id": id})
}

type UpdatePlanInput struct {
	Name        *string  `json:"name"`
	Price       *float64 `json:"price"`
	MaxUsers    *int     `json:"max_users"`
	MaxProducts *int     `json:"max_products"`
	Description *string  `json:"description"`
	IsActive    *bool    `json:"is_active"`
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	var in UpdatePlanInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if in.Name != nil {
		updates["name"] = *in.Name
	}
	if in.Price != nil {
		updates["price"] = *in.Price
	}
	if in.MaxUsers != nil {
		updates["max_users"] = *in.MaxUsers
	}
	if in.MaxProducts != nil {
		updates["max_products"] = *in.MaxProducts
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.IsActive != nil {
		updates["is_active"] = *in.IsActive
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	res := h.db.Table("plans").Where("id = ?", id).Updates(updates)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeletePlan(c *gin.Context) {
	id := c.Param("id")
	res := h.db.Table("plans").Where("id = ?", id).Update("is_active", false)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "plan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "plan deactivated"})
}

// ─── Dashboard Stats ────────────────────────────────────────────────────────

func (h *Handler) Stats(c *gin.Context) {
	var totalTenants int64
	h.db.Table("tenants").Where("status != 'deleted'").Count(&totalTenants)

	var activeTenants int64
	h.db.Table("tenants").Where("status = 'active'").Count(&activeTenants)

	var totalRevenue float64
	h.db.Raw("SELECT COALESCE(SUM(amount), 0) FROM subscription_orders WHERE status = 'approved'").Scan(&totalRevenue)

	var pendingOrders int64
	h.db.Table("subscription_orders").Where("status IN ('pending_payment','pending_review')").Count(&pendingOrders)

	var expiringIn7 int64
	h.db.Raw(`SELECT COUNT(*) FROM subscriptions WHERE status IN ('active','trial') AND end_date BETWEEN NOW() AND NOW() + INTERVAL '7 days'`).Scan(&expiringIn7)

	c.JSON(http.StatusOK, gin.H{
		"total_tenants":    totalTenants,
		"active_tenants":   activeTenants,
		"total_revenue":    totalRevenue,
		"pending_orders":   pendingOrders,
		"expiring_in_7d":   expiringIn7,
	})
}

// ─── Subscription Order Management ──────────────────────────────────────────

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

// ─── Revenue Report ─────────────────────────────────────────────────────────

type revenueRow struct {
	ID          string  `json:"id"`
	TenantName  string  `json:"tenant_name"`
	PlanName    string  `json:"plan_name"`
	Amount      float64 `json:"amount"`
	ApprovedAt  string  `json:"approved_at"`
	Invoice     string  `json:"invoice"`
}

func (h *Handler) RevenueReport(c *gin.Context) {
	from := c.DefaultQuery("from", "")
	to := c.DefaultQuery("to", "")

	query := `
		SELECT o.id, t.name AS tenant_name, p.name AS plan_name, o.amount,
		       TO_CHAR(o.reviewed_at, 'YYYY-MM-DD HH24:MI') AS approved_at,
		       o.invoice_number AS invoice
		FROM subscription_orders o
		JOIN tenants t ON t.id = o.tenant_id
		JOIN plans p ON p.id = o.plan_id
		WHERE o.status = 'approved'
	`
	args := []interface{}{}
	if from != "" {
		query += " AND o.reviewed_at >= ?"
		args = append(args, from)
	}
	if to != "" {
		query += " AND o.reviewed_at < ?::date + INTERVAL '1 day'"
		args = append(args, to)
	}
	query += " ORDER BY o.reviewed_at DESC"

	var rows []revenueRow
	if err := h.db.Raw(query, args...).Scan(&rows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get revenue report"})
		return
	}

	var totalRevenue float64
	var planBreakdown []struct {
		PlanName   string  `json:"plan_name"`
		OrderCount int     `json:"order_count"`
		Total      float64 `json:"total"`
	}

	summaryQuery := `
		SELECT p.name AS plan_name, COUNT(*) AS order_count, SUM(o.amount) AS total
		FROM subscription_orders o
		JOIN plans p ON p.id = o.plan_id
		WHERE o.status = 'approved'
	`
	summaryArgs := []interface{}{}
	if from != "" {
		summaryQuery += " AND o.reviewed_at >= ?"
		summaryArgs = append(summaryArgs, from)
	}
	if to != "" {
		summaryQuery += " AND o.reviewed_at < ?::date + INTERVAL '1 day'"
		summaryArgs = append(summaryArgs, to)
	}
	summaryQuery += " GROUP BY p.name ORDER BY total DESC"

	h.db.Raw(summaryQuery, summaryArgs...).Scan(&planBreakdown)
	for _, pb := range planBreakdown {
		totalRevenue += pb.Total
	}

	avgPerOrder := 0.0
	if len(rows) > 0 {
		avgPerOrder = math.Round(totalRevenue/float64(len(rows))*100) / 100
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":         rows,
		"total_revenue":  totalRevenue,
		"total_orders":   len(rows),
		"avg_per_order":  fmt.Sprintf("%.0f", avgPerOrder),
		"plan_breakdown": planBreakdown,
	})
}
