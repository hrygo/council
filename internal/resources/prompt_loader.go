package resources

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// AgentConfig holds the model configuration from prompt front matter.
type AgentConfig struct {
	Name        string  `yaml:"name"`
	Provider    string  `yaml:"provider"`
	Model       string  `yaml:"model"`
	Temperature float64 `yaml:"temperature"`
	MaxTokens   int     `yaml:"max_tokens"`
	TopP        float64 `yaml:"top_p"`
}

// AgentPrompt represents a parsed prompt file with config and content.
type AgentPrompt struct {
	Config  AgentConfig
	Content string
}

// LoadPrompt reads and parses a prompt file from the embedded filesystem.
// The file must have YAML front matter delimited by "---".
func LoadPrompt(filename string) (*AgentPrompt, error) {
	data, err := PromptFiles.ReadFile("prompts/" + filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompt file %s: %w", filename, err)
	}

	// Split YAML front matter from markdown content
	// Format: ---\nyaml content\n---\nmarkdown content
	parts := bytes.SplitN(data, []byte("---"), 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid prompt format in %s: missing YAML front matter", filename)
	}

	var config AgentConfig
	if err := yaml.Unmarshal(parts[1], &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML front matter in %s: %w", filename, err)
	}

	return &AgentPrompt{
		Config:  config,
		Content: string(bytes.TrimSpace(parts[2])),
	}, nil
}

// LoadAllPrompts loads all default agent prompts.
func LoadAllPrompts() (map[string]*AgentPrompt, error) {
	files := []string{
		"system_affirmative.md",
		"system_negative.md",
		"system_adjudicator.md",
		"system_surgeon.md",
	}

	prompts := make(map[string]*AgentPrompt)
	for _, file := range files {
		prompt, err := LoadPrompt(file)
		if err != nil {
			return nil, err
		}
		// Use filename without extension as key
		key := file[:len(file)-3] // Remove ".md"
		prompts[key] = prompt
	}

	return prompts, nil
}
