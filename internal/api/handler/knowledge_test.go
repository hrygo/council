package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/stretchr/testify/assert"
)

func TestKnowledgeHandler_GetSessionKnowledge(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	// Create handler with mock service
	memService := memory.NewService(nil, nil, nil)
	handler := NewKnowledgeHandler(memService)
	
	// Setup router
	r := gin.New()
	r.GET("/api/v1/sessions/:sessionID/knowledge", handler.GetSessionKnowledge)
	
	t.Run("Get all knowledge items", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		
		assert.Equal(t, 5, resp.Total)
		assert.Equal(t, 50, resp.Limit)
		assert.Equal(t, 0, resp.Offset)
		assert.Len(t, resp.Items, 5)
	})
	
	t.Run("Filter by memory layer", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge?layer=sandboxed", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		
		assert.Equal(t, 2, resp.Total)
		for _, item := range resp.Items {
			assert.Equal(t, "sandboxed", item.MemoryLayer)
		}
	})
	
	t.Run("Pagination", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge?limit=2&offset=1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		
		assert.Equal(t, 5, resp.Total)
		assert.Equal(t, 2, resp.Limit)
		assert.Equal(t, 1, resp.Offset)
		assert.Len(t, resp.Items, 2)
	})
	
	t.Run("Filter by working memory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge?layer=working", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		
		assert.Equal(t, 2, resp.Total)
		for _, item := range resp.Items {
			assert.Equal(t, "working", item.MemoryLayer)
		}
	})
	
	t.Run("Filter by long-term memory", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/sessions/session-1/knowledge?layer=long-term", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var resp KnowledgeResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		
		assert.Equal(t, 1, resp.Total)
		for _, item := range resp.Items {
			assert.Equal(t, "long-term", item.MemoryLayer)
		}
	})
}
