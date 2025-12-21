# WebSocket 调试报告

> **日期**: 2025-12-21  
> **版本**: v0.14.0 → v0.15.0  
> **报告类型**: Bug 复盘与根因分析  
> **调试时长**: ~30 分钟

---

## 一、问题概述

会议室功能异常：用户启动会话后，前端无法接收和显示 WebSocket 消息，UI 状态停留在 IDLE。

---

## 二、Bug 清单与深度根因分析

### Bug 1: WebSocket 消息静默丢失

| 维度         | 详情                                     |
| ------------ | ---------------------------------------- |
| **现象**     | 后端正常发送消息，前端 Router 无任何处理 |
| **修复位置** | `internal/core/workflow/context.go:10`   |
| **修复内容** | JSON tag `"type"` → `"event"`            |

**问题代码**:
```go
// 后端 (Go)
type StreamEvent struct {
    Type string `json:"type"` // ❌ 序列化为 {"type": "..."}
}
```

```typescript
// 前端 (TypeScript)
interface WSMessage {
    event: WSEventType; // ❌ 期望 {"event": "..."}
}
```

**深度根因**:

1. **缺乏 API Contract 文档**: 前后端分别参考了不同的命名惯例：
   - Go/后端: 参考 Server-Sent Events 规范使用 `type`
   - TypeScript/前端: 参考 Socket.IO 风格使用 `event`

2. **无自动化类型同步**: 未使用 `tygo` 等工具从 Go struct 生成 TypeScript interface

3. **测试覆盖盲区**: 单元测试分别测试后端发送和前端接收，但**未进行端到端消息格式验证**

**影响范围**: 所有 WebSocket 消息（token_stream, node_state_change, error 等）

---

### Bug 2: node_id 使用 Agent UUID 而非 Graph Node ID

| 维度         | 详情                                              |
| ------------ | ------------------------------------------------- |
| **现象**     | 消息显示为 `94cc1d02-4d8e-53f3-8891-...`          |
| **修复位置** | `internal/core/workflow/nodes/agent.go:26,88,113` |
| **修复内容** | event data 中使用 `a.NodeID` 替代 `a.AgentID`     |

**问题代码**:
```go
// 修复前
stream <- StreamEvent{
    Data: map[string]interface{}{"node_id": a.AgentID}, // Agent 数据库 UUID
}
```

**深度根因**:

1. **ID 语义混淆**: 
   - `AgentID`: 数据库中 Agent 实体的主键（UUID）
   - `NodeID`: 工作流图中节点的逻辑标识符（如 "agent_affirmative"）
   
   开发时**未明确区分这两个概念**，直接使用 `AgentID` 作为事件中的标识。

2. **命名不规范**: 
   - `AgentProcessor.AgentID` 命名暗示了它是"Agent 的 ID"，但实际调用场景需要的是"当前节点的 ID"
   
3. **缺少设计评审**: 数据流设计：
   ```
   Template → NodeFactory → AgentProcessor → StreamEvent → Frontend
   ```
   在这条链路中，NodeFactory 只传递了 `agent_id`（属性），未传递 `node.ID`（图结构 ID）

**影响范围**: 所有 Agent 节点的消息显示

---

### Bug 3: NodeStateSnapshot 缺少 name/type 字段

| 维度         | 详情                                                                |
| ------------ | ------------------------------------------------------------------- |
| **现象**     | 即使 node_id 正确，仍显示 "Unknown Node"                            |
| **修复位置** | `frontend/src/types/session.ts:59`, `stores/useSessionStore.ts:107` |
| **修复内容** | 扩展接口，initSession 时存储完整元数据                              |

**问题代码**:
```typescript
// 修复前
interface NodeStateSnapshot {
    id: string;
    status: NodeStatus;
    // ❌ 缺少 name 和 type
}
```

**深度根因**:

1. **类型设计未考虑 UI 需求**: 
   最初设计 `NodeStateSnapshot` 时，只关注"运行时状态跟踪"，忽略了 UI 渲染需要显示元数据

2. **数据丢失**: 
   `initSession` 接收完整的 `{id, name, type}`，但存储时只保留 `{id, status}`

3. **信息传递断层**:
   ```
   Template.graph.nodes → initSession(nodes) → currentSession.nodes
   ```
   传入有 name，存储时丢弃

**影响范围**: 所有消息组的节点名称显示

---

### Bug 4: LLM Model 降级硬编码

| 维度         | 详情                                            |
| ------------ | ----------------------------------------------- |
| **现象**     | 使用 Gemini provider 时尝试调用不存在的 `gpt-4` |
| **修复位置** | `internal/core/workflow/nodes/agent.go:70`      |
| **修复内容** | `"gpt-4"` → `a.Registry.GetDefaultModel()`      |

**问题代码**:
```go
if req.Model == "" {
    req.Model = "gpt-4" // ❌ 硬编码
}
```

**深度根因**:

1. **开发阶段假设**: 早期以 OpenAI 为主要 provider，`gpt-4` 是"安全默认值"

