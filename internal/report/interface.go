package report

import (
	"time"

	"gorm.io/gorm"
)

type ServiceInterface interface {
	DB() *gorm.DB
	SalesSummary(tenantID string, fromDate, toDate time.Time) (*SalesSummaryResult, error)
	TopProducts(tenantID string, fromDate, toDate time.Time) ([]TopProductRow, error)
	InventorySummary(tenantID string) ([]InventorySummaryRow, error)
	SalesDaily(tenantID string, fromDate, toDate time.Time) ([]SalesDailyRow, error)
	SalesTransactions(tenantID string, fromDate, toDate time.Time, limit, offset int) ([]SalesTransactionRow, error)
	CountSalesTransactions(tenantID string, fromDate, toDate time.Time) (int64, error)
	PaymentsReport(tenantID string, fromDate, toDate time.Time) ([]PaymentRow, error)
	SalesHourly(tenantID string, date time.Time) ([]SalesHourlyRow, error)
	ProfitReport(tenantID string, fromDate, toDate time.Time) (ProfitSummary, []ProfitRow, error)
	CashiersReport(tenantID string, fromDate, toDate time.Time) ([]CashierRow, error)
	ShiftsReport(tenantID string, fromDate, toDate time.Time) ([]ShiftReportRow, error)
	SalesCompare(tenantID string, curFrom, curTo, prevFrom, prevTo time.Time) (*CompareResult, error)
	ProductMargin(tenantID string, fromDate, toDate time.Time) ([]MarginRow, error)
}

var _ ServiceInterface = (*Service)(nil)
