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
