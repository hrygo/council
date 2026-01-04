# ID命名规范

<cite>
**本文档引用文件**  
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md)
- [session.go](file://internal/core/workflow/session.go)
- [types.go](file://internal/core/workflow/types.go)
- [entity.go](file://internal/core/agent/entity.go)
- [entity.go](file://internal/core/group/entity.go)
- [workflow.go](file://internal/api/handler/workflow.go)
- [hub.go](file://internal/api/ws/hub.go)
- [agent.go](file://internal/core/workflow/nodes/agent.go)
- [001_v2_schema_init.up.sql](file://internal/infrastructure/db/migrations/001_v2_schema_init.up.sql)
- [20251226_architecture_standardization_report.md](file://docs/reports/20251226_architecture_standardization_report.md)
</cite>

## 目录
1. [概述](#概述)
2. [命名约定](#命名约定)
3. [ID类型区分](#id类型区分)
4. [命名模式](#命名模式)
5. [跨层传递规则](#跨层传递规则)
6. [常见错误](#常见错误)
7. [参考实现](#参考实现)
8. [数据库Schema规范](#数据库schema规范)
9. [更新日志](#更新日志)

## 概述

本文档定义了Council项目中ID字段的命名规范和语义约定。该规范旨在解决因命名歧义导致的技术债务，强化前后端数据契约的一致性，并消除类型相关的编译错误。

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L1-L6)

## 命名约定

### 1. 实体 ID

| 字段名 | 类型 | 语义 | 示例 |
|--------|------|------|------|
| `ID` | uuid.UUID | 实体主键（数据库） | `550e8400-e29b-41d4-a716-446655440000` |
| `id` (JSON) | string | 实体 ID（API） | `550e8400-e29b-41d4-a716-446655440000` |

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L9-L15)

### 2. 关联 ID

| 字段名 | 类型 | 语义 | 用途 |
|--------|------|------|------|
| `SessionID` | string | 会话 ID | 标识工作流执行实例 |
| `WorkflowID` | string | 工作流定义 ID | 引用工作流模板 |
| `GroupID` | string/uuid.UUID | 群组 ID | 标识 Agent 所属群组 |
| `AgentID` | string/uuid.UUID | Agent ID | 引用数据库中的 Agent |
| `NodeID` | string | 工作流节点 ID | 工作流图中的节点标识 |

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L16-L25)

### 3. 特殊 ID

| 字段名 | 类型 | 语义 | 说明 |
|--------|------|------|------|
| `node_id` (JSON) | string | WebSocket 消息节点 ID | 事件关联的节点 |
| `StartNodeID` | string | 起始节点 ID | 工作流入口节点 |

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L26-L32)

## ID类型区分

### UUID vs String

**UUID 类型** (`uuid.UUID`):
- 用于数据库主键
- 全局唯一
- Go 结构体字段

**String 类型** (`string`):
- 用于工作流图节点标识
- API JSON 字段
- WebSocket 消息字段

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L33-L46)

## 命名模式

### 节点 ID 命名

工作流图中的节点 ID 建议使用语义化命名：

```
{类型}_{名称}[_{序号}]
```

**示例**:
- `agent_affirmative` - 正方辩手节点
- `agent_negative` - 反方辩手节点
- `vote_1` - 第一次投票节点
- `loop_optimization` - 优化循环节点
- `human_review_final` - 最终人工审核节点

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L69-L84)

### 会话 ID 命名

会话 ID 使用 UUID v4：

```go
sessionID := uuid.New().String()
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L85-L92)

## 跨层传递规则

### 1. API → 工作流引擎

```go
// API Handler
type ExecuteRequest struct {
    WorkflowID string `json:"workflow_id"`  // 工作流模板 ID
    Inputs     map[string]interface{} `json:"inputs"`
}

// 传递到引擎
session := &workflow.Session{
    ID:         uuid.New().String(),  // 生成新会话 ID
    WorkflowID: req.WorkflowID,       // 传递工作流 ID
    Inputs:     req.Inputs,
}
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L94-L111)

### 2. 引擎 → 节点处理器

```go
// 引擎调用节点
func (e *Engine) executeNode(nodeID string) {
    node := e.Graph.Nodes[nodeID]  // nodeID 是图中的节点 ID
    
    // 如果是 Agent 节点，需要解析 AgentID
    if node.Type == NodeTypeAgent {
        agentID := node.Properties["agent_id"].(string)  // 从配置读取 Agent UUID
        // 使用 agentID 查询数据库
    }
}
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L112-L125)

### 3. 节点 → WebSocket

```go
// 节点发送事件
stream <- workflow.StreamEvent{
    NodeID: nodeID,  // 使用工作流图中的节点 ID
    Data: map[string]interface{}{
        "session_id": ctx.SessionID,  // 会话 ID
    },
}
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L127-L137)

## 常见错误

### ❌ 错误示例 1: 混淆 NodeID 和 AgentID

```go
// 错误：将 NodeID 当作 AgentID 使用
agent, err := repo.GetAgent(nodeID)  // nodeID 不是 UUID
```

**正确做法**:

```go
// 从节点配置读取 AgentID
agentID := node.Properties["agent_id"].(string)
agent, err := repo.GetAgent(uuid.MustParse(agentID))
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L140-L155)

### ❌ 错误示例 2: ID 类型不一致

```go
// 错误：字符串 ID 直接赋值给 UUID 字段
agent := &Agent{
    ID: "some-string-id",  // 类型错误
}
```

**正确做法**:

```go
agent := &Agent{
    ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
}
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L156-L171)

### ❌ 错误示例 3: 节点 ID 命名不明确

```go
// 不好的命名
nodeID := "node1"
nodeID := "n_1"

// 好的命名
nodeID := "agent_affirmative"
nodeID := "vote_round_1"
```

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L173-L183)

## 参考实现

### 会话实体实现

```go
// Session represents a single execution instance of a workflow
type Session struct {
    ID        string `json:"session_uuid"`
    Graph     *GraphDefinition
    Status    SessionStatus
    StartTime time.Time
    EndTime   time.Time
    Inputs    map[string]interface{}
    Outputs   map[string]interface{}
    Error     error
}
```

**Section sources**
- [session.go](file://internal/core/workflow/session.go#L23-L42)

### 工作流图定义

```go
// GraphDefinition represents the static definition of a workflow
type GraphDefinition struct {
    ID          string           `json:"workflow_uuid"`
    Name        string           `json:"name"`
    Description string           `json:"description"`
    Nodes       map[string]*Node `json:"nodes"`
    StartNodeID string           `json:"start_node_id"`
}
```

**Section sources**
- [types.go](file://internal/core/workflow/types.go#L44-L51)

### Agent实体实现

```go
// Agent represents an AI persona.
type Agent struct {
    ID            uuid.UUID    `json:"agent_uuid" db:"agent_uuid"`
    Name          string       `json:"name" db:"name"`
    Avatar        *string      `json:"avatar" db:"avatar"`
    Description   *string      `json:"description" db:"description"`
    PersonaPrompt string       `json:"persona_prompt" db:"persona_prompt"`
    ModelConfig   ModelConfig  `json:"model_config" db:"model_config"`
    Capabilities  Capabilities `json:"capabilities" db:"capabilities"`
    CreatedAt     time.Time    `json:"created_at" db:"created_at"`
    UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}
```

**Section sources**
- [entity.go](file://internal/core/agent/entity.go#L9-L20)

### 群组实体实现

```go
// Group represents a collaboration group (Project/Context).
type Group struct {
    ID                uuid.UUID   `json:"group_uuid" db:"group_uuid"`
    Name              string      `json:"name" db:"name"`
    Icon              *string     `json:"icon" db:"icon"`
    SystemPrompt      *string     `json:"system_prompt" db:"system_prompt"`
    DefaultAgentUUIDs []uuid.UUID `json:"default_agent_uuids" db:"default_agent_uuids"`
    CreatedAt         time.Time   `json:"created_at" db:"created_at"`
    UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}
```

**Section sources**
- [entity.go](file://internal/core/group/entity.go#L9-L19)

## 数据库schema规范

### 标识符标准化规范

| ID 类型    | 后缀规范 | 语义                             | 示例                                      |
| :--------- | :------- | :------------------------------- | :---------------------------------------- |
| **主键**   | `_uuid`  | 数据库全局唯一主键 (Primary Key) | `user_uuid`, `session_uuid`, `agent_uuid` |
| **逻辑ID** | `_id`    | 上下文特定或逻辑标识符           | `node_id`, `provider_id`, `message_id`    |

**Section sources**
- [20251226_architecture_standardization_report.md](file://docs/reports/20251226_architecture_standardization_report.md#L14-L19)

### 数据库迁移示例

```sql
-- 1. Groups Table
CREATE TABLE groups (
    group_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    icon VARCHAR(256),
    system_prompt TEXT,
    default_agent_uuids JSONB DEFAULT '[]',
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

-- 6. Sessions Table
CREATE TABLE sessions (
    session_uuid UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_uuid UUID REFERENCES groups(group_uuid) ON DELETE CASCADE,
    workflow_uuid UUID REFERENCES workflows(workflow_uuid),
    status VARCHAR(32) DEFAULT 'pending',
    proposal JSONB,
    started_at TIMESTAMPTZ,
    ended_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
```

**Section sources**
- [001_v2_schema_init.up.sql](file://internal/infrastructure/db/migrations/001_v2_schema_init.up.sql#L8-L72)

## 更新日志

- **2025-12-26**: 初始版本

**Section sources**
- [id-naming-convention.md](file://docs/guide/id-naming-convention.md#L239-L242)