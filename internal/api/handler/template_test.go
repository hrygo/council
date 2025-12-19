package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestTemplateHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mocks.TemplateMockRepository{
		Templates: []workflow.Template{
			{ID: "1", Name: "T1"},
			{ID: "2", Name: "T2"},
		},
	}
	h := NewTemplateHandler(repo)
	r := gin.New()
	r.GET("/templates", h.List)

	req, _ := http.NewRequest(http.MethodGet, "/templates", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp []workflow.Template
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if len(resp) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(resp))
	}
}

func TestTemplateHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mocks.TemplateMockRepository{}
	h := NewTemplateHandler(repo)
	r := gin.New()
	r.POST("/templates", h.Create)

	tpl := workflow.Template{
		Name: "New Template",
	}
	body, _ := json.Marshal(tpl)
	req, _ := http.NewRequest(http.MethodPost, "/templates", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	if len(repo.Templates) != 1 {
		t.Errorf("Expected 1 template in repo, got %d", len(repo.Templates))
	}
}

func TestTemplateHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	repo := &mocks.TemplateMockRepository{
		Templates: []workflow.Template{{ID: "123"}},
	}
	h := NewTemplateHandler(repo)
	r := gin.New()
	r.DELETE("/templates/:id", h.Delete)

	req, _ := http.NewRequest(http.MethodDelete, "/templates/123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if len(repo.Templates) != 0 {
		t.Errorf("Expected 0 templates in repo, got %d", len(repo.Templates))
	}
}
