package cli

import (
	"bytes"
	"testing"

	"github.com/hrygo/dialecta/internal/config"
)

func TestNewRunner(t *testing.T) {
	cfg := config.New()
	runner := NewRunner(cfg, true)

	if runner == nil {
		t.Fatal("NewRunner() returned nil")
	}
	if runner.cfg != cfg {
		t.Error("NewRunner() did not set cfg correctly")
	}
	if !runner.stream {
		t.Error("NewRunner() did not set stream correctly")
	}
	if runner.ui == nil {
		t.Error("NewRunner() should create default UI")
	}
	if runner.input == nil {
		t.Error("NewRunner() should create default InputReader")
	}
	if runner.executor == nil {
		t.Error("NewRunner() should create executor")
	}
}

func TestNewRunnerWithOptions(t *testing.T) {
	cfg := config.New()
	ui := NewUI(&bytes.Buffer{}, &bytes.Buffer{})
	input := NewInputReader(&bytes.Buffer{}, &bytes.Buffer{})

	runner := NewRunnerWithOptions(cfg, false, ui, input)

	if runner == nil {
		t.Fatal("NewRunnerWithOptions() returned nil")
	}
	if runner.ui != ui {
		t.Error("NewRunnerWithOptions() did not set ui correctly")
	}
	if runner.input != input {
		t.Error("NewRunnerWithOptions() did not set input correctly")
	}
	if runner.stream != false {
		t.Error("NewRunnerWithOptions() did not set stream correctly")
	}
}

func TestRunner_StreamMode(t *testing.T) {
	cfg := config.New()

	streamRunner := NewRunner(cfg, true)
	if !streamRunner.stream {
		t.Error("Runner with stream=true should have stream mode enabled")
	}

	nonStreamRunner := NewRunner(cfg, false)
	if nonStreamRunner.stream {
		t.Error("Runner with stream=false should have stream mode disabled")
	}
}

func TestSetupContext(t *testing.T) {
	ctx, cancel := SetupContext()
	defer cancel()

	if ctx == nil {
		t.Fatal("SetupContext() returned nil context")
	}
	if cancel == nil {
		t.Fatal("SetupContext() returned nil cancel function")
	}

	// Context should not be done initially
	select {
	case <-ctx.Done():
		t.Error("Context should not be done immediately")
	default:
		// Expected
	}

	// Cancel should work
	cancel()
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Error("Context should be done after cancel")
	}
}
