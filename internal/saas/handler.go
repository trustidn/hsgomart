package saas

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetSettings(c *gin.Context) {
	s, err := h.svc.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get settings"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *Handler) UpdateSettings(c *gin.Context) {
	var in UpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	s, err := h.svc.UpdateSettings(in)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update settings"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func (h *Handler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "logo file required"})
		return
	}
	url, err := h.svc.UploadLogo(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"logo_url": url})
}

// PublicInfo returns SaaS name, logo, and bank info (no auth required)
func (h *Handler) PublicInfo(c *gin.Context) {
	s, err := h.svc.GetSettings()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"saas_name": "HSMart POS"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"saas_name":    s.SaasName,
		"logo_url":     s.LogoURL,
		"tagline":      s.Tagline,
		"bank_name":    s.BankName,
		"bank_account": s.BankAccount,
		"bank_holder":  s.BankHolder,
	})
}
