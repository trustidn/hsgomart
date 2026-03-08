package inventory

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

// MovementResponse for GET /api/inventory/movements.
type MovementResponse struct {
	ProductName string `json:"product_name"`
	Type        string `json:"type"`
	Quantity    int    `json:"quantity"`
	StockAfter  int    `json:"stock_after"` // running balance after this movement
	Reference   string `json:"reference"`
	Reason      string `json:"reason"`
	CreatedAt   string `json:"created_at"`
}

// ListMovements handles GET /api/inventory/movements (optional: product_id, type, from_date, to_date, limit, page).
func (h *Handler) ListMovements(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	productID := c.Query("product_id")
	movementType := c.Query("type")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	rows, total, err := h.service.ListMovementRowsPaginated(tenantID, productID, movementType, fromDate, toDate, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list movements"})
		return
	}
	res := make([]MovementResponse, 0, len(rows))
	for _, r := range rows {
		createdAt := ""
		if !r.CreatedAt.IsZero() {
			createdAt = r.CreatedAt.Format("2006-01-02T15:04:05")
		}
		res = append(res, MovementResponse{
			ProductName: r.ProductName,
			Type:        r.Type,
			Quantity:    r.Quantity,
			StockAfter:  r.StockAfter,
			Reference:   r.Reference,
			Reason:      r.Reason,
			CreatedAt:   createdAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"movements": res, "total": total})
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

// LowStock returns products with stock at or below threshold (GET /api/inventory/low-stock).
func (h *Handler) LowStock(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	list, err := h.service.GetLowStockProducts(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get low stock list"})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) Expiring(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days <= 0 {
		days = 30
	}
	list, err := h.service.GetExpiringProducts(tenantID, days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get expiring products"})
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
	Reason    string `json:"reason"` // adjustment reason for audit (e.g. expired product, damaged item)
}

// AdjustStock applies a stock reduction only (POST /api/products/:id/adjust-stock). Quantity must be <= 0.
// Stock increase is only via Purchase.
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

	err := h.service.AdjustStock(tenantID, productID, in.Quantity, in.Type, in.Reference, in.Reason)
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err == ErrInsufficientStock {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err == ErrAdjustOnlyDecrease {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to adjust stock"})
		return
	}

	stock, _ := h.service.GetStock(tenantID, productID)
	c.JSON(http.StatusOK, gin.H{"product_id": productID, "stock": stock})
}
