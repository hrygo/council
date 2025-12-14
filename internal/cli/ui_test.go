package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/debate"
)

func TestNewUI(t *testing.T) {
	var out, errOut bytes.Buffer
	ui := NewUI(&out, &errOut)

	if ui == nil {
		t.Fatal("NewUI() returned nil")
	}
	if ui.out != &out {
		t.Error("NewUI() did not set out correctly")
	}
	if ui.err != &errOut {
		t.Error("NewUI() did not set err correctly")
	}
}

func TestDefaultUI(t *testing.T) {
	ui := DefaultUI()
	if ui == nil {
		t.Fatal("DefaultUI() returned nil")
	}
}

func TestUI_PrintBanner(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintBanner()

	output := out.String()
	// Check for banner elements (ASCII art uses block characters)
	if !strings.Contains(output, "██") {
		t.Error("PrintBanner() should contain ASCII art block characters")
	}
	if !strings.Contains(output, "Multi-Persona") && !strings.Contains(output, "AI") {
		t.Error("PrintBanner() should contain 'Multi-Persona' or 'AI'")
	}
}

func TestUI_PrintConfig(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})
	cfg := config.New()

	ui.PrintConfig(cfg)

	output := out.String()
	// Check for configuration elements
	if !strings.Contains(output, "PRO") && !strings.Contains(output, "正方") {
		t.Error("PrintConfig() should contain 'PRO' or '正方'")
	}
	if !strings.Contains(output, "CON") && !strings.Contains(output, "反方") {
		t.Error("PrintConfig() should contain 'CON' or '反方'")
	}
	if !strings.Contains(output, "ADJ") && !strings.Contains(output, "裁决") {
		t.Error("PrintConfig() should contain 'ADJ' or '裁决'")
	}
}

func TestUI_PrintDebating(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintDebating()

	output := out.String()
	// Check for debating status indicators
	if !strings.Contains(output, "DEBATE") && !strings.Contains(output, "正方") {
		t.Error("PrintDebating() should contain debating status")
	}
}

func TestUI_PrintComplete(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintComplete()

	output := out.String()
	if !strings.Contains(output, "COMPLETE") && !strings.Contains(output, "完成") {
		t.Error("PrintComplete() should contain completion message")
	}
}

func TestUI_PrintError(t *testing.T) {
	var errOut bytes.Buffer
	ui := NewUI(&bytes.Buffer{}, &errOut)

	ui.PrintError("test error")

	output := errOut.String()
	if !strings.Contains(output, "test error") {
		t.Error("PrintError() should contain the error message")
	}
	if !strings.Contains(output, "ERROR") && !strings.Contains(output, "⚠") {
		t.Error("PrintError() should contain error indicator")
	}
}

func TestUI_PrintWarning(t *testing.T) {
	var errOut bytes.Buffer
	ui := NewUI(&bytes.Buffer{}, &errOut)

	ui.PrintWarning("test warning")

	output := errOut.String()
	if !strings.Contains(output, "test warning") {
		t.Error("PrintWarning() should contain the warning message")
	}
}

func TestUI_PrintSectionHeader(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintSectionHeader("Test Title", "🔵", ColorBlue)

	output := out.String()
	if !strings.Contains(output, "Test Title") {
		t.Error("PrintSectionHeader() should contain the title")
	}
	if !strings.Contains(output, "🔵") {
		t.Error("PrintSectionHeader() should contain the icon")
	}
	if !strings.Contains(output, "━") {
		t.Error("PrintSectionHeader() should contain separator line")
	}
}

func TestUI_PrintProHeader(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintProHeader()

	output := out.String()
	if !strings.Contains(output, "AFFIRMATIVE") && !strings.Contains(output, "正方") {
		t.Error("PrintProHeader() should contain pro title")
	}
	if !strings.Contains(output, "🟢") {
		t.Error("PrintProHeader() should contain green icon")
	}
}

func TestUI_PrintConHeader(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintConHeader()

	output := out.String()
	if !strings.Contains(output, "NEGATIVE") && !strings.Contains(output, "反方") {
		t.Error("PrintConHeader() should contain con title")
	}
	if !strings.Contains(output, "🔴") {
		t.Error("PrintConHeader() should contain red icon")
	}
}

func TestUI_PrintJudgeHeader(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintJudgeHeader()

	output := out.String()
	if !strings.Contains(output, "ADJUDICATOR") && !strings.Contains(output, "裁决") {
		t.Error("PrintJudgeHeader() should contain judge title")
	}
	if !strings.Contains(output, "⚖") {
		t.Error("PrintJudgeHeader() should contain judge icon")
	}
}

func TestUI_PrintResult(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	result := &debate.Result{
		ProArgument: "pro content",
		ConArgument: "con content",
		Verdict:     "verdict content",
	}

	ui.PrintResult(result)

	output := out.String()
	if !strings.Contains(output, "pro content") {
		t.Error("PrintResult() should contain pro argument")
	}
	if !strings.Contains(output, "con content") {
		t.Error("PrintResult() should contain con argument")
	}
	if !strings.Contains(output, "verdict content") {
		t.Error("PrintResult() should contain verdict")
	}
}

func TestUI_Print(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.Print("test content")

	if out.String() != "test content" {
		t.Errorf("Print() = %q, want %q", out.String(), "test content")
	}
}

func TestUI_Println(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.Println("test content")

	if out.String() != "test content\n" {
		t.Errorf("Println() = %q, want %q", out.String(), "test content\n")
	}
}

func TestUI_PrintDivider(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintDivider()

	output := out.String()
	if !strings.Contains(output, "─") {
		t.Error("PrintDivider() should contain divider character")
	}
}

func TestUI_PrintInfo(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintInfo("info message")

	output := out.String()
	if !strings.Contains(output, "info message") {
		t.Error("PrintInfo() should contain the message")
	}
	if !strings.Contains(output, "◈") {
		t.Error("PrintInfo() should contain info icon")
	}
}

func TestUI_PrintSuccess(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintSuccess("success message")

	output := out.String()
	if !strings.Contains(output, "success message") {
		t.Error("PrintSuccess() should contain the message")
	}
	if !strings.Contains(output, "✓") {
		t.Error("PrintSuccess() should contain success icon")
	}
}

func TestUI_PrintThinking(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintThinking("TestAgent")

	output := out.String()
	if !strings.Contains(output, "TestAgent") {
		t.Error("PrintThinking() should contain agent name")
	}
	if !strings.Contains(output, "processing") {
		t.Error("PrintThinking() should contain processing message")
	}
}

func TestColorConstants(t *testing.T) {
	// Verify color constants are defined
	if ColorReset == "" {
		t.Error("ColorReset should not be empty")
	}
	if ColorGreen == "" {
		t.Error("ColorGreen should not be empty")
	}
	if ColorRed == "" {
		t.Error("ColorRed should not be empty")
	}
	if ColorYellow == "" {
		t.Error("ColorYellow should not be empty")
	}
	if ColorBlue == "" {
		t.Error("ColorBlue should not be empty")
	}
	if ColorCyan == "" {
		t.Error("ColorCyan should not be empty")
	}
	if ColorBold == "" {
		t.Error("ColorBold should not be empty")
	}
	// Test new color constants
	if ColorBrightCyan == "" {
		t.Error("ColorBrightCyan should not be empty")
	}
	if ColorDim == "" {
		t.Error("ColorDim should not be empty")
	}
}
