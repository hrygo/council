package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/workflow"
)

// WorkflowMgmtHandler handles CRUD operations for workflows (templates and user graphs)
type WorkflowMgmtHandler struct {
	// We need a repository. For MVP we might use a simple in-memory or file-based one
	// But ideally we should define a Repository interface in core/workflow
	// For now, let's assume we use the existing AgentRepo pattern or similar.
	// Since we haven't created a WorkflowRepo yet, I'll define a simple interface here or use a map for MVP if verified.
	// But to be professional, we should probably add a repo.
}

// In-Memory storage for MVP until DB is ready
var workflowStore = map[string]*workflow.GraphDefinition{
	"code-review": {
		ID:          "code-review",
		Name:        "Code Review",
		Description: "Standard parallel code review",
		Nodes: map[string]*workflow.Node{
			"start": {ID: "start", Type: workflow.NodeTypeStart, NextIDs: []string{"parallel_review"}},
			"parallel_review": {
				ID: "parallel_review", Type: workflow.NodeTypeParallel,
				NextIDs: []string{"security_agent", "performance_agent"},
			},
			"security_agent": {
				ID: "security_agent", Type: workflow.NodeTypeAgent, NextIDs: []string{"merge"},
				Properties: map[string]interface{}{"role": "security"},
			},
			"performance_agent": {
				ID: "performance_agent", Type: workflow.NodeTypeAgent, NextIDs: []string{"merge"},
				Properties: map[string]interface{}{"role": "performance"},
			},
			"merge": {ID: "merge", Type: workflow.NodeTypeEnd},
		},
		StartNodeID: "start",
	},
}

func NewWorkflowMgmtHandler() *WorkflowMgmtHandler {
	return &WorkflowMgmtHandler{}
}

type ListWorkflowsResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

// List returns available workflows
func (h *WorkflowMgmtHandler) List(c *gin.Context) {
	// TODO: Filter by type query param
	var list []ListWorkflowsResponse
	for _, w := range workflowStore {
		list = append(list, ListWorkflowsResponse{
			ID:        w.ID,
			Name:      w.Name,
			UpdatedAt: time.Now().Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, list)
}

// Get returns a specific workflow definition
func (h *WorkflowMgmtHandler) Get(c *gin.Context) {
	id := c.Param("id")
	if w, ok := workflowStore[id]; ok {
		c.JSON(http.StatusOK, w)
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
}

// Create saves a new workflow
func (h *WorkflowMgmtHandler) Create(c *gin.Context) {
	var req workflow.GraphDefinition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	workflowStore[req.ID] = &req
	c.JSON(http.StatusCreated, req)
}

// Update updates an existing workflow
func (h *WorkflowMgmtHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req workflow.GraphDefinition
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, ok := workflowStore[id]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}

	req.ID = id
	workflowStore[id] = &req
	c.JSON(http.StatusOK, req)
}

type GenerateWorkflowRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

// Generate creates a workflow from natural language
func (h *WorkflowMgmtHandler) Generate(c *gin.Context) {
	var req GenerateWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mock response for MVP
	// In real impl, we call LLM here
	mockGraph := workflow.GraphDefinition{
		ID:          uuid.New().String(),
		Name:        "Generated: " + req.Prompt,
		Description: "AI Generated workflow based on: " + req.Prompt,
		StartNodeID: "start",
		Nodes: map[string]*workflow.Node{
			"start": {ID: "start", Type: workflow.NodeTypeStart, NextIDs: []string{"end"}},
			"end":   {ID: "end", Type: workflow.NodeTypeEnd},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"graph":       mockGraph,
		"explanation": "Generated a simple start-end workflow for demo purposes.",
	})
}
