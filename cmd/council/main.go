package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/api/handler"
	"github.com/hrygo/council/internal/api/ws"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/infrastructure/cache"
	"github.com/hrygo/council/internal/infrastructure/db"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/persistence"
	"github.com/hrygo/council/internal/pkg/config"
)

func main() {
	fmt.Println("The Council Backend is starting...")

	cfg := config.Load()

	// Initialize Database
	if err := db.Init(context.Background(), cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis
	// Assuming no password and default DB 0 for now as per config
	if err := cache.Init(cfg.RedisURL, "", 0); err != nil {
		log.Printf("Warning: Redis initialization failed: %v", err)
		// Proceeding without Redis? or fail? Implementation plan implies we need it.
		// For MVP, maybe log fatal? TDD says "Strict Three-Tier Memory", so Redis is likely required.
		// Let's log fatal to be safe and ensure infra is up.
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer cache.Close()

	r := gin.Default()

	pool := db.GetPool()

	// Repositories
	groupRepo := persistence.NewGroupRepository(pool)
	agentRepo := persistence.NewAgentRepository(pool)
	workflowRepo := persistence.NewWorkflowRepository(pool)

	// LLM Provider
	llmRouter := llm.NewRouter()
	// Map config.LLMConfig to llm.LLMConfig
	llmCfg := llm.LLMConfig{
		Type:    cfg.LLM.Provider,
		APIKey:  cfg.LLM.APIKey,
		BaseURL: cfg.LLM.BaseURL,
		Model:   cfg.LLM.Model,
	}
	llmProvider, err := llmRouter.GetLLMProvider(llmCfg)
	if err != nil {
		log.Fatalf("Failed to initialize LLM provider: %v", err)
	}

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Services
	embedder, ok := llmProvider.(llm.Embedder)
	if !ok {
		log.Fatalf("Selected LLM provider does not support embeddings")
	}
	memoryService := memory.NewService(embedder, pool, cache.GetClient())

	// Handlers
	groupHandler := handler.NewGroupHandler(groupRepo)
	agentHandler := handler.NewAgentHandler(agentRepo)
	workflowHandler := handler.NewWorkflowHandler(hub, agentRepo, llmProvider)
	workflowMgmtHandler := handler.NewWorkflowMgmtHandler(workflowRepo, llmProvider)
	templateRepo := persistence.NewTemplateRepository(pool)
	templateHandler := handler.NewTemplateHandler(templateRepo)
	memoryHandler := handler.NewMemoryHandler(memoryService)

	// Routes
	r.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	api := r.Group("/api/v1")
	{
		// Groups
		api.POST("/groups", groupHandler.Create)
		api.GET("/groups", groupHandler.List)
		api.GET("/groups/:id", groupHandler.Get)
		api.PUT("/groups/:id", groupHandler.Update)
		api.DELETE("/groups/:id", groupHandler.Delete)

		// Agents
		api.POST("/agents", agentHandler.Create)
		api.GET("/agents", agentHandler.List)
		api.GET("/agents/:id", agentHandler.Get)
		api.PUT("/agents/:id", agentHandler.Update)
		api.DELETE("/agents/:id", agentHandler.Delete)

		// Workflows
		api.POST("/workflows/execute", workflowHandler.Execute)
		api.POST("/sessions/:id/control", workflowHandler.Control)
		api.POST("/sessions/:id/signal", workflowHandler.Signal)
		api.POST("/sessions/:id/review", workflowHandler.Review)

		// Workflow Management
		api.GET("/workflows", workflowMgmtHandler.List)
		api.GET("/workflows/:id", workflowMgmtHandler.Get)
		api.POST("/workflows", workflowMgmtHandler.Create)
		api.PUT("/workflows/:id", workflowMgmtHandler.Update)
		api.POST("/workflows/generate", workflowMgmtHandler.Generate)
		api.POST("/workflows/estimate", workflowMgmtHandler.EstimateCost)

		// Templates
		api.GET("/templates", templateHandler.List)
		api.POST("/templates", templateHandler.Create)
		api.DELETE("/templates/:id", templateHandler.Delete)

		// Memory
		api.POST("/memory/ingest", memoryHandler.Ingest)
		api.POST("/memory/query", memoryHandler.Query)
	}

	fmt.Printf("Server listening on :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
