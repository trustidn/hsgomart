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

// UpdateCategory updates a category by tenant_id and id.
func UpdateCategory(db *gorm.DB, tenantID, categoryID string, updates map[string]interface{}) error {
	res := db.Model(&Category{}).Where("tenant_id = ? AND id = ?", tenantID, categoryID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteCategory deletes a category by tenant_id and id.
func DeleteCategory(db *gorm.DB, tenantID, categoryID string) error {
	res := db.Where("tenant_id = ? AND id = ?", tenantID, categoryID).Delete(&Category{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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

// ListBarcodesByProductID returns all barcodes for a product (product_id only; caller ensures tenant).
func ListBarcodesByProductID(db *gorm.DB, productID string) ([]ProductBarcode, error) {
	var list []ProductBarcode
	err := db.Where("product_id = ?", productID).Find(&list).Error
	return list, err
}

// UpdateProduct updates a product by tenant_id and id (updates only non-zero/non-empty as needed).
func UpdateProduct(db *gorm.DB, tenantID, productID string, updates map[string]interface{}) error {
	res := db.Model(&Product{}).Where("tenant_id = ? AND id = ?", tenantID, productID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteProduct deletes a product by tenant_id and id.
func DeleteProduct(db *gorm.DB, tenantID, productID string) error {
	res := db.Where("tenant_id = ? AND id = ?", tenantID, productID).Delete(&Product{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
