package shift

import "time"

const StatusOpen = "open"
const StatusClosed = "closed"

type CashierShift struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    string     `gorm:"type:uuid;not null;index"`
	UserID      string     `gorm:"type:uuid;not null;index"`
	OpeningCash float64    `gorm:"column:opening_cash;type:numeric;not null"`
	ClosingCash *float64   `gorm:"column:closing_cash;type:numeric"`
	OpenedAt    time.Time  `gorm:"column:opened_at;not null"`
	ClosedAt    *time.Time `gorm:"column:closed_at"`
	Status      string     `gorm:"type:varchar(20);default:open"`
}

func (CashierShift) TableName() string {
	return "cashier_shifts"
}
