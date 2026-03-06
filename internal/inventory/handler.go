package inventory

import (
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

// MovementResponse for GET /api/inventory/movements.
type MovementResponse struct {
	ProductName string `json:"product_name"`
	Type        string `json:"type"`
	Quantity    int    `json:"quantity"`
	Reference   string `json:"reference"`
	CreatedAt   string `json:"created_at"`
}

// ListMovements handles GET /api/inventory/movements (optional query: product_id).
func (h *Handler) ListMovements(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	productID := c.Query("product_id")
	rows, err := h.service.ListMovementRows(tenantID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list movements"})
		return
	}
	res := make([]MovementResponse, 0, len(rows))
	for _, r := range rows {
		createdAt := ""
		if !r.CreatedAt.IsZero() {
			createdAt = r.CreatedAt.Format(time.RFC3339)
		}
		res = append(res, MovementResponse{
			ProductName: r.ProductName,
			Type:        r.Type,
			Quantity:    r.Quantity,
			Reference:   r.Reference,
			CreatedAt:   createdAt,
		})
	}
	c.JSON(http.StatusOK, res)
}

// List returns all inventory rows for the tenant (product_id and stock per product).
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	list, err := h.service.ListInventoryByTenant(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list inventory"})
		return
	}

	c.JSON(http.StatusOK, list)
}

// GetStock returns current stock for a product (GET /api/products/:id/stock).
func (h *Handler) GetStock(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id required"})
		return
	}

	stock, err := h.service.GetStock(tenantID, productID)
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"product_id": productID, "stock": stock})
}

type AdjustStockInput struct {
	Quantity  int    `json:"quantity" binding:"required"`
	Type      string `json:"type" binding:"required"` // adjustment, purchase, sale, return
	Reference string `json:"reference"`
}

// AdjustStock applies a stock adjustment (POST /api/products/:id/adjust-stock).
func (h *Handler) AdjustStock(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id required"})
		return
	}

	var in AdjustStockInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.AdjustStock(tenantID, productID, in.Quantity, in.Type, in.Reference)
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err == ErrInsufficientStock {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to adjust stock"})
		return
	}

	stock, _ := h.service.GetStock(tenantID, productID)
	c.JSON(http.StatusOK, gin.H{"product_id": productID, "stock": stock})
}
