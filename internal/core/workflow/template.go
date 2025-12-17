package workflow

import (
	"time"
)

type TemplateCategory string

const (
	TemplateCategoryCodeReview    TemplateCategory = "code_review"
	TemplateCategoryBusinessPlan  TemplateCategory = "business_plan"
	TemplateCategoryQuickDecision TemplateCategory = "quick_decision"
	TemplateCategoryCustom        TemplateCategory = "custom"
	TemplateCategoryOther         TemplateCategory = "other"
)

type Template struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Category    TemplateCategory `json:"category"`
	IsSystem    bool             `json:"is_system"`
	Graph       GraphDefinition  `json:"graph"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
