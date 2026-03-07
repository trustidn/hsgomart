package opname

import (
	"errors"
	"time"

	"github.com/trustidn/hsmart-saas/internal/inventory"
	"gorm.io/gorm"
)

var (
	ErrOpnameNotFound    = errors.New("stock opname not found")
	ErrOpnameNotDraft    = errors.New("stock opname is not in draft status")
	ErrOpnameAlreadyDone = errors.New("stock opname already completed")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) StartOpname(tenantID, userID string) (*StockOpname, error) {
	op := &StockOpname{
		TenantID: tenantID,
		UserID:   userID,
		Status:   "draft",
	}
	if err := s.db.Create(op).Error; err != nil {
		return nil, err
	}
	return op, nil
}

type SubmitItemInput struct {
	ProductID   string `json:"product_id" binding:"required"`
	ActualStock int    `json:"actual_stock" binding:"min=0"`
	Notes       string `json:"notes"`
}

func (s *Service) SubmitItems(tenantID, opnameID string, items []SubmitItemInput) (*StockOpname, error) {
	var op StockOpname
	if err := s.db.Where("id = ? AND tenant_id = ?", opnameID, tenantID).First(&op).Error; err != nil {
		return nil, ErrOpnameNotFound
	}
	if op.Status != "draft" {
		return nil, ErrOpnameNotDraft
	}

	s.db.Where("opname_id = ?", opnameID).Delete(&StockOpnameItem{})

	for _, it := range items {
		inv, _ := inventory.GetInventoryByProduct(s.db, tenantID, it.ProductID)
		sysStock := 0
		if inv != nil {
			sysStock = inv.Stock
		}
		diff := it.ActualStock - sysStock
		item := &StockOpnameItem{
			OpnameID:    opnameID,
			ProductID:   it.ProductID,
			SystemStock: sysStock,
			ActualStock: it.ActualStock,
			Difference:  diff,
			Notes:       it.Notes,
		}
		if err := s.db.Create(item).Error; err != nil {
			return nil, err
		}
	}

	op.Status = "submitted"
	s.db.Model(&op).Update("status", "submitted")
	return &op, nil
}

func (s *Service) ApproveOpname(tenantID, opnameID string) (*StockOpname, error) {
	var op StockOpname
	if err := s.db.Where("id = ? AND tenant_id = ?", opnameID, tenantID).First(&op).Error; err != nil {
		return nil, ErrOpnameNotFound
	}
	if op.Status == "completed" {
		return nil, ErrOpnameAlreadyDone
	}

	var items []StockOpnameItem
	s.db.Where("opname_id = ?", opnameID).Find(&items)

	return &op, s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if item.Difference == 0 {
				continue
			}
			inv, _ := inventory.GetInventoryByProduct(tx, tenantID, item.ProductID)
			if inv == nil {
				continue
			}
			newStock := item.ActualStock
			if err := inventory.UpdateStock(tx, tenantID, item.ProductID, newStock); err != nil {
				return err
			}
			refID := opnameID
			m := &inventory.StockMovement{
				TenantID:    tenantID,
				ProductID:   item.ProductID,
				Type:        "opname",
				Quantity:    item.Difference,
				Reference:   "Stock opname adjustment",
				ReferenceID: &refID,
			}
			if err := inventory.CreateMovement(tx, m); err != nil {
				return err
			}
		}
		now := time.Now()
		return tx.Model(&StockOpname{}).Where("id = ?", opnameID).Updates(map[string]interface{}{
			"status":       "completed",
			"completed_at": now,
		}).Error
	})
}

func (s *Service) GetOpname(tenantID, opnameID string) (*StockOpname, []StockOpnameItem, error) {
	var op StockOpname
	if err := s.db.Where("id = ? AND tenant_id = ?", opnameID, tenantID).First(&op).Error; err != nil {
		return nil, nil, ErrOpnameNotFound
	}
	var items []StockOpnameItem
	s.db.Where("opname_id = ?", opnameID).Find(&items)
	return &op, items, nil
}

func (s *Service) ListOpnames(tenantID string, limit, offset int) ([]StockOpname, error) {
	var list []StockOpname
	q := s.db.Where("tenant_id = ?", tenantID).Order("created_at DESC")
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset > 0 {
		q = q.Offset(offset)
	}
	err := q.Find(&list).Error
	return list, err
}
