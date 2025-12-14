package cli

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/term"
)

// DualPrinter handles vertical split streaming output (Top/Bottom boxes)
type DualPrinter struct {
	mu         sync.Mutex
	out        *os.File
	termWidth  int
	termHeight int

	proLines []string
	conLines []string

	// Configuration
	boxHeight    int // Height of content area inside the box
	contentWidth int // Width of text area

	// State
	isStarted bool
}

func NewDualPrinter(out *os.File) *DualPrinter {
	w, h, err := term.GetSize(int(out.Fd()))
	if err != nil {
		w, h = 80, 24
	}

	// Calculate Box Height based on terminal height
	// We need 2 boxes. Each box has BorderTop + Content + BorderBot = Content + 2 lines.
	// Plus 1 spacer line. Total overhead = 2*2 + 1 = 5 lines.
	// Available for content = H - 5.
	// content per box = (H - 5) / 2.

	targetH := (h - 6) / 2
	if targetH < 5 {
		targetH = 5 // Minimum viable height
	}
	if targetH > 15 {
		targetH = 15 // Cap maximum height to avoid taking over huge screens completely
	}

	// Content width = Width - 4 (Border chars + padding)
	cW := w - 4
	if cW < 20 {
		cW = 20
	}

	return &DualPrinter{
		out:          out,
		termWidth:    w,
		termHeight:   h,
		boxHeight:    targetH,
		contentWidth: cW,
		proLines:     []string{""},
		conLines:     []string{""},
	}
}

func (dp *DualPrinter) Start() {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	if dp.isStarted {
		return
	}

	// Reserve screen real estate by printing newlines
	// Total Height = (BoxHeight + 2) * 2 + 1
	totalLines := (dp.boxHeight+2)*2 + 1
	fmt.Fprint(dp.out, strings.Repeat("\n", totalLines))

	// Draw initial empty boxes
	dp.drawBox(true)  // Draw Pro
	dp.drawBox(false) // Draw Con

	dp.isStarted = true
}

func (dp *DualPrinter) UpdatePro(text string) {
	dp.update(true, text)
}

func (dp *DualPrinter) UpdateCon(text string) {
	dp.update(false, text)
}

func (dp *DualPrinter) update(isPro bool, text string) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	runes := []rune(text)
	for _, r := range runes {
		dp.appendRune(isPro, r)
	}

	// Only redraw the modified box
	dp.drawBox(isPro)
}

func (dp *DualPrinter) appendRune(isPro bool, r rune) {
	var lines *[]string
	if isPro {
		lines = &dp.proLines
	} else {
		lines = &dp.conLines
	}

	// Get current line
	idx := len(*lines) - 1
	line := (*lines)[idx]

	if r == '\n' {
		*lines = append(*lines, "")
		return
	}

	rw := runeWidth(r)
	curW := stringWidth(line)

	if curW+rw > dp.contentWidth {
		*lines = append(*lines, string(r))
	} else {
		(*lines)[idx] += string(r)
	}
}

