package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// InputReader handles reading material input from various sources
type InputReader struct {
	stdin io.Reader
	out   io.Writer
}

// NewInputReader creates a new InputReader
func NewInputReader(stdin io.Reader, out io.Writer) *InputReader {
	return &InputReader{
		stdin: stdin,
		out:   out,
	}
}

// DefaultInputReader creates an InputReader using os.Stdin and os.Stdout
func DefaultInputReader() *InputReader {
	return NewInputReader(os.Stdin, os.Stdout)
}

// ReadFile reads material from a file
func (r *InputReader) ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	return string(data), nil
}

// ReadStdin reads material from stdin (pipe mode)
func (r *InputReader) ReadStdin() (string, error) {
	data, err := io.ReadAll(r.stdin)
	if err != nil {
		return "", fmt.Errorf("读取标准输入失败: %w", err)
	}
	return string(data), nil
}

// ReadInteractive reads material interactively from the user
// The user can finish input by entering two consecutive empty lines
func (r *InputReader) ReadInteractive() (string, error) {
	fmt.Fprintf(r.out, "%s📝 请输入待分析的材料（输入两个空行结束）:%s\n", ColorCyan, ColorReset)
	fmt.Fprintln(r.out, strings.Repeat("─", 40))

	var lines []string
	scanner := bufio.NewScanner(r.stdin)
	emptyCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			emptyCount++
			if emptyCount >= 2 {
				break
			}
		} else {
			emptyCount = 0
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取输入失败: %w", err)
	}

	// Trim trailing empty lines
	for len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return strings.Join(lines, "\n"), nil
}

// ReadMaterial reads material based on the input mode
// - If interactive is true, reads interactively
// - If source is "-", reads from stdin
// - Otherwise reads from the file at the given path
func (r *InputReader) ReadMaterial(source string, interactive bool) (string, error) {
	if interactive {
		return r.ReadInteractive()
	}
	if source == "-" {
		return r.ReadStdin()
	}
	return r.ReadFile(source)
}

// ValidateMaterial checks if the material is valid (non-empty)
func ValidateMaterial(material string) error {
	if strings.TrimSpace(material) == "" {
		return fmt.Errorf("材料内容为空")
	}
	return nil
}
