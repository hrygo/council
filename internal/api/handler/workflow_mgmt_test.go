package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/workflow"
)

func TestWorkflowMgmtHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewWorkflowMgmtHandler()
	r := gin.New()
	r.GET("/workflows", h.List)

	req, _ := http.NewRequest(http.MethodGet, "/workflows", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []ListWorkflowsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// We expect at least the hardcoded "code-review" workflow
	if len(resp) < 1 {
		t.Errorf("Expected at least 1 workflow, got %d", len(resp))
	}
}

func TestWorkflowMgmtHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewWorkflowMgmtHandler()
	r := gin.New()
	r.GET("/workflows/:id", h.Get)

	t.Run("Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/workflows/code-review", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		var resp workflow.GraphDefinition
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if resp.ID != "code-review" {
			t.Errorf("Expected ID code-review, got %s", resp.ID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/workflows/non-existent", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestWorkflowMgmtHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewWorkflowMgmtHandler()
	r := gin.New()
	r.POST("/workflows", h.Create)

	newWF := workflow.GraphDefinition{
		Name: "New Workflow",
		Nodes: map[string]*workflow.Node{
			"start": {ID: "start", Type: workflow.NodeTypeStart},
		},
		StartNodeID: "start",
	}

	body, _ := json.Marshal(newWF)
	req, _ := http.NewRequest(http.MethodPost, "/workflows", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp workflow.GraphDefinition
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp.ID == "" {
		t.Error("Expected ID to be generated")
	}
}

func TestWorkflowMgmtHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewWorkflowMgmtHandler()
	r := gin.New()
	r.PUT("/workflows/:id", h.Update)

	// Pre-seed a workflow to update
	id := uuid.New().String()
	workflowStore[id] = &workflow.GraphDefinition{
		ID:          id,
		Name:        "Original",
		StartNodeID: "start",
	}

	t.Run("Success", func(t *testing.T) {
		updatePayload := workflow.GraphDefinition{
			Name:        "Updated Name",
			StartNodeID: "start",
		}
		body, _ := json.Marshal(updatePayload)
		req, _ := http.NewRequest(http.MethodPut, "/workflows/"+id, bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
		var resp workflow.GraphDefinition
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatal(err)
		}
		if resp.Name != "Updated Name" {
			t.Errorf("Expected Name updated, got %s", resp.Name)
		}
		if resp.ID != id {
			t.Errorf("Expected ID to persist, got %s", resp.ID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/workflows/invalid", bytes.NewBufferString("{}"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestWorkflowMgmtHandler_Generate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := NewWorkflowMgmtHandler()
	r := gin.New()
	r.POST("/workflows/generate", h.Generate)

	payload := GenerateWorkflowRequest{Prompt: "Review code"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/workflows/generate", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if _, ok := resp["graph"]; !ok {
		t.Error("Expected graph in response")
	}
	if _, ok := resp["explanation"]; !ok {
		t.Error("Expected explanation in response")
	}
}
