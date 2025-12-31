package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/api/handler"
	"github.com/hrygo/council/internal/api/ws"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/infrastructure/cache"
	"github.com/hrygo/council/internal/infrastructure/db"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/persistence"
	"github.com/hrygo/council/internal/pkg/config"
	"github.com/hrygo/council/internal/resources"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	fmt.Println("The Council Backend is starting...")

	cfg := config.Load()

	// Initialize Database
	if err := db.Init(context.Background(), cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()
	pool := db.GetPool()

	// Seed default data (Agents, Groups, Workflows)
	seeder := resources.NewSeeder(pool, cfg)
	if err := seeder.SeedAll(context.Background()); err != nil {
		log.Printf("Warning: Failed to seed default data: %v", err)
	} else {
		log.Println("Default data seeded successfully")
	}

	// Initialize Redis
	if err := cache.Init(cfg.RedisURL, "", 0); err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}
	defer cache.Close()

	r := gin.Default()

	// LLM Registry
	registry := llm.NewRegistry(cfg)

	// Core Services
	// Map config.EmbeddingConfig to llm.EmbeddingConfig
	embedCfg := llm.EmbeddingConfig{
		Type:    cfg.Embedding.Provider, // Mapped from Provider to Type
		APIKey:  cfg.Embedding.APIKey,
		BaseURL: cfg.Embedding.BaseURL,
		Model:   cfg.Embedding.Model,
	}
	embedder, err := registry.NewEmbedder(embedCfg)
	if err != nil {
		log.Fatalf("Failed to initialize embedder: %v", err)
	}
	memoryService := memory.NewService(embedder, pool, cache.GetClient())

	// WebSocket Hub
	hub := ws.NewHub()
	go hub.Run()

	// Repositories
	groupRepo := persistence.NewGroupRepository(pool)
	agentRepo := persistence.NewAgentRepository(pool)
	workflowRepo := persistence.NewWorkflowRepository(pool)
	templateRepo := persistence.NewTemplateRepository(pool)
	sessionRepo := persistence.NewSessionRepository(pool)
	fileRepo := persistence.NewSessionFileRepository(pool)

	// Handlers
	agentHandler := handler.NewAgentHandler(agentRepo)
	groupHandler := handler.NewGroupHandler(groupRepo)
	templateHandler := handler.NewTemplateHandler(templateRepo)
	memoryHandler := handler.NewMemoryHandler(memoryService)
	knowledgeHandler := handler.NewKnowledgeHandler(memoryService, sessionRepo)
	workflowMgmtHandler := handler.NewWorkflowMgmtHandler(workflowRepo, registry)
	llmHandler := handler.NewLLMHandler(cfg, pool)

	// WorkflowHandler dependency injection
	workflowHandler := handler.NewWorkflowHandler(
		hub,
		agentRepo,
		registry,
		memoryService,
		sessionRepo,
		fileRepo,
		workflowRepo,
	)

	// Routes
	r.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	api := r.Group("/api/v1")
	{
		// Agents
		api.POST("/agents", agentHandler.Create)
		api.GET("/agents", agentHandler.List)
		api.GET("/agents/:id", agentHandler.Get)
		api.PUT("/agents/:id", agentHandler.Update)
		api.DELETE("/agents/:id", agentHandler.Delete)

		// Groups
		api.POST("/groups", groupHandler.Create)
		api.GET("/groups", groupHandler.List)
		api.GET("/groups/:id", groupHandler.Get)
		api.PUT("/groups/:id", groupHandler.Update)
		api.DELETE("/groups/:id", groupHandler.Delete)

		// Workflows Management
		api.GET("/workflows", workflowMgmtHandler.List)
		api.GET("/workflows/:id", workflowMgmtHandler.Get)
		api.POST("/workflows", workflowMgmtHandler.Create)
		api.PUT("/workflows/:id", workflowMgmtHandler.Update)
		api.POST("/workflows/generate", workflowMgmtHandler.Generate)
		api.POST("/workflows/estimate", workflowMgmtHandler.EstimateCost)

		// Workflows Execution
		api.POST("/workflows/execute", workflowHandler.Execute)
		api.GET("/sessions/:id", workflowHandler.GetSession)
		api.POST("/sessions/:id/control", workflowHandler.Control)

		api.POST("/sessions/:id/signal", workflowHandler.Signal)
		api.POST("/sessions/:id/review", workflowHandler.Review)
		api.GET("/sessions/:id/files", workflowHandler.ListFiles)
		api.GET("/sessions/:id/files/history", workflowHandler.GetFileHistory)

		// Templates
		api.GET("/templates", templateHandler.List)
		api.POST("/templates", templateHandler.Create)
		api.DELETE("/templates/:id", templateHandler.Delete)

		// Memory
		api.POST("/memory/ingest", memoryHandler.Ingest)
		api.POST("/memory/query", memoryHandler.Query)

		// Knowledge (Session-specific)
		api.GET("/sessions/:id/knowledge", knowledgeHandler.GetSessionKnowledge)

		// LLM Options
		api.GET("/llm/providers", llmHandler.GetProviderOptions)
	}

	// Graceful Shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		fmt.Printf("Server listening on :%s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
