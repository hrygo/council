package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/api/ws"
	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/middleware"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type WorkflowHandler struct {
	Hub       *ws.Hub
	AgentRepo agent.Repository
	LLM       llm.LLMProvider
}

func NewWorkflowHandler(hub *ws.Hub, agentRepo agent.Repository, llmProvider llm.LLMProvider) *WorkflowHandler {
	return &WorkflowHandler{
		Hub:       hub,
		AgentRepo: agentRepo,
		LLM:       llmProvider,
	}
}

type ExecuteRequest struct {
	Graph *workflow.GraphDefinition `json:"graph"`
	Input map[string]interface{}    `json:"input"`
}

func (h *WorkflowHandler) Execute(c *gin.Context) {
	var req ExecuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For MVP, we pass explicit Graph definition.
	// Later we can lookup by ID.

	// Create Session
	session := workflow.NewSession(req.Graph, req.Input)
	session.Start(context.Background())

	// Create Engine
	engine := workflow.NewEngine(session)

	// Configure Factory
	engine.NodeFactory = nodes.NewNodeFactory(nodes.NodeDependencies{
		LLM:       h.LLM,
		AgentRepo: h.AgentRepo,
	})

	// Inject Middleware
	// Try to get Embedder from LLM Provider
	var embedder llm.Embedder
	if e, ok := h.LLM.(llm.Embedder); ok {
		embedder = e
	}
	memService := memory.NewService(embedder) // Using default global pools
	engine.Middlewares = []workflow.Middleware{
		middleware.NewCircuitBreaker(10),           // Logic Circuit Breaker (Depth > 10)
		middleware.NewFactCheckTrigger(),           // Anti-Hallucination
		middleware.NewMemoryMiddleware(memService), // Memory Persistence
	}

	// Run in Goroutine
	go func() {
		defer session.Complete()

		// Bridge Engine Stream -> WS Hub
		// We need to modify Engine to allow tapping or we just read from the stream channel
		// Engine exposes StreamChannel
		go func() {
			for event := range engine.StreamChannel {
				// Augment event with SessionID?
				event.Data["session_id"] = session.ID
				h.Hub.Broadcast(event)
			}
		}()

		engine.Run(session.Context())
		close(engine.StreamChannel)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"session_id": session.ID,
		"status":     "started",
	})
}
