package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestAgentHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewAgentMockRepository()
	h := NewAgentHandler(repo)
	r := gin.New()
	r.POST("/agents", h.Create)

	t.Run("Success", func(t *testing.T) {
		a := &agent.Agent{
			Name: "Test Agent",
		}
		body, _ := json.Marshal(a)
		req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp agent.Agent
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if resp.Name != "Test Agent" {
			t.Errorf("Expected name 'Test Agent', got '%s'", resp.Name)
		}
		if resp.ID == uuid.Nil {
			t.Errorf("Expected ID to be set")
		}
	})

	t.Run("CreateError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		a := &agent.Agent{Name: "Fail"}
		body, _ := json.Marshal(a)
		req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/agents", bytes.NewBufferString("{invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestAgentHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewAgentMockRepository()
	h := NewAgentHandler(repo)
	r := gin.New()
	r.GET("/agents/:id", h.Get)

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/agents/invalid-uuid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	existingAgent := &agent.Agent{ID: uuid.New(), Name: "Bond"}
	if err := repo.Create(context.Background(), existingAgent); err != nil {
		t.Fatalf("Failed to create mock agent: %v", err)
	}

	t.Run("Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/agents/"+existingAgent.ID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/agents/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestAgentHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewAgentMockRepository()
	h := NewAgentHandler(repo)
	r := gin.New()
	r.GET("/agents", h.List)

	if err := repo.Create(context.Background(), &agent.Agent{Name: "A1"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(context.Background(), &agent.Agent{Name: "A2"}); err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/agents", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []*agent.Agent
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("Expected 2 agents, got %d", len(resp))
	}

	t.Run("RepoError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		req, _ := http.NewRequest(http.MethodGet, "/agents", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})
}

func TestAgentHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewAgentMockRepository()
	h := NewAgentHandler(repo)
	r := gin.New()
	r.PUT("/agents/:id", h.Update)

	id := uuid.New()
	if err := repo.Create(context.Background(), &agent.Agent{ID: id, Name: "Old"}); err != nil {
		t.Fatal(err)
	}

	updated := &agent.Agent{Name: "New"}
	body, _ := json.Marshal(updated)
	req, _ := http.NewRequest(http.MethodPut, "/agents/"+id.String(), bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	a, _ := repo.GetByID(context.Background(), id)
	if a.Name != "New" {
		t.Errorf("Expected name 'New', got '%s'", a.Name)
	}

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/agents/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		id := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, "/agents/"+id.String(), bytes.NewBufferString("{inv"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		id := uuid.New()
		a := &agent.Agent{Name: "X"}
		body, _ := json.Marshal(a)
		req, _ := http.NewRequest(http.MethodPut, "/agents/"+id.String(), bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})
}

func TestAgentHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewAgentMockRepository()
	h := NewAgentHandler(repo)
	r := gin.New()
	r.DELETE("/agents/:id", h.Delete)

	id := uuid.New()
	if err := repo.Create(context.Background(), &agent.Agent{ID: id}); err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodDelete, "/agents/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	_, err := repo.GetByID(context.Background(), id)
	if err == nil {
		t.Errorf("Expected agent to be deleted")
	}

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/agents/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		id := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, "/agents/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})
}
