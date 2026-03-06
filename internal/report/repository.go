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
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Stock       int     `json:"stock"`
	CostPrice   float64 `json:"cost_price"`
}

// GetInventorySummary returns all products for the tenant with stock (0 if no inventory row yet).
func GetInventorySummary(db *gorm.DB, tenantID string) ([]InventorySummaryRow, error) {
	var rows []InventorySummaryRow
	err := db.Table("products").
		Select("products.id as product_id, products.name as product_name, COALESCE(inventories.stock, 0) as stock, COALESCE(products.cost_price, 0) as cost_price").
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

// GetSalesTransactions returns each transaction in the date range for the tenant.
func GetSalesTransactions(db *gorm.DB, tenantID string, fromDate, toDate time.Time) ([]SalesTransactionRow, error) {
	var rows []SalesTransactionRow
	err := db.Table("transactions").
		Select("transactions.id as id, to_char(transactions.created_at, 'YYYY-MM-DD HH24:MI') as created_at, transactions.total_amount as total_amount, COALESCE(users.email, users.name, '') as cashier").
		Joins("LEFT JOIN users ON users.id = transactions.user_id AND users.tenant_id = transactions.tenant_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Order("transactions.created_at DESC").
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
func GetProfitReport(db *gorm.DB, tenantID string, fromDate, toDate time.Time) (summary ProfitSummary, rows []ProfitRow, err error) {
	err = db.Table("transaction_items").
		Select("products.name as product_name, SUM(transaction_items.quantity) as quantity_sold, SUM(transaction_items.subtotal) as revenue, SUM(products.cost_price * transaction_items.quantity) as cost").
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
		Select("COALESCE(users.email, users.name, 'Unknown') as cashier, COUNT(*) as transactions, COALESCE(SUM(transactions.total_amount), 0) as revenue").
		Joins("LEFT JOIN users ON users.id = transactions.user_id AND users.tenant_id = transactions.tenant_id").
		Where("transactions.tenant_id = ? AND transactions.status = ?", tenantID, "completed").
		Where("transactions.created_at >= ? AND transactions.created_at <= ?", fromDate, toDate).
		Group("transactions.user_id, users.email, users.name").
		Order("revenue DESC").
		Scan(&rows).Error
	return rows, err
}
