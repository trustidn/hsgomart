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

// LowStockRow holds a product with stock at or below threshold.
type LowStockRow struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Stock       int    `json:"stock"`
	Threshold   int    `json:"threshold"`
}

// GetLowStockProducts returns products where stock <= low_stock_threshold for the tenant.
func GetLowStockProducts(db *gorm.DB, tenantID string) ([]LowStockRow, error) {
	var list []LowStockRow
	err := db.Table("products").
		Select("products.id as product_id, products.name as product_name, COALESCE(inventories.stock, 0) as stock, COALESCE(products.low_stock_threshold, 10) as threshold").
		Joins("INNER JOIN inventories ON inventories.product_id = products.id AND inventories.tenant_id = products.tenant_id").
		Where("products.tenant_id = ?", tenantID).
		Where("COALESCE(inventories.stock, 0) <= COALESCE(products.low_stock_threshold, 10)").
		Order("inventories.stock ASC").
		Scan(&list).Error
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

// MovementRow holds a movement with product name and running balance for listing.
type MovementRow struct {
	ProductName string
	Type        string
	Quantity    int
	Reference   string
	Reason      string
	CreatedAt   time.Time
	StockAfter  int // running balance: stock level after this movement (cumulative sum per product)
}

// ListMovementRows returns movements with product name for the tenant (productID empty = all).
func ListMovementRows(db *gorm.DB, tenantID, productID string) ([]MovementRow, error) {
	return listMovementRows(db, tenantID, productID, "", "", "", 0, 0)
}

// listMovementRows with optional filters and limit/offset (limit 0 = no limit).
// Uses a window function to compute stock_after (running balance) per product.
// movementType, fromDate, toDate empty = no filter. Dates in YYYY-MM-DD; range is inclusive.
func listMovementRows(db *gorm.DB, tenantID, productID, movementType, fromDate, toDate string, limit, offset int) ([]MovementRow, error) {
	// Window: SUM(quantity) OVER (PARTITION BY product_id ORDER BY created_at) = running balance after each movement.
	sql := `WITH running AS (
  SELECT sm.tenant_id, sm.product_id, sm.type, sm.quantity, sm.reference, COALESCE(sm.reason,'') AS reason, sm.created_at,
         p.name AS product_name,
         (SUM(sm.quantity) OVER (PARTITION BY sm.product_id ORDER BY sm.created_at))::integer AS stock_after
  FROM stock_movements sm
  LEFT JOIN products p ON p.id = sm.product_id
  WHERE sm.tenant_id = ?
)
SELECT product_name, type, quantity, reference, reason, created_at, stock_after
FROM running
WHERE 1=1`
	args := []interface{}{tenantID}

	if productID != "" {
		sql += " AND product_id = ?"
		args = append(args, productID)
	}
	if movementType != "" {
		sql += " AND type = ?"
		args = append(args, movementType)
	}
	if fromDate != "" {
		sql += " AND created_at >= ?"
		args = append(args, fromDate+"T00:00:00Z")
	}
	if toDate != "" {
		sql += " AND created_at <= ?"
		args = append(args, toDate+"T23:59:59.999Z")
	}
	sql += " ORDER BY created_at DESC"
	if limit > 0 {
		sql += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		sql += " OFFSET ?"
		args = append(args, offset)
	}
	var list []MovementRow
	err := db.Raw(sql, args...).Scan(&list).Error
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

type ExpiringRow struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	BatchID     string `json:"batch_id"`
	Remaining   int    `json:"remaining"`
	ExpiredAt   string `json:"expired_at"`
}

func GetExpiringProducts(db *gorm.DB, tenantID string, days int) ([]ExpiringRow, error) {
	var list []ExpiringRow
	err := db.Raw(`
		SELECT ib.id AS batch_id, p.id AS product_id, p.name AS product_name,
		       ib.remaining_quantity AS remaining, TO_CHAR(ib.expired_at, 'YYYY-MM-DD') AS expired_at
		FROM inventory_batches ib
		JOIN products p ON p.id = ib.product_id
		WHERE p.tenant_id = ?
		  AND ib.remaining_quantity > 0
		  AND ib.expired_at IS NOT NULL
		  AND ib.expired_at <= CURRENT_DATE + INTERVAL '1 day' * ?
		ORDER BY ib.expired_at ASC
	`, tenantID, days).Scan(&list).Error
	return list, err
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
