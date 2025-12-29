# SPEC-1100: Foundation Refactors (Relational VFS)

| Metadata     | Value                                |
| :----------- | :----------------------------------- |
| **Title**    | Foundation Refactors: Relational VFS |
| **Author**   | Antigravity                          |
| **Status**   | Approved (V3 Audit)                  |
| **Priority** | P0 (Blocker)                         |
| **Sprint**   | Sprint 11                            |

## 1. Goal
Implement a **Scalable, Concurrency-Safe Virtual File System (VFS)** using a relational database approach, replacing the rejected JSONB/Context proposal.

## 2. Changes

### 2.1 Configuration: Workspace Root
**File**: `internal/pkg/config/config.go`
-   **Add**: `WorkspaceRoot` (Legacy support/Default to `./workspace` but largely superseded by VFS for Agents).

### 2.2 LLM Infrastructure: Tool Types
**File**: `internal/infrastructure/llm/llm.go`
-   **Update**: Add `Tools` and `ToolCalls` support to `CompletionRequest` and `Message` structs (Standard Re-Act Protocol).

### 2.3 Persistence: Relational VFS (The Core Change)
**File**: `internal/infrastructure/db/migrations/006_create_session_files.up.sql`

```sql
CREATE TABLE session_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID REFERENCES sessions(session_uuid) ON DELETE CASCADE,
    path VARCHAR(255) NOT NULL,
    version INT NOT NULL,
    content TEXT NOT NULL,
    author VARCHAR(64), -- e.g., "agent:surgeon", "user"
    reason VARCHAR(255), -- e.g., "Fixing compilation error"
    created_at TIMESTAMPTZ DEFAULT NOW(),
    -- Composite Unique Key ensures strict versioning and prevents race conditions on numeric versions
    -- Note: Application logic must handle allocating the next version number atomically or robustly.
    UNIQUE(session_id, path, version)
);
CREATE INDEX idx_session_files_lookup ON session_files(session_id, path);
```

### 2.4 Repository Layer
**File**: `internal/core/workflow/file_repository.go` & `infrastructure/persistence/file_repository.go`

```go
type SessionFileRepository interface {
    // AddVersion creates a new version. Application should determine 'version' = max+1.
    // Returns error if version conflict (Optimistic Locking).
    AddVersion(ctx context.Context, sessionID, path, content, author, reason string) (version int, err error)
    
    // GetLatest returns the highest version of a file.
    GetLatest(ctx context.Context, sessionID, path string) (*FileEntity, error)
    
    // ListFiles returns map of Path -> Latest FileEntity
    ListFiles(ctx context.Context, sessionID string) (map[string]*FileEntity, error)
}
```

### 2.5 Session Integration (Dependency Injection)
**File**: `internal/core/workflow/session.go`

The `Session` struct is primarily a Data Entity, but for Tools to work elegantly, we will inject the Repository interface into it (acting as a Domain Service facade).

```go
type Session struct {
    // ... existing fields
    FileRepo SessionFileRepository `json:"-"` // Injected at runtime
}

// WriteFile is a helper that delegates to FileRepo
func (s *Session) WriteFile(path, content, author, reason string) (int, error) {
    if s.FileRepo == nil {
        return 0, fmt.Errorf("file repository not injected")
    }
    return s.FileRepo.AddVersion(s.Context(), s.ID, path, content, author, reason)
}
```

### 2.6 API Layer
**Endpoint**: `GET /api/v1/sessions/:id/files`
-   Returns list of files with their latest content.
**Endpoint**: `GET /api/v1/sessions/:id/files/:path/history`
-   Returns version history for a specific file (for Diff View).

## 3. TDD Strategy

### 3.1 Tests
-   [ ] **Persistence Test**: `TestFileRepo_Concurrency`
    -   Spawn 10 goroutines trying to add version for `main.go`.
    -   Ensure versions are monotonic or handle conflicts gracefully.
-   [ ] **Session Test**: `TestSession_WriteFile`
    -   Mock `FileRepo`.
    -   Verify delegation works.

## 4. Implementation Steps
1.  **Migration**: Create SQL migration.
2.  **Repo**: Implement `SessionFileRepository` (Postgres).
3.  **Domain**: Update `Session` struct and `NewEngine` (to inject Repo).
4.  **LLM**: Update `llm.go` types.
5.  **API**: Add Endpoints.
