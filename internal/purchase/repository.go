package purchase

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func CreatePurchase(db *gorm.DB, p *Purchase) error {
	return db.Create(p).Error
}

func CreatePurchaseItem(db *gorm.DB, item *PurchaseItem) error {
	return db.Create(item).Error
}

func CreateInventoryBatch(db *gorm.DB, b *InventoryBatch) error {
	return db.Create(b).Error
}

// ListPurchasesByTenant returns purchases for the tenant, newest first.
func ListPurchasesByTenant(db *gorm.DB, tenantID string) ([]Purchase, error) {
	var list []Purchase
	err := db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&list).Error
	return list, err
}

// PurchaseProductNameRow is used to load product names per purchase for list view.
type PurchaseProductNameRow struct {
	PurchaseID  string
	ProductName string
}

// GetProductNamesByPurchaseIDs returns a map of purchase_id -> list of product names (order preserved per item).
func GetProductNamesByPurchaseIDs(db *gorm.DB, purchaseIDs []string) (map[string][]string, error) {
	if len(purchaseIDs) == 0 {
		return map[string][]string{}, nil
	}
	var rows []PurchaseProductNameRow
	err := db.Table("purchase_items").
		Select("purchase_items.purchase_id as purchase_id, COALESCE(products.name, '') as product_name").
		Joins("LEFT JOIN products ON products.id = purchase_items.product_id").
		Where("purchase_items.purchase_id IN ?", purchaseIDs).
		Order("purchase_items.purchase_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make(map[string][]string)
	for _, r := range rows {
		out[r.PurchaseID] = append(out[r.PurchaseID], r.ProductName)
	}
	return out, nil
}

// ExistsPurchaseByTenantAndInvoice returns true if a purchase with this tenant_id and invoice_number exists.
func ExistsPurchaseByTenantAndInvoice(db *gorm.DB, tenantID, invoiceNumber string) (bool, error) {
	if invoiceNumber == "" {
		return false, nil
	}
	var count int64
	err := db.Model(&Purchase{}).Where("tenant_id = ? AND invoice_number = ?", tenantID, invoiceNumber).Count(&count).Error
	return count > 0, err
}

// GetPurchaseByID returns one purchase by ID if it belongs to the tenant.
func GetPurchaseByID(db *gorm.DB, tenantID, id string) (*Purchase, error) {
	var p Purchase
	err := db.Where("tenant_id = ? AND id = ?", tenantID, id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPurchaseItems returns all items for a purchase (caller ensures purchase belongs to tenant).
func ListPurchaseItems(db *gorm.DB, purchaseID string) ([]PurchaseItem, error) {
	var list []PurchaseItem
	err := db.Where("purchase_id = ?", purchaseID).Find(&list).Error
	return list, err
}

// PurchaseItemRow adds product name for API response.
type PurchaseItemRow struct {
	PurchaseItem
	ProductName string `json:"product_name"`
}

// ListPurchaseItemRows returns items with product name for a purchase.
func ListPurchaseItemRows(db *gorm.DB, purchaseID string) ([]PurchaseItemRow, error) {
	var list []PurchaseItemRow
	err := db.Table("purchase_items").
		Select("purchase_items.*, COALESCE(products.name, '') as product_name").
		Joins("LEFT JOIN products ON products.id = purchase_items.product_id").
		Where("purchase_items.purchase_id = ?", purchaseID).
		Scan(&list).Error
	return list, err
}

// FindAvailableBatches returns batches for the product with remaining_quantity > 0, ordered by created_at ASC (FIFO).
// Uses FOR UPDATE to prevent double-selling under concurrent POS checkouts.
func FindAvailableBatches(db *gorm.DB, productID string) ([]InventoryBatch, error) {
	var list []InventoryBatch
	err := db.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("product_id = ? AND remaining_quantity > 0", productID).
		Order("created_at ASC").
		Find(&list).Error
	return list, err
}

// UpdateBatchRemaining sets remaining_quantity for the batch.
func UpdateBatchRemaining(db *gorm.DB, batchID string, remaining int) error {
	return db.Model(&InventoryBatch{}).Where("id = ?", batchID).
		Update("remaining_quantity", remaining).Error
}
