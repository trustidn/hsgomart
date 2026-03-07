package inventory

type ServiceInterface interface {
	GetStock(tenantID, productID string) (int, error)
	AdjustStock(tenantID, productID string, quantity int, movementType, reference, reason string) error
	IncreaseStock(tenantID, productID string, quantity int, movementType, reference string) error
	DecreaseStock(tenantID, productID string, quantity int, movementType, reference string) error
	ListInventoryByTenant(tenantID string) ([]Inventory, error)
	ListMovements(tenantID, productID string) ([]StockMovement, error)
	ListMovementRows(tenantID, productID string) ([]MovementRow, error)
	ListMovementRowsPaginated(tenantID, productID, movementType, fromDate, toDate string, limit, offset int) ([]MovementRow, int64, error)
	GetLowStockProducts(tenantID string) ([]LowStockRow, error)
	GetExpiringProducts(tenantID string, days int) ([]ExpiringRow, error)
}

var _ ServiceInterface = (*Service)(nil)
