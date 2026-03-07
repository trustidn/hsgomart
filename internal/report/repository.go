package report

import (
	"time"

	"gorm.io/gorm"
)

// SalesSummaryResult holds sales report aggregates.
type SalesSummaryResult struct {
	TotalTransactions int     `json:"total_transactions"`
	TotalSales        float64 `json:"total_sales"`
}

// GetSalesSummary returns total transaction count and total sales amount for the tenant in the date range.
func GetSalesSummary(db *gorm.DB, tenantID string, fromDate, toDate time.Time) (*SalesSummaryResult, error) {
	var result SalesSummaryResult
	err := db.Table("transactions").
		Select("COUNT(*) as total_transactions, COALESCE(SUM(total_amount), 0) as total_sales").
		Where("tenant_id = ? AND status = ?", tenantID, "completed").
		Where("created_at >= ? AND created_at <= ?", fromDate, toDate).
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// TopProductRow holds one row of top products by sales.
type TopProductRow struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	QuantitySold int    `json:"quantity_sold"`
	Revenue     float64 `json:"revenue"`
}

// GetTopProducts returns products with quantity sold and revenue in the date range (tenant-scoped).
func GetTopProducts(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]TopProductRow, error) {
	var rows []TopProductRow
	err := db.Table("transaction_items").
		Select("products.id as product_id, products.name as product_name, SUM(transaction_items.quantity) as quantity_sold, SUM(transaction_items.subtotal) as revenue").
		Joins("INNER JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("INNER JOIN products ON products.id = transaction_items.product_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Group("products.id, products.name").
		Order("revenue DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// InventorySummaryRow holds one row of inventory summary.
type InventorySummaryRow struct {
	ProductID      string  `json:"product_id"`
	ProductName    string  `json:"product_name"`
	Stock          int     `json:"stock"`
	CostPrice      float64 `json:"cost_price"`       // product reference cost
	InventoryValue float64 `json:"inventory_value"`  // SUM(remaining_quantity * cost_price) from batches
}

// GetInventorySummary returns all products for the tenant with stock and inventory value from batches.
func GetInventorySummary(db *gorm.DB, tenantID string) ([]InventorySummaryRow, error) {
	var rows []InventorySummaryRow
	err := db.Table("products").
		Select(`products.id as product_id, products.name as product_name,
			COALESCE(inventories.stock, 0) as stock,
			COALESCE(products.cost_price, 0) as cost_price,
			COALESCE((SELECT SUM(ib.remaining_quantity * ib.cost_price) FROM inventory_batches ib WHERE ib.product_id = products.id), 0) as inventory_value`).
		Joins("LEFT JOIN inventories ON inventories.product_id = products.id AND inventories.tenant_id = products.tenant_id").
		Where("products.tenant_id = ?", tenantID).
		Order("products.name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// SalesDailyRow holds one day of sales.
type SalesDailyRow struct {
	Date                string  `json:"date"`
	TotalTransactions   int     `json:"total_transactions"`
	TotalSales          float64 `json:"total_sales"`
}

// GetSalesDaily returns daily sales breakdown for the tenant in the date range.
func GetSalesDaily(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]SalesDailyRow, error) {
	var rows []SalesDailyRow
	err := db.Table("transactions").
		Select("to_char(created_at, 'YYYY-MM-DD') as date, COUNT(*) as total_transactions, COALESCE(SUM(total_amount), 0) as total_sales").
		Where("tenant_id = ? AND status = ?", tenantID, "completed").
		Where("created_at >= ? AND created_at <= ?", fromDate, toDate).
		Group("to_char(created_at, 'YYYY-MM-DD')").
		Order("date DESC").
		Scan(&rows).Error
	return rows, err
}

// SalesTransactionRow holds one transaction for detail list.
type SalesTransactionRow struct {
	ID          string  `json:"id"`
	CreatedAt   string  `json:"created_at"`
	TotalAmount float64 `json:"total_amount"`
	Cashier     string  `json:"cashier"`
}

// GetSalesTransactions returns paginated transactions. Use limit=0 for no limit (export).
func GetSalesTransactions(db *gorm.DB, tenantID string, fromDate, toDate time.Time, limit, offset int) ([]SalesTransactionRow, error) {
	var rows []SalesTransactionRow
	q := db.Table("transactions").
		Select("transactions.id as id, to_char(transactions.created_at, 'YYYY-MM-DD HH24:MI') as created_at, transactions.total_amount as total_amount, COALESCE(NULLIF(TRIM(users.name), ''), users.email, '') as cashier").
		Joins("LEFT JOIN users ON users.id = transactions.user_id AND users.tenant_id = transactions.tenant_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Order("transactions.created_at DESC")
	if limit > 0 {
		q = q.Limit(limit).Offset(offset)
	}
	err := q.Scan(&rows).Error
	return rows, err
}

// CountSalesTransactions returns total count of transactions in the date range.
func CountSalesTransactions(db *gorm.DB, tenantID string, fromDate, toDate time.Time) (int64, error) {
	var count int64
	err := db.Table("transactions").
		Where("tenant_id = ? AND status = ?", tenantID, "completed").
		Where("created_at >= ? AND created_at <= ?", fromDate, toDate).
		Count(&count).Error
	return count, err
}

// PaymentRow holds payment method aggregates.
type PaymentRow struct {
	Method       string  `json:"method"`
	Transactions int     `json:"transactions"`
	Revenue      float64 `json:"revenue"`
}

// GetPaymentsReport returns per payment-method count and revenue for the tenant in the date range.
func GetPaymentsReport(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]PaymentRow, error) {
	var rows []PaymentRow
	err := db.Table("payments").
		Select("payments.method as method, COUNT(DISTINCT payments.transaction_id) as transactions, COALESCE(SUM(payments.amount), 0) as revenue").
		Joins("INNER JOIN transactions ON transactions.id = payments.transaction_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Group("payments.method").
		Order("revenue DESC").
		Scan(&rows).Error
	return rows, err
}

// SalesHourlyRow holds one hour of sales for a given date.
type SalesHourlyRow struct {
	Hour       int     `json:"hour"`
	Transactions int   `json:"transactions"`
	Revenue    float64 `json:"revenue"`
}

// GetSalesHourly returns hourly breakdown for the tenant on the given date (YYYY-MM-DD).
func GetSalesHourly(db *gorm.DB, tenantID string, date time.Time) ([]SalesHourlyRow, error) {
	start := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	end := start.Add(24*time.Hour - time.Second)
	var rows []SalesHourlyRow
	err := db.Table("transactions").
		Select("EXTRACT(HOUR FROM created_at)::int as hour, COUNT(*) as transactions, COALESCE(SUM(total_amount), 0) as revenue").
		Where("tenant_id = ? AND status = ?", tenantID, "completed").
		Where("created_at >= ? AND created_at <= ?", start, end).
		Group("EXTRACT(HOUR FROM created_at)").
		Order("hour").
		Scan(&rows).Error
	return rows, err
}

// ProfitRow holds product-level profit (revenue - cost).
type ProfitRow struct {
	ProductName  string  `json:"product_name"`
	QuantitySold int     `json:"quantity_sold"`
	Revenue      float64 `json:"revenue"`
	Cost         float64 `json:"cost"`
	Profit       float64 `json:"profit"`
}

// ProfitSummary holds totals for profit report.
type ProfitSummary struct {
	Revenue float64 `json:"revenue"`
	Cost    float64 `json:"cost"`
	Profit  float64 `json:"profit"`
}

// GetProfitReport returns product-level profit and summary for the tenant in the date range.
// Cost uses transaction_items.cogs (FIFO batch cost) when present; falls back to product cost for legacy rows.
func GetProfitReport(db *gorm.DB, tenantID string, fromDate, toDate time.Time) (summary ProfitSummary, rows []ProfitRow, err error) {
	err = db.Table("transaction_items").
		Select("products.name as product_name, SUM(transaction_items.quantity) as quantity_sold, SUM(transaction_items.subtotal) as revenue, SUM(COALESCE(NULLIF(transaction_items.cogs, 0), products.cost_price * transaction_items.quantity)) as cost").
		Joins("INNER JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("INNER JOIN products ON products.id = transaction_items.product_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Group("products.id, products.name").
		Order("revenue DESC").
		Scan(&rows).Error
	if err != nil {
		return summary, nil, err
	}
	for i := range rows {
		rows[i].Profit = rows[i].Revenue - rows[i].Cost
		summary.Revenue += rows[i].Revenue
		summary.Cost += rows[i].Cost
		summary.Profit += rows[i].Profit
	}
	return summary, rows, nil
}

// CashierRow holds cashier performance.
type CashierRow struct {
	Cashier      string  `json:"cashier"`
	Transactions int     `json:"transactions"`
	Revenue      float64 `json:"revenue"`
}

// GetCashiersReport returns per-cashier transaction count and revenue for the tenant in the date range.
func GetCashiersReport(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]CashierRow, error) {
	var rows []CashierRow
	err := db.Table("transactions").
		Select("COALESCE(NULLIF(TRIM(users.name), ''), users.email, 'Unknown') as cashier, COUNT(*) as transactions, COALESCE(SUM(transactions.total_amount), 0) as revenue").
		Joins("LEFT JOIN users ON users.id = transactions.user_id AND users.tenant_id = transactions.tenant_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Group("transactions.user_id, users.name, users.email").
		Order("revenue DESC").
		Scan(&rows).Error
	return rows, err
}

// ShiftReportRow holds one row for the shifts report (cash reconciliation).
type ShiftReportRow struct {
	Date       string  `json:"date"`
	Cashier    string  `json:"cashier"`
	Opening    float64 `json:"opening"`
	Sales      float64 `json:"sales"`
	Expected   float64 `json:"expected"`
	Actual     float64 `json:"actual"`
	Difference float64 `json:"difference"`
}

// GetShiftsReport returns shift reconciliation rows for the tenant (closed shifts in date range by closed_at).
func GetShiftsReport(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]ShiftReportRow, error) {
	type shiftRow struct {
		ID          string
		UserID      string
		OpeningCash float64
		ClosingCash *float64
		OpenedAt    time.Time
		ClosedAt    *time.Time
		CashierName string
	}
	var list []shiftRow
	err := db.Table("cashier_shifts").
		Select("cashier_shifts.id, cashier_shifts.user_id, cashier_shifts.opening_cash, cashier_shifts.closing_cash, cashier_shifts.opened_at, cashier_shifts.closed_at, COALESCE(NULLIF(TRIM(users.name), ''), users.email, '') as cashier_name").
		Joins("LEFT JOIN users ON users.id = cashier_shifts.user_id AND users.tenant_id = cashier_shifts.tenant_id").
		Where("cashier_shifts.tenant_id = ? AND cashier_shifts.status = ?", tenantID, "closed").
		Where("cashier_shifts.closed_at >= ? AND cashier_shifts.closed_at <= ?", fromDate, toDate).
		Order("cashier_shifts.closed_at DESC").
		Scan(&list).Error
	if err != nil {
		return nil, err
	}
	result := make([]ShiftReportRow, 0, len(list))
	for _, s := range list {
		if s.ClosedAt == nil {
			continue
		}
		var cashSales float64
		_ = db.Table("payments").
			Select("COALESCE(SUM(payments.amount), 0)").
			Joins("INNER JOIN transactions ON transactions.id = payments.transaction_id").
			Where("transactions.tenant_id = ? AND payments.method = ?", tenantID, "cash").
			Where("payments.created_at >= ? AND payments.created_at <= ?", s.OpenedAt, *s.ClosedAt).
			Scan(&cashSales).Error
		expected := s.OpeningCash + cashSales
		actual := 0.0
		if s.ClosingCash != nil {
			actual = *s.ClosingCash
		}
		dateStr := s.ClosedAt.Format("2006-01-02 15:04")
		result = append(result, ShiftReportRow{
			Date:       dateStr,
			Cashier:    s.CashierName,
			Opening:    s.OpeningCash,
			Sales:      cashSales,
			Expected:   expected,
			Actual:     actual,
			Difference: actual - expected,
		})
	}
	return result, nil
}

type MarginRow struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Revenue     float64 `json:"revenue"`
	COGS        float64 `json:"cogs"`
	Margin      float64 `json:"margin"`
	MarginPct   float64 `json:"margin_pct"`
}

func GetProductMargin(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]MarginRow, error) {
	var rows []MarginRow
	err := db.Raw(`
		SELECT ti.product_id,
		       p.name AS product_name,
		       SUM(ti.subtotal) AS revenue,
		       SUM(ti.cogs) AS cogs,
		       SUM(ti.subtotal) - SUM(ti.cogs) AS margin,
		       CASE WHEN SUM(ti.subtotal) > 0
		            THEN ROUND(((SUM(ti.subtotal) - SUM(ti.cogs)) / SUM(ti.subtotal) * 100)::numeric, 2)
		            ELSE 0
		       END AS margin_pct
		FROM transaction_items ti
		JOIN transactions t ON t.id = ti.transaction_id
		JOIN products p ON p.id = ti.product_id
		WHERE t.tenant_id = ? AND t.status = 'completed'
		  AND t.created_at >= ? AND t.created_at <= ?
		GROUP BY ti.product_id, p.name
		ORDER BY margin DESC
	`, tenantID, fromDate, toDate).Scan(&rows).Error
	return rows, err
}
