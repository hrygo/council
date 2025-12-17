package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/persistence"
)

type TemplateHandler struct {
	Repo *persistence.TemplateRepository
}

func NewTemplateHandler(repo *persistence.TemplateRepository) *TemplateHandler {
	return &TemplateHandler{Repo: repo}
}

func (h *TemplateHandler) List(c *gin.Context) {
	templates, err := h.Repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, templates)
}

func (h *TemplateHandler) Create(c *gin.Context) {
	var req struct {
		Name        string                    `json:"name"`
		Description string                    `json:"description"`
		Category    workflow.TemplateCategory `json:"category"`
		Graph       workflow.GraphDefinition  `json:"graph"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tpl := workflow.Template{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		IsSystem:    false, // User created
		Graph:       req.Graph,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.Repo.Create(c.Request.Context(), &tpl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tpl)
}

func (h *TemplateHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.Repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
