package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(c *gin.Context) {
	var in RegisterInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.service.Register(in)
	if err != nil {
		if err == ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err == utils.ErrWeakPassword {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed"})
		return
	}

	c.JSON(http.StatusCreated, TokenResponse{Token: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (h *Handler) Login(c *gin.Context) {
	var in LoginInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.service.Login(in)
	if err != nil {
		if err == ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (h *Handler) Refresh(c *gin.Context) {
	var in RefreshInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pair, err := h.service.Refresh(in.RefreshToken)
	if err != nil {
		if err == ErrInvalidRefreshToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "refresh failed"})
		return
	}

	c.JSON(http.StatusOK, TokenResponse{Token: pair.AccessToken, RefreshToken: pair.RefreshToken})
}

func (h *Handler) Profile(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	tenantID, _ := utils.GetTenantID(c)
	role, _ := utils.GetUserRole(c)

	email, err := h.service.UserEmail(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":   userID,
		"tenant_id": tenantID,
		"email":     email,
		"role":      role,
	})
}
