package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
)

type KnowledgeHandler struct {
	memoryManager memory.MemoryManager
	sessionRepo   workflow.SessionRepository
}

func NewKnowledgeHandler(memoryManager memory.MemoryManager, sessionRepo workflow.SessionRepository) *KnowledgeHandler {
	return &KnowledgeHandler{
		memoryManager: memoryManager,
		sessionRepo:   sessionRepo,
	}
}

// KnowledgeItem represents a knowledge item displayed in the UI
type KnowledgeItem struct {
	ID              string    `json:"knowledge_uuid"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
	Content         string    `json:"content"`
	MemoryLayer     string    `json:"memory_layer"`
	RelevanceScore  int       `json:"relevance_score"`
	SourceMessageID string    `json:"source_message_uuid,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

// KnowledgeResponse is the API response for knowledge list
type KnowledgeResponse struct {
	Items  []KnowledgeItem `json:"items"`
	Total  int             `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
}

// GetSessionKnowledge handles GET /api/v1/sessions/:sessionID/knowledge
func (h *KnowledgeHandler) GetSessionKnowledge(c *gin.Context) {
	sessionID := c.Param("sessionID")

	// Parse query parameters
	memoryLayer := c.DefaultQuery("layer", "all")
	searchQuery := c.DefaultQuery("q", "")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Retrieve knowledge from memory service
	// We need GroupID for the session
	var groupID string
	// Strategy: check active engines first
	if engine := activeEngines[sessionID]; engine != nil {
		groupID, _ = engine.Session.Inputs["group_id"].(string)
	}
	// If not found (or inactive), check DB
	if groupID == "" && h.sessionRepo != nil {
		sEntity, err := h.sessionRepo.Get(c.Request.Context(), sessionID)
		if err == nil {
			groupID = sEntity.GroupID
		}
	}

	rawItems, err := h.memoryManager.Retrieve(c.Request.Context(), searchQuery, groupID, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Map raw memory items to KnowledgeItem DTO
	var items []KnowledgeItem
	for i, ri := range rawItems {
		items = append(items, KnowledgeItem{
			ID:             strconv.Itoa(i + 1),
			Title:          "Memory Fragment",
			Content:        ri.Content,
			Summary:        ri.Content, // Use content as summary for now
			MemoryLayer:    ri.Source,
			RelevanceScore: int(ri.Score * 5), // Map 0-1 to 1-5
			CreatedAt:      time.Now(),
		})
	}

	// Filter by layer if not "all"
	var filtered []KnowledgeItem
	if memoryLayer != "all" {
		for _, it := range items {
			if it.MemoryLayer == memoryLayer {
				filtered = append(filtered, it)
			}
		}
	} else {
		filtered = items
	}
	items = filtered

	// Apply pagination
	start := offset
	if start > len(items) {
		start = len(items)
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	pagedItems := items[start:end]

	c.JSON(http.StatusOK, KnowledgeResponse{
		Items:  pagedItems,
		Total:  len(items),
		Limit:  limit,
		Offset: offset,
	})
}
