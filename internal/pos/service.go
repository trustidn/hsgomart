package pos

import (
	"errors"
	"math"

	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/product"
	"github.com/trustidn/hsmart-saas/internal/purchase"
	"github.com/trustidn/hsmart-saas/internal/user"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidItems      = errors.New("items cannot be empty")
	ErrPaymentShort      = errors.New("total payment is less than amount due")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type CheckoutItemInput struct {
	ProductID     string  `json:"product_id" binding:"required"`
	Quantity      int     `json:"quantity" binding:"required,min=1"`
	DiscountType  string  `json:"discount_type"`
	DiscountValue float64 `json:"discount_value"`
}

type PaymentInput struct {
	Method string  `json:"method" binding:"required"`
	Amount float64 `json:"amount" binding:"required,min=0"`
}

type CheckoutInput struct {
	Items         []CheckoutItemInput `json:"items" binding:"required,min=1,dive"`
	Payments      []PaymentInput      `json:"payments"`
	PaymentMethod string              `json:"payment_method"`
	PaidAmount    float64             `json:"paid_amount"`
}

type CheckoutResult struct {
	TransactionID  string  `json:"transaction_id"`
	Cashier        string  `json:"cashier"`
	Total          float64 `json:"total"`
	DiscountTotal  float64 `json:"discount_total"`
	Change         float64 `json:"change"`
}

func (s *Service) Checkout(tenantID, userID string, in CheckoutInput) (*CheckoutResult, error) {
	if len(in.Items) == 0 {
		return nil, ErrInvalidItems
	}

	payments := in.Payments
	if len(payments) == 0 && in.PaymentMethod != "" {
		payments = []PaymentInput{{Method: in.PaymentMethod, Amount: in.PaidAmount}}
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var totalAmount, totalDiscount float64
	type itemCalc struct {
		product   *product.Product
		quantity  int
		price     float64
		subtotal  float64
		discType  string
		discValue float64
		discAmt   float64
	}
	var calcs []itemCalc

	for _, it := range in.Items {
		p, err := product.FindProductByID(tx, tenantID, it.ProductID)
		if err != nil {
			tx.Rollback()
			if err == gorm.ErrRecordNotFound {
				return nil, ErrProductNotFound
			}
			return nil, err
		}
		inv, err := inventory.GetInventoryByProduct(tx, tenantID, it.ProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if inv == nil || inv.Stock < it.Quantity {
			tx.Rollback()
			return nil, ErrInsufficientStock
		}

		lineTotal := p.SellPrice * float64(it.Quantity)
		var discAmt float64
		if it.DiscountType == "percent" && it.DiscountValue > 0 {
			discAmt = math.Round(lineTotal*it.DiscountValue) / 100
		} else if it.DiscountType == "fixed" && it.DiscountValue > 0 {
			discAmt = it.DiscountValue * float64(it.Quantity)
		}
		if discAmt > lineTotal {
			discAmt = lineTotal
		}

		subtotal := lineTotal - discAmt
		totalAmount += subtotal
		totalDiscount += discAmt

		calcs = append(calcs, itemCalc{
			product:   p,
			quantity:  it.Quantity,
			price:     p.SellPrice,
			subtotal:  subtotal,
			discType:  it.DiscountType,
			discValue: it.DiscountValue,
			discAmt:   discAmt,
		})
	}

	var totalPaid float64
	for _, p := range payments {
		totalPaid += p.Amount
	}
	if totalPaid < totalAmount {
		tx.Rollback()
		return nil, ErrPaymentShort
	}

	t := &Transaction{
		TenantID:       tenantID,
		UserID:         userID,
		TotalAmount:    totalAmount,
		DiscountAmount: totalDiscount,
		Status:         "completed",
	}
	if err := CreateTransaction(tx, t); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, calc := range calcs {
		cogs, err := deductBatchesFIFO(tx, calc.product.ID, calc.quantity)
		if err != nil {
			if err == ErrInsufficientStock {
				fallbackCost := calc.product.LastPurchasePrice
				if fallbackCost <= 0 {
					fallbackCost = calc.product.CostPrice
				}
				cogs = fallbackCost * float64(calc.quantity)
				if cogs < 0 {
					cogs = 0
				}
			} else {
				tx.Rollback()
				return nil, err
			}
		}
		unitCost := 0.0
		if calc.quantity > 0 {
			unitCost = cogs / float64(calc.quantity)
		}

		if err := inventory.DecreaseStock(tx, tenantID, calc.product.ID, calc.quantity); err != nil {
			tx.Rollback()
			return nil, err
		}
		refID := t.ID
		m := &inventory.StockMovement{
			TenantID:    tenantID,
			ProductID:   calc.product.ID,
			Type:        inventory.MovementTypeSale,
			Quantity:    calc.quantity,
			Reference:   "POS sale",
			ReferenceID: &refID,
		}
		if err := inventory.CreateMovement(tx, m); err != nil {
			tx.Rollback()
			return nil, err
		}

		item := &TransactionItem{
			TransactionID:  t.ID,
			ProductID:      calc.product.ID,
			Price:          calc.price,
			Quantity:       calc.quantity,
			Subtotal:       calc.subtotal,
			UnitCost:       unitCost,
			Cogs:           cogs,
			DiscountType:   calc.discType,
			DiscountValue:  calc.discValue,
			DiscountAmount: calc.discAmt,
		}
		if err := CreateTransactionItem(tx, item); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	for _, p := range payments {
		pay := &Payment{
			TransactionID: t.ID,
			Method:        p.Method,
			Amount:        p.Amount,
		}
		if err := CreatePayment(tx, pay); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	change := totalPaid - totalAmount
	if change < 0 {
		change = 0
	}

	cashierName := ""
	if u, err := user.FindUserByID(s.db, tenantID, userID); err == nil {
		cashierName = u.Name
		if cashierName == "" {
			cashierName = u.Email
		}
	}

	return &CheckoutResult{
		TransactionID: t.ID,
		Cashier:       cashierName,
		Total:         totalAmount,
		DiscountTotal: totalDiscount,
		Change:        change,
	}, nil
}

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
