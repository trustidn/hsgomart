package refund

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRefund(c *gin.Context) {
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

	var in RefundInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.CreateRefund(tenantID, userID, in)
	if err != nil {
		switch err {
		case ErrTransactionNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrAlreadyRefunded:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case ErrRefundExceedsTotal:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "refund failed"})
		}
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *Handler) ListRefunds(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	list, err := h.service.ListRefunds(tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list refunds"})
		return
	}
	c.JSON(http.StatusOK, list)
}
