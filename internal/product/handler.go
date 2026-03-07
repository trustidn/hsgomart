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

func (h *Handler) UpdateCategory(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id required"})
		return
	}
	var in struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cat, err := h.service.UpdateCategory(tenantID, categoryID, in.Name)
	if err != nil {
		if err == ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update category"})
		return
	}
	c.JSON(http.StatusOK, cat)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	categoryID := c.Param("id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category id required"})
		return
	}
	if err := h.service.DeleteCategory(tenantID, categoryID); err != nil {
		if err == ErrCategoryNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete category"})
		return
	}
	c.Status(http.StatusNoContent)
}

// productResponse for GET /api/products list and GET /api/products/:id (includes category_id, cost_price, barcode for edit).
type productResponse struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	SKU               string   `json:"sku"`
	CategoryID        *string  `json:"category_id,omitempty"`
	CostPrice         float64  `json:"cost_price"`
	SellPrice         float64  `json:"sell_price"`
	Unit              string   `json:"unit"`
	Status            string   `json:"status"`
	Barcode           string   `json:"barcode,omitempty"`
	LowStockThreshold int      `json:"low_stock_threshold"`
}

func toProductResponse(p *Product) productResponse {
	return toProductResponseWithBarcode(p, "")
}

func toProductResponseWithBarcode(p *Product, barcode string) productResponse {
	unit := p.Unit
	if unit == "" {
		unit = "pcs"
	}
	return productResponse{
		ID:                p.ID,
		Name:              p.Name,
		SKU:               p.SKU,
		CategoryID:        p.CategoryID,
		CostPrice:         p.CostPrice,
		SellPrice:         p.SellPrice,
		Unit:              unit,
		Status:            p.Status,
		Barcode:           barcode,
		LowStockThreshold: p.LowStockThreshold,
	}
}

// ProductByBarcodeResponse for GET /api/products/barcode/:barcode
type ProductByBarcodeResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	SellPrice float64 `json:"sell_price"`
}

func (h *Handler) GetProductByBarcode(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	barcode := c.Param("barcode")
	if barcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "barcode required"})
		return
	}
	p, err := h.service.FindProductByBarcode(tenantID, barcode)
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product"})
		return
	}
	c.JSON(http.StatusOK, ProductByBarcodeResponse{
		ID:        p.ID,
		Name:      p.Name,
		SellPrice: p.SellPrice,
	})
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
	barcodes, _ := h.service.GetProductBarcodes(tenantID, productID)
	firstBarcode := ""
	if len(barcodes) > 0 {
		firstBarcode = barcodes[0]
	}
	c.JSON(http.StatusOK, toProductResponseWithBarcode(p, firstBarcode))
}

func (h *Handler) UpdateProduct(c *gin.Context) {
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
	var in UpdateProductInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	catID := in.CategoryID
	if in.CategoryID != nil && *in.CategoryID == "" {
		catID = nil
	}
	p, err := h.service.UpdateProduct(tenantID, productID, UpdateProductInput{
		Name:              in.Name,
		SKU:               in.SKU,
		CategoryID:        catID,
		CostPrice:         in.CostPrice,
		SellPrice:         in.SellPrice,
		Unit:              in.Unit,
		Status:            in.Status,
		LowStockThreshold: in.LowStockThreshold,
	})
	if err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product"})
		return
	}
	c.JSON(http.StatusOK, toProductResponse(p))
}

func (h *Handler) DeleteProduct(c *gin.Context) {
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
	if err := h.service.DeleteProduct(tenantID, productID); err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete product"})
		return
	}
	c.Status(http.StatusNoContent)
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

func (h *Handler) DeleteBarcode(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	productID := c.Param("id")
	barcode := c.Param("barcode")
	if productID == "" || barcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id and barcode required"})
		return
	}
	if err := h.service.DeleteBarcode(tenantID, productID, barcode); err != nil {
		if err == ErrProductNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "barcode not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete barcode"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListBarcodes(c *gin.Context) {
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
	barcodes, err := h.service.GetProductBarcodes(tenantID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list barcodes"})
		return
	}
	c.JSON(http.StatusOK, barcodes)
}
