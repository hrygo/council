package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool" // Added import for pgxpool
)

type LLMHandler struct {
	Config *config.Config
	DB     *pgxpool.Pool // Added DB field
}

func NewLLMHandler(cfg *config.Config, db *pgxpool.Pool) *LLMHandler { // Modified signature
	return &LLMHandler{
		Config: cfg,
		DB:     db, // Initialized DB field
	}
}

type ProviderOption struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Icon   string   `json:"icon"`
	Models []string `json:"models"`
}

// GetProviderOptions returns the available LLM providers and their models.
func (h *LLMHandler) GetProviderOptions(c *gin.Context) {
	// 1. Fetch all providers from DB
	rows, err := h.DB.Query(c.Request.Context(),
		"SELECT provider_id, name, icon FROM llm_providers ORDER BY sort_order ASC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch providers"})
		return
	}
	defer rows.Close()

	var dbOptions []ProviderOption
	for rows.Next() {
		var p ProviderOption
		if err := rows.Scan(&p.ID, &p.Name, &p.Icon); err != nil {
			continue
		}
		dbOptions = append(dbOptions, p)
	}

	// 2. Prepare options with filtering logic
	var options []ProviderOption
	activeProviders := make(map[string]bool)

	// Check which valid API keys are heavily confirmed to optimize SiliconFlow
	// (We still list all providers, but use this for the exclusion logic)
	if h.Config.DeepSeekKey != "" {
		activeProviders["deepseek"] = true
	}
	if h.Config.DashScopeKey != "" {
		activeProviders["dashscope"] = true
	}

	for _, p := range dbOptions {
		// Fetch models for this provider (All stored models)
		modelRows, err := h.DB.Query(c.Request.Context(),
			"SELECT model_id FROM llm_models WHERE provider_id = $1 AND is_mainstream = true ORDER BY sort_order ASC", p.ID)
		if err != nil {
			continue
		}

		var models []string
		for modelRows.Next() {
			var mID string
			if err := modelRows.Scan(&mID); err == nil {
				models = append(models, mID)
			}
		}
		modelRows.Close()

		// 3. Apply Filtering Logic for SiliconFlow
		// Rule: If "deepseek" is active (has key), filter out "deepseek-ai/*" from SiliconFlow
		// Rule: If "dashscope" is active, filter out "Qwen/*" from SiliconFlow
		if p.ID == "siliconflow" {
			var filteredModels []string
			for _, m := range models {
				// Filter DeepSeek duplicates
				if activeProviders["deepseek"] && strings.Contains(strings.ToLower(m), "deepseek") {
					continue
				}
				// Filter Qwen duplicates if dashscope is enabled
				if activeProviders["dashscope"] && strings.Contains(strings.ToLower(m), "qwen") {
					continue
				}
				filteredModels = append(filteredModels, m)
			}
			models = filteredModels
		}

		p.Models = models
		options = append(options, p)
	}

	c.JSON(http.StatusOK, gin.H{"providers": options})
}
