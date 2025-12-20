-- Add updated_at to workflow_templates
ALTER TABLE workflow_templates ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();

-- Add updated_at to sessions
ALTER TABLE sessions ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
