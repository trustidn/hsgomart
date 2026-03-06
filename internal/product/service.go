package product

import (
	"errors"

	"github.com/trustidn/hsmart-saas/internal/subscription"
	"gorm.io/gorm"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrBarcodeExists     = errors.New("barcode already registered")
	ErrCategoryNotFound  = errors.New("category not found")
)

// PlanLimitChecker is used to enforce plan max_products (e.g. subscription.Service).
type PlanLimitChecker interface {
	CheckSubscription(tenantID string) (*subscription.SubscriptionWithPlan, error)
}

type Service struct {
	db               *gorm.DB
	planLimitChecker PlanLimitChecker
}

func NewService(db *gorm.DB, planLimitChecker PlanLimitChecker) *Service {
	return &Service{db: db, planLimitChecker: planLimitChecker}
}

func (s *Service) CreateCategory(tenantID string, name string) (*Category, error) {
	c := &Category{
		TenantID: tenantID,
		Name:     name,
	}
	if err := CreateCategory(s.db, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListCategories(tenantID string) ([]Category, error) {
	return ListCategoriesByTenant(s.db, tenantID)
}

func (s *Service) UpdateCategory(tenantID, categoryID string, name string) (*Category, error) {
	if err := UpdateCategory(s.db, tenantID, categoryID, map[string]interface{}{"name": name}); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	var c Category
	if err := s.db.Where("tenant_id = ? AND id = ?", tenantID, categoryID).First(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Service) DeleteCategory(tenantID, categoryID string) error {
	if err := DeleteCategory(s.db, tenantID, categoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	return nil
}

func (s *Service) CreateProduct(tenantID string, in CreateProductInput) (*Product, error) {
	if s.planLimitChecker != nil {
		subWithPlan, err := s.planLimitChecker.CheckSubscription(tenantID)
		if err != nil {
			return nil, err
		}
		var count int64
		s.db.Model(&Product{}).Where("tenant_id = ?", tenantID).Count(&count)
		if int(count) >= subWithPlan.Plan.MaxProducts {
			return nil, subscription.ErrPlanLimitReached
		}
	}

	p := &Product{
		TenantID:  tenantID,
		Name:      in.Name,
		SKU:       in.SKU,
		CostPrice: in.CostPrice,
		SellPrice: in.SellPrice,
		Status:    "active",
	}
	if in.CategoryID != nil && *in.CategoryID != "" {
		p.CategoryID = in.CategoryID
	}
	if in.Status != "" {
		p.Status = in.Status
	}
	if err := CreateProduct(s.db, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) ListProducts(tenantID string) ([]Product, error) {
	return ListProductsByTenant(s.db, tenantID)
}

func (s *Service) GetProduct(tenantID, productID string) (*Product, error) {
	p, err := FindProductByID(s.db, tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *Service) AddBarcode(tenantID, productID string, barcode string) (*ProductBarcode, error) {
	_, err := FindProductByID(s.db, tenantID, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	var existing ProductBarcode
	if err := s.db.Where("barcode = ?", barcode).First(&existing).Error; err == nil {
		return nil, ErrBarcodeExists
	}

	pb := &ProductBarcode{
		ProductID: productID,
		Barcode:   barcode,
	}
	if err := AddBarcodeToProduct(s.db, pb); err != nil {
		return nil, err
	}
	return pb, nil
}

type CreateProductInput struct {
	Name       string   `json:"name" binding:"required"`
	SKU        string   `json:"sku"`
	CategoryID *string  `json:"category_id"`
	CostPrice  float64  `json:"cost_price"`
	SellPrice  float64  `json:"sell_price" binding:"required"`
	Status     string   `json:"status"`
}

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

type AddBarcodeInput struct {
	Barcode string `json:"barcode" binding:"required"`
}
