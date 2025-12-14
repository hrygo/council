// Package cli provides command-line interface functionality for Dialecta.
// This package contains UI components, input/output handling, and CLI-specific
// utilities that can be replaced or extended for other interaction modes
// (e.g., Web API, GUI, TUI).
package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/debate"
)

// ANSI color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorRed    = "\033[31m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// UI handles all user interface output for the CLI
type UI struct {
	out io.Writer
	err io.Writer
}

// NewUI creates a new UI with the specified output writers
func NewUI(out, err io.Writer) *UI {
	return &UI{
		out: out,
		err: err,
	}
}

// DefaultUI creates a UI using stdout and stderr
func DefaultUI() *UI {
	return NewUI(os.Stdout, os.Stderr)
}

// PrintBanner prints the application banner
func (u *UI) PrintBanner() {
	fmt.Fprintf(u.out, "\n%s╔══════════════════════════════════════════════════════════════╗%s\n", ColorCyan, ColorReset)
	fmt.Fprintf(u.out, "%s║           🎭 Dialecta - 多角色辩论系统                        ║%s\n", ColorCyan, ColorReset)
	fmt.Fprintf(u.out, "%s╚══════════════════════════════════════════════════════════════╝%s\n\n", ColorCyan, ColorReset)
}

// PrintConfig prints the configuration info
func (u *UI) PrintConfig(cfg *config.Config) {
	fmt.Fprintf(u.out, "%s📋 配置信息%s\n", ColorBold, ColorReset)
	fmt.Fprintf(u.out, "   正方: %s/%s\n", cfg.ProRole.Provider, cfg.ProRole.Model)
	fmt.Fprintf(u.out, "   反方: %s/%s\n", cfg.ConRole.Provider, cfg.ConRole.Model)
	fmt.Fprintf(u.out, "   裁决: %s/%s\n\n", cfg.JudgeRole.Provider, cfg.JudgeRole.Model)
}

// PrintDebating prints the debating status message
func (u *UI) PrintDebating() {
	fmt.Fprintf(u.out, "%s⏳ 正反方并行辩论中...%s\n\n", ColorYellow, ColorReset)
}

// PrintComplete prints the completion message
func (u *UI) PrintComplete() {
	fmt.Fprintf(u.out, "\n%s✅ 辩论完成%s\n", ColorGreen, ColorReset)
}

// PrintError prints an error message
func (u *UI) PrintError(message string) {
	fmt.Fprintf(u.err, "%s❌ %s%s\n", ColorRed, message, ColorReset)
}

// PrintWarning prints a warning message
func (u *UI) PrintWarning(message string) {
	fmt.Fprintf(u.err, "\n%s⚠️ %s%s\n", ColorYellow, message, ColorReset)
}

// PrintSectionHeader prints a section header with the given title, icon and color
func (u *UI) PrintSectionHeader(title, icon, color string) {
	fmt.Fprintf(u.out, "\n%s%s%s %s%s\n", ColorBold, color, icon, title, ColorReset)
	fmt.Fprintln(u.out, strings.Repeat("─", 60))
}

// PrintProHeader prints the affirmative (pro) section header
func (u *UI) PrintProHeader() {
	u.PrintSectionHeader("正方论述 (The Affirmative)", "🟢", ColorGreen)
}

// PrintConHeader prints the negative (con) section header
func (u *UI) PrintConHeader() {
	u.PrintSectionHeader("反方论述 (The Negative)", "🔴", ColorRed)
}

// PrintJudgeHeader prints the adjudicator (judge) section header
func (u *UI) PrintJudgeHeader() {
	u.PrintSectionHeader("裁决方报告 (The Adjudicator)", "⚖️", ColorBlue)
}

// PrintResult prints the complete debate result (non-streaming mode)
func (u *UI) PrintResult(result *debate.Result) {
	u.PrintProHeader()
	fmt.Fprintln(u.out, result.ProArgument)

	u.PrintConHeader()
	fmt.Fprintln(u.out, result.ConArgument)

	u.PrintJudgeHeader()
	fmt.Fprintln(u.out, result.Verdict)
}

// Print writes content to the output
func (u *UI) Print(content string) {
	fmt.Fprint(u.out, content)
}

// Println writes content to the output with a newline
func (u *UI) Println(content string) {
	fmt.Fprintln(u.out, content)
}
