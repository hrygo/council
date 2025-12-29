-- Down Migration for 001_v2_schema_init

DROP TABLE IF EXISTS session_files;
DROP TABLE IF EXISTS llm_models;
DROP TABLE IF EXISTS llm_providers;
DROP TABLE IF EXISTS quarantine_logs;
DROP TABLE IF EXISTS memories;
DROP TABLE IF EXISTS session_messages;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS workflow_templates;
DROP TABLE IF EXISTS workflows;
DROP TABLE IF EXISTS group_agents;
DROP TABLE IF EXISTS agents;
DROP TABLE IF EXISTS groups;
DROP EXTENSION IF EXISTS vector;
