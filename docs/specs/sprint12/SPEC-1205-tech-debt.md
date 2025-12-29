# SPEC-1205: Technical Debt & Naming Standardization

| Metadata     | Value                                   |
| :----------- | :-------------------------------------- |
| **Title**    | Technical Debt & Naming Standardization |
| **Type**     | Refactor                                |
| **Sprint**   | Sprint 12                               |
| **Priority** | P1                                      |

## 1. Context
During Sprint 11 audit, we identified several naming inconsistencies and technical debt items accumulated during the "ID Refactor" and "VFS Implementation" phases. To ensure long-term maintainability and adherence to `GEMINI.md` standards, these must be addressed.

## 2. Problem Statement

### 2.1 ID Naming Inconsistency
- **Workflow**: The API (`GraphDefinition`) uses `workflow_id`, while the Database uses `workflow_uuid`.
- **Frontend**: The frontend codebase largely uses `workflow_id` but `group_uuid`. This mix creates cognitive load.
- **Goal**: Standardize on `*_uuid` for all Entity IDs across DB, API, and Frontend.

### 2.2 Implicit Parameter Coupling
- **Problem**: `WorkflowHandler.Execute` implicitly reads `req.Input["group_id"]`.
- **Risk**: Since `req.Input` is `map[string]interface{}`, a frontend rename to `group_uuid` (which happened in some places) would cause a silent failure (nil value) in the backend handler.

### 2.3 JSON Compliance
- **Problem**: `Session.ContextData` lacks a JSON tag, defaulting to PascalCase (`ContextData`) in serialization, violating the `snake_case` rule.

## 3. Implementation Plan

### 3.1 Backend: Struct Tags & Logic
1.  **Session Struct**: Add `json:"context_data"` to `ContextData`.
2.  **Workflow Structs**:
    - Update `GraphDefinition` tag: `ID` field -> `json:"workflow_uuid"` (breaking change).
    - Update `ListWorkflowsResponse`: Ensure `workflow_uuid` is consistently used.
3.  **Handler Logic**:
    - Update `WorkflowHandler.Execute`:
        - Explicitly check for `group_uuid` in `Input`.
        - Add fallback to `group_id` for backward compatibility during migration, then log a deprecation warning.

### 3.2 Frontend: Type Sync
1.  Run `tygo generate` to update TypeScript definitions.
2.  Refactor Frontend code to rely on generated types.
3.  Search & Replace `workflow_id` -> `workflow_uuid` in:
    - `useSessionStore.ts`
    - `SessionStarter.tsx`
    - `WorkflowEditor.tsx`

### 3.3 Verification
- Run `make verify-types` to ensure clean sync.
- Execute E2E tests to confirm no regressions in binding.

## 4. Acceptance Criteria
- [ ] `Session.ContextData` is serialized as `context_data`.
- [ ] `GraphDefinition` serialization uses `workflow_uuid`.
- [ ] Frontend strictly uses `workflow_uuid` and `group_uuid`.
- [ ] No direct `Input["group_id"]` string casting without nil checks in Backend.
