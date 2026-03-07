package purchase

import "time"

type Purchase struct {
	ID            string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID      string    `gorm:"type:uuid;not null;index"`
	SupplierName  string    `gorm:"type:varchar(255)"`
	InvoiceNumber string    `gorm:"type:varchar(255)"`
	Notes         string    `gorm:"type:text"`
	TotalAmount   float64   `gorm:"column:total_amount;type:numeric;not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
}

func (Purchase) TableName() string {
	return "purchases"
}

type PurchaseItem struct {
	ID         string  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	PurchaseID string  `gorm:"type:uuid;not null;index"`
	ProductID  string  `gorm:"type:uuid;not null;index"`
	Quantity   int     `gorm:"not null"`
	CostPrice  float64 `gorm:"column:cost_price;type:numeric;not null"`
	Subtotal   float64 `gorm:"type:numeric;not null"`
}

func (PurchaseItem) TableName() string {
	return "purchase_items"
}

type InventoryBatch struct {
	ID                string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID         string    `gorm:"type:uuid;not null;index"`
	PurchaseItemID    string    `gorm:"type:uuid;not null;index"`
	Quantity          int       `gorm:"not null"`
	RemainingQuantity int       `gorm:"not null"`
	CostPrice         float64   `gorm:"column:cost_price;type:numeric;not null"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
}

func (InventoryBatch) TableName() string {
	return "inventory_batches"
}
