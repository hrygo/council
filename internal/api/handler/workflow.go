package handler

import (
	"context"
	"log"
	"net/http"
	"time"

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
	Hub           *ws.Hub
	AgentRepo     agent.Repository
	Registry      *llm.Registry
	MemoryManager memory.MemoryManager
	SessionRepo   workflow.SessionRepository
	FileRepo      workflow.SessionFileRepository
	WorkflowRepo  workflow.Repository
}

func NewWorkflowHandler(hub *ws.Hub, agentRepo agent.Repository, registry *llm.Registry, memManager memory.MemoryManager, sessionRepo workflow.SessionRepository, fileRepo workflow.SessionFileRepository, workflowRepo workflow.Repository) *WorkflowHandler {
	return &WorkflowHandler{
		Hub:           hub,
		AgentRepo:     agentRepo,
		Registry:      registry,
		MemoryManager: memManager,
		SessionRepo:   sessionRepo,
		FileRepo:      fileRepo,
		WorkflowRepo:  workflowRepo,
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
	session.SetFileRepository(h.FileRepo)

	// Persist Session
	groupID, _ := req.Input["group_uuid"].(string)
	if groupID == "" {
		groupID, _ = req.Input["group_id"].(string) // Backward compatibility
		if groupID != "" {
			log.Printf("[Workflow] Warning: Client used deprecated 'group_id' parameter. Please update to 'group_uuid'.")
		}
	}
	workflowID := ""
	if req.Graph != nil {
		workflowID = req.Graph.ID
		// Auto-persist dynamic workflow if it doesn't exist
		// This prevents 404 errors when frontend tries to fetch workflow details later.
		if workflowID != "" && h.WorkflowRepo != nil {
			if _, err := h.WorkflowRepo.Get(c.Request.Context(), workflowID); err != nil {
				// Assume duplicate/exists or not found. If error, try to create.
				// A simple way is to just try Create and ignore specific "already exists" error,
				// but Get() error usually means not found or DB error.
				// Let's try Create.
				if req.Graph.Name == "" {
					req.Graph.Name = "Dynamic Workflow " + workflowID[:8]
				}
				if err := h.WorkflowRepo.Create(c.Request.Context(), req.Graph); err != nil {
					log.Printf("[Workflow] Warning: Failed to auto-persist dynamic workflow %s: %v", workflowID, err)
				} else {
					log.Printf("[Workflow] Auto-persisted dynamic workflow %s", workflowID)
				}
			}
		}
	}
	session.Start(context.Background())

	if err := h.SessionRepo.Create(c.Request.Context(), session, groupID, workflowID); err != nil {
		log.Printf("[Workflow] Failed to persist session: %v", err)
		// We continue anyway for MVP but ideally fail here
	}

	// Create Engine
	engine := workflow.NewEngine(session)
	activeEngines[session.ID] = engine

	// Inject CouncilMergeStrategy for Council workflows (SPEC-1206)
	// This aggregates agent_output from parallel branches into aggregated_outputs
	engine.MergeStrategy = &workflow.CouncilMergeStrategy{}

	// Configure Factory
	engine.NodeFactory = nodes.NewNodeFactory(nodes.NodeDependencies{
		Registry:      h.Registry,
		AgentRepo:     h.AgentRepo,
		MemoryManager: h.MemoryManager,
		Session:       session,
	})

	// First, create memService as it's a dependency for NodeDependencies now.
	// Note: We use global getters here for simplicity, but ideally these would be in WorkflowHandler
	engine.Middlewares = []workflow.Middleware{
		middleware.NewCircuitBreaker(10),                // Logic Circuit Breaker (Depth > 10)
		middleware.NewFactCheckTrigger(),                // Anti-Hallucination
		middleware.NewMemoryMiddleware(h.MemoryManager), // Memory Persistence
	}

	// Run in Goroutine
	go func() {
		log.Printf("[Workflow] Starting execution for session %s", session.ID)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[Workflow] PANIC in session %s: %v", session.ID, r)
			}
			session.Complete()
			log.Printf("[Workflow] Session %s completed", session.ID)
		}()

		// Bridge Engine Stream -> WS Hub
		// We need to modify Engine to allow tapping or we just read from the stream channel
		// Engine exposes StreamChannel
		go func() {
			for event := range engine.StreamChannel {
				// Augment event with SessionID?
				event.Data["session_uuid"] = session.ID

				h.Hub.Broadcast(event)
			}
		}()

		engine.Run(session.Context())

		// Emit completion event
		engine.StreamChannel <- workflow.StreamEvent{
			Type:      "execution:completed",
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"status": "completed"},
		}

		close(engine.StreamChannel)
	}()

	c.JSON(http.StatusAccepted, gin.H{
		"session_uuid": session.ID,
		"status":       "started",
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

	engine := h.getEngine(id) // Helper method
	if engine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found or not active"})
		return
	}
	session := engine.Session

	switch req.Action {
	case "pause":
		session.Pause()
	case "resume":
		session.Resume()
	case "stop":
		session.Stop()
	}

	c.JSON(http.StatusOK, gin.H{
		"session_uuid": id,
		"status":       session.Status,
		"action":       req.Action,
	})
}

// Global active engine map (Protected by mutex in real impl)
var activeEngines = map[string]*workflow.Engine{}

func (h *WorkflowHandler) getEngine(sessionID string) *workflow.Engine {
	return activeEngines[sessionID]
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

	engine := h.getEngine(id)
	if engine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}
	session := engine.Session

	if err := session.SendSignal(req.NodeID, req.Payload); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "signal_sent"})
}

type ReviewRequest struct {
	NodeID string                 `json:"node_id" binding:"required"`
	Action string                 `json:"action" binding:"required,oneof=approve reject modify"`
	Data   map[string]interface{} `json:"data"` // Optional patch data
}

func (h *WorkflowHandler) Review(c *gin.Context) {
	id := c.Param("id")
	var req ReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	engine := h.getEngine(id)
	if engine == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found or not active"})
		return
	}

	// Logic for Human Review Resume
	// We need to fetch the node to confirm it is suspended?
	// The Engine.ResumeNode checks this.

	// Construct output payload based on Action
	output := map[string]interface{}{
		"review_action": req.Action,
		"reviewer":      "human", // Placeholder
		"timestamp":     c.GetHeader("Date"),
	}
	if req.Data != nil {
		for k, v := range req.Data {
			output[k] = v
		}
	}

	if err := engine.ResumeNode(c.Request.Context(), req.NodeID, output); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "resumed"})
}

func (h *WorkflowHandler) GetSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.SessionRepo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}
