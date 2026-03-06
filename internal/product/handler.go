package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/internal/subscription"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListCategories(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	list, err := h.service.ListCategories(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list categories"})
		return
	}

	c.JSON(http.StatusOK, list)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	var in CreateCategoryInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cat, err := h.service.CreateCategory(tenantID, in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, cat)
}

// productResponse for GET /api/products list (id, name, sku, sell_price, status)
type productResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	SKU       string  `json:"sku"`
	SellPrice float64 `json:"sell_price"`
	Status    string  `json:"status"`
}

func toProductResponse(p *Product) productResponse {
	return productResponse{
		ID:        p.ID,
		Name:      p.Name,
		SKU:       p.SKU,
		SellPrice: p.SellPrice,
		Status:    p.Status,
	}
}

func (h *Handler) ListProducts(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	list, err := h.service.ListProducts(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	res := make([]productResponse, 0, len(list))
	for i := range list {
		res = append(res, toProductResponse(&list[i]))
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateProduct(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	var in CreateProductInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.service.CreateProduct(tenantID, in)
	if err != nil {
		if err == subscription.ErrPlanLimitReached {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "plan limit reached"})
			return
		}
		if err == subscription.ErrSubscriptionRequired || err == subscription.ErrSubscriptionExpired {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, toProductResponse(p))
}

func (h *Handler) GetProduct(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id required"})
		return
	}

	p, err := h.service.GetProduct(tenantID, productID)
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}

	c.JSON(http.StatusOK, toProductResponse(p))
}

func (h *Handler) AddBarcode(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id required"})
		return
	}

	var in AddBarcodeInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pb, err := h.service.AddBarcode(tenantID, productID, in.Barcode)
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		if err == ErrBarcodeExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add barcode"})
		return
	}

	c.JSON(http.StatusCreated, pb)
}
