package tools

import (
	"context"

	"github.com/hrygo/council/internal/core/workflow"
)

// Tool defines an executable capability
type Tool interface {
	Name() string
	Description() string
	Parameters() map[string]interface{}
	Execute(ctx context.Context, args map[string]interface{}) (string, error)
}

// SessionAwareTool is a tool that requires access to the workflow session (e.g. VFS).
type SessionAwareTool interface {
	Tool
	ExecuteWithSession(ctx context.Context, session *workflow.Session, args map[string]interface{}) (string, error)
}

// Registry manages available tools
type Registry struct {
	tools map[string]Tool
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]Tool),
	}
}

func (r *Registry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

func (r *Registry) GetTool(name string) (Tool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

func (r *Registry) ListTools() []Tool {
	var list []Tool
	for _, t := range r.tools {
		list = append(list, t)
	}
	return list
}
