package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestWorkflowHandler_Control(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	repo := mocks.NewAgentMockRepository()
	// We don't need real WS or LLM for this test as we only test Control logic
	sessionRepo := mocks.NewSessionMockRepository()
	h := NewWorkflowHandler(nil, repo, nil, nil, sessionRepo)
	r := gin.New()
	r.POST("/sessions/:id/control", h.Control)

	// Create a dummy session and register it
	graph := &workflow.GraphDefinition{
		ID: "test-graph",
		Nodes: map[string]*workflow.Node{
			"start": {ID: "start", Type: workflow.NodeTypeStart},
		},
		StartNodeID: "start",
	}
	session := workflow.NewSession(graph, nil)
	session.Start(context.Background())

	// Manually inject into activeEngines for testing
	engine := workflow.NewEngine(session)
	activeEngines[session.ID] = engine

	t.Run("Pause", func(t *testing.T) {
		payload := ControlRequest{Action: "pause"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/control", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		if session.Status != workflow.SessionPaused {
			t.Errorf("Expected session to be Paused, got %s", session.Status)
		}
	})

	t.Run("Resume", func(t *testing.T) {
		payload := ControlRequest{Action: "resume"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/control", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		if session.Status != workflow.SessionRunning {
			t.Errorf("Expected session to be Running, got %s", session.Status)
		}
	})

	t.Run("Stop", func(t *testing.T) {
		payload := ControlRequest{Action: "stop"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/control", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}

		// Stop sets status to Failed in current implementation
		if session.Status != workflow.SessionFailed && session.Status != workflow.SessionCancelled {
			t.Errorf("Expected session to be Failed/Cancelled, got %s", session.Status)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		payload := ControlRequest{Action: "pause"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/invalid-id/control", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected 404, got %d", w.Code)
		}
	})

	t.Run("InvalidAction", func(t *testing.T) {
		payload := ControlRequest{Action: "invalid"}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/control", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected 400, got %d", w.Code)
		}
	})
}

func TestWorkflowHandler_Signal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewAgentMockRepository()
	sessionRepo := mocks.NewSessionMockRepository()
	h := NewWorkflowHandler(nil, repo, nil, nil, sessionRepo)
	r := gin.New()
	r.POST("/sessions/:id/signal", h.Signal)

	// Setup Session
	session := workflow.NewSession(nil, nil)
	session.Start(context.Background())
	engine := workflow.NewEngine(session)
	activeEngines[session.ID] = engine

	// Manually open a signal channel
	ch := session.GetSignalChannel("wait_node")

	// Consume signal in background
	go func() {
		<-ch
	}()

	payload := SignalRequest{
		NodeID:  "wait_node",
		Payload: map[string]interface{}{"decision": "approve"},
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/signal", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	// Test Invalid Node
	payload.NodeID = "invalid_node"
	body, _ = json.Marshal(payload)
	req, _ = http.NewRequest(http.MethodPost, "/sessions/"+session.ID+"/signal", bytes.NewBuffer(body))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict { // or 500/409 depending on impl
		t.Errorf("Expected 409, got %d", w.Code)
	}
}
