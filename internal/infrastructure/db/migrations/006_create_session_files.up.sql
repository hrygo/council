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
