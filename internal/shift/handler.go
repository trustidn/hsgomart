package shift

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// OpenShift handles POST /api/shifts/open (cashier opens a shift).
func (h *Handler) OpenShift(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok || strings.TrimSpace(tenantID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	userID, ok := utils.GetUserID(c)
	if !ok || strings.TrimSpace(userID) == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context required"})
		return
	}
	var in OpenShiftInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.OpenShift(tenantID, userID, in)
	if err != nil {
		if err == ErrShiftAlreadyOpen {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Printf("[shift] OpenShift error: %v", err)
		msg := err.Error()
		if len(msg) > 200 {
			msg = msg[:200] + "..."
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open shift: " + msg})
		return
	}
	c.JSON(http.StatusOK, result)
}

// CloseShift handles POST /api/shifts/close (cashier closes current shift). Body may include shift_id if multiple; we close by current.
func (h *Handler) CloseShift(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context required"})
		return
	}
	var in CloseShiftInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Close current shift for this user
	active, err := h.svc.GetCurrentShift(tenantID, userID)
	if err != nil || active == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no active shift to close"})
		return
	}
	result, err := h.svc.CloseShift(tenantID, userID, active.ID, in)
	if err != nil {
		if err == ErrShiftNotFound || err == ErrShiftAlreadyClosed {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to close shift"})
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetCurrentShift handles GET /api/shifts/current.
func (h *Handler) GetCurrentShift(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user context required"})
		return
	}
	sh, err := h.svc.GetCurrentShift(tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get shift"})
		return
	}
	if sh == nil {
		c.JSON(http.StatusOK, gin.H{"shift": nil})
		return
	}
	openedAt := sh.OpenedAt.Format("2006-01-02T15:04:05Z07:00")
	c.JSON(http.StatusOK, gin.H{
		"shift": gin.H{
			"id":           sh.ID,
			"opening_cash": sh.OpeningCash,
			"opened_at":    openedAt,
			"status":       sh.Status,
		},
	})
}

// ListShifts handles GET /api/shifts (owner: list all shifts for tenant).
func (h *Handler) ListShifts(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if offset < 0 {
		offset = 0
	}
	list, err := h.svc.ListShifts(tenantID, limit, offset)
	if err != nil {
		log.Printf("[shift] ListShifts error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list shifts"})
		return
	}
	log.Printf("[shift] ListShifts tenant_id=%s count=%d", tenantID, len(list))
	rows := make([]gin.H, 0, len(list))
	for _, s := range list {
		openedAt := s.OpenedAt.Format("2006-01-02T15:04:05Z07:00")
		row := gin.H{
			"id":           s.ID,
			"user_id":      s.UserID,
			"opening_cash": s.OpeningCash,
			"opened_at":    openedAt,
			"status":       s.Status,
		}
		if s.ClosingCash != nil {
			row["closing_cash"] = *s.ClosingCash
		}
		if s.ClosedAt != nil {
			row["closed_at"] = s.ClosedAt.Format("2006-01-02T15:04:05Z07:00")
		}
		rows = append(rows, row)
	}
	c.JSON(http.StatusOK, gin.H{"shifts": rows})
}