func (dp *DualPrinter) drawBox(isPro bool) {
	// 1. Save Cursor position (Assuming we are at the bottom because we always restore downstream)
	// Actually, standard practice is:
	// We allocated N lines. We assume cursor is at the line AFTER the last box.

	// Calculate Up Distance
	// Layout (Bottom Up):
	// [Cursor]
	// Con Border Bot (1)
	// Con Content (H)
	// Con Header (1)
	// Spacer (1)
	// Pro Border Bot (1)
	// Pro Content (H)
	// Pro Header (1)

	var upLines int
	boxTotal := dp.boxHeight + 2

	if isPro {
		// To redraw Pro, we go to top of Pro.
		// Dist = ConTotal + Spacer + ProTotal
		// Start of Pro is at Top.
		// But we draw from top down.
		// Start line of Pro is (ConTotal + Spacer + ProTotal) lines UP from current.
		upLines = boxTotal + 1 + boxTotal
	} else {
		// Con
		upLines = boxTotal
	}

	// Move Up
	fmt.Fprintf(dp.out, "\033[%dA", upLines)

	// Draw Box
	// Header
	var title, borderColor, titleColor string
	if isPro {
		title = "🟢 正方 (PRO)"
		borderColor = ColorBrightGreen
		titleColor = ColorBrightGreen // Or use Bold
	} else {
		title = "🔴 反方 (CON)"
		borderColor = ColorBrightRed
		titleColor = ColorBrightRed
	}

	// Header Line: ┌─ TITLE ──...─┐
	// Architecture: ┌─ (2) + Space (1) + Title (W) + Space (1) + Dashes (N) + ┐ (1)
	// Total Width used explicitly = 2 + 1 + W + 1 + 1 = 5 + W.
	// Remaining for dashes = TermWidth - (5 + W).

	headerBarLen := dp.termWidth - (5 + stringWidth(title))
	if headerBarLen < 0 {
		headerBarLen = 0
	}

	fmt.Fprintf(dp.out, "\r%s┌─ %s%s %s┐%s\n",
		borderColor,
		titleColor, title,
		borderColor+strings.Repeat("─", headerBarLen),
		ColorReset,
	)

	// Content with Clipping/Scrolling
	lines := dp.conLines
	if isPro {
		lines = dp.proLines
	}

	// Determine visible window (Tail)
	startIdx := 0
	if len(lines) > dp.boxHeight {
		startIdx = len(lines) - dp.boxHeight
	}
	visibleLines := lines[startIdx:]

	// Draw Content Rows
	for i := 0; i < dp.boxHeight; i++ {
		text := ""
		if i < len(visibleLines) {
			text = visibleLines[i]
		}

		// Fill width
		padding := dp.contentWidth - stringWidth(text)
		if padding < 0 {
			padding = 0
		} // Should not happen due to wrap logic

		fmt.Fprintf(dp.out, "\r%s│ %s%s%s%s%s │%s\n",
			borderColor,
			ColorReset,
			text,
			strings.Repeat(" ", padding),
			borderColor, // Restore border color for │
			"",          // extra arg
			ColorReset,
		)
	}

	// Footer Line
	fmt.Fprintf(dp.out, "\r%s└%s┘%s",
		borderColor,
		strings.Repeat("─", dp.termWidth-2),
		ColorReset,
	)

	// Restore Down
	// We just printed (1 + H + 1) lines. We moved Up `upLines`.
	// Current pos is at end of Footer.
	// We need to move down by remaining distance.
	// Total Up was `upLines`. We printed `boxTotal`.
	// Remaining Down = upLines - boxTotal. NOT QUITE.
	// Because printing moves cursor down.
	// Initial: Y
	// Move Up X: Y-X
	// Print N: Y-X+N
	// Target: Y.
	// Need to move down: Y - (Y-X+N) = X - N.

	downNeeded := upLines - boxTotal
	if downNeeded > 0 {
		fmt.Fprintf(dp.out, "\n\033[%dB", downNeeded)
		// Note: \n moves down 1. So (downNeeded - 1)?
		// Actually, Printf("\n") already is 1 line.
		// Let's use Move Down explicitly without newline to be safe?
		// No, the last Fprintf didn't verify newline at end?
		// "└...┘" -> No newline at end in my code.
		// So current is on same line as footer.
		// So we actually outputted `boxTotal - 1` newlines. (Header, H content lines).
		// Footer is just printed, cursor at end of footer line.

		// So we are at Box Bottom Line.
		// If isPro: we are at Pro Bottom. We need to jump over Spacer (1) + Con (BoxTotal).
		// Distance = 1 + BoxTotal.
		fmt.Fprintf(dp.out, "\n\033[%dB", downNeeded-1) // \n takes 1 line
	} else {
		// If Con (downNeeded == 0)
		// We are at Con Footer.
		// Cursor should be at line BELOW Con Footer.
		fmt.Fprint(dp.out, "\n")
	}

	// Reset column
	fmt.Fprint(dp.out, "\r")
}

// Helpers retained
func stringWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

func runeWidth(r rune) int {
	if utf8.RuneLen(r) > 1 {
		return 2
	}
	return 1
}
