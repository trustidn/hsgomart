package tenant

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db        *gorm.DB
	uploadDir string
}

func NewService(db *gorm.DB) *Service {
	dir := "./uploads/logos"
	os.MkdirAll(dir, 0755)
	return &Service{db: db, uploadDir: dir}
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
