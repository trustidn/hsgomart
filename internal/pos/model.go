package pos

import "time"

type Transaction struct {
	ID             string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID       string    `gorm:"type:uuid;not null;index"`
	UserID         string    `gorm:"type:uuid;not null;index"`
	TotalAmount    float64   `gorm:"column:total_amount;type:numeric;not null"`
	DiscountAmount float64   `gorm:"column:discount_amount;type:numeric;default:0"`
	Status         string    `gorm:"type:varchar(50);default:completed"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionItem struct {
	ID             string  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TransactionID  string  `gorm:"type:uuid;not null;index"`
	ProductID      string  `gorm:"type:uuid;not null;index"`
	Price          float64 `gorm:"type:numeric;not null"`
	Quantity       int     `gorm:"not null"`
	Subtotal       float64 `gorm:"type:numeric;not null"`
	UnitCost       float64 `gorm:"column:unit_cost;type:numeric"`
	Cogs           float64 `gorm:"type:numeric;default:0"`
	DiscountType   string  `gorm:"column:discount_type;type:varchar(10)"`
	DiscountValue  float64 `gorm:"column:discount_value;type:numeric;default:0"`
	DiscountAmount float64 `gorm:"column:discount_amount;type:numeric;default:0"`
}

func (TransactionItem) TableName() string {
	return "transaction_items"
}

type Payment struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TransactionID string    `gorm:"type:uuid;not null;index"`
	Method        string    `gorm:"type:varchar(50);not null"`
	Amount        float64   `gorm:"type:numeric;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (Payment) TableName() string {
	return "payments"
}
