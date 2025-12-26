package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

// WorkflowMgmtHandler handles CRUD operations for workflows
// WorkflowMgmtHandler handles CRUD operations for workflows
type WorkflowMgmtHandler struct {
	Repo     workflow.Repository
	Registry *llm.Registry
}

func NewWorkflowMgmtHandler(repo workflow.Repository, registry *llm.Registry) *WorkflowMgmtHandler {
	return &WorkflowMgmtHandler{
		Repo:     repo,
		Registry: registry,
	}
}

type ListWorkflowsResponse struct {
	ID        string `json:"workflow_uuid"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

// List returns available workflows
func (h *WorkflowMgmtHandler) List(c *gin.Context) {
	ctx := c.Request.Context()
	workflows, err := h.Repo.List(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var list []ListWorkflowsResponse
	for _, w := range workflows {
		list = append(list, ListWorkflowsResponse{
			ID:        w.ID,
			Name:      w.Name,
			UpdatedAt: w.UpdatedAt.Format(time.RFC3339),
		})
	}
	c.JSON(http.StatusOK, list)
}

// Get returns a specific workflow definition
func (h *WorkflowMgmtHandler) Get(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	w, err := h.Repo.Get(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
		return
	}
	c.JSON(http.StatusOK, w)
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

	ctx := c.Request.Context()
	if err := h.Repo.Create(ctx, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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

	req.ID = id
	ctx := c.Request.Context()
	if err := h.Repo.Update(ctx, &req); err != nil {
		if strings.Contains(err.Error(), "workflow not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, req)
}

// Generate creates a workflow from natural language
func (h *WorkflowMgmtHandler) Generate(c *gin.Context) {
	var req struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	systemPrompt := `You are an expert Workflow Designer.
Your goal is to generate a valid JSON GraphDefinition based on the user's request.
Ref:
type GraphDefinition struct {
    ID          string              json:"workflow_id"
    Name        string              json:"name"
    Description string              json:"description"
    StartNodeID string              json:"start_node_id"
    Nodes       map[string]Node     json:"nodes"
}
type Node struct {
    ID         string                 json:"node_id"
    Type       NodeType               json:"type" // start, end, agent, llm, tool, parallel, sequence
    Name       string                 json:"name"
    NextIDs    []string               json:"next_ids,omitempty"
    Properties map[string]interface{} json:"properties,omitempty"
}
Output STRICT JSON only.`

	ctx := c.Request.Context()
	model := os.Getenv("LLM_MODEL")
	if model == "" {
		model = "gemini-1.5-flash" // Fallback
	}

	// Resolve Provider
	provider, err := h.Registry.GetLLMProvider("default")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get default LLM provider"})
		return
	}

	resp, err := provider.Generate(ctx, &llm.CompletionRequest{
		Messages: []llm.Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: req.Prompt},
		},
		Temperature: 0.2,
		Model:       model,
	})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "LLM generation failed: " + err.Error()})
		return
	}

	// Strip markdown fences if present
	content := resp.Content
	if strings.Contains(content, "```json") {
		content = strings.ReplaceAll(content, "```json", "")
		content = strings.ReplaceAll(content, "```", "")
	} else if strings.Contains(content, "```") {
		content = strings.ReplaceAll(content, "```", "")
	}
	content = strings.TrimSpace(content)

	var graph workflow.GraphDefinition
	if err := json.Unmarshal([]byte(content), &graph); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse generated workflow", "raw": content})
		return
	}

	// Ensure ID is empty for draft so frontend knows to POST (Create) instead of PUT (Update)
	graph.ID = ""

	c.JSON(http.StatusOK, gin.H{
		"graph":       graph,
		"explanation": "Generated workflow based on your prompt.",
	})
}

// EstimateCost calculates the estimated cost of a workflow
func (h *WorkflowMgmtHandler) EstimateCost(c *gin.Context) {
	// Support both: POST with JSON body (draft workflow) OR GET /:id (saved workflow)
	// Spec says: POST /api/v1/workflows/:id/estimate (if ID exists)
	// But usually we want to estimate *before* saving edits.
	// Let's support POST on collection /workflows/estimate with body.

	// If ID is passed in URL:
	id := c.Param("id")
	var graph workflow.GraphDefinition

	if id != "" {
		g, err := h.Repo.Get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workflow not found"})
			return
		}
		graph = *g
	} else {
		// Expect body
		if err := c.ShouldBindJSON(&graph); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	estimate := workflow.EstimateWorkflowCost(&graph)
	c.JSON(http.StatusOK, estimate)
}
