package opname

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/trustidn/hsmart-saas/pkg/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Start(c *gin.Context) {
	tenantID, _ := utils.GetTenantID(c)
	userID, _ := utils.GetUserID(c)
	op, err := h.service.StartOpname(tenantID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start opname"})
		return
	}
	c.JSON(http.StatusCreated, op)
}

func (h *Handler) SubmitItems(c *gin.Context) {
	tenantID, _ := utils.GetTenantID(c)
	opnameID := c.Param("id")
	var items []SubmitItemInput
	if err := c.ShouldBindJSON(&items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	op, err := h.service.SubmitItems(tenantID, opnameID, items)
	if err != nil {
		if err == ErrOpnameNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, op)
}

func (h *Handler) Approve(c *gin.Context) {
	tenantID, _ := utils.GetTenantID(c)
	opnameID := c.Param("id")
	op, err := h.service.ApproveOpname(tenantID, opnameID)
	if err != nil {
		if err == ErrOpnameNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, op)
}

func (h *Handler) Delete(c *gin.Context) {
	tenantID, _ := utils.GetTenantID(c)
	opnameID := c.Param("id")
	err := h.service.DeleteOpname(tenantID, opnameID)
	if err != nil {
		if err == ErrOpnameNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == ErrOpnameNotDraft {
			c.JSON(http.StatusBadRequest, gin.H{"error": "hanya opname berstatus draft yang dapat dihapus"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus opname"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "opname deleted"})
}

func (h *Handler) Get(c *gin.Context) {
	tenantID, _ := utils.GetTenantID(c)
	opnameID := c.Param("id")
	op, items, err := h.service.GetOpname(tenantID, opnameID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	type itemWithName struct {
		StockOpnameItem
		ProductName string `json:"product_name"`
	}
	enriched := make([]itemWithName, 0, len(items))
	for _, it := range items {
		name := it.ProductID
		var pName string
		if err := h.service.DB().Raw("SELECT name FROM products WHERE id = ?", it.ProductID).Scan(&pName).Error; err == nil && pName != "" {
			name = pName
		}
		enriched = append(enriched, itemWithName{StockOpnameItem: it, ProductName: name})
	}
	c.JSON(http.StatusOK, gin.H{"opname": op, "items": enriched})
}

func (h *Handler) List(c *gin.Context) {
	tenantID, _ := utils.GetTenantID(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	list, err := h.service.ListOpnames(tenantID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list opnames"})
		return
	}
	c.JSON(http.StatusOK, list)
}
