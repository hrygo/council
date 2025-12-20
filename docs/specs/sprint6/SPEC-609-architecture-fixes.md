# SPEC-609: Architecture Defect Remediation

> **优先级**: P0  
> **类型**: Bug Fix / Architecture  
> **预估工时**: 6h  
> **来源**: v5 架构审计报告

## 1. 概述

本 SPEC 记录 Sprint 6 架构审计发现的关键缺陷，确保在开发迭代中得到修复。

---

## 2. P0 Critical Defects

### 2.1 Defect-1: NodeTypeMemoryRetrieval 未定义

**位置**: `internal/core/workflow/types.go`

**现状**: types.go 中不存在 `memory_retrieval` 节点类型

**修复**:

```go
// internal/core/workflow/types.go (Line 39 后添加)
NodeTypeMemoryRetrieval NodeType = "memory_retrieval" // NEW
```

```go
// internal/core/workflow/nodes/factory.go (添加 case)
case workflow.NodeTypeMemoryRetrieval:
    memSvc := deps.MemoryService // 需要添加到 NodeDependencies
    return &MemoryRetrievalProcessor{
        MemoryService: memSvc,
        MaxResults:    int(node.Properties["max_results"].(float64)),
    }, nil
```

**验收**: `go build` 通过，Optimize 流程可启动

---

### 2.2 Defect-2: ID 类型不匹配 (UUID vs String)

**位置**: 
- `internal/infrastructure/db/migrations/001_init_schema.up.sql`
- `internal/resources/seeder.go`

**现状**: Schema 定义 `id UUID`，但 Seeder 插入字符串 ID

**修复方案 A** (推荐): 修改 Schema 为 VARCHAR

```sql
-- 002_fix_id_type.up.sql
ALTER TABLE agents ALTER COLUMN id TYPE VARCHAR(64);
ALTER TABLE groups ALTER COLUMN id TYPE VARCHAR(64);
ALTER TABLE workflow_templates ALTER COLUMN id TYPE VARCHAR(64);
```

**修复方案 B**: 在 Seeder 中生成 UUID

```go
agentID := uuid.NewSHA1(uuid.NameSpaceOID, []byte("system_affirmative")).String()
```

**验收**: `make seed` 执行成功

---

### 2.3 Defect-3: Loop 节点参数名不一致

**位置**:
- `internal/core/workflow/nodes/factory.go`
- `internal/resources/seeder.go` (optimizeWorkflowGraph)

**现状**: 
- Workflow 使用 `exit_on_score`
- Factory 读取 `exit_condition`

**修复**: 统一使用 `exit_on_score`

```go
// factory.go Line 57-62
case workflow.NodeTypeLoop:
    maxRounds, _ := node.Properties["max_rounds"].(float64)
    exitOnScore, _ := node.Properties["exit_on_score"].(float64) // Fixed
    return &LoopProcessor{
        MaxRounds:   int(maxRounds),
        ExitOnScore: int(exitOnScore), // Fixed
    }, nil
```

**验收**: Loop 可根据 score >= 90 自动退出

---

### 2.4 Defect-4: Parallel 节点缺少 Join 逻辑

**位置**: `internal/core/workflow/engine.go`

**现状**: Parallel 分支各自执行完成后，没有等待所有分支完成再继续

**问题分析**:
```
parallel_debate
   ├─► agent_affirmative ─► agent_adjudicator
   └─► agent_negative ───► agent_adjudicator
```
两个分支独立完成，各自触发 `agent_adjudicator`，导致执行两次。

**修复方案**: 添加隐式 Join 节点或使用 sync.WaitGroup

```go
// engine.go
if node.Type == NodeTypeParallel {
    var wg sync.WaitGroup
    results := make(map[string]interface{})
    
    for _, branchID := range node.NextIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            // 执行分支
            result := e.executeNode(ctx, id, input)
            results[id] = result
        }(branchID)
    }
    
    wg.Wait()
    // 合并 results 后继续
    return e.findJoinNode(node), mergedResults
}
```

**验收**: Adjudicator 只执行一次，且收到双方论点

---

### 2.5 Defect-5: workflow_templates 缺少 updated_at 列

**位置**: `internal/infrastructure/db/migrations/001_init_schema.up.sql`

**现状**: 表定义缺少 `updated_at` 列

**修复**:

```sql
-- 002_fix_workflow_templates.up.sql
ALTER TABLE workflow_templates ADD COLUMN updated_at TIMESTAMPTZ DEFAULT NOW();
```

或修改 Seeder SQL 移除 `updated_at`:

```go
INSERT INTO workflow_templates (id, name, description, graph_definition, created_at)
VALUES ($1, $2, $3, $4::jsonb, NOW())
```

**验收**: Seeder 执行成功

---

## 3. P1 Medium Defects

### 3.1 Defect-6: 数据流传递机制未定义

**问题**: Aff/Neg 的输出如何传递给 Adjudicator?

**修复方案**: 在 SPEC-603 中增加数据流规范

```json
"agent_adjudicator": {
    "properties": {
        "agent_id": "system_adjudicator",
        "input_sources": ["agent_affirmative", "agent_negative"]  // NEW
    }
}
```

Engine 在执行 agent 节点前，检查 `input_sources` 并合并之前节点的输出。

---

### 3.2 Defect-7: AgentConfig 缺少 TopP 字段

**位置**: `internal/resources/prompt_loader.go`

**修复**:

```go
type AgentConfig struct {
    Name         string          `yaml:"name"`
    Provider     string          `yaml:"provider"`
    Model        string          `yaml:"model"`
    Temperature  float64         `yaml:"temperature"`
    MaxTokens    int             `yaml:"max_tokens"`
    TopP         float64         `yaml:"top_p"`       // NEW
    Capabilities map[string]bool `yaml:"capabilities"`
}
```

---

## 4. 验收标准

- [ ] `go build ./...` 通过
- [ ] `make migrate` 成功
- [ ] `make seed` 成功（或服务启动后自动 seed）
- [ ] `council_optimize` 流程可完整执行
- [ ] Adjudicator 只执行一次且收到双方输入
- [ ] Loop 可根据 score 自动退出

## 5. 依赖

- 需先完成 SPEC-607 (Memory Retrieval Node) 的 types/factory 扩展
- Migration 顺序：Schema Fix → Seeder

## 6. 风险

| 风险                          | 缓解                        |
| ----------------------------- | --------------------------- |
| Schema 变更影响现有数据       | 使用 `ALTER TABLE` 而非重建 |
| Engine 并发修改可能引入新 bug | 增加单元测试覆盖            |
