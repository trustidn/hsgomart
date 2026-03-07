package report

import (
	"net/http"
	"strconv"
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

// SalesDaily handles GET /api/reports/sales/daily?from=&to=
func (h *Handler) SalesDaily(c *gin.Context) {
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
	rows, err := h.service.SalesDaily(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sales daily"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// SalesTransactions handles GET /api/reports/sales/transactions?from=&to=&page=1&limit=20 (limit=0 = all)
func (h *Handler) SalesTransactions(c *gin.Context) {
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
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	limit := 20
	if q := c.Query("limit"); q != "" {
		if v, _ := strconv.Atoi(q); v >= 0 {
			limit = v
		}
	}
	offset := 0
	if limit > 0 {
		offset = (page - 1) * limit
	}
	total, err := h.service.CountSalesTransactions(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count transactions"})
		return
	}
	rows, err := h.service.SalesTransactions(tenantID, from, to, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sales transactions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
		"rows":  rows,
	})
}

// SalesHourly handles GET /api/reports/sales/hourly?date=YYYY-MM-DD
func (h *Handler) SalesHourly(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date required (YYYY-MM-DD)"})
		return
	}
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}
	rows, err := h.service.SalesHourly(tenantID, date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get hourly sales"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// PaymentsReport handles GET /api/reports/payments?from=&to=
func (h *Handler) PaymentsReport(c *gin.Context) {
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
	rows, err := h.service.PaymentsReport(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payments report"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

// ProfitReport handles GET /api/reports/profit?from=&to=
func (h *Handler) ProfitReport(c *gin.Context) {
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
	summary, rows, err := h.service.ProfitReport(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profit report"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"rows":    rows,
	})
}

// CashiersReport handles GET /api/reports/cashiers?from=&to=
func (h *Handler) CashiersReport(c *gin.Context) {
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
	rows, err := h.service.CashiersReport(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cashiers report"})
		return
	}
	c.JSON(http.StatusOK, rows)
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

// InventorySummaryResponse for API (product_id, product_name, stock, cost_price, inventory_value from batches).
type InventorySummaryResponse struct {
	ProductID      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	Stock          int     `json:"stock"`
	CostPrice      float64 `json:"cost_price"`
	InventoryValue float64 `json:"inventory_value"`
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
			ProductID:      r.ProductID,
			ProductName:    r.ProductName,
			Stock:          r.Stock,
			CostPrice:      r.CostPrice,
			InventoryValue: r.InventoryValue,
		})
	}
	c.JSON(http.StatusOK, res)
}

// ShiftsReport handles GET /api/reports/shifts?from=&to=
func (h *Handler) ShiftsReport(c *gin.Context) {
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
	rows, err := h.service.ShiftsReport(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get shifts report"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) SalesCompare(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	curFrom, _ := time.Parse("2006-01-02", c.Query("current_from"))
	curTo, _ := time.Parse("2006-01-02", c.Query("current_to"))
	prevFrom, _ := time.Parse("2006-01-02", c.Query("previous_from"))
	prevTo, _ := time.Parse("2006-01-02", c.Query("previous_to"))

	if curFrom.IsZero() || curTo.IsZero() {
		now := time.Now()
		curFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		curTo = now
		prevFrom = curFrom.AddDate(0, -1, 0)
		prevTo = curFrom.Add(-time.Second)
	}
	if prevFrom.IsZero() {
		prevFrom = curFrom.AddDate(0, -1, 0)
		prevTo = curFrom.Add(-time.Second)
	}
	curTo = curTo.Add(24*time.Hour - time.Second)
	prevTo = prevTo.Add(24*time.Hour - time.Second)

	result, err := h.service.SalesCompare(tenantID, curFrom, curTo, prevFrom, prevTo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compare sales"})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) ProductMargin(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	from, to, err := parseDateRange(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}
	rows, err := h.service.ProductMargin(tenantID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product margin"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) GetReceipt(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	txnID := c.Param("id")

	type receiptItem struct {
		ProductName string  `json:"product_name"`
		Quantity    int     `json:"quantity"`
		Price       float64 `json:"price"`
		Subtotal    float64 `json:"subtotal"`
	}

	var txn struct {
		ID          string  `json:"id"`
		TotalAmount float64 `json:"total_amount"`
		CreatedAt   string  `json:"created_at"`
		Cashier     string  `json:"cashier"`
	}
	err := h.service.DB().Raw(`
		SELECT t.id, t.total_amount, TO_CHAR(t.created_at, 'YYYY-MM-DD HH24:MI') AS created_at,
		       COALESCE(NULLIF(TRIM(u.name), ''), u.email) AS cashier
		FROM transactions t LEFT JOIN users u ON u.id = t.user_id
		WHERE t.id = ? AND t.tenant_id = ?
	`, txnID, tenantID).Scan(&txn).Error
	if err != nil || txn.ID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	var items []receiptItem
	h.service.DB().Raw(`
		SELECT p.name AS product_name, ti.quantity, ti.price, ti.subtotal
		FROM transaction_items ti JOIN products p ON p.id = ti.product_id
		WHERE ti.transaction_id = ?
	`, txnID).Scan(&items)

	var payments []struct {
		Method string  `json:"method"`
		Amount float64 `json:"amount"`
	}
	h.service.DB().Raw(`SELECT method, amount FROM payments WHERE transaction_id = ?`, txnID).Scan(&payments)

	c.JSON(http.StatusOK, gin.H{
		"transaction": txn,
		"items":       items,
		"payments":    payments,
	})
}
