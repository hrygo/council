package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListFiles returns the latest state of the VFS for a session.
func (h *WorkflowHandler) ListFiles(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id required"})
		return
	}

	files, err := h.FileRepo.ListFiles(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}

// GetFileHistory returns all versions of a specific file.
func (h *WorkflowHandler) GetFileHistory(c *gin.Context) {
	sessionID := c.Param("id")
	path := c.Query("path") // Use query param for path to avoid URL encoding issues with forward slashes
	if sessionID == "" || path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "session_id and path required"})
		return
	}

	versions, err := h.FileRepo.ListVersions(c.Request.Context(), sessionID, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list file versions: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"history": versions})
}
