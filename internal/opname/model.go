package opname

import "time"

type StockOpname struct {
	ID          string     `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID    string     `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID      string     `gorm:"type:uuid;not null" json:"user_id"`
	Status      string     `gorm:"type:varchar(50);default:draft" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

func (StockOpname) TableName() string {
	return "stock_opnames"
}

type StockOpnameItem struct {
	ID          string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	OpnameID    string `gorm:"type:uuid;not null;index" json:"opname_id"`
	ProductID   string `gorm:"type:uuid;not null" json:"product_id"`
	SystemStock int    `json:"system_stock"`
	ActualStock int    `json:"actual_stock"`
	Difference  int    `json:"difference"`
	Notes       string `gorm:"type:text" json:"notes"`
}

func (StockOpnameItem) TableName() string {
	return "stock_opname_items"
}
