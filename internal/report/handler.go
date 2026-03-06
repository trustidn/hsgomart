package report

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

// parseDateRange reads from and to query params (YYYY-MM-DD). Defaults to start/end of current month.
func parseDateRange(c *gin.Context) (from, to time.Time, err error) {
	now := time.Now()
	fromStr := c.DefaultQuery("from", "")
	toStr := c.DefaultQuery("to", "")

	if fromStr != "" && toStr != "" {
		from, err = time.Parse("2006-01-02", fromStr)
		if err != nil {
			return from, to, err
		}
		to, err = time.Parse("2006-01-02", toStr)
		if err != nil {
			return from, to, err
		}
		// Include full to date (end of day)
		to = to.Add(24*time.Hour - time.Second)
		return from, to, nil
	}

	// Default: current month
	from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	to = now
	return from, to, nil
}

// SalesSummary handles GET /api/reports/sales?from=2026-01-01&to=2026-01-31
func (h *Handler) SalesSummary(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD for from and to"})
		return
	}

	result, err := h.service.SalesSummary(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sales summary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_transactions": result.TotalTransactions,
		"total_sales":        result.TotalSales,
	})
}

// TopProductsResponse for API (product_name, quantity_sold, revenue per spec).
type TopProductsResponse struct {
	ProductName  string  `json:"product_name"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
}

// TopProducts handles GET /api/reports/products?from=2026-01-01&to=2026-01-31
func (h *Handler) TopProducts(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD for from and to"})
		return
	}

	rows, err := h.service.TopProducts(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get top products"})
		return
	}

	res := make([]TopProductsResponse, 0, len(rows))
	for _, r := range rows {
		res = append(res, TopProductsResponse{
			ProductName:  r.ProductName,
			QuantitySold: r.QuantitySold,
			Revenue:      r.Revenue,
		})
	}
	c.JSON(http.StatusOK, res)
}

// InventorySummaryResponse for API (product_id for adjust-stock, product_name, stock).
type InventorySummaryResponse struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
}

// InventorySummary handles GET /api/reports/inventory
func (h *Handler) InventorySummary(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	rows, err := h.service.InventorySummary(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get inventory summary"})
		return
	}

	res := make([]InventorySummaryResponse, 0, len(rows))
	for _, r := range rows {
		res = append(res, InventorySummaryResponse{
			ProductID:   r.ProductID,
			ProductName: r.ProductName,
			Stock:       r.Stock,
		})
	}
	c.JSON(http.StatusOK, res)
}
