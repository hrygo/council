package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/memory"
)

type KnowledgeHandler struct {
	memoryService *memory.Service
}

func NewKnowledgeHandler(memoryService *memory.Service) *KnowledgeHandler {
	return &KnowledgeHandler{
		memoryService: memoryService,
	}
}

// KnowledgeItem represents a knowledge item displayed in the UI
type KnowledgeItem struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
	Content         string    `json:"content"`
	MemoryLayer     string    `json:"memory_layer"`
	RelevanceScore  int       `json:"relevance_score"`
	SourceMessageID string    `json:"source_message_id,omitempty"`
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
	// For MVP, we return mock data
	// TODO: Implement actual retrieval from memory service
	items := h.mockKnowledgeItems(sessionID, memoryLayer, searchQuery)

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

// mockKnowledgeItems returns mock knowledge items for development
// TODO: Replace with actual memory service integration
func (h *KnowledgeHandler) mockKnowledgeItems(sessionID, layer, query string) []KnowledgeItem {
	allItems := []KnowledgeItem{
		{
			ID:             "k1",
			Title:          "工作流执行上下文",
			Summary:        "当前工作流正在执行序列节点，已完成初始化步骤",
			Content:        "详细内容：工作流引擎已启动，当前处于序列执行模式...",
			MemoryLayer:    "sandboxed",
			RelevanceScore: 5,
			CreatedAt:      time.Now().Add(-5 * time.Minute),
		},
		{
			ID:             "k2",
			Title:          "Agent 配置信息",
			Summary:        "Adjudicator agent 使用 gpt-4o 模型",
			Content:        "配置详情：Adjudicator 角色设定、Prompt 模板...",
			MemoryLayer:    "working",
			RelevanceScore: 4,
			CreatedAt:      time.Now().Add(-10 * time.Minute),
		},
		{
			ID:             "k3",
			Title:          "历史辩论记录",
			Summary:        "上次会议讨论了 AI 伦理问题",
			Content:        "辩论摘要：正方提出..., 反方认为...",
			MemoryLayer:    "long-term",
			RelevanceScore: 3,
			CreatedAt:      time.Now().Add(-2 * time.Hour),
		},
		{
			ID:             "k4",
			Title:          "用户输入上下文",
			Summary:        "用户提出关于模型选择策略的问题",
			Content:        "用户问题：如何根据任务复杂度选择合适的 LLM？",
			MemoryLayer:    "sandboxed",
			RelevanceScore: 4,
			CreatedAt:      time.Now().Add(-3 * time.Minute),
		},
		{
			ID:             "k5",
			Title:          "知识库检索结果",
			Summary:        "找到 5 条与当前话题相关的文档",
			Content:        "检索到的文档包括：模型选择指南、成本优化建议...",
			MemoryLayer:    "working",
			RelevanceScore: 5,
			CreatedAt:      time.Now().Add(-1 * time.Minute),
		},
	}

	// Filter by layer
	var filtered []KnowledgeItem
	for _, item := range allItems {
		if layer == "all" || item.MemoryLayer == layer {
			filtered = append(filtered, item)
		}
	}

	// TODO: Implement search filtering based on query
	if query != "" {
		// Simple mock: return all for now
		// Real implementation would filter by title/summary/content
	}

	return filtered
}
