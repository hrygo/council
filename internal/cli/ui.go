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
	"time"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/debate"
)

// ANSI color and style codes for terminal output
const (
	// Reset
	ColorReset = "\033[0m"

	// Regular colors
	ColorBlack   = "\033[30m"
	ColorRed     = "\033[31m"
	ColorGreen   = "\033[32m"
	ColorYellow  = "\033[33m"
	ColorBlue    = "\033[34m"
	ColorMagenta = "\033[35m"
	ColorCyan    = "\033[36m"
	ColorWhite   = "\033[37m"

	// Bright/Bold colors
	ColorBrightBlack   = "\033[90m"
	ColorBrightRed     = "\033[91m"
	ColorBrightGreen   = "\033[92m"
	ColorBrightYellow  = "\033[93m"
	ColorBrightBlue    = "\033[94m"
	ColorBrightMagenta = "\033[95m"
	ColorBrightCyan    = "\033[96m"
	ColorBrightWhite   = "\033[97m"

	// Styles
	ColorBold      = "\033[1m"
	ColorDim       = "\033[2m"
	ColorItalic    = "\033[3m"
	ColorUnderline = "\033[4m"

	// Background colors
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

// Gradient characters for visual effects
var gradientChars = []string{"░", "▒", "▓", "█"}

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

// PrintBanner prints a futuristic AI-themed application banner
func (u *UI) PrintBanner() {
	fmt.Fprintln(u.out)

	// ASCII Art definition
	asciiArt := []string{
		"██████╗ ██╗ █████╗ ██╗     ███████╗ ██████╗████████╗ █████╗",
		"██╔══██╗██║██╔══██╗██║     ██╔════╝██╔════╝╚══██╔══╝██╔══██╗",
		"██║  ██║██║███████║██║     █████╗  ██║        ██║   ███████║",
		"██║  ██║██║██╔══██║██║     ██╔══╝  ██║        ██║   ██╔══██║",
		"██████╔╝██║██║  ██║███████╗███████╗╚██████╗   ██║   ██║  ██║",
		"╚═════╝ ╚═╝╚═╝  ╚═╝╚══════╝╚══════╝ ╚═════╝   ╚═╝   ╚═╝  ╚═╝",
	}

	// Layout constants
	borderTop := "╭─────────────────────────────────────────────────────────────╮" // 61 dashes
	borderBot := "╰─────────────────────────────────────────────────────────────╯"
	padding := "  "
	// Content width logic:
	// Border is 63 chars wide (1 + 61 + 1)
	// Inner content space is 61 chars
	// ASCII Art is ~55 chars wide visually
	// Padding left: 3 chars ("   ")
	// Padding right: 3 chars ("   ")
	// Total: 3 + 55 + 3 = 61 matches inner width

	// 1. Top Border
	fmt.Fprintf(u.out, "%s%s%s%s\n", padding, ColorBrightCyan, borderTop, ColorReset)

	// Helper to print a line with correct coloring
	printLine := func(content string, color string) {
		// Left Border: Padding + Cyan │ + Reset
		fmt.Fprintf(u.out, "%s%s│%s", padding, ColorBrightCyan, ColorReset)

		// Content: BgBlack + 3 spaces + Color Content + BgBlack 3 spaces + Reset
		// Note: We avoid ColorReset inside the content to keep BgBlack active.
		// Instead, we just set the color we want.
		// But wait, 'content' is printed with 'color'. After content, we need to ensure right padding has BgBlack.
		// And we need to ensure 'content' printing doesn't reset bg.
		fmt.Fprintf(u.out, "%s   %s%s   %s", BgBlack, color, content, ColorReset)

		// Right Border: Cyan │ + Reset
		fmt.Fprintf(u.out, "%s│%s\n", ColorBrightCyan, ColorReset)
	}

	// 2. ASCII Art
	for i := 0; i < 2; i++ {
		printLine(asciiArt[i], ColorBrightMagenta)
	}
	for i := 2; i < 4; i++ {
		printLine(asciiArt[i], ColorBrightCyan)
	}
	for i := 4; i < 6; i++ {
		printLine(asciiArt[i], ColorBrightBlue)
	}

	// 3. Spacing line
	fmt.Fprintf(u.out, "%s%s│%s%s%s%s%s│%s\n",
		padding, ColorBrightCyan, ColorReset,
		BgBlack, strings.Repeat(" ", 61), ColorReset,
		ColorBrightCyan, ColorReset)

	// 4. Metadata Footer
	// Line 1: "  ◆ Multi-Persona AI Debate System            ► v1.0.0     "
	// We use colors directly without intermediate resets to maintain BgBlack
	// Let's rewrite the format string to be simpler and safer
	// Line 1
	fmt.Fprintf(u.out, "%s%s│%s%s  %s◆ Multi-Persona AI Debate System           %s► v1.0.0      %s%s│%s\n",
		padding, ColorBrightCyan, ColorReset,
		BgBlack,
		ColorBrightWhite,
		ColorDim,
		ColorReset, // End of content
		ColorBrightCyan, ColorReset)

	// Line 2
	// "  ◆ Powered by DeepSeek × Gemini × Qwen                     "
	fmt.Fprintf(u.out, "%s%s│%s%s  %s◆ Powered by DeepSeek × Gemini × Qwen                     %s%s│%s\n",
		padding, ColorBrightCyan, ColorReset,
		BgBlack,
		ColorBrightYellow,
		ColorReset,
		ColorBrightCyan, ColorReset)

	// 5. Bottom Border
	fmt.Fprintf(u.out, "%s%s%s%s\n\n", padding, ColorBrightCyan, borderBot, ColorReset)
}

// PrintConfig prints the configuration info with a modern card-style layout
func (u *UI) PrintConfig(cfg *config.Config) {
	fmt.Fprintf(u.out, "%s%s┌─ 🧠 AI Configuration ─────────────────────────────────────────┐%s\n", ColorBrightBlue, ColorBold, ColorReset)

	// Pro role
	fmt.Fprintf(u.out, "%s│%s  %s▹ PRO%s  %s%-12s%s │ %s%s%s\n",
		ColorBrightBlue, ColorReset,
		ColorBrightGreen, ColorReset,
		ColorBold, cfg.ProRole.Provider, ColorReset,
		ColorDim, cfg.ProRole.Model, ColorReset)

	// Con role
	fmt.Fprintf(u.out, "%s│%s  %s▹ CON%s  %s%-12s%s │ %s%s%s\n",
		ColorBrightBlue, ColorReset,
		ColorBrightRed, ColorReset,
		ColorBold, cfg.ConRole.Provider, ColorReset,
		ColorDim, cfg.ConRole.Model, ColorReset)

	// Judge role
	fmt.Fprintf(u.out, "%s│%s  %s▹ ADJ%s  %s%-12s%s │ %s%s%s\n",
		ColorBrightBlue, ColorReset,
		ColorBrightYellow, ColorReset,
		ColorBold, cfg.JudgeRole.Provider, ColorReset,
		ColorDim, cfg.JudgeRole.Model, ColorReset)

	fmt.Fprintf(u.out, "%s%s└───────────────────────────────────────────────────────────────┘%s\n\n", ColorBrightBlue, ColorBold, ColorReset)
}

// PrintDebating prints the debating status with animated-style indicators
func (u *UI) PrintDebating() {
	fmt.Fprintf(u.out, "%s%s◉ INITIATING PARALLEL DEBATE SEQUENCE...%s\n", ColorBrightYellow, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s  ├─ 🟢 正方 Agent: Generating affirmative arguments...%s\n", ColorDim, ColorReset)
	fmt.Fprintf(u.out, "%s  └─ 🔴 反方 Agent: Generating counter-arguments...%s\n\n", ColorDim, ColorReset)
}

// PrintComplete prints the completion message with a success indicator
func (u *UI) PrintComplete() {
	fmt.Fprintln(u.out)
	fmt.Fprintf(u.out, "%s%s╔══════════════════════════════════════════════════════════════╗%s\n", ColorBrightGreen, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s%s║                    ✓ DEBATE COMPLETE                          ║%s\n", ColorBrightGreen, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s%s╚══════════════════════════════════════════════════════════════╝%s\n", ColorBrightGreen, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%sSession ended at %s%s\n", ColorDim, time.Now().Format("2006-01-02 15:04:05"), ColorReset)
}

// PrintError prints an error message with a distinctive style
func (u *UI) PrintError(message string) {
	fmt.Fprintf(u.err, "\n%s%s┌─ ⚠ ERROR ─────────────────────────────────────────────────────┐%s\n", ColorBrightRed, ColorBold, ColorReset)
	fmt.Fprintf(u.err, "%s│%s %s%s\n", ColorBrightRed, ColorReset, message, ColorReset)
	fmt.Fprintf(u.err, "%s%s└───────────────────────────────────────────────────────────────┘%s\n\n", ColorBrightRed, ColorBold, ColorReset)
}

// PrintWarning prints a warning message
func (u *UI) PrintWarning(message string) {
	fmt.Fprintf(u.err, "\n%s%s⚡ %s%s\n", ColorBrightYellow, ColorBold, message, ColorReset)
}

// PrintSectionHeader prints a section header with the given title, icon and color
func (u *UI) PrintSectionHeader(title, icon, color string) {
	fmt.Fprintln(u.out)
	fmt.Fprintf(u.out, "%s%s%s ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ %s%s\n", color, ColorBold, icon, icon, ColorReset)
	fmt.Fprintf(u.out, "%s%s  %s%s\n", color, ColorBold, title, ColorReset)
	fmt.Fprintf(u.out, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", color, ColorReset)
}

// PrintProHeader prints the affirmative (pro) section header
func (u *UI) PrintProHeader() {
	fmt.Fprintln(u.out)
	fmt.Fprintf(u.out, "%s%s🟢 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 🟢%s\n", ColorBrightGreen, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s%s   AFFIRMATIVE ARGUMENT │ 正方论述%s\n", ColorBrightGreen, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorGreen, ColorReset)
}

// PrintConHeader prints the negative (con) section header
func (u *UI) PrintConHeader() {
	fmt.Fprintln(u.out)
	fmt.Fprintf(u.out, "%s%s🔴 ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 🔴%s\n", ColorBrightRed, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s%s   NEGATIVE ARGUMENT │ 反方论述%s\n", ColorBrightRed, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorRed, ColorReset)
}

// PrintJudgeHeader prints the adjudicator (judge) section header
func (u *UI) PrintJudgeHeader() {
	fmt.Fprintln(u.out)
	fmt.Fprintf(u.out, "%s%s⚖️  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ ⚖️%s\n", ColorBrightYellow, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s%s   ADJUDICATOR'S VERDICT │ 裁决报告%s\n", ColorBrightYellow, ColorBold, ColorReset)
	fmt.Fprintf(u.out, "%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", ColorYellow, ColorReset)
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

// PrintDivider prints a subtle divider line
func (u *UI) PrintDivider() {
	fmt.Fprintf(u.out, "%s%s%s\n", ColorDim, strings.Repeat("─", 66), ColorReset)
}

// PrintInfo prints an info message
func (u *UI) PrintInfo(message string) {
	fmt.Fprintf(u.out, "%s%s◈ %s%s\n", ColorBrightCyan, ColorBold, message, ColorReset)
}

// PrintSuccess prints a success message
func (u *UI) PrintSuccess(message string) {
	fmt.Fprintf(u.out, "%s%s✓ %s%s\n", ColorBrightGreen, ColorBold, message, ColorReset)
}

// PrintThinking prints a "thinking" indicator for AI processing
func (u *UI) PrintThinking(agentName string) {
	fmt.Fprintf(u.out, "%s◐ %s is processing...%s\n", ColorDim, agentName, ColorReset)
}
