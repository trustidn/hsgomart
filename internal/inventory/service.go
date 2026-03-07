package inventory

import (
	"errors"

	"github.com/trustidn/hsmart-saas/internal/product"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrInsufficientStock   = errors.New("insufficient stock")
	ErrAdjustOnlyDecrease  = errors.New("penambahan stok hanya melalui Purchase; Adjust Stock hanya untuk pengurangan (quantity negatif)")
)

const (
	MovementTypePurchase   = "purchase"
	MovementTypeSale       = "sale"
	MovementTypeAdjustment = "adjustment"
	MovementTypeReturn     = "return"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// ensureProductBelongsToTenant returns ErrProductNotFound if product is not in tenant.
func (s *Service) ensureProductBelongsToTenant(tenantID, productID string) error {
	_, err := product.FindProductByID(s.db, tenantID, productID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

// GetStock returns current stock from the inventory table. Returns 0 if no inventory row exists. Verifies product belongs to tenant.
// All stock checks (e.g. POS) must use this or GetInventoryByProduct; purchases only add to inventory, they do not replace it.
func (s *Service) GetStock(tenantID, productID string) (int, error) {
	if err := s.ensureProductBelongsToTenant(tenantID, productID); err != nil {
		return 0, err
	}
	inv, err := GetInventoryByProduct(s.db, tenantID, productID)
	if err != nil {
		return 0, err
	}
	if inv == nil {
		return 0, nil
	}
	return inv.Stock, nil
}

// AdjustStock applies a delta (negative only). Creates an adjustment movement record.
// Policy: stock increase only via Purchase; Adjust Stock is for corrections/reductions only.
func (s *Service) AdjustStock(tenantID, productID string, quantity int, movementType, reference string) error {
	if quantity > 0 {
		return ErrAdjustOnlyDecrease
	}
	if err := s.ensureProductBelongsToTenant(tenantID, productID); err != nil {
		return err
	}
	inv, err := GetInventoryByProduct(s.db, tenantID, productID)
	if err != nil {
		return err
	}
	if inv == nil {
		inv = &Inventory{TenantID: tenantID, ProductID: productID, Stock: 0}
		if err := CreateInventory(s.db, inv); err != nil {
			return err
		}
	}
	newStock := inv.Stock + quantity
	if newStock < 0 {
		return ErrInsufficientStock
	}
	if err := UpdateStock(s.db, tenantID, productID, newStock); err != nil {
		return err
	}
	m := &StockMovement{
		TenantID:  tenantID,
		ProductID: productID,
		Type:      movementType,
		Quantity:  quantity,
		Reference: reference,
	}
	return CreateMovement(s.db, m)
}

// IncreaseStock adds quantity and creates a movement record.
func (s *Service) IncreaseStock(tenantID, productID string, quantity int, movementType, reference string) error {
	if err := s.ensureProductBelongsToTenant(tenantID, productID); err != nil {
		return err
	}
	if err := IncreaseStock(s.db, tenantID, productID, quantity); err != nil {
		return err
	}
	m := &StockMovement{
		TenantID:  tenantID,
		ProductID: productID,
		Type:      movementType,
		Quantity:  quantity,
		Reference: reference,
	}
	return CreateMovement(s.db, m)
}

// DecreaseStock subtracts quantity and creates a movement record. Returns ErrInsufficientStock if not enough.
func (s *Service) DecreaseStock(tenantID, productID string, quantity int, movementType, reference string) error {
	if err := s.ensureProductBelongsToTenant(tenantID, productID); err != nil {
		return err
	}
	inv, err := GetInventoryByProduct(s.db, tenantID, productID)
	if err != nil {
		return err
	}
	if inv == nil || inv.Stock < quantity {
		return ErrInsufficientStock
	}
	if err := DecreaseStock(s.db, tenantID, productID, quantity); err != nil {
		return err
	}
	m := &StockMovement{
		TenantID:  tenantID,
		ProductID: productID,
		Type:      movementType,
		Quantity:  quantity,
		Reference: reference,
	}
	return CreateMovement(s.db, m)
}

// ListInventoryByTenant returns all inventory rows for the tenant.
func (s *Service) ListInventoryByTenant(tenantID string) ([]Inventory, error) {
	return ListInventoriesByTenant(s.db, tenantID)
}

// ListMovements returns movements for tenant, optionally for one product.
func (s *Service) ListMovements(tenantID, productID string) ([]StockMovement, error) {
	return ListMovements(s.db, tenantID, productID)
}

// ListMovementRows returns movements with product name for the tenant.
func (s *Service) ListMovementRows(tenantID, productID string) ([]MovementRow, error) {
	return ListMovementRows(s.db, tenantID, productID)
}

// ListMovementRowsPaginated returns movements and total count for pagination.
func (s *Service) ListMovementRowsPaginated(tenantID, productID string, limit, offset int) ([]MovementRow, int64, error) {
	return ListMovementRowsPaginated(s.db, tenantID, productID, limit, offset)
}
