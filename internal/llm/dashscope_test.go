package llm

import (
	"os"
	"strings"
	"testing"
)

func TestNewDashScopeClient(t *testing.T) {
	tests := []struct {
		name      string
		apiKey    string
		cfg       Config
		wantModel string
	}{
		{
			name:      "with custom model",
			apiKey:    "test-key",
			cfg:       Config{Model: "qwen-max", Temperature: 0.5, MaxTokens: 2048},
			wantModel: "qwen-max",
		},
		{
			name:      "with empty model - uses default",
			apiKey:    "test-key",
			cfg:       Config{Model: "", Temperature: 0.8, MaxTokens: 4096},
			wantModel: "qwen-plus",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewDashScopeClient(tt.apiKey, tt.cfg)

			if client == nil {
				t.Fatal("NewDashScopeClient() returned nil")
			}

			if client.apiKey != tt.apiKey {
				t.Errorf("client.apiKey = %v, want %v", client.apiKey, tt.apiKey)
			}

			if client.cfg.Model != tt.wantModel {
				t.Errorf("client.cfg.Model = %v, want %v", client.cfg.Model, tt.wantModel)
			}

			if client.http == nil {
				t.Error("client.http is nil")
			}
		})
	}
}

func TestDashScopeClient_Implements_Client(t *testing.T) {
	// Verify DashScopeClient implements Client interface
	var _ Client = (*DashScopeClient)(nil)
}

func TestDashScopeClient_ConfigPreserved(t *testing.T) {
	cfg := Config{
		Model:       "qwen-turbo",
		Temperature: 0.9,
		MaxTokens:   2048,
	}

	client := NewDashScopeClient("api-key", cfg)

	if client.cfg.Temperature != 0.9 {
		t.Errorf("Temperature = %v, want %v", client.cfg.Temperature, 0.9)
	}
	if client.cfg.MaxTokens != 2048 {
		t.Errorf("MaxTokens = %v, want %v", client.cfg.MaxTokens, 2048)
	}
}

func TestDashScopeClient_HandleStream(t *testing.T) {
	// Test the handleStream function with mock data
	streamData := `data: {"choices":[{"delta":{"content":"你好"}}]}
data: {"choices":[{"delta":{"content":"世界"}}]}
data: {"choices":[{"delta":{"content":"！"}}]}
data: [DONE]
`

	client := &DashScopeClient{}

	var chunks []string
	onChunk := func(s string) {
		chunks = append(chunks, s)
	}

	result, err := client.handleStream(strings.NewReader(streamData), onChunk)

	if err != nil {
		t.Errorf("handleStream() error = %v", err)
	}

	if result != "你好世界！" {
		t.Errorf("handleStream() result = %v, want %v", result, "你好世界！")
	}

	if len(chunks) != 3 {
		t.Errorf("Callback was called %d times, want 3", len(chunks))
	}
}

func TestDashScopeClient_HandleStream_EmptyChoices(t *testing.T) {
	streamData := `data: {"choices":[]}
data: {"choices":[{"delta":{"content":"test"}}]}
data: [DONE]
`

	client := &DashScopeClient{}
	result, err := client.handleStream(strings.NewReader(streamData), nil)

	if err != nil {
		t.Errorf("handleStream() error = %v", err)
	}

	if result != "test" {
		t.Errorf("handleStream() result = %v, want %v", result, "test")
	}
}

func TestDashScopeClient_HandleStream_InvalidJSON(t *testing.T) {
	streamData := `data: invalid json
data: {"choices":[{"delta":{"content":"valid"}}]}
data: [DONE]
`

	client := &DashScopeClient{}
	result, err := client.handleStream(strings.NewReader(streamData), nil)

	if err != nil {
		t.Errorf("handleStream() error = %v", err)
	}

	// Should skip invalid JSON and continue
	if result != "valid" {
		t.Errorf("handleStream() result = %v, want %v", result, "valid")
	}
}

func TestDashScopeClient_HandleStream_NilCallback(t *testing.T) {
	streamData := `data: {"choices":[{"delta":{"content":"test"}}]}
data: [DONE]
`

	client := &DashScopeClient{}
	result, err := client.handleStream(strings.NewReader(streamData), nil)

	if err != nil {
		t.Errorf("handleStream() error = %v", err)
	}

	if result != "test" {
		t.Errorf("handleStream() result = %v, want %v", result, "test")
	}
}

// Integration test helper - skipped unless API key is set
func TestDashScopeClient_Integration(t *testing.T) {
	apiKey := os.Getenv("DASHSCOPE_API_KEY")
	if apiKey == "" {
		t.Skip("DASHSCOPE_API_KEY not set, skipping integration test")
	}

	// This test would make a real API call - only run if explicitly enabled
	t.Skip("Skipping integration test to avoid API costs")
}
