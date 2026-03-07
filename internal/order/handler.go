package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type CreateOrderInput struct {
	PlanID int `json:"plan_id" binding:"required"`
}

func (h *Handler) CreateOrder(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	var in CreateOrderInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.CreateOrder(tenantID, in.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *Handler) UploadPaymentProof(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	orderID := c.Param("id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id required"})
		return
	}
	file, err := c.FormFile("payment_proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "payment_proof file required"})
		return
	}
	order, err := h.svc.UploadPaymentProof(tenantID, orderID, file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *Handler) ListOrders(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	orders, err := h.svc.ListOrders(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list orders"})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *Handler) GetOrder(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	orderID := c.Param("id")
	order, err := h.svc.GetOrder(tenantID, orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}
