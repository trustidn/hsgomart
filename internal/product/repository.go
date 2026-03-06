package product

import (
	"gorm.io/gorm"
)

// CreateCategory inserts a category. Caller must set TenantID.
func CreateCategory(db *gorm.DB, c *Category) error {
	return db.Create(c).Error
}

// ListCategoriesByTenant returns all categories for the tenant (tenant isolation).
func ListCategoriesByTenant(db *gorm.DB, tenantID string) ([]Category, error) {
	var list []Category
	err := db.Where("tenant_id = ?", tenantID).Find(&list).Error
	return list, err
}

// CreateProduct inserts a product. Caller must set TenantID.
func CreateProduct(db *gorm.DB, p *Product) error {
	return db.Create(p).Error
}

// ListProductsByTenant returns all products for the tenant (tenant isolation).
func ListProductsByTenant(db *gorm.DB, tenantID string) ([]Product, error) {
	var list []Product
	err := db.Where("tenant_id = ?", tenantID).Find(&list).Error
	return list, err
}

// FindProductByID returns one product by ID scoped by tenant_id.
func FindProductByID(db *gorm.DB, tenantID, productID string) (*Product, error) {
	var p Product
	err := db.Where("tenant_id = ? AND id = ?", tenantID, productID).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// AddBarcodeToProduct inserts a barcode for a product. Product must belong to tenant (caller ensures).
func AddBarcodeToProduct(db *gorm.DB, pb *ProductBarcode) error {
	return db.Create(pb).Error
}

// FindProductByBarcode returns the product that has this barcode, scoped by tenant_id.
func FindProductByBarcode(db *gorm.DB, tenantID, barcode string) (*Product, error) {
	var p Product
	err := db.Model(&Product{}).
		Joins("INNER JOIN product_barcodes ON product_barcodes.product_id = products.id").
		Where("products.tenant_id = ? AND product_barcodes.barcode = ?", tenantID, barcode).
		First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}
