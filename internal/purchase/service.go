package purchase

import (
	"errors"

	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/product"
	"gorm.io/gorm"
)

var (
	ErrInvalidItems    = errors.New("items cannot be empty")
	ErrProductInvalid  = errors.New("product not found or does not belong to tenant")
	ErrInvalidQtyCost  = errors.New("quantity and cost_price must be positive")
	ErrDuplicateInvoice = errors.New("invoice number already exists for this tenant")
)

type CreatePurchaseItemInput struct {
	ProductID string  `json:"product_id" binding:"required"`
	Quantity  int     `json:"quantity" binding:"required,min=1"`
	CostPrice float64 `json:"cost_price" binding:"required,min=0"`
}

type CreatePurchaseInput struct {
	SupplierName  string                   `json:"supplier_name"`
	InvoiceNumber string                   `json:"invoice_number"`
	Items         []CreatePurchaseItemInput `json:"items" binding:"required,min=1,dive"`
}

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreatePurchase creates a purchase with items, inventory batches, updates stock, and stock movements in one transaction.
func (s *Service) CreatePurchase(tenantID string, in CreatePurchaseInput) (*Purchase, error) {
	if len(in.Items) == 0 {
		return nil, ErrInvalidItems
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalAmount float64
	for _, it := range in.Items {
		if it.Quantity <= 0 || it.CostPrice <= 0 {
			tx.Rollback()
			return nil, ErrInvalidQtyCost
		}
		_, err := product.FindProductByID(tx, tenantID, it.ProductID)
		if err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, ErrProductInvalid
			}
			return nil, err
		}
		subtotal := it.CostPrice * float64(it.Quantity)
		totalAmount += subtotal
	}

	// Prevent duplicate invoice per tenant (when invoice number is set)
	if in.InvoiceNumber != "" {
		exists, err := ExistsPurchaseByTenantAndInvoice(tx, tenantID, in.InvoiceNumber)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if exists {
			tx.Rollback()
			return nil, ErrDuplicateInvoice
		}
	}

	p := &Purchase{
		TenantID:      tenantID,
		SupplierName:  in.SupplierName,
		InvoiceNumber: in.InvoiceNumber,
		TotalAmount:   totalAmount,
	}
	if err := CreatePurchase(tx, p); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, it := range in.Items {
		subtotal := it.CostPrice * float64(it.Quantity)
		item := &PurchaseItem{
			PurchaseID: p.ID,
			ProductID:  it.ProductID,
			Quantity:   it.Quantity,
			CostPrice:  it.CostPrice,
			Subtotal:   subtotal,
		}
		if err := CreatePurchaseItem(tx, item); err != nil {
			tx.Rollback()
			return nil, err
		}

		batch := &InventoryBatch{
			ProductID:         it.ProductID,
			PurchaseItemID:    item.ID,
			Quantity:          it.Quantity,
			RemainingQuantity: it.Quantity,
			CostPrice:         it.CostPrice,
		}
		if err := CreateInventoryBatch(tx, batch); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := inventory.IncreaseStock(tx, tenantID, it.ProductID, it.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}

		ref := "purchase " + p.InvoiceNumber
		if ref == "purchase " {
			ref = "purchase " + p.ID
		}
		m := &inventory.StockMovement{
			TenantID:    tenantID,
			ProductID:   it.ProductID,
			Type:        inventory.MovementTypePurchase,
			Quantity:    it.Quantity,
			Reference:   ref,
			ReferenceID: p.ID,
		}
		if err := inventory.CreateMovement(tx, m); err != nil {
			tx.Rollback()
			return nil, err
		}
		// Update product last_purchase_price for reference in UI
		if err := product.UpdateLastPurchasePrice(tx, it.ProductID, it.CostPrice); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	return p, nil
}

// ListPurchases returns all purchases for the tenant.
func (s *Service) ListPurchases(tenantID string) ([]Purchase, error) {
	return ListPurchasesByTenant(s.db, tenantID)
}

// GetPurchase returns one purchase by ID with its items.
func (s *Service) GetPurchase(tenantID, purchaseID string) (*Purchase, []PurchaseItem, error) {
	p, err := GetPurchaseByID(s.db, tenantID, purchaseID)
	if err != nil {
		return nil, nil, err
	}
	items, err := ListPurchaseItems(s.db, p.ID)
	if err != nil {
		return nil, nil, err
	}
	return p, items, nil
}

// GetPurchaseWithItemRows returns one purchase by ID with items including product names.
func (s *Service) GetPurchaseWithItemRows(tenantID, purchaseID string) (*Purchase, []PurchaseItemRow, error) {
	p, err := GetPurchaseByID(s.db, tenantID, purchaseID)
	if err != nil {
		return nil, nil, err
	}
	rows, err := ListPurchaseItemRows(s.db, p.ID)
	if err != nil {
		return nil, nil, err
	}
	return p, rows, nil
}
