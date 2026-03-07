package inventory

import "time"

type Inventory struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  string    `gorm:"type:uuid;not null;uniqueIndex:idx_inv_tenant_product"`
	ProductID string    `gorm:"type:uuid;not null;uniqueIndex:idx_inv_tenant_product"`
	Stock     int       `gorm:"default:0"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (Inventory) TableName() string {
	return "inventories"
}

type StockMovement struct {
	ID          string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID    string    `gorm:"type:uuid;not null;index"`
	ProductID   string    `gorm:"type:uuid;not null;index"`
	Type        string    `gorm:"type:varchar(50);not null"` // purchase, sale, adjustment, return
	Quantity    int       `gorm:"not null"`
	Reference   string    `gorm:"type:varchar(255)"`
	Reason      string    `gorm:"type:varchar(255)"`             // adjustment reason for audit (e.g. expired, damaged)
	ReferenceID *string   `gorm:"type:uuid;column:reference_id"` // purchase_id or transaction_id for traceability; NULL for adjustment
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (StockMovement) TableName() string {
	return "stock_movements"
}
