package order

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/trustidn/hsmart-saas/internal/subscription"
	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	uploadDir string
}

func NewService(db *gorm.DB) *Service {
	dir := "./uploads/payments"
	os.MkdirAll(dir, 0755)
	return &Service{db: db, uploadDir: dir}
}

func (s *Service) CreateOrder(tenantID string, planID int) (*SubscriptionOrder, error) {
	var plan subscription.Plan
	if err := s.db.Where("id = ?", planID).First(&plan).Error; err != nil {
		return nil, fmt.Errorf("plan not found")
	}

	// Check for existing pending order for the same plan
	var existing SubscriptionOrder
	err := s.db.Where("tenant_id = ? AND plan_id = ? AND status IN ?", tenantID, planID,
		[]string{"pending_payment", "pending_review"}).First(&existing).Error
	if err == nil {
		return nil, fmt.Errorf("you already have a pending order for this plan")
	}

	invoiceNumber := fmt.Sprintf("INV-%s-%04d", time.Now().Format("20060102"), time.Now().UnixMilli()%10000)

	order := SubscriptionOrder{
		TenantID:      tenantID,
		PlanID:        planID,
		Amount:        plan.Price,
		Status:        "pending_payment",
		InvoiceNumber: invoiceNumber,
	}
	if err := s.db.Create(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

var allowedProofExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".pdf": true,
}

func (s *Service) UploadPaymentProof(tenantID, orderID string, header *multipart.FileHeader) (*SubscriptionOrder, error) {
	var order SubscriptionOrder
	if err := s.db.Where("id = ? AND tenant_id = ?", orderID, tenantID).First(&order).Error; err != nil {
		return nil, fmt.Errorf("order not found")
	}
	if order.Status != "pending_payment" && order.Status != "rejected" {
		return nil, fmt.Errorf("cannot upload proof for order with status: %s", order.Status)
	}

	if header.Size > 5*1024*1024 {
		return nil, fmt.Errorf("file too large, max 5MB")
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedProofExts[ext] {
		return nil, fmt.Errorf("invalid file type")
	}

	src, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	filename := fmt.Sprintf("%s_%s_%d%s", tenantID, orderID[:8], time.Now().UnixMilli(), ext)
	dstPath := filepath.Join(s.uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"payment_proof_url": "/uploads/payments/" + filename,
		"status":            "pending_review",
		"paid_at":           now,
	}
	if err := s.db.Model(&order).Updates(updates).Error; err != nil {
		return nil, err
	}

	order.PaymentProofURL = "/uploads/payments/" + filename
	order.Status = "pending_review"
	order.PaidAt = &now
	return &order, nil
}

func (s *Service) ListOrders(tenantID string) ([]SubscriptionOrder, error) {
	var orders []SubscriptionOrder
	err := s.db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (s *Service) GetOrder(tenantID, orderID string) (*SubscriptionOrder, error) {
	var order SubscriptionOrder
	if err := s.db.Where("id = ? AND tenant_id = ?", orderID, tenantID).First(&order).Error; err != nil {
		return nil, err
	}
	return &order, nil
}
