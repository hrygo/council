-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Groups Table
CREATE TABLE groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    icon VARCHAR(256),
    system_prompt TEXT,
    default_agent_ids JSONB DEFAULT '[]',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Agents Table
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(64) NOT NULL,
    avatar VARCHAR(256),
    description VARCHAR(512),
    persona_prompt TEXT NOT NULL,
    model_config JSONB NOT NULL DEFAULT '{"provider": "openai", "model": "gpt-4o", "temperature": 0.7}',
    capabilities JSONB DEFAULT '{"web_search": true, "search_provider": "tavily", "code_execution": false}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Group-Agent Association Table
CREATE TABLE group_agents (
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id) ON DELETE CASCADE,
    sort_order INT DEFAULT 0,
    PRIMARY KEY (group_id, agent_id)
);

-- Workflows Table
CREATE TABLE workflows (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    name VARCHAR(128) NOT NULL,
    graph_definition JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Workflow Templates Table
CREATE TABLE workflow_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    graph_definition JSONB NOT NULL,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Sessions Table
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    workflow_id UUID REFERENCES workflows(id),
    status VARCHAR(32) DEFAULT 'pending', -- pending, running, paused, completed, error
    proposal JSONB, -- {text: string, files: string[]}
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ
);

-- Session Messages Table
CREATE TABLE session_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(id) ON DELETE CASCADE,
    node_id VARCHAR(64) NOT NULL,
    agent_id UUID REFERENCES agents(id),
    content TEXT NOT NULL,
    token_count INT DEFAULT 0,
    is_thinking BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_session_messages_session ON session_messages(session_id);

-- Memories Table (Vector Store)
CREATE TABLE memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id), -- NULL means group memory
    session_id UUID REFERENCES sessions(id),
    content TEXT NOT NULL,
    embedding VECTOR(1536), -- 1536 dim for text-embedding-ada-002
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_memories_group ON memories(group_id);
-- Note: ivfflat index creation might fail if table is empty on some pg versions without data, 
-- but usually it's fine to create. Ideally HNSW is better for performance but ivfflat is in schema.
-- I'll stick to TDD recommendation but verify syntax.
CREATE INDEX idx_memories_embedding ON memories USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
