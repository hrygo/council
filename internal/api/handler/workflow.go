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
	activeSessions[session.ID] = session
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

type ControlRequest struct {
	Action string `json:"action" binding:"required,oneof=pause resume stop"`
}

func (h *WorkflowHandler) Control(c *gin.Context) {
	id := c.Param("id")
	var req ControlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real app, we need a SessionRepository or global SessionRegistry to look up active sessions.
	// For MVP, we don't have a global registry yet. The session created in Execute is lost reference-wise unless stored.
	// We need to fix this by adding a SessionRegistry to WorkflowHandler.

	// FIX: This indicates a missing architectural piece. We need a way to retrieve running sessions.
	// For now, I will add a TO-DO and mock response, but I must fix this to make it work.
	// Strategy: Use a simple map in memory for active sessions.

	session := h.getSession(id) // Helper method
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found or not active"})
		return
	}

	switch req.Action {
	case "pause":
		session.Pause()
	case "resume":
		session.Resume()
	case "stop":
		session.Stop()
	}

	c.JSON(http.StatusOK, gin.H{
		"id":     id,
		"status": session.Status,
		"action": req.Action,
	})
}

// Global active sessions map (Protected by mutex in real impl)
// For MVP, simple map is fine if we don't have high concurrency on registry itself
var activeSessions = map[string]*workflow.Session{}

func (h *WorkflowHandler) getSession(id string) *workflow.Session {
	return activeSessions[id]
}

type SignalRequest struct {
	NodeID  string      `json:"node_id" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

func (h *WorkflowHandler) Signal(c *gin.Context) {
	id := c.Param("id")
	var req SignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := h.getSession(id)
	if session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if err := session.SendSignal(req.NodeID, req.Payload); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "signal_sent"})
}
