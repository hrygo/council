package debate

import (
	"strings"
	"sync"
)

// StreamParser processes the stream to extract One-Liner early
type StreamParser struct {
	mu           sync.Mutex
	buffer       strings.Builder
	oneLiner     string
	fullBody     string
	oneLinerSent bool
	delimiter    string
}

func NewStreamParser(delimiter string) *StreamParser {
	return &StreamParser{
		delimiter: delimiter,
	}
}

// Feed adds a chunk and returns the One-Liner if it just became available
func (p *StreamParser) Feed(chunk string) (string, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.buffer.WriteString(chunk)
	current := p.buffer.String()

	// Heuristic: Look for the start of Full Argument section
	if !p.oneLinerSent {
		if strings.Contains(current, p.delimiter) {
			// Extract everything before this marker
			parts := strings.SplitN(current, p.delimiter, 2)
			if len(parts) > 0 {
				// Clean up the One-Liner section
				p.oneLiner = extractOneLinerContent(parts[0])
				p.oneLinerSent = true
				return p.oneLiner, true
			}
		}
	}

	return "", false
}

// Finalize parses the full buffer at the end to ensure everything is captured
func (p *StreamParser) Finalize() {
	p.mu.Lock()
	defer p.mu.Unlock()

	fullText := p.buffer.String()

	// Split by markers
	// 1. One-Liner
	// 2. Full Argument

	parts := strings.Split(fullText, p.delimiter)
	if len(parts) >= 2 {
		p.oneLiner = extractOneLinerContent(parts[0])
		p.fullBody = strings.TrimSpace(parts[1])
	} else {
		// Fallback if model failed format slightly
		p.fullBody = fullText
	}
}

func extractOneLinerContent(raw string) string {
	// Raw contains: "## 💡 One-Liner\n(Content)\n..."
	// We remove the header "## 💡 One-Liner"
	// And remove any parentheses if present in prompt placeholder

	lines := strings.Split(raw, "\n")
	var cleaned []string
	headerFound := false

	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if strings.Contains(trim, "## 💡 One-Liner") {
			headerFound = true
			continue
		}
		if !headerFound {
			continue // Skip everything before header
		}
		if trim == "" {
			continue
		}

		// Remove placeholder text like "(在此处写下一句...)"
		if strings.HasPrefix(trim, "(") && strings.HasSuffix(trim, ")") {
			continue
		}

		cleaned = append(cleaned, trim)
	}

	return strings.Join(cleaned, " ")
}
