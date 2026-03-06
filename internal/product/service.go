package product

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrBarcodeExists  = errors.New("barcode already registered")
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
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

func (s *Service) CreateProduct(tenantID string, in CreateProductInput) (*Product, error) {
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
