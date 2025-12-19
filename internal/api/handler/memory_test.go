package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/infrastructure/mocks"
)

func TestMemoryHandler_Ingest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockManager := &mocks.MemoryMockManager{}
	h := NewMemoryHandler(mockManager)
	r := gin.New()
	r.POST("/memory/ingest", h.Ingest)

	reqBody := IngestRequest{
		GroupID: "g1",
		Content: "test content",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/memory/ingest", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestMemoryHandler_Query(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockManager := &mocks.MemoryMockManager{}
	h := NewMemoryHandler(mockManager)
	r := gin.New()
	r.POST("/memory/query", h.Query)

	reqBody := QueryRequest{
		GroupID: "g1",
		Query:   "test query",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest(http.MethodPost, "/memory/query", bytes.NewBuffer(body))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
