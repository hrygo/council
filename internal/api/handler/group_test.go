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
	"github.com/hrygo/council/internal/core/group"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestGroupHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewGroupMockRepository()
	h := NewGroupHandler(repo)
	r := gin.New()
	r.POST("/groups", h.Create)

	t.Run("Success", func(t *testing.T) {
		g := &group.Group{
			Name: "Test Group",
		}
		body, _ := json.Marshal(g)
		req, _ := http.NewRequest(http.MethodPost, "/groups", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", w.Code)
		}

		var resp group.Group
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}
		if resp.Name != "Test Group" {
			t.Errorf("Expected name 'Test Group', got '%s'", resp.Name)
		}
		if resp.ID == uuid.Nil {
			t.Errorf("Expected ID to be set")
		}
	})

	t.Run("CreateError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		g := &group.Group{Name: "Fail"}
		body, _ := json.Marshal(g)
		req, _ := http.NewRequest(http.MethodPost, "/groups", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil // Reset
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/groups", bytes.NewBufferString("{invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})
}

func TestGroupHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewGroupMockRepository()
	h := NewGroupHandler(repo)
	r := gin.New()
	r.GET("/groups/:id", h.Get)

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/groups/invalid-uuid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	existingGroup := &group.Group{ID: uuid.New(), Name: "Existing"}
	if err := repo.Create(context.Background(), existingGroup); err != nil {
		t.Fatalf("Failed to create mock group: %v", err)
	}

	t.Run("Found", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/groups/"+existingGroup.ID.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/groups/"+uuid.New().String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", w.Code)
		}
	})
}

func TestGroupHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewGroupMockRepository()
	h := NewGroupHandler(repo)
	r := gin.New()
	r.GET("/groups", h.List)

	if err := repo.Create(context.Background(), &group.Group{Name: "G1"}); err != nil {
		t.Fatal(err)
	}
	if err := repo.Create(context.Background(), &group.Group{Name: "G2"}); err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodGet, "/groups", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []*group.Group
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(resp))
	}

	t.Run("RepoError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		req, _ := http.NewRequest(http.MethodGet, "/groups", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})
}

func TestGroupHandler_Update(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewGroupMockRepository()
	h := NewGroupHandler(repo)
	r := gin.New()
	r.PUT("/groups/:id", h.Update)

	id := uuid.New()
	if err := repo.Create(context.Background(), &group.Group{ID: id, Name: "Old"}); err != nil {
		t.Fatal(err)
	}

	updated := &group.Group{Name: "New"}
	body, _ := json.Marshal(updated)
	req, _ := http.NewRequest(http.MethodPut, "/groups/"+id.String(), bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	g, _ := repo.GetByID(context.Background(), id)
	if g.Name != "New" {
		t.Errorf("Expected name 'New', got '%s'", g.Name)
	}

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPut, "/groups/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		id := uuid.New()
		req, _ := http.NewRequest(http.MethodPut, "/groups/"+id.String(), bytes.NewBufferString("{inv"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		id := uuid.New()
		g := &group.Group{Name: "X"}
		body, _ := json.Marshal(g)
		req, _ := http.NewRequest(http.MethodPut, "/groups/"+id.String(), bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})
}

func TestGroupHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := mocks.NewGroupMockRepository()
	h := NewGroupHandler(repo)
	r := gin.New()
	r.DELETE("/groups/:id", h.Delete)

	id := uuid.New()
	if err := repo.Create(context.Background(), &group.Group{ID: id}); err != nil {
		t.Fatal(err)
	}

	req, _ := http.NewRequest(http.MethodDelete, "/groups/"+id.String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status 204, got %d", w.Code)
	}

	_, err := repo.GetByID(context.Background(), id)
	if err == nil {
		t.Errorf("Expected group to be deleted")
	}

	t.Run("InvalidID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodDelete, "/groups/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", w.Code)
		}
	})

	t.Run("RepoError", func(t *testing.T) {
		repo.Err = context.DeadlineExceeded
		id := uuid.New()
		req, _ := http.NewRequest(http.MethodDelete, "/groups/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", w.Code)
		}
		repo.Err = nil
	})
}
