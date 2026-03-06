package pos

import (
	"errors"

	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/product"
	"github.com/trustidn/hsmart-saas/internal/purchase"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidItems      = errors.New("items cannot be empty")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CheckoutItemInput struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

type CheckoutInput struct {
	Items         []CheckoutItemInput `json:"items" binding:"required,min=1,dive"`
	PaymentMethod string              `json:"payment_method" binding:"required"`
	PaidAmount    float64             `json:"paid_amount" binding:"required,min=0"`
}

type CheckoutResult struct {
	TransactionID string  `json:"transaction_id"`
	Total         float64 `json:"total"`
	Change        float64 `json:"change"`
}

// Checkout runs a sale in a single DB transaction: validate products, check stock, decrease inventory,
// create movements, then create transaction, items, and payment. Rolls back on any error.
func (s *Service) Checkout(tenantID, userID string, in CheckoutInput) (*CheckoutResult, error) {
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
	var itemsToInsert []TransactionItem

	for _, it := range in.Items {
		// 1. Find product (ensures product belongs to tenant)
		p, err := product.FindProductByID(tx, tenantID, it.ProductID)
		if err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, ErrProductNotFound
			}
			return nil, err
		}

		// 2. Check stock
		inv, err := inventory.GetInventoryByProduct(tx, tenantID, it.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if inv == nil || inv.Stock < it.Quantity {
			tx.Rollback()
			return nil, ErrInsufficientStock
		}

		// 3. FIFO: deduct from batches and compute COGS
		cogs, err := deductBatchesFIFO(tx, it.ProductID, it.Quantity)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// 4. Calculate subtotal (revenue)
		price := p.SellPrice
		subtotal := price * float64(it.Quantity)
		totalAmount += subtotal

		// 5. Decrease inventory stock
		if err := inventory.DecreaseStock(tx, tenantID, it.ProductID, it.Quantity); err != nil {
			tx.Rollback()
			return nil, err
		}

		// 6. Create stock movement (type = sale)
		m := &inventory.StockMovement{
			TenantID:  tenantID,
			ProductID: it.ProductID,
			Type:      inventory.MovementTypeSale,
			Quantity:  it.Quantity,
			Reference: "POS sale",
		}
		if err := inventory.CreateMovement(tx, m); err != nil {
			tx.Rollback()
			return nil, err
		}

		itemsToInsert = append(itemsToInsert, TransactionItem{
			ProductID: it.ProductID,
			Price:     price,
			Quantity:  it.Quantity,
			Subtotal:  subtotal,
			Cogs:      cogs,
		})
	}

	// Insert transaction
	t := &Transaction{
		TenantID:    tenantID,
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "completed",
	}
	if err := CreateTransaction(tx, t); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Insert transaction items
	for i := range itemsToInsert {
		itemsToInsert[i].TransactionID = t.ID
		if err := CreateTransactionItem(tx, &itemsToInsert[i]); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Insert payment
	pay := &Payment{
		TransactionID: t.ID,
		Method:        in.PaymentMethod,
		Amount:        in.PaidAmount,
	}
	if err := CreatePayment(tx, pay); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	change := in.PaidAmount - totalAmount
	if change < 0 {
		change = 0
	}

	return &CheckoutResult{
		TransactionID: t.ID,
		Total:         totalAmount,
		Change:        change,
	}, nil
}

// deductBatchesFIFO deducts quantity from oldest batches first, updates remaining_quantity, returns total COGS.
func deductBatchesFIFO(tx *gorm.DB, productID string, quantity int) (float64, error) {
	batches, err := purchase.FindAvailableBatches(tx, productID)
	if err != nil {
		return 0, err
	}
	var cogs float64
	remaining := quantity
	for _, b := range batches {
		if remaining <= 0 {
			break
		}
		deduct := b.RemainingQuantity
		if deduct > remaining {
			deduct = remaining
		}
		cogs += b.CostPrice * float64(deduct)
		newRemaining := b.RemainingQuantity - deduct
		if err := purchase.UpdateBatchRemaining(tx, b.ID, newRemaining); err != nil {
			return 0, err
		}
		remaining -= deduct
	}
	if remaining > 0 {
		return 0, ErrInsufficientStock
	}
	return cogs, nil
}
