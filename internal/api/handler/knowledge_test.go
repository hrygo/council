package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/stretchr/testify/assert"
)

func TestKnowledgeHandler_GetSessionKnowledge(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create handler with mock service
	mockManager := &mocks.MemoryMockManager{}
	sessionRepo := mocks.NewSessionMockRepository()
	handler := NewKnowledgeHandler(mockManager, sessionRepo)

	// Setup router
	r := gin.New()
	r.GET("/api/v1/sessions/:sessionID/knowledge", handler.GetSessionKnowledge)

	t.Run("Get all knowledge items", func(t *testing.T) {
		mockManager.RetrieveResult = []memory.ContextItem{
			{Content: "item 1", Source: "cold", Score: 0.9},
			{Content: "item 2", Source: "hot", Score: 0.8},
		}

		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, 2, resp.Total)
		assert.Len(t, resp.Items, 2)
	})

	t.Run("Filter by memory layer", func(t *testing.T) {
		mockManager.RetrieveResult = []memory.ContextItem{
			{Content: "item 1", Source: "cold", Score: 0.9},
			{Content: "item 2", Source: "hot", Score: 0.8},
		}

		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge?layer=cold", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, 1, resp.Total)
		assert.Equal(t, "cold", resp.Items[0].MemoryLayer)
	})
}
