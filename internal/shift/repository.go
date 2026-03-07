package shift

import (
	"time"

	"gorm.io/gorm"
)

func Create(db *gorm.DB, s *CashierShift) error {
	return db.Create(s).Error
}

func GetActiveByUser(db *gorm.DB, tenantID, userID string) (*CashierShift, error) {
	var list []CashierShift
	err := db.Where("tenant_id = ? AND user_id = ? AND status = ?", tenantID, userID, StatusOpen).Limit(1).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return &list[0], nil
}

func GetByID(db *gorm.DB, tenantID, shiftID string) (*CashierShift, error) {
	var s CashierShift
	err := db.Where("tenant_id = ? AND id = ?", tenantID, shiftID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func CloseShift(db *gorm.DB, shiftID string, closingCash float64, closedAt time.Time) error {
	return db.Model(&CashierShift{}).Where("id = ?", shiftID).
		Updates(map[string]interface{}{
			"closing_cash": closingCash,
			"closed_at":    closedAt,
			"status":       StatusClosed,
		}).Error
}

func ListShifts(db *gorm.DB, tenantID string, limit, offset int) ([]CashierShift, error) {
	var list []CashierShift
	q := db.Table("cashier_shifts").Where("tenant_id = ?", tenantID).Order("opened_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Find(&list).Error
	return list, err
}

// SumCashPaymentsBetween returns sum of payment amounts where method = 'cash' and created_at in [from, to].
func SumCashPaymentsBetween(db *gorm.DB, tenantID string, from, to time.Time) (float64, error) {
	var sum float64
	err := db.Table("payments").
		Select("COALESCE(SUM(payments.amount), 0)").
		Joins("INNER JOIN transactions ON transactions.id = payments.transaction_id").
		Where("transactions.tenant_id = ? AND payments.method = ?", tenantID, "cash").
		Where("payments.created_at >= ? AND payments.created_at <= ?", from, to).
		Scan(&sum).Error
	return sum, err
}
