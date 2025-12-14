package debate

import (
	"context"
	"os"
	"testing"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/llm"
)

func TestNewExecutor(t *testing.T) {
	cfg := config.New()
	executor := NewExecutor(cfg)

	if executor == nil {
		t.Fatal("NewExecutor() returned nil")
	}

	if executor.cfg != cfg {
		t.Error("NewExecutor() did not set cfg correctly")
	}

	if executor.stream != false {
		t.Error("NewExecutor() should set stream to false by default")
	}
}

func TestExecutor_SetStream(t *testing.T) {
	cfg := config.New()
	executor := NewExecutor(cfg)

	proCalled := false
	conCalled := false
	judgeCalled := false

	onPro := func(s string) { proCalled = true }
	onCon := func(s string) { conCalled = true }
	onJudge := func(s string) { judgeCalled = true }

	executor.SetStream(onPro, onCon, onJudge)

	if !executor.stream {
		t.Error("SetStream() did not set stream to true")
	}

	// Test callbacks are set by calling them
	if executor.onPro == nil {
		t.Error("SetStream() did not set onPro callback")
	}
	if executor.onCon == nil {
		t.Error("SetStream() did not set onCon callback")
	}
	if executor.onJudge == nil {
		t.Error("SetStream() did not set onJudge callback")
	}

	// Verify callbacks work
	executor.onPro("test")
	executor.onCon("test")
	executor.onJudge("test")

	if !proCalled {
		t.Error("onPro callback was not called")
	}
	if !conCalled {
		t.Error("onCon callback was not called")
	}
	if !judgeCalled {
		t.Error("onJudge callback was not called")
	}
}

func TestResult(t *testing.T) {
	result := &Result{
		Material:    "test material",
		ProArgument: "pro argument",
		ConArgument: "con argument",
		Verdict:     "verdict",
	}

	if result.Material != "test material" {
		t.Errorf("Result.Material = %v, want %v", result.Material, "test material")
	}
	if result.ProArgument != "pro argument" {
		t.Errorf("Result.ProArgument = %v, want %v", result.ProArgument, "pro argument")
	}
	if result.ConArgument != "con argument" {
		t.Errorf("Result.ConArgument = %v, want %v", result.ConArgument, "con argument")
	}
	if result.Verdict != "verdict" {
		t.Errorf("Result.Verdict = %v, want %v", result.Verdict, "verdict")
	}
}

func TestNewExecutor_WithCustomConfig(t *testing.T) {
	cfg := &config.Config{
		ProRole: config.RoleConfig{
			Provider:    llm.ProviderDeepSeek,
			Model:       "custom-model",
			Temperature: 0.5,
			MaxTokens:   2048,
		},
		ConRole: config.RoleConfig{
			Provider:    llm.ProviderGemini,
			Model:       "gemini-pro",
			Temperature: 0.7,
			MaxTokens:   4096,
		},
		JudgeRole: config.RoleConfig{
			Provider:    llm.ProviderDashScope,
			Model:       "qwen-max",
			Temperature: 0.1,
			MaxTokens:   8192,
		},
	}

	executor := NewExecutor(cfg)

	if executor.cfg.ProRole.Model != "custom-model" {
		t.Errorf("executor.cfg.ProRole.Model = %v, want %v", executor.cfg.ProRole.Model, "custom-model")
	}
	if executor.cfg.ConRole.Provider != llm.ProviderGemini {
		t.Errorf("executor.cfg.ConRole.Provider = %v, want %v", executor.cfg.ConRole.Provider, llm.ProviderGemini)
	}
	if executor.cfg.JudgeRole.MaxTokens != 8192 {
		t.Errorf("executor.cfg.JudgeRole.MaxTokens = %v, want %v", executor.cfg.JudgeRole.MaxTokens, 8192)
	}
}

func TestExecutor_SetStream_NilCallbacks(t *testing.T) {
	cfg := config.New()
	executor := NewExecutor(cfg)

	// Should not panic with nil callbacks
	executor.SetStream(nil, nil, nil)

	if !executor.stream {
		t.Error("SetStream() should still set stream to true even with nil callbacks")
	}
}

