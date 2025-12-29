# SPEC-1101: Core Tool Infrastructure & Safety Sandbox

| Metadata     | Value                                     |
| :----------- | :---------------------------------------- |
| **Title**    | Core Tool Infrastructure & Safety Sandbox |
| **Author**   | Antigravity                               |
| **Status**   | Draft                                     |
| **Priority** | P0                                        |
| **Sprint**   | Sprint 11                                 |

## 1. Goal
Upgrade the Core `AgentProcessor` to support **Tool Calling (Function Calling)** capabilities, enabling Agents to interact with the external world (FileSystem, API, etc.) securely.

## 2. Architecture

### 2.1 Interface Definition (`internal/core/workflow/nodes/tools/types.go`)

```go
package tools

import (
    "context"
)

// Tool defines an executable capability
type Tool interface {
    Name() string
    Description() string
    Parameters() map[string]interface{} // JSON Schema
    Execute(ctx context.Context, args map[string]interface{}) (string, error)
}

// Registry manages available tools
type Registry interface {
    GetTool(name string) (Tool, error)
    Register(tool Tool)
}
```

### 2.2 Execution Flow (Re-Act Pattern)

1.  `AgentProcessor` receives input.
2.  Passes `tools` definition to `LLMProvider`.
3.  **LLM Response**:
    -   *Case A*: Text content -> Finish.
    -   *Case B*: `FunctionCall` detected -> **Pause**.
4.  **Tool Execution**:
    -   Validate Tool Name and Args.
    -   Execute `tool.Execute()`.
    -   Capture output (or error).
5.  **Recursion**:
    -   Append `ToolOutput` message to conversation history.
    -   **Re-call LLM** with updated history.
    -   Repeat until LLM stops calling tools or `MaxIterations` reached.

## 3. Session-Aware Tool Execution (VFS Proxy)
Instead of a Physical Sandbox, the Tool implementation must have access to the `Session` context to call `WriteFile`.

```go
// Define a new interface for Session-Aware Tools
type SessionAwareTool interface {
    Tool
    ExecuteWithSession(ctx context.Context, session *workflow.Session, args map[string]interface{}) (string, error)
}

type WriteFileTool struct {}

func (t *WriteFileTool) ExecuteWithSession(ctx context.Context, session *workflow.Session, args map[string]interface{}) (string, error) {
    path := args["path"].(string)
    content := args["content"].(string)
    // Virtual Write
    ver, err := session.WriteFile(path, content, "agent", "tool_call")
    return fmt.Sprintf("File %s written (Version %d)", path, ver), err
}
```

This requires upgrading `AgentProcessor` to fetch `Session` from context (or struct) and check if tool implements `SessionAwareTool`.

## 4. TDD Strategy

### 4.1 Mock Strategy
We do NOT need real LLMs for testing. We will mock the `LLMProvider` interface.

**Scenario: Auto-Fix**
1.  **Mock LLM Response 1**: content="", tool_calls=[{name="write_file", args={path="main.go", content="..."}}]
2.  **System**: Executes `write_file`.
3.  **Mock LLM Response 2**: content="Fix applied."

### 4.2 Test Cases (`agent_tool_test.go`)
-   [ ] `TestAgent_ToolCall_HappyPath`: Verify request -> tool -> re-request loop.
-   [ ] `TestAgent_ToolCall_MaxIterations`: Ensure infinite loops (Agent calling tool forever) are stopped (limit: 5).
-   [ ] `TestVFS_Integration`: Agent calls `write_file` -> valid VFS version created in Session.

## 5. Implementation Steps
1.  Define Interfaces (`Tool`, `SessionAwareTool`).
2.  Implement `WriteFileTool` (VFS).
3.  Refactor `AgentProcessor` to support Tool Protocol.
4.  Implement `TestAgent_ToolCall` (Mock).
