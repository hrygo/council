package nodes

import (
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
	"github.com/hrygo/council/internal/infrastructure/mocks"
	"github.com/hrygo/council/internal/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestGenericNodeFactory(t *testing.T) {
	cfg := &config.Config{}
	registry := llm.NewRegistry(cfg)

	// Register default provider to avoid factory failures
	mockLLM := &llm.MockProvider{}
	registry.RegisterProvider("default", mockLLM)

	agentRepo := mocks.NewAgentMockRepository()
	memoryManager := &mocks.MemoryMockManager{}

	factory := NewGenericNodeFactory(registry, agentRepo, memoryManager)

	testCases := []struct {
		nodeType workflow.NodeType
		nodeID   string
		props    map[string]interface{}
		wantErr  bool
	}{
		{workflow.NodeTypeStart, "start", nil, false},
		{workflow.NodeTypeEnd, "end", nil, false},
		{workflow.NodeTypeAgent, "agent", map[string]interface{}{"agent_uuid": "a1"}, false},
		{workflow.NodeTypeAgent, "agent-fail", nil, true}, // Missing agent_id
		{workflow.NodeTypeVote, "vote", nil, false},
		{workflow.NodeTypeLoop, "loop", nil, false},
		{workflow.NodeTypeFactCheck, "factcheck", nil, false},
		{workflow.NodeTypeHumanReview, "human", nil, false},
		{workflow.NodeTypeMemoryRetrieval, "memory", nil, false},
		{"unknown", "unknown", nil, true},
	}

	for _, tc := range testCases {
		t.Run(string(tc.nodeType), func(t *testing.T) {
			node := &workflow.Node{
				ID:         tc.nodeID,
				Type:       tc.nodeType,
				Properties: tc.props,
			}
			processor, err := factory.CreateNode(node, workflow.FactoryDeps{})
			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, processor)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, processor)
			}
		})
	}
}
