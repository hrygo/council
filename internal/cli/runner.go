package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/debate"
)

// Runner orchestrates the CLI application execution
type Runner struct {
	ui       *UI
	input    *InputReader
	cfg      *config.Config
	stream   bool
	executor *debate.Executor
}

// NewRunner creates a new CLI runner
func NewRunner(cfg *config.Config, stream bool) *Runner {
	return &Runner{
		ui:       DefaultUI(),
		input:    DefaultInputReader(),
		cfg:      cfg,
		stream:   stream,
		executor: debate.NewExecutor(cfg),
	}
}

// RunWithOptions creates a runner with custom UI and input reader
func NewRunnerWithOptions(cfg *config.Config, stream bool, ui *UI, input *InputReader) *Runner {
	return &Runner{
		ui:       ui,
		input:    input,
		cfg:      cfg,
		stream:   stream,
		executor: debate.NewExecutor(cfg),
	}
}

// Run executes the debate with the given material
func (r *Runner) Run(ctx context.Context, material string) error {
	// Validate material
	if err := ValidateMaterial(material); err != nil {
		return err
	}

	// Print banner and config
	r.ui.PrintBanner()
	r.ui.PrintConfig(r.cfg)

	if r.stream {
		return r.runStreaming(ctx, material)
	}
	return r.runNonStreaming(ctx, material)
}

// runStreaming executes the debate in streaming mode using sequential display
func (r *Runner) runStreaming(ctx context.Context, material string) error {
	r.ui.PrintDebating()

	var (
		mu           sync.Mutex
		conBuffer    strings.Builder
		proFinished  bool
		proStarted   bool
		conStarted   bool
		judgeStarted bool
	)

	r.executor.SetStream(
		// Pro Callback
		func(chunk string, done bool) {
			mu.Lock()
			defer mu.Unlock()

			if done {
				proFinished = true
				// Flush any buffered Con content immediately
				if conBuffer.Len() > 0 {
					if !conStarted {
						r.ui.Println("")
						r.ui.PrintConHeader()
						conStarted = true
					}
					r.ui.Print(conBuffer.String())
					conBuffer.Reset()
				}
				return
			}

			if !proStarted {
				r.ui.PrintProHeader()
				proStarted = true
			}
			r.ui.Print(chunk)
		},

		// Con Callback
		func(chunk string, done bool) {
			mu.Lock()
			defer mu.Unlock()

			if done {
				return
			}

			if proFinished {
				// Direct print if Pro is done
				if !conStarted {
					r.ui.Println("")
					r.ui.PrintConHeader()
					conStarted = true
				}
				r.ui.Print(chunk)
			} else {
				// Buffer if Pro is running
				conBuffer.WriteString(chunk)
			}
		},

		// Judge Callback
		func(chunk string, done bool) {
			if done {
				return
			}

			if !judgeStarted {
				r.ui.Println("")
				r.ui.PrintJudgeHeader()
				judgeStarted = true
			}
			r.ui.Print(chunk)
		},
	)

	_, err := r.executor.Execute(ctx, material)
	if err != nil {
		return fmt.Errorf("执行失败: %w", err)
	}

	r.ui.PrintComplete()

	return nil
}

// runNonStreaming executes the debate in non-streaming mode
func (r *Runner) runNonStreaming(ctx context.Context, material string) error {
	r.ui.PrintDebating()

	result, err := r.executor.Execute(ctx, material)
	if err != nil {
		return fmt.Errorf("执行失败: %w", err)
	}

	r.ui.PrintResult(result)
	r.ui.PrintComplete()

	return nil
}

// SetupContext creates a context that can be cancelled by interrupt signals
func SetupContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		ui := DefaultUI()
		ui.PrintWarning("中断信号接收，正在取消...")
		cancel()
	}()

	return ctx, cancel
}
