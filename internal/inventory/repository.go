package inventory

import (
	"time"

	"gorm.io/gorm"
)

// GetInventoryByProduct returns the inventory row for tenant+product. Nil if not found.
// This table (inventories.stock) is the single source of truth for stock; purchases only increase it, POS and others read from it.
func GetInventoryByProduct(db *gorm.DB, tenantID, productID string) (*Inventory, error) {
	var inv Inventory
	err := db.Where("tenant_id = ? AND product_id = ?", tenantID, productID).First(&inv).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

// CreateInventory inserts a new inventory row (e.g. stock 0). Caller must set TenantID, ProductID, Stock.
func CreateInventory(db *gorm.DB, inv *Inventory) error {
	return db.Create(inv).Error
}

// UpdateStock updates the stock and updated_at for the given tenant+product row.
func UpdateStock(db *gorm.DB, tenantID, productID string, stock int) error {
	return db.Model(&Inventory{}).
		Where("tenant_id = ? AND product_id = ?", tenantID, productID).
		Updates(map[string]interface{}{"stock": stock, "updated_at": gorm.Expr("CURRENT_TIMESTAMP")}).Error
}

// IncreaseStock adds quantity to current stock. Does not create movement; caller creates movement.
func IncreaseStock(db *gorm.DB, tenantID, productID string, quantity int) error {
	res := db.Model(&Inventory{}).
		Where("tenant_id = ? AND product_id = ?", tenantID, productID).
		Update("stock", gorm.Expr("stock + ?", quantity))
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		// Row may not exist; create with stock = quantity
		inv := &Inventory{TenantID: tenantID, ProductID: productID, Stock: quantity}
		return CreateInventory(db, inv)
	}
	return nil
}

// DecreaseStock subtracts quantity from current stock. Caller must ensure sufficient stock.
func DecreaseStock(db *gorm.DB, tenantID, productID string, quantity int) error {
	return db.Model(&Inventory{}).
		Where("tenant_id = ? AND product_id = ?", tenantID, productID).
		Update("stock", gorm.Expr("stock - ?", quantity)).Error
}

// CreateMovement inserts a stock movement record.
func CreateMovement(db *gorm.DB, m *StockMovement) error {
	return db.Create(m).Error
}

// ListInventoriesByTenant returns all inventory rows for the tenant.
func ListInventoriesByTenant(db *gorm.DB, tenantID string) ([]Inventory, error) {
	var list []Inventory
	err := db.Where("tenant_id = ?", tenantID).Find(&list).Error
	return list, err
}

// ListMovements returns movements for tenant, optionally filtered by productID (empty = all).
func ListMovements(db *gorm.DB, tenantID, productID string) ([]StockMovement, error) {
	var list []StockMovement
	q := db.Where("tenant_id = ?", tenantID)
	if productID != "" {
		q = q.Where("product_id = ?", productID)
	}
	err := q.Order("created_at DESC").Find(&list).Error
	return list, err
}

// MovementRow holds a movement with product name for listing (used by ListMovementRows).
type MovementRow struct {
	ProductName string
	Type        string
	Quantity    int
	Reference   string
	Reason      string
	CreatedAt   time.Time
}

// ListMovementRows returns movements with product name for the tenant (productID empty = all).
func ListMovementRows(db *gorm.DB, tenantID, productID string) ([]MovementRow, error) {
	return listMovementRows(db, tenantID, productID, "", "", "", 0, 0)
}

// listMovementRows with optional filters and limit/offset (limit 0 = no limit).
// movementType, fromDate, toDate empty = no filter. Dates in YYYY-MM-DD; range is inclusive.
func listMovementRows(db *gorm.DB, tenantID, productID, movementType, fromDate, toDate string, limit, offset int) ([]MovementRow, error) {
	var list []MovementRow
	q := db.Table("stock_movements").
		Select("products.name as product_name, stock_movements.type as type, stock_movements.quantity as quantity, stock_movements.reference as reference, COALESCE(stock_movements.reason, '') as reason, stock_movements.created_at as created_at").
		Joins("LEFT JOIN products ON products.id = stock_movements.product_id").
		Where("stock_movements.tenant_id = ?", tenantID)
	if productID != "" {
		q = q.Where("stock_movements.product_id = ?", productID)
	}
	if movementType != "" {
		q = q.Where("stock_movements.type = ?", movementType)
	}
	if fromDate != "" {
		q = q.Where("stock_movements.created_at >= ?", fromDate+"T00:00:00Z")
	}
	if toDate != "" {
		q = q.Where("stock_movements.created_at <= ?", toDate+"T23:59:59.999Z")
	}
	q = q.Order("stock_movements.created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Scan(&list).Error
	return list, err
}

// CountMovements returns total count with same filters as listMovementRows.
func CountMovements(db *gorm.DB, tenantID, productID, movementType, fromDate, toDate string) (int64, error) {
	var n int64
	q := db.Model(&StockMovement{}).Where("tenant_id = ?", tenantID)
	if productID != "" {
		q = q.Where("product_id = ?", productID)
	}
	if movementType != "" {
		q = q.Where("type = ?", movementType)
	}
	if fromDate != "" {
		q = q.Where("created_at >= ?", fromDate+"T00:00:00Z")
	}
	if toDate != "" {
		q = q.Where("created_at <= ?", toDate+"T23:59:59.999Z")
	}
	err := q.Count(&n).Error
	return n, err
}

// ListMovementRowsPaginated returns movements and total count for pagination.
func ListMovementRowsPaginated(db *gorm.DB, tenantID, productID, movementType, fromDate, toDate string, limit, offset int) ([]MovementRow, int64, error) {
	list, err := listMovementRows(db, tenantID, productID, movementType, fromDate, toDate, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := CountMovements(db, tenantID, productID, movementType, fromDate, toDate)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
