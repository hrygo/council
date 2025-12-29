# SPEC-1102: System Surgeon Agent

| Metadata     | Value                              |
| :----------- | :--------------------------------- |
| **Title**    | System Surgeon Agent (The Builder) |
| **Author**   | Antigravity                        |
| **Status**   | Draft                              |
| **Priority** | P0                                 |
| **Sprint**   | Sprint 11                          |

## 1. Goal
Introduce a specialized role `system_surgeon` equipped with **Write** permissions, responsible for physically applying code/document changes based on the Adjudicator's verdict.

## 2. Agent Persona (`internal/resources/prompts/system_surgeon.md`)

```markdown
### Role
You are the **System Surgeon**, a precise engineering agent responsible for applying patches.
Your goal is to translate the "Mending Orders" from the Adjudicator into **Virtual File Changes** using the `write_file` tool.
All changes are applied to the Session's Virtual File System (VFS), creating a new version draft.

### Constraints
1.  **Precision**: Only touch the files explicitly mentioned.
2.  **Safety**: You are operating on a Virtual Draft. Do not worry about breaking the production server, but aim for correctness.
3.  **Minimalism**: Just call the tool.
```

## 3. Workflow Integration

### 3.1 Graph Update (`council_optimize`)
The workflow needs a **Conditional Switch** based on the Adjudicator's JSON output.

-   **Node: `agent_adjudicator`**
    -   Output: `{ "verdict": "apply_fix", "mending_orders": "..." }`
-   **Link**:
    -   If `verdict == "apply_fix"` -> **Go To `agent_surgeon`**.
    -   If `verdict == "discuss"` -> **Go To `human_review`** (or skip Surgeon).

*Note: For V1, we simply insert Surgeon before Loop Check.*

`Adjudicator` -> `Surgeon` -> `Loop Decision`

## 4. TDD Strategy

### 4.1 Test Cases
-   [ ] `TestSurgeon_Prompt`: Verify seeder loads the prompt correctly.
-   [ ] `TestSurgeon_Graph`: Verify the `council_optimize` reference in `seeder.go` includes the new node connection.

## 5. Implementation Steps
1.  Create `system_surgeon.md`.
2.  Update `Seeder.SeedAgents` to include Surgeon (with `capabilities: {file_access: true}`).
3.  Update `Seeder.SeedWorkflows` to re-wire the graph.
