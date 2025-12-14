package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/llm"
)

// Options holds the parsed command-line options
type Options struct {
	ProProvider   string
	ProModel      string
	ConProvider   string
	ConModel      string
	JudgeProvider string
	JudgeModel    string
	Stream        bool
	Interactive   bool
	Source        string // file path, "-" for stdin, or empty for no source
}

// ParseFlags parses command-line flags and returns Options
func ParseFlags() *Options {
	opts := &Options{}

	flag.StringVar(&opts.ProProvider, "pro-provider", "deepseek", "Provider for affirmative (deepseek, gemini, dashscope)")
	flag.StringVar(&opts.ProModel, "pro-model", "", "Model for affirmative")
	flag.StringVar(&opts.ConProvider, "con-provider", "dashscope", "Provider for negative (deepseek, gemini, dashscope)")
	flag.StringVar(&opts.ConModel, "con-model", "", "Model for negative")
	flag.StringVar(&opts.JudgeProvider, "judge-provider", "gemini", "Provider for adjudicator (deepseek, gemini, dashscope)")
	flag.StringVar(&opts.JudgeModel, "judge-model", "", "Model for adjudicator")
	flag.BoolVar(&opts.Stream, "stream", true, "Enable streaming output")
	flag.BoolVar(&opts.Interactive, "interactive", false, "Interactive mode - enter material via stdin")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `%sDialecta - Multi-Persona Debate System%s

Usage:
  dialecta [options] <file>       Analyze material from file
  dialecta [options] -            Read material from stdin (pipe)
  dialecta --interactive          Interactive mode

Providers:
  deepseek   - DeepSeek API (DEEPSEEK_API_KEY)
  gemini     - Google Gemini (GEMINI_API_KEY or GOOGLE_API_KEY)
  dashscope  - Alibaba DashScope/Qwen (DASHSCOPE_API_KEY)

Examples:
  dialecta proposal.md
  cat plan.txt | dialecta -
  echo "我们应该启动AI创业项目" | dialecta -
  dialecta --judge-provider deepseek --judge-model deepseek-chat proposal.md

Options:
`, ColorBold, ColorReset)
		flag.PrintDefaults()
	}

	flag.Parse()

	// Get source from remaining arguments
	if flag.NArg() > 0 {
		opts.Source = flag.Arg(0)
	}

	return opts
}

// ApplyToConfig applies the options to a config
func (opts *Options) ApplyToConfig(cfg *config.Config) {
	if p, err := llm.ParseProvider(opts.ProProvider); err == nil {
		cfg.ProRole.Provider = p
	}
	if opts.ProModel != "" {
		cfg.ProRole.Model = opts.ProModel
	}

	if p, err := llm.ParseProvider(opts.ConProvider); err == nil {
		cfg.ConRole.Provider = p
	}
	if opts.ConModel != "" {
		cfg.ConRole.Model = opts.ConModel
	}

	if p, err := llm.ParseProvider(opts.JudgeProvider); err == nil {
		cfg.JudgeRole.Provider = p
	}
	if opts.JudgeModel != "" {
		cfg.JudgeRole.Model = opts.JudgeModel
	}
}

// NeedsHelp returns true if help should be shown (no source and not interactive)
func (opts *Options) NeedsHelp() bool {
	return opts.Source == "" && !opts.Interactive
}
