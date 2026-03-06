package product

import "time"

type Category struct {
	ID        string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID  string    `gorm:"type:uuid;not null;index"`
	Name      string    `gorm:"type:varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Category) TableName() string {
	return "categories"
}

type Product struct {
	ID         string    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TenantID   string    `gorm:"type:uuid;not null;index"`
	CategoryID *string   `gorm:"type:uuid;index"`
	Name       string    `gorm:"type:varchar(255)"`
	SKU        string    `gorm:"column:sku;type:varchar(100)"`
	CostPrice  float64   `gorm:"column:cost_price;type:numeric"`
	SellPrice  float64   `gorm:"column:sell_price;type:numeric"`
	Status     string    `gorm:"type:varchar(50);default:active"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
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
