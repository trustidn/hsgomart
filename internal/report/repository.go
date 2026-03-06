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
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
}

// GetInventorySummary returns product id, name and current stock for the tenant.
func GetInventorySummary(db *gorm.DB, tenantID string) ([]InventorySummaryRow, error) {
	var rows []InventorySummaryRow
	err := db.Table("inventories").
		Select("inventories.product_id as product_id, products.name as product_name, inventories.stock as stock").
		Joins("INNER JOIN products ON products.id = inventories.product_id").
		Where("inventories.tenant_id = ?", tenantID).
		Order("products.name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	return rows, nil
}
