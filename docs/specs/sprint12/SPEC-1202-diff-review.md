# SPEC-1202: Advanced Human Review (Diff-Based)

| Metadata     | Value                    |
| :----------- | :----------------------- |
| **Title**    | Advanced Human Review UI |
| **Type**     | Frontend Feature         |
| **Sprint**   | Sprint 12                |
| **Priority** | P0                       |

## 1. Goal
Upgrade the generic `HumanReviewModal` to support **Code-Aware Reviews**. When the Surgeon proposes changes, the human must see the *Impact* before approving.

## 2. Problem
Currently, the Review Modal only shows text prompts. For the Surgeon, the output is a *File Change* (Tool Call). We need to render this specific payload.

## 3. Logic
1.  **Detection**: The `HumanReviewProcessor` receives input from the previous node.
    - If previous node was `agent_surgeon` (or output contains `file_diff`), the Modal enters "Code Review Mode".
2.  **Display**:
    - Instead of just "Approve/Reject", show the **Diff** of the proposed change.
    - Use the VFS Explorer component (SPEC-1201) in "Draft Mode" (comparing current Head vs Proposed).
3.  **Actions**:
    - **Approve**: Applies the change (Surgeon tool executes).
    - **Reject**: Feedback sent back to Surgeon.
    - **Modify**: Human edits the code directly in the Modal (Override).

## 4. Implementation details
- **Payload**: The `HumanReview` node payload must include the `tool_call` arguments (path, content).
- **Frontend**:
    - Parse `payload.tool_calls` from the WebSocket event.
    - If `write_file` detected, render Code Editor.

## 5. Acceptance Criteria
- [ ] Review Modal identifies Surgeon's `write_file` intent.
- [ ] Modal displays Code Diff (Current vs New Content).
- [ ] User can edit the "New Content" before approving.
