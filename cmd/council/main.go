package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/api/handler"
	"github.com/hrygo/council/internal/api/ws"
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

	r := gin.Default()

	pool := db.GetPool()

	// Repositories
	groupRepo := persistence.NewGroupRepository(pool)
	agentRepo := persistence.NewAgentRepository(pool)

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

	// Handlers
	groupHandler := handler.NewGroupHandler(groupRepo)
	agentHandler := handler.NewAgentHandler(agentRepo)
	workflowHandler := handler.NewWorkflowHandler(hub, agentRepo, llmProvider)
	workflowMgmtHandler := handler.NewWorkflowMgmtHandler()

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

		// Workflow Management
		api.GET("/workflows", workflowMgmtHandler.List)
		api.GET("/workflows/:id", workflowMgmtHandler.Get)
		api.POST("/workflows", workflowMgmtHandler.Create)
		api.PUT("/workflows/:id", workflowMgmtHandler.Update)
		api.POST("/workflows/generate", workflowMgmtHandler.Generate)
	}

	fmt.Printf("Server listening on :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
