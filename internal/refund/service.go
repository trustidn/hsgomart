package refund

import (
	"errors"

	"github.com/trustidn/hsmart-saas/internal/inventory"
	"github.com/trustidn/hsmart-saas/internal/pos"
	"github.com/trustidn/hsmart-saas/pkg/utils"
	"gorm.io/gorm"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrAlreadyRefunded     = errors.New("transaction already refunded")
	ErrRefundExceedsTotal  = errors.New("refund amount exceeds transaction total")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

type RefundInput struct {
	TransactionID string `json:"transaction_id" binding:"required"`
	Reason        string `json:"reason"`
}

type RefundResult struct {
	RefundID string  `json:"refund_id"`
	Amount   float64 `json:"amount"`
	Status   string  `json:"status"`
}

func (s *Service) CreateRefund(tenantID, userID string, in RefundInput) (*RefundResult, error) {
	var txn pos.Transaction
	if err := s.db.Where("id = ? AND tenant_id = ?", in.TransactionID, tenantID).First(&txn).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, err
	}

	if txn.Status == "refunded" {
		return nil, ErrAlreadyRefunded
	}

	var existingRefundTotal float64
	s.db.Model(&Refund{}).Where("transaction_id = ? AND tenant_id = ?", in.TransactionID, tenantID).
		Select("COALESCE(SUM(amount), 0)").Scan(&existingRefundTotal)

	refundAmount := txn.TotalAmount - existingRefundTotal
	if refundAmount <= 0 {
		return nil, ErrAlreadyRefunded
	}

	return s.processRefund(tenantID, userID, in, &txn, refundAmount)
}

func (s *Service) processRefund(tenantID, userID string, in RefundInput, txn *pos.Transaction, amount float64) (*RefundResult, error) {
	var result *RefundResult

	err := s.db.Transaction(func(tx *gorm.DB) error {
		r := &Refund{
			TenantID:      tenantID,
			TransactionID: in.TransactionID,
			UserID:        userID,
			Amount:        amount,
			Reason:        in.Reason,
			Status:        "completed",
		}
		if err := tx.Create(r).Error; err != nil {
			return err
		}

		var items []pos.TransactionItem
		if err := tx.Where("transaction_id = ?", in.TransactionID).Find(&items).Error; err != nil {
			return err
		}

		for _, item := range items {
			if err := inventory.IncreaseStock(tx, tenantID, item.ProductID, item.Quantity); err != nil {
				return err
			}
			refID := r.ID
			m := &inventory.StockMovement{
				TenantID:    tenantID,
				ProductID:   item.ProductID,
				Type:        "refund",
				Quantity:    item.Quantity,
				Reference:   "Refund: " + in.Reason,
				ReferenceID: &refID,
			}
			if err := inventory.CreateMovement(tx, m); err != nil {
				return err
			}
		}

		refID := r.ID
		pay := &pos.Payment{
			TransactionID: in.TransactionID,
			Method:        "refund",
			Amount:        -amount,
		}
		if err := tx.Create(pay).Error; err != nil {
			return err
		}

		if err := tx.Model(&pos.Transaction{}).Where("id = ?", in.TransactionID).Update("status", "refunded").Error; err != nil {
			return err
		}

		utils.LogAudit(tx, tenantID, userID, "refund", "transaction", in.TransactionID, map[string]interface{}{
			"refund_id": refID,
			"amount":    amount,
			"reason":    in.Reason,
		})

		result = &RefundResult{
			RefundID: r.ID,
			Amount:   amount,
			Status:   "completed",
		}
		return nil
	})

	return result, err
}

func (s *Service) ListRefunds(tenantID string, limit, offset int) ([]Refund, error) {
	return ListRefundsByTenant(s.db, tenantID, limit, offset)
}
