package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/memory"
)

type MemoryHandler struct {
	Service *memory.Service
}

func NewMemoryHandler(service *memory.Service) *MemoryHandler {
	return &MemoryHandler{Service: service}
}

type IngestRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *MemoryHandler) Ingest(c *gin.Context) {
	var req IngestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trigger Promotion (Long-Term Memory)
	// In MVP we expose this directly. In production, this might be async.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := h.Service.Promote(ctx, req.GroupID, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ingested"})
}

type QueryRequest struct {
	GroupID string `json:"group_id" binding:"required"`
	Query   string `json:"query" binding:"required"`
}

func (h *MemoryHandler) Query(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	results, err := h.Service.Retrieve(ctx, req.Query, req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
