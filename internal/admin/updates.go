package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type updateRow struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

func (h *Handler) ListUpdates(c *gin.Context) {
	rows := make([]updateRow, 0)
	err := h.db.Raw(`
		SELECT id, title, content, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM platform_updates ORDER BY created_at DESC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list updates"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) ListRecentUpdates(c *gin.Context) {
	rows := make([]updateRow, 0)
	err := h.db.Raw(`
		SELECT id, title, content, TO_CHAR(created_at, 'YYYY-MM-DD HH24:MI') AS created_at
		FROM platform_updates ORDER BY created_at DESC LIMIT 5
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list updates"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

type createUpdateInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *Handler) CreateUpdate(c *gin.Context) {
	var in createUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var id int
	err := h.db.Raw(`INSERT INTO platform_updates (title, content) VALUES (?, ?) RETURNING id`, in.Title, in.Content).Scan(&id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create update"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "created"})
}

type editUpdateInput struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func (h *Handler) EditUpdate(c *gin.Context) {
	id := c.Param("id")
	var in editUpdateInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Content != nil {
		updates["content"] = *in.Content
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	res := h.db.Table("platform_updates").Where("id = ?", id).Updates(updates)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "update not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteUpdate(c *gin.Context) {
	id := c.Param("id")
	res := h.db.Exec("DELETE FROM platform_updates WHERE id = ?", id)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "update not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
