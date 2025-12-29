package tools

import (
	"context"
	"fmt"

	"github.com/hrygo/council/internal/core/workflow"
)

type WriteFileTool struct{}

func (t *WriteFileTool) Name() string {
	return "write_file"
}

func (t *WriteFileTool) Description() string {
	return "Write content to a file in the virtual file system. If the file exists, it creates a new version."
}

func (t *WriteFileTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Relative path of the file",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "Full content of the file",
			},
			"reason": map[string]interface{}{
				"type":        "string",
				"description": "Reason for the change (e.g. 'Fixing syntax error')",
			},
		},
		"required": []string{"path", "content"},
	}
}

func (t *WriteFileTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return "", fmt.Errorf("session required")
}

func (t *WriteFileTool) ExecuteWithSession(ctx context.Context, session *workflow.Session, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path is required")
	}
	content, ok := args["content"].(string)
	if !ok {
		return "", fmt.Errorf("content is required")
	}
	reason, _ := args["reason"].(string)
	if reason == "" {
		reason = "tool_call"
	}

	ver, err := session.WriteFile(path, content, "agent", reason)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("File %s written (Version %d)", path, ver), nil
}

type ReadFileTool struct{}

func (t *ReadFileTool) Name() string {
	return "read_file"
}

func (t *ReadFileTool) Description() string {
	return "Read the content of a file from the virtual file system (latest version)."
}

func (t *ReadFileTool) Parameters() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"path": map[string]interface{}{
				"type":        "string",
				"description": "Relative path of the file",
			},
		},
		"required": []string{"path"},
	}
}

func (t *ReadFileTool) Execute(ctx context.Context, args map[string]interface{}) (string, error) {
	return "", fmt.Errorf("session required")
}

func (t *ReadFileTool) ExecuteWithSession(ctx context.Context, session *workflow.Session, args map[string]interface{}) (string, error) {
	path, ok := args["path"].(string)
	if !ok {
		return "", fmt.Errorf("path is required")
	}

	// Try to get from VFS first
	file, err := session.GetLatestFile(path)
	if err == nil {
		return file.Content, nil
	}

	// Fallback? Spec didn't explicitly say. But generally System Surgeon works on VFS.
	// We can return error if not found.
	return "", fmt.Errorf("file not found: %s", path)
}
