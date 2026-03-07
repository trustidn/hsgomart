package product

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	TenantID  string         `gorm:"type:uuid;not null;index" json:"tenant_id,omitempty"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Category) TableName() string {
	return "categories"
}

type Product struct {
	ID                 string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID           string    `gorm:"type:uuid;not null;index"`
	CategoryID         *string   `gorm:"type:uuid;index"`
	Name               string    `gorm:"type:varchar(255)"`
	SKU                string    `gorm:"column:sku;type:varchar(100)"`
	CostPrice          float64   `gorm:"column:cost_price;type:numeric"`
	LastPurchasePrice  float64   `gorm:"column:last_purchase_price;type:numeric"`
	SellPrice            float64   `gorm:"column:sell_price;type:numeric"`
	Unit                 string    `gorm:"column:unit;type:varchar(20);default:pcs"`
	LowStockThreshold    int       `gorm:"column:low_stock_threshold;default:10"`
	Status    string         `gorm:"type:varchar(50);default:active"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Product) TableName() string {
	return "products"
}

type ProductBarcode struct {
	ID        string `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProductID string `gorm:"type:uuid;not null;index"`
	Barcode   string `gorm:"type:varchar(100);uniqueIndex"`
}

func (ProductBarcode) TableName() string {
	return "product_barcodes"
}
