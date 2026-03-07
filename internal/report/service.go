package report

import (
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func (s *Service) DB() *gorm.DB {
	return s.db
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

// SalesTransactions returns paginated transactions. Limit 0 = all (for export).
func (s *Service) SalesTransactions(tenantID string, fromDate, toDate time.Time, limit, offset int) ([]SalesTransactionRow, error) {
	return GetSalesTransactions(s.db, tenantID, fromDate, toDate, limit, offset)
}

// CountSalesTransactions returns total count of transactions in the date range.
func (s *Service) CountSalesTransactions(tenantID string, fromDate, toDate time.Time) (int64, error) {
	return CountSalesTransactions(s.db, tenantID, fromDate, toDate)
}

// PaymentsReport returns payment method report for the tenant in the date range.
func (s *Service) PaymentsReport(tenantID string, fromDate, toDate time.Time) ([]PaymentRow, error) {
	return GetPaymentsReport(s.db, tenantID, fromDate, toDate)
}

// SalesHourly returns hourly sales for the tenant on the given date.
func (s *Service) SalesHourly(tenantID string, date time.Time) ([]SalesHourlyRow, error) {
	return GetSalesHourly(s.db, tenantID, date)
}

// ProfitReport returns profit summary and product rows for the tenant in the date range.
func (s *Service) ProfitReport(tenantID string, fromDate, toDate time.Time) (ProfitSummary, []ProfitRow, error) {
	return GetProfitReport(s.db, tenantID, fromDate, toDate)
}

// CashiersReport returns cashier performance for the tenant in the date range.
func (s *Service) CashiersReport(tenantID string, fromDate, toDate time.Time) ([]CashierRow, error) {
	return GetCashiersReport(s.db, tenantID, fromDate, toDate)
}

// ShiftsReport returns shift reconciliation report for the tenant in the date range.
func (s *Service) ShiftsReport(tenantID string, fromDate, toDate time.Time) ([]ShiftReportRow, error) {
	return GetShiftsReport(s.db, tenantID, fromDate, toDate)
}

type CompareResult struct {
	Current  PeriodStats `json:"current"`
	Previous PeriodStats `json:"previous"`
	ChangePct float64    `json:"change_pct"`
}

type PeriodStats struct {
	TotalSales        float64 `json:"total_sales"`
	TotalTransactions int     `json:"total_transactions"`
}

func (s *Service) SalesCompare(tenantID string, curFrom, curTo, prevFrom, prevTo time.Time) (*CompareResult, error) {
	cur, err := GetSalesSummary(s.db, tenantID, curFrom, curTo)
	if err != nil {
		return nil, err
	}
	prev, err := GetSalesSummary(s.db, tenantID, prevFrom, prevTo)
	if err != nil {
		return nil, err
	}
	changePct := 0.0
	if prev.TotalSales > 0 {
		changePct = ((cur.TotalSales - prev.TotalSales) / prev.TotalSales) * 100
	}
	return &CompareResult{
		Current:  PeriodStats{TotalSales: cur.TotalSales, TotalTransactions: cur.TotalTransactions},
		Previous: PeriodStats{TotalSales: prev.TotalSales, TotalTransactions: prev.TotalTransactions},
		ChangePct: changePct,
	}, nil
}

func (s *Service) ProductMargin(tenantID string, fromDate, toDate time.Time) ([]MarginRow, error) {
	return GetProductMargin(s.db, tenantID, fromDate, toDate)
}