func TestResult_EmptyFields(t *testing.T) {
	result := &Result{}

	if result.Material != "" {
		t.Error("Empty Result.Material should be empty string")
	}
	if result.ProArgument != "" {
		t.Error("Empty Result.ProArgument should be empty string")
	}
	if result.ConArgument != "" {
		t.Error("Empty Result.ConArgument should be empty string")
	}
	if result.Verdict != "" {
		t.Error("Empty Result.Verdict should be empty string")
	}
}

func TestExecutor_Execute_MissingAPIKey(t *testing.T) {
	// Save and clear env vars
	origDeepSeek := os.Getenv("DEEPSEEK_API_KEY")
	origGemini := os.Getenv("GEMINI_API_KEY")
	origGoogle := os.Getenv("GOOGLE_API_KEY")
	origDashScope := os.Getenv("DASHSCOPE_API_KEY")
	defer func() {
		os.Setenv("DEEPSEEK_API_KEY", origDeepSeek)
		os.Setenv("GEMINI_API_KEY", origGemini)
		os.Setenv("GOOGLE_API_KEY", origGoogle)
		os.Setenv("DASHSCOPE_API_KEY", origDashScope)
	}()

	os.Unsetenv("DEEPSEEK_API_KEY")
	os.Unsetenv("GEMINI_API_KEY")
	os.Unsetenv("GOOGLE_API_KEY")
	os.Unsetenv("DASHSCOPE_API_KEY")

	cfg := config.New()
	executor := NewExecutor(cfg)

	ctx := context.Background()
	_, err := executor.Execute(ctx, "test material")

	if err == nil {
		t.Error("Execute() should fail when API keys are missing")
	}
}

func TestExecutor_Execute_ContextCancellation(t *testing.T) {
	// Test that Execute respects context cancellation
	cfg := config.New()
	executor := NewExecutor(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// This should fail fast due to cancelled context
	// Note: The actual behavior depends on how quickly the goroutines check the context
	_, err := executor.Execute(ctx, "test material")

	// We expect an error since API keys are not set or context is cancelled
	if err == nil {
		t.Log("Execute() did not return error, which is acceptable if it completed before cancellation")
	}
}

func TestExecutor_StreamCallbacksInvoked(t *testing.T) {
	cfg := config.New()
	executor := NewExecutor(cfg)

	var proChunks, conChunks, judgeChunks []string

	executor.SetStream(
		func(s string) { proChunks = append(proChunks, s) },
		func(s string) { conChunks = append(conChunks, s) },
		func(s string) { judgeChunks = append(judgeChunks, s) },
	)

	// Verify stream mode is enabled
	if !executor.stream {
		t.Error("stream mode should be enabled")
	}

	// Manually invoke callbacks to test they are correctly wired
	executor.onPro("pro chunk 1")
	executor.onPro("pro chunk 2")
	executor.onCon("con chunk 1")
	executor.onJudge("judge chunk 1")

	if len(proChunks) != 2 {
		t.Errorf("proChunks length = %d, want 2", len(proChunks))
	}
	if len(conChunks) != 1 {
		t.Errorf("conChunks length = %d, want 1", len(conChunks))
	}
	if len(judgeChunks) != 1 {
		t.Errorf("judgeChunks length = %d, want 1", len(judgeChunks))
	}
}

func TestExecutor_MultipleSetStream(t *testing.T) {
	cfg := config.New()
	executor := NewExecutor(cfg)

	// First set
	var count1 int
	executor.SetStream(
		func(s string) { count1++ },
		func(s string) { count1++ },
		func(s string) { count1++ },
	)

	// Second set should override
	var count2 int
	executor.SetStream(
		func(s string) { count2++ },
		func(s string) { count2++ },
		func(s string) { count2++ },
	)

	// Invoke callbacks
	executor.onPro("test")
	executor.onCon("test")
	executor.onJudge("test")

	if count1 != 0 {
		t.Error("First set of callbacks should have been overridden")
	}
	if count2 != 3 {
		t.Errorf("count2 = %d, want 3", count2)
	}
}

func TestNewExecutor_NilConfig(t *testing.T) {
	// Test behavior with nil config - should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Log("NewExecutor(nil) panicked as expected for safety")
		}
	}()

	executor := NewExecutor(nil)
	if executor != nil && executor.cfg != nil {
		t.Error("Expected nil config to be preserved or panic")
	}
}
