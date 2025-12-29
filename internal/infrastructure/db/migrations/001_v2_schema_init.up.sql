-- Squashed Migration: 001_v2_schema_init (replaces 001-006)
-- Date: 2025-12-29
-- Content: Core Entities (Standardized UUIDs), LLM Config, VFS Support

-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- 1. Groups Table
CREATE TABLE groups (
    group_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    icon VARCHAR(256),
    system_prompt TEXT,
    default_agent_uuids JSONB DEFAULT '[]', -- Was default_agent_ids
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 2. Agents Table
CREATE TABLE agents (
    agent_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) NOT NULL,
    avatar VARCHAR(256),
    description VARCHAR(512),
    persona_prompt TEXT NOT NULL,
    model_config JSONB NOT NULL DEFAULT '{"provider": "deepseek", "model": "deepseek-chat", "temperature": 0.7}',
    capabilities JSONB DEFAULT '{"web_search": true, "search_provider": "tavily", "code_execution": false}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Group-Agent Association Table
CREATE TABLE group_agents (
    group_uuid UUID REFERENCES groups(group_uuid) ON DELETE CASCADE,
    agent_uuid UUID REFERENCES agents(agent_uuid) ON DELETE CASCADE,
    sort_order INT DEFAULT 0,
    PRIMARY KEY (group_uuid, agent_uuid)
);

-- 4. Workflows Table
CREATE TABLE workflows (
    workflow_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_uuid UUID REFERENCES groups(group_uuid) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    graph_definition JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Workflow Templates Table
CREATE TABLE workflow_templates (
    template_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    graph_definition JSONB NOT NULL,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 6. Sessions Table
CREATE TABLE sessions (
    session_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_uuid UUID REFERENCES groups(group_uuid) ON DELETE CASCADE,
    workflow_uuid UUID REFERENCES workflows(workflow_uuid),
    status VARCHAR(32) DEFAULT 'pending', -- pending, running, paused, completed, error
    proposal JSONB, -- {text: string, files: string[]}
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. Session Messages Table
CREATE TABLE session_messages (
    message_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_uuid UUID REFERENCES sessions(session_uuid) ON DELETE CASCADE,
    node_id VARCHAR(64) NOT NULL,
    agent_uuid UUID REFERENCES agents(agent_uuid),
    content TEXT NOT NULL,
    token_count INT DEFAULT 0,
    is_thinking BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_session_messages_session ON session_messages(session_uuid);

-- 8. Memories Table (Vector Store)
CREATE TABLE memories (
    memory_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_uuid UUID REFERENCES groups(group_uuid) ON DELETE CASCADE,
    agent_uuid UUID REFERENCES agents(agent_uuid), -- NULL means group memory
    session_uuid UUID REFERENCES sessions(session_uuid),
    content TEXT NOT NULL,
    embedding VECTOR(1536), -- 1536 dim for text-embedding-ada-002
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_memories_group ON memories(group_uuid);
-- HNSW is preferred for production, but ivfflat is compatible with current test env setup
CREATE INDEX idx_memories_embedding ON memories USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);

-- 9. Quarantine Logs (Audit)
CREATE TABLE quarantine_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID, -- Keeping original name for now or assuming standardized? Let's stabilize to _uuid for consistency if allowed, or keep consistent with code. Code likely uses ID struct fields which map to DB columns. Let's assume code expects session_id if not updated, but Migration 005 renamed *everything*. So usage of session_uuid is safer.
    -- Wait, Migration 005 didn't rename Quarantine Logs columns? Let's check 005 content.
    -- 005 content: "ALTER TABLE sessions RENAME...", it didn't touch Quarantine.
    -- However, for consistency, let's use session_uuid.
    session_uuid UUID,
    content TEXT,
    raw_metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_quarantine_logs_session ON quarantine_logs(session_uuid);

-- 10. LLM Configuration
CREATE TABLE llm_providers (
    provider_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    icon TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE llm_models (
    model_id TEXT PRIMARY KEY,
    provider_id TEXT NOT NULL REFERENCES llm_providers(provider_id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    is_mainstream BOOLEAN DEFAULT FALSE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
CREATE INDEX idx_llm_models_provider ON llm_models(provider_id);

-- 11. VFS (Session Files)
CREATE TABLE session_files (
    file_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_uuid UUID REFERENCES sessions(session_uuid) ON DELETE CASCADE,
    path VARCHAR(255) NOT NULL,
    version INT NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(64),
    reason VARCHAR(255),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(session_uuid, path, version)
);
CREATE INDEX idx_session_files_lookup ON session_files(session_uuid, path);