2. **未遵循 DRY 原则**: `Registry` 已有 `GetDefaultModel()` 方法，但未被使用

3. **缺乏配置驱动**: 默认值应来自环境变量/配置文件，而非代码

**影响范围**: 所有未配置 Model 的 Agent 节点

---

### Bug 5: Session 状态永远 IDLE

| 维度         | 详情                                             |
| ------------ | ------------------------------------------------ |
| **现象**     | 工作流运行中，UI 仍显示 "IDLE" 状态              |
| **修复位置** | `frontend/src/hooks/useWebSocketRouter.ts:39`    |
| **修复内容** | 收到第一个 `running` 节点时自动更新 session 状态 |

**问题代码**:
```typescript
case 'node_state_change':
    // ❌ 只更新节点状态，未触发 session 状态变更
    workflowStore.updateNodeStatus(data.node_id, data.status);
```

**深度根因**:

1. **缺少显式状态事件**: 后端未发送 `execution:started` 事件

2. **隐式 vs 显式**: 
   - 前端期望显式的 `session.status = running` 事件
   - 后端假设前端会从节点状态推断 session 状态

3. **状态机边界不清**: `initSession` 设置 `status: 'idle'`，但未定义何时转为 `running`

**影响范围**: 所有会话的状态指示器

---

## 三、根因分类与统计

| 根因类别             | Bug 数量 | 具体案例                           |
| -------------------- | -------- | ---------------------------------- |
| **前后端合约不一致** | 2        | Bug 1 (JSON字段), Bug 5 (状态事件) |
| **ID/语义混淆**      | 1        | Bug 2                              |
| **数据模型不完整**   | 1        | Bug 3                              |
| **硬编码/缺少抽象**  | 1        | Bug 4                              |

---

## 四、改进建议与规约更新

### 4.1 立即执行项

| 项目                                | 负责   | 优先级 |
| ----------------------------------- | ------ | ------ |
| 使用 `tygo` 生成 WebSocket 消息类型 | 后端   | P0     |
| 建立 `docs/api-contract.md`         | 架构   | P0     |
| CI 添加类型一致性检查               | DevOps | P1     |

### 4.2 规约更新建议

#### GEMINI.md 新增条款

```markdown
## 前后端协作规约

1. **JSON 字段命名**: 前后端必须使用相同字段名，优先采用前端命名惯例
2. **类型同步**: 所有共享类型必须通过 tygo 自动生成，禁止手动维护
3. **ID 语义**: 
   - 数据库 ID 使用 `{entity}_uuid` 命名 (如 `agent_uuid`)
   - 逻辑 ID 使用 `{context}_id` 命名 (如 `node_id`)
4. **默认值获取**: 禁止硬编码，必须通过配置或工厂方法
```

### 4.3 研发计划更新

在 Sprint 8 后增加技术债务清理任务：
- [ ] WebSocket 消息类型自动生成
- [ ] 端到端消息格式测试
- [ ] ID 命名规范审查

---

## 五、预防清单

| 检查项                 | 适用时机         | 检查方法         |
| ---------------------- | ---------------- | ---------------- |
| 前后端 JSON 字段名一致 | 新增 API/WS 消息 | tygo 生成 + diff |
| ID 参数语义明确        | 跨层传递 ID      | 代码审查         |
| 类型包含 UI 所需字段   | 新增 Store 类型  | 设计评审         |
| 默认值来自配置         | 新增配置项       | Lint 规则        |
| 状态转换显式定义       | 新增状态机       | 文档审查         |

---

## 六、附录

### 调试时间线

```
11:44 - 用户报告会议室异常
11:45 - 检查 WebSocket 连接状态 → 连接正常
11:46 - 对比前后端消息定义 → 发现 type/event 不匹配
11:47 - 修复 StreamEvent JSON tag
11:49 - 确认消息开始显示，但节点名称错误
11:50 - 追踪 AgentProcessor.AgentID → 应为 NodeID
11:52 - 修复 AgentProcessor，添加 NodeID 字段
11:55 - 修复 NodeStateSnapshot 存储 name/type
12:00 - 修复 LLM Model 降级逻辑
12:00 - 用户记录完整待办清单
12:14 - 所有测试通过
12:15 - Release v0.15.0
```

### 修改文件清单

| 文件                                       | 修改类型                  |
| ------------------------------------------ | ------------------------- |
| `internal/core/workflow/context.go`        | JSON tag 修复             |
| `internal/core/workflow/nodes/agent.go`    | 添加 NodeID, 修复事件发送 |
| `internal/core/workflow/nodes/factory.go`  | 传递 NodeID               |
| `internal/infrastructure/llm/router.go`    | 添加 GetDefaultModel      |
| `frontend/src/types/session.ts`            | 扩展 NodeStateSnapshot    |
| `frontend/src/stores/useSessionStore.ts`   | 存储完整节点信息          |
| `frontend/src/hooks/useWebSocketRouter.ts` | 自动更新 session 状态     |
| `frontend/vite.config.ts`                  | 修复 WS proxy             |
