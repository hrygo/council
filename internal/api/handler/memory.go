package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/memory"
)

type MemoryHandler struct {
	Manager memory.MemoryManager
}

func NewMemoryHandler(manager memory.MemoryManager) *MemoryHandler {
	return &MemoryHandler{Manager: manager}
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

	if err := h.Manager.Promote(ctx, req.GroupID, req.Content); err != nil {
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

	sessionID := c.Query("session_id")
	results, err := h.Manager.Retrieve(ctx, req.Query, req.GroupID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"results": results})
}
