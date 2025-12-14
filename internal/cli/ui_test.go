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
	if !strings.Contains(output, "Dialecta") {
		t.Error("PrintBanner() should contain 'Dialecta'")
	}
	if !strings.Contains(output, "🎭") {
		t.Error("PrintBanner() should contain emoji")
	}
}

func TestUI_PrintConfig(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})
	cfg := config.New()

	ui.PrintConfig(cfg)

	output := out.String()
	if !strings.Contains(output, "配置信息") {
		t.Error("PrintConfig() should contain '配置信息'")
	}
	if !strings.Contains(output, "正方") {
		t.Error("PrintConfig() should contain '正方'")
	}
	if !strings.Contains(output, "反方") {
		t.Error("PrintConfig() should contain '反方'")
	}
	if !strings.Contains(output, "裁决") {
		t.Error("PrintConfig() should contain '裁决'")
	}
}

func TestUI_PrintDebating(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintDebating()

	output := out.String()
	if !strings.Contains(output, "正反方并行辩论中") {
		t.Error("PrintDebating() should contain status message")
	}
}

func TestUI_PrintComplete(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintComplete()

	output := out.String()
	if !strings.Contains(output, "辩论完成") {
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
	if !strings.Contains(output, "❌") {
		t.Error("PrintError() should contain error icon")
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
	if !strings.Contains(output, "⚠️") {
		t.Error("PrintWarning() should contain warning icon")
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
	if !strings.Contains(output, "─") {
		t.Error("PrintSectionHeader() should contain separator line")
	}
}

func TestUI_PrintProHeader(t *testing.T) {
	var out bytes.Buffer
	ui := NewUI(&out, &bytes.Buffer{})

	ui.PrintProHeader()

	output := out.String()
	if !strings.Contains(output, "正方论述") {
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
	if !strings.Contains(output, "反方论述") {
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
	if !strings.Contains(output, "裁决方报告") {
		t.Error("PrintJudgeHeader() should contain judge title")
	}
	if !strings.Contains(output, "⚖️") {
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
}
