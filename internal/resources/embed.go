// Package resources provides embedded resources for default data seeding.
package resources

import "embed"

// PromptFiles contains embedded prompt markdown files for default agents.
//
//go:embed prompts/*.md
var PromptFiles embed.FS
