package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type docRow struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	SortOrder   int    `json:"sort_order"`
	IsPublished bool   `json:"is_published"`
	UpdatedAt   string `json:"updated_at"`
}

func (h *Handler) ListDocumentation(c *gin.Context) {
	rows := make([]docRow, 0)
	err := h.db.Raw(`
		SELECT id, title, content, sort_order, is_published,
		       TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at
		FROM documentation ORDER BY sort_order ASC, id ASC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list documentation"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

func (h *Handler) ListPublishedDocumentation(c *gin.Context) {
	rows := make([]docRow, 0)
	err := h.db.Raw(`
		SELECT id, title, content, sort_order, is_published,
		       TO_CHAR(updated_at, 'YYYY-MM-DD HH24:MI') AS updated_at
		FROM documentation WHERE is_published = true ORDER BY sort_order ASC, id ASC
	`).Scan(&rows).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list documentation"})
		return
	}
	c.JSON(http.StatusOK, rows)
}

type createDocInput struct {
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	SortOrder   int    `json:"sort_order"`
	IsPublished *bool  `json:"is_published"`
}

func (h *Handler) CreateDocumentation(c *gin.Context) {
	var in createDocInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	published := true
	if in.IsPublished != nil {
		published = *in.IsPublished
	}
	var id int
	err := h.db.Raw(`
		INSERT INTO documentation (title, content, sort_order, is_published, updated_at)
		VALUES (?, ?, ?, ?, ?) RETURNING id
	`, in.Title, in.Content, in.SortOrder, published, time.Now()).Scan(&id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create documentation"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "created"})
}

type updateDocInput struct {
	Title       *string `json:"title"`
	Content     *string `json:"content"`
	SortOrder   *int    `json:"sort_order"`
	IsPublished *bool   `json:"is_published"`
}

func (h *Handler) UpdateDocumentation(c *gin.Context) {
	id := c.Param("id")
	var in updateDocInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]interface{}{"updated_at": time.Now()}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Content != nil {
		updates["content"] = *in.Content
	}
	if in.SortOrder != nil {
		updates["sort_order"] = *in.SortOrder
	}
	if in.IsPublished != nil {
		updates["is_published"] = *in.IsPublished
	}
	res := h.db.Table("documentation").Where("id = ?", id).Updates(updates)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "documentation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *Handler) DeleteDocumentation(c *gin.Context) {
	id := c.Param("id")
	res := h.db.Exec("DELETE FROM documentation WHERE id = ?", id)
	if res.Error != nil || res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "documentation not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
