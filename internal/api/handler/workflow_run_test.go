package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/api/ws"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestWorkflowHandler_Execute(t *testing.T) {
	gin.SetMode(gin.TestMode)
	hub := ws.NewHub()
	go hub.Run()

	repo := mocks.NewAgentMockRepository()
	mockLLM := &llm.MockProvider{}
	h := NewWorkflowHandler(hub, repo, mockLLM)

	router := gin.New()
	router.POST("/execute", h.Execute)

	graph := &workflow.GraphDefinition{
		ID: "test-wf",
		Nodes: map[string]*workflow.Node{
			"start": {ID: "start", Type: workflow.NodeTypeStart, NextIDs: []string{"end"}},
			"end":   {ID: "end", Type: workflow.NodeTypeEnd},
		},
		StartNodeID: "start",
	}

	payload := ExecuteRequest{
		Graph: graph,
		Input: map[string]interface{}{"query": "hello"},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/execute", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Errorf("Expected 202, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["status"] != "started" || resp["session_id"] == "" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestWorkflowHandler_Review(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewWorkflowHandler(nil, nil, nil)

	router := gin.New()
	router.POST("/sessions/:id/review", h.Review)

	// Setup active session in a suspended node
	graph := &workflow.GraphDefinition{
		ID: "test-wf",
		Nodes: map[string]*workflow.Node{
			"review": {ID: "review", Type: workflow.NodeTypeHumanReview},
		},
	}
	session := workflow.NewSession(graph, nil)
	engine := workflow.NewEngine(session)
	activeEngines[session.ID] = engine

	// To simulate suspended, we manually set status in map
	engine.Status["review"] = workflow.StatusSuspended

	payload := ReviewRequest{
		NodeID: "review",
		Action: "approve",
		Data:   map[string]interface{}{"comment": "ok"},
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/sessions/"+session.ID+"/review", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
