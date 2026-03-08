package saas

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
	dir := "./uploads/saas"
	os.MkdirAll(dir, 0755)
	return &Service{db: db, uploadDir: dir}
}

func (s *Service) GetSettings() (*Settings, error) {
	var settings Settings
	if err := s.db.First(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

type UpdateInput struct {
	SaasName     *string `json:"saas_name"`
	Tagline      *string `json:"tagline"`
	BankName     *string `json:"bank_name"`
	BankAccount  *string `json:"bank_account"`
	BankHolder   *string `json:"bank_holder"`
	ContactEmail   *string `json:"contact_email"`
	ContactPhone   *string `json:"contact_phone"`
	WhatsappNumber *string `json:"whatsapp_number"`
}

func (s *Service) UpdateSettings(in UpdateInput) (*Settings, error) {
	updates := map[string]interface{}{}
	if in.SaasName != nil {
		updates["saas_name"] = strings.TrimSpace(*in.SaasName)
	}
	if in.Tagline != nil {
		updates["tagline"] = strings.TrimSpace(*in.Tagline)
	}
	if in.BankName != nil {
		updates["bank_name"] = strings.TrimSpace(*in.BankName)
	}
	if in.BankAccount != nil {
		updates["bank_account"] = strings.TrimSpace(*in.BankAccount)
	}
	if in.BankHolder != nil {
		updates["bank_holder"] = strings.TrimSpace(*in.BankHolder)
	}
	if in.ContactEmail != nil {
		updates["contact_email"] = strings.TrimSpace(*in.ContactEmail)
	}
	if in.ContactPhone != nil {
		updates["contact_phone"] = strings.TrimSpace(*in.ContactPhone)
	}
	if in.WhatsappNumber != nil {
		updates["whatsapp_number"] = strings.TrimSpace(*in.WhatsappNumber)
	}
	if len(updates) > 0 {
		if err := s.db.Model(&Settings{}).Where("id = 1").Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetSettings()
}

func (s *Service) UploadLogo(header *multipart.FileHeader) (string, error) {
	if header.Size > 2*1024*1024 {
		return "", fmt.Errorf("file too large, max 2MB")
	}
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".png": true, ".jpg": true, ".jpeg": true, ".webp": true}
	if !allowed[ext] {
		return "", fmt.Errorf("invalid file type")
	}

	src, err := header.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := fmt.Sprintf("saas_logo_%d%s", time.Now().UnixMilli(), ext)
	dstPath := filepath.Join(s.uploadDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	logoURL := "/uploads/saas/" + filename
	s.db.Model(&Settings{}).Where("id = 1").Update("logo_url", logoURL)
	return logoURL, nil
}
