package purchase

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
	"gorm.io/gorm"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Create handles POST /api/purchases
func (h *Handler) Create(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	var in CreatePurchaseInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := h.svc.CreatePurchase(tenantID, in)
	if err != nil {
		switch err {
		case ErrInvalidItems, ErrProductInvalid, ErrInvalidQtyCost, ErrDuplicateInvoice:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create purchase"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":             p.ID,
		"supplier_name":  p.SupplierName,
		"invoice_number": p.InvoiceNumber,
		"notes":          p.Notes,
		"total_amount":   p.TotalAmount,
		"created_at":     p.CreatedAt,
	})
}

// List handles GET /api/purchases
func (h *Handler) List(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}

	list, err := h.svc.ListPurchases(tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list purchases"})
		return
	}

	rows := make([]gin.H, 0, len(list))
	for _, p := range list {
		rows = append(rows, gin.H{
			"id":             p.ID,
			"supplier_name":  p.SupplierName,
			"invoice_number": p.InvoiceNumber,
			"notes":          p.Notes,
			"product_names":  p.ProductNames,
			"total_amount":   p.TotalAmount,
			"created_at":     p.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"purchases": rows})
}

// GetByID handles GET /api/purchases/:id
func (h *Handler) GetByID(c *gin.Context) {
	tenantID, ok := utils.GetTenantID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "purchase id required"})
		return
	}

	p, items, err := h.svc.GetPurchaseWithItemRows(tenantID, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "purchase not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get purchase"})
		return
	}

	itemList := make([]gin.H, 0, len(items))
	for _, it := range items {
		itemList = append(itemList, gin.H{
			"product_name": it.ProductName,
			"quantity":     it.Quantity,
			"cost_price":   it.CostPrice,
			"subtotal":     it.Subtotal,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"purchase": gin.H{
			"id":             p.ID,
			"supplier_name":  p.SupplierName,
			"invoice_number": p.InvoiceNumber,
			"notes":          p.Notes,
			"total_amount":   p.TotalAmount,
			"created_at":     p.CreatedAt,
		},
		"items": itemList,
	})
}
