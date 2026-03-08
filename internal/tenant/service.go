package tenant

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidConfirmationCode = errors.New("invalid confirmation code")
	ErrInvalidPassword         = errors.New("invalid password")
)

const resetConfirmationCode = "RESET-SEMUA-DATA"

type Service struct {
	db        *gorm.DB
	uploadDir string
}

func NewService(db *gorm.DB) *Service {
	dir := "./uploads/logos"
	os.MkdirAll(dir, 0755)
	return &Service{db: db, uploadDir: dir}
}

type ResetDataInput struct {
	ConfirmationCode string `json:"confirmation_code" binding:"required"`
	Password         string `json:"password" binding:"required"`
}

func (s *Service) ResetData(tenantID, userID string, in ResetDataInput) error {
	if in.ConfirmationCode != resetConfirmationCode {
		return ErrInvalidConfirmationCode
	}

	var user struct {
		PasswordHash string
	}
	if err := s.db.Table("users").Select("password_hash").
		Where("id = ? AND tenant_id = ? AND role = 'owner'", userID, tenantID).
		Scan(&user).Error; err != nil || user.PasswordHash == "" {
		return ErrInvalidPassword
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password)); err != nil {
		return ErrInvalidPassword
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		deletes := []string{
			"DELETE FROM payments WHERE transaction_id IN (SELECT id FROM transactions WHERE tenant_id = ?)",
			"DELETE FROM transaction_items WHERE transaction_id IN (SELECT id FROM transactions WHERE tenant_id = ?)",
			"DELETE FROM transactions WHERE tenant_id = ?",
			"DELETE FROM refunds WHERE tenant_id = ?",
			"DELETE FROM stock_movements WHERE tenant_id = ?",
			"DELETE FROM stock_opname_items WHERE opname_id IN (SELECT id FROM stock_opnames WHERE tenant_id = ?)",
			"DELETE FROM stock_opnames WHERE tenant_id = ?",
			"DELETE FROM inventory_batches WHERE product_id IN (SELECT id FROM products WHERE tenant_id = ?)",
			"DELETE FROM purchase_items WHERE purchase_id IN (SELECT id FROM purchases WHERE tenant_id = ?)",
			"DELETE FROM purchases WHERE tenant_id = ?",
			"DELETE FROM inventories WHERE tenant_id = ?",
			"DELETE FROM product_barcodes WHERE product_id IN (SELECT id FROM products WHERE tenant_id = ?)",
			"DELETE FROM products WHERE tenant_id = ?",
			"DELETE FROM categories WHERE tenant_id = ?",
			"DELETE FROM cashier_shifts WHERE tenant_id = ?",
			"DELETE FROM audit_logs WHERE tenant_id = ?",
		}
		for _, q := range deletes {
			if err := tx.Exec(q, tenantID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) GetProfile(tenantID string) (*Tenant, error) {
	var t Tenant
	if err := s.db.Where("id = ?", tenantID).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

type UpdateProfileInput struct {
	Name        *string `json:"name"`
	Phone       *string `json:"phone"`
	Address     *string `json:"address"`
	Description *string `json:"description"`
}

func (s *Service) UpdateProfile(tenantID string, in UpdateProfileInput) (*Tenant, error) {
	updates := map[string]interface{}{}
	if in.Name != nil {
		updates["name"] = strings.TrimSpace(*in.Name)
	}
	if in.Phone != nil {
		updates["phone"] = strings.TrimSpace(*in.Phone)
	}
	if in.Address != nil {
		updates["address"] = strings.TrimSpace(*in.Address)
	}
	if in.Description != nil {
		updates["description"] = strings.TrimSpace(*in.Description)
	}
	if len(updates) == 0 {
		return s.GetProfile(tenantID)
	}

	if err := s.db.Model(&Tenant{}).Where("id = ?", tenantID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetProfile(tenantID)
}

var allowedExts = map[string]bool{
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true,
}

func (s *Service) UploadLogo(tenantID string, header *multipart.FileHeader) (string, error) {
	if header.Size > 2*1024*1024 {
		return "", fmt.Errorf("file too large, max 2MB")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("invalid file type, allowed: png, jpg, jpeg, webp")
	}

	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := fmt.Sprintf("%s_%d%s", tenantID, time.Now().UnixMilli(), ext)
	dstPath := filepath.Join(s.uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	logoURL := "/uploads/logos/" + filename
	if err := s.db.Model(&Tenant{}).Where("id = ?", tenantID).Update("logo_url", logoURL).Error; err != nil {
		return "", err
	}

	return logoURL, nil
}
