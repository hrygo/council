# ID 命名规范文档

## 概述

本文档定义了 Council 项目中 ID 字段的命名规范和语义约定。

## 命名约定

### 1. 实体 ID

| 字段名 | 类型 | 语义 | 示例 |
|--------|------|------|------|
| `ID` | uuid.UUID | 实体主键（数据库） | `550e8400-e29b-41d4-a716-446655440000` |
| `id` (JSON) | string | 实体 ID（API） | `550e8400-e29b-41d4-a716-446655440000` |

### 2. 关联 ID

| 字段名 | 类型 | 语义 | 用途 |
|--------|------|------|------|
| `SessionID` | string | 会话 ID | 标识工作流执行实例 |
| `WorkflowID` | string | 工作流定义 ID | 引用工作流模板 |
| `GroupID` | string/uuid.UUID | 群组 ID | 标识 Agent 所属群组 |
| `AgentID` | string/uuid.UUID | Agent ID | 引用数据库中的 Agent |
| `NodeID` | string | 工作流节点 ID | 工作流图中的节点标识 |

### 3. 特殊 ID

| 字段名 | 类型 | 语义 | 说明 |
|--------|------|------|------|
| `node_id` (JSON) | string | WebSocket 消息节点 ID | 事件关联的节点 |
| `StartNodeID` | string | 起始节点 ID | 工作流入口节点 |

## ID 类型区分

### UUID vs String

**UUID 类型** (`uuid.UUID`):
- 用于数据库主键
- 全局唯一
- Go 结构体字段

**String 类型** (`string`):
- 用于工作流图节点标识
- API JSON 字段
- WebSocket 消息字段

### 示例对比

```go
// Agent 实体（数据库）
type Agent struct {
    ID   uuid.UUID  // 数据库主键，UUID 类型
    Name string
}

// Agent 节点（工作流图）
type AgentNode struct {
    NodeID  string     // 工作流节点 ID，字符串类型（如 "agent_affirmative"）
    AgentID string     // 引用 Agent.ID，字符串类型（UUID 转字符串）
}

// WebSocket 消息
type StreamEvent struct {
    NodeID string  // 节点 ID，对应工作流图中的 NodeID
}
```

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

### 会话 ID 命名

会话 ID 使用 UUID v4：

```go
sessionID := uuid.New().String()
```

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

### ❌ 错误示例 3: 节点 ID 命名不明确

```go
// 不好的命名
nodeID := "node1"
nodeID := "n_1"

// 好的命名
nodeID := "agent_affirmative"
nodeID := "vote_round_1"
```

## 检查清单

开发新功能时，请检查：

- [ ] ID 字段类型正确（UUID vs String）
- [ ] ID 语义明确（是节点 ID 还是实体 ID？）
- [ ] 跨层传递时类型转换正确
- [ ] 节点 ID 使用语义化命名
- [ ] WebSocket 消息包含正确的 ID 字段
- [ ] 错误消息中包含 ID 信息便于调试

## 参考代码

### ID 生成

```go
// 生成 UUID
import "github.com/google/uuid"

// 新建实体
entityID := uuid.New()

// 字符串转 UUID
entityID, err := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")

// UUID 转字符串
idStr := entityID.String()
```

### ID 验证

```go
// 验证 UUID 格式
func isValidUUID(id string) bool {
    _, err := uuid.Parse(id)
    return err == nil
}
```

### 类型转换

```go
// interface{} → string ID
if id, ok := data["id"].(string); ok {
    // 使用 id
}

// string ID → UUID
entityID, err := uuid.Parse(stringID)
if err != nil {
    return fmt.Errorf("invalid UUID: %w", err)
}
```

## 更新日志

- **2025-12-26**: 初始版本
