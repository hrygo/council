package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/hrygo/council/internal/pkg/config"
)

func TestWorkflowMgmtHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := &mocks.WorkflowMockRepository{}
	handler := NewWorkflowMgmtHandler(mockRepo, nil)

	mockRepo.ListFunc = func(ctx context.Context) ([]*workflow.WorkflowEntity, error) {
		return []*workflow.WorkflowEntity{
			{ID: "w1", Name: "Workflow 1", UpdatedAt: time.Now()},
		}, nil
	}

	router := gin.New()
	router.GET("/workflows", handler.List)

	req, _ := http.NewRequest("GET", "/workflows", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []ListWorkflowsResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) != 1 || resp[0].ID != "w1" {
		t.Errorf("Unexpected response: %+v", resp)
	}
}

func TestWorkflowMgmtHandler_Get(t *testing.T) {
	mockRepo := &mocks.WorkflowMockRepository{}
	handler := NewWorkflowMgmtHandler(mockRepo, nil)

	mockRepo.GetFunc = func(ctx context.Context, id string) (*workflow.GraphDefinition, error) {
		if id == "w1" {
			return &workflow.GraphDefinition{ID: "w1", Name: "Workflow 1"}, nil
		}
		return nil, nil // Error case usually handled by Get returning err
	}

	router := gin.New()
	router.GET("/workflows/:id", handler.Get)

	t.Run("Found", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/workflows/w1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestWorkflowMgmtHandler_Create(t *testing.T) {
	mockRepo := &mocks.WorkflowMockRepository{}
	handler := NewWorkflowMgmtHandler(mockRepo, nil)

	router := gin.New()
	router.POST("/workflows", handler.Create)

	payload := workflow.GraphDefinition{Name: "New Workflow"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "/workflows", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d", w.Code)
	}
}

func TestWorkflowMgmtHandler_Update(t *testing.T) {
	mockRepo := &mocks.WorkflowMockRepository{}
	handler := NewWorkflowMgmtHandler(mockRepo, nil)

	mockRepo.UpdateFunc = func(ctx context.Context, graph *workflow.GraphDefinition) error {
		return nil
	}

	router := gin.New()
	router.PUT("/workflows/:id", handler.Update)

	payload := workflow.GraphDefinition{Name: "Updated Workflow"}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", "/workflows/w1", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestWorkflowMgmtHandler_EstimateCost(t *testing.T) {
	mockRepo := &mocks.WorkflowMockRepository{}
	handler := NewWorkflowMgmtHandler(mockRepo, nil)

	router := gin.New()
	router.POST("/workflows/estimate", handler.EstimateCost)
	router.POST("/workflows/:id/estimate", handler.EstimateCost)

	t.Run("Draft", func(t *testing.T) {
		graph := workflow.GraphDefinition{
			Nodes: map[string]*workflow.Node{
				"n1": {ID: "n1", Type: "agent", Properties: map[string]interface{}{"model": "gpt-4"}},
			},
		}
		body, _ := json.Marshal(graph)
		req, _ := http.NewRequest("POST", "/workflows/estimate", bytes.NewReader(body))
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})

	t.Run("Saved", func(t *testing.T) {
		mockRepo.GetFunc = func(ctx context.Context, id string) (*workflow.GraphDefinition, error) {
			return &workflow.GraphDefinition{ID: id, Nodes: map[string]*workflow.Node{}}, nil
		}
		req, _ := http.NewRequest("POST", "/workflows/w1/estimate", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", w.Code)
		}
	})
}

func TestWorkflowMgmtHandler_Generate(t *testing.T) {
	mockLLM := &llm.MockProvider{
		GenerateResponse: &llm.CompletionResponse{
			Content: `{"id": "gen-1", "name": "Generated"}`,
		},
	}

	cfg := &config.Config{}
	registry := llm.NewRegistry(cfg)
	registry.RegisterProvider("default", mockLLM)

	handler := NewWorkflowMgmtHandler(nil, registry)

	router := gin.New()
	router.POST("/workflows/generate", handler.Generate)

	body, _ := json.Marshal(map[string]string{"prompt": "make a workflow"})
	req, _ := http.NewRequest("POST", "/workflows/generate", bytes.NewReader(body))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}
