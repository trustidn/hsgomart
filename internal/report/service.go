package report

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// SalesSummary returns sales summary for the tenant in the date range.
func (s *Service) SalesSummary(tenantID string, fromDate, toDate time.Time) (*SalesSummaryResult, error) {
	return GetSalesSummary(s.db, tenantID, fromDate, toDate)
}

// TopProducts returns top products by revenue for the tenant in the date range.
func (s *Service) TopProducts(tenantID string, fromDate, toDate time.Time) ([]TopProductRow, error) {
	return GetTopProducts(s.db, tenantID, fromDate, toDate)
}

// InventorySummary returns current inventory summary for the tenant.
func (s *Service) InventorySummary(tenantID string) ([]InventorySummaryRow, error) {
	return GetInventorySummary(s.db, tenantID)
}

// SalesDaily returns daily sales breakdown for the tenant in the date range.
func (s *Service) SalesDaily(tenantID string, fromDate, toDate time.Time) ([]SalesDailyRow, error) {
	return GetSalesDaily(s.db, tenantID, fromDate, toDate)
}

// SalesTransactions returns each transaction in the date range for the tenant.
func (s *Service) SalesTransactions(tenantID string, fromDate, toDate time.Time) ([]SalesTransactionRow, error) {
	return GetSalesTransactions(s.db, tenantID, fromDate, toDate)
}

// ProfitReport returns profit summary and product rows for the tenant in the date range.
func (s *Service) ProfitReport(tenantID string, fromDate, toDate time.Time) (ProfitSummary, []ProfitRow, error) {
	return GetProfitReport(s.db, tenantID, fromDate, toDate)
}

// CashiersReport returns cashier performance for the tenant in the date range.
func (s *Service) CashiersReport(tenantID string, fromDate, toDate time.Time) ([]CashierRow, error) {
	return GetCashiersReport(s.db, tenantID, fromDate, toDate)
}
