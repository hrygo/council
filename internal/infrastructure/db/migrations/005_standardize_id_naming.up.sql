-- Rename primary keys to standardized {entity}_uuid
-- Groups
ALTER TABLE groups RENAME COLUMN id TO group_uuid;
-- Agents
ALTER TABLE agents RENAME COLUMN id TO agent_uuid;
-- Workflows
ALTER TABLE workflows RENAME COLUMN id TO workflow_uuid;
ALTER TABLE workflows RENAME COLUMN group_id TO group_uuid;
-- Workflow Templates
ALTER TABLE workflow_templates RENAME COLUMN id TO template_uuid;
-- Sessions
ALTER TABLE sessions RENAME COLUMN id TO session_uuid;
ALTER TABLE sessions RENAME COLUMN group_id TO group_uuid;
ALTER TABLE sessions RENAME COLUMN workflow_id TO workflow_uuid;
-- Session Messages
ALTER TABLE session_messages RENAME COLUMN id TO message_uuid;
ALTER TABLE session_messages RENAME COLUMN session_id TO session_uuid;
ALTER TABLE session_messages RENAME COLUMN agent_id TO agent_uuid;
-- Group Agents (Association table columns)
ALTER TABLE group_agents RENAME COLUMN group_id TO group_uuid;
ALTER TABLE group_agents RENAME COLUMN agent_id TO agent_uuid;
-- Memories
ALTER TABLE memories RENAME COLUMN id TO memory_uuid;
ALTER TABLE memories RENAME COLUMN group_id TO group_uuid;
ALTER TABLE memories RENAME COLUMN agent_id TO agent_uuid;
ALTER TABLE memories RENAME COLUMN session_id TO session_uuid;
-- LLM Tables (Logical IDs)
ALTER TABLE llm_providers RENAME COLUMN id TO provider_id;
ALTER TABLE llm_models RENAME COLUMN id TO model_id;
ALTER TABLE llm_models RENAME COLUMN provider_id TO provider_id; -- No change but for completeness

-- Migrate existing JSONB data
UPDATE workflow_templates 
SET graph_definition = jsonb_set(graph_definition, '{workflow_id}', graph_definition->'id') - 'id';

UPDATE workflows
SET graph_definition = jsonb_set(graph_definition, '{workflow_id}', graph_definition->'id') - 'id';

-- Group Default Agent IDs
ALTER TABLE groups RENAME COLUMN default_agent_ids TO default_agent_uuids;
