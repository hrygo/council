# 🔴 终极无情架构审计报告 v5

**审计员**: 首席架构破坏者  
**日期**: 2024-12-20  
**态度**: 恶意挑刺 + 代码级验证  
**结论**: **❌ 方案存在多个致命缺陷，无法直接实施**

---

## 💀 CRITICAL DEFECTS (P0 - 必须修复)

### Defect 1: `NodeTypeMemoryRetrieval` 不存在

**位置**: `internal/core/workflow/types.go` (Line 28-40)

**现有 NodeType**:
```go
const (
    NodeTypeStart       NodeType = "start"
    NodeTypeEnd         NodeType = "end"
    NodeTypeAgent       NodeType = "agent"
    NodeTypeLLM         NodeType = "llm"
    NodeTypeTool        NodeType = "tool"
    NodeTypeParallel    NodeType = "parallel"
    NodeTypeSequence    NodeType = "sequence"
    NodeTypeVote        NodeType = "vote"
    NodeTypeLoop        NodeType = "loop"
    NodeTypeFactCheck   NodeType = "fact_check"
    NodeTypeHumanReview NodeType = "human_review"
)
```

**缺失**: `NodeTypeMemoryRetrieval NodeType = "memory_retrieval"`

**SPEC-603 Workflow**:
```json
{
  "type": "memory_retrieval",  // ← 这个类型不存在！
  ...
}
```

**后果**: 
1. `factory.go` 的 switch 会落入 `default` 分支
2. 返回 `fmt.Errorf("unsupported node type: memory_retrieval")`
3. **整个 Optimize 流程无法启动**

---

### Defect 2: 数据库 Schema 使用 UUID，Seeder 插入 String

**位置**: `internal/infrastructure/db/migrations/001_init_schema.up.sql` (Line 16-17)

```sql
CREATE TABLE agents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),  -- UUID 类型
    ...
);
```

**Seeder 代码** (`internal/resources/seeder.go` Line 52):
```go
_, err = s.db.ExecContext(ctx, `
    INSERT INTO agents (id, name, ...)
    VALUES ($1, $2, ...)
`, agentID, ...)  // agentID = "system_affirmative" (STRING!)
```

**后果**: PostgreSQL 会报错 `invalid input syntax for type uuid: "system_affirmative"`

---

### Defect 3: Loop 节点参数名不匹配

**SPEC-603 Workflow** (Line 121-124):
```json
"properties": {
    "max_rounds": 5,
    "exit_on_score": 90  // ← 参数名
}
```

**factory.go** (Line 57-58):
```go
maxRounds, _ := node.Properties["max_rounds"].(float64)
exitCond, _ := node.Properties["exit_condition"].(string) // ← 读的是 exit_condition!
```

**后果**: `exit_on_score` 永远不会被读取，Loop 无法根据分数退出。

---

### Defect 4: Parallel 节点 DAG 逻辑错误

**SPEC-603 Workflow** (Line 77-95):
```json
"parallel_debate": {
    "type": "parallel",
    "next_ids": ["agent_affirmative", "agent_negative"]  // 分支
},
"agent_affirmative": {
    "next_ids": ["agent_adjudicator"]  // 汇聚
},
"agent_negative": {
    "next_ids": ["agent_adjudicator"]  // 汇聚
}
```

这意味着 `agent_affirmative` 和 `agent_negative` 都独立指向 `agent_adjudicator`。

**Engine 如何知道要等待两个分支都完成？**

查看 `engine.go` (Line 80):
```go
if node.Type == NodeTypeParallel {
    // 并发执行分支...
}
```

**问题**: 没有看到 "等待所有分支完成后再继续" 的 Join 逻辑。  
**后果**: `agent_adjudicator` 可能在只有一个分支完成时就被触发两次。

---

### Defect 5: workflow_templates 缺少 `updated_at` 字段

**Schema** (Line 47-54):
```sql
CREATE TABLE workflow_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(128) NOT NULL,
    description VARCHAR(512),
    graph_definition JSONB NOT NULL,
    is_system BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
    -- 没有 updated_at！
);
```

**Seeder** (Line 240-241):
```go
INSERT INTO workflow_templates (id, name, description, graph_definition, created_at, updated_at)
                                                                                    ^^^^^^^^^^
```

**后果**: SQL 执行失败 `column "updated_at" does not exist`

---

## 🟡 MEDIUM DEFECTS (P1)

### Defect 6: Adjudicator 输入缺少双方论点

`agent_affirmative` 和 `agent_negative` 的输出如何传递给 `agent_adjudicator`？

Workflow 定义只有 `next_ids`，没有定义数据流。`agent_adjudicator` 如何知道要读取哪些输入？

---

### Defect 7: top_p 在 Prompt YAML 中但不在 AgentConfig

**Prompt Front Matter**:
```yaml
top_p: 0.95
```

**SPEC-608 AgentConfig** (已更新版):
```go
type AgentConfig struct {
    Name         string          `yaml:"name"`
    Provider     string          `yaml:"provider"`
    Model        string          `yaml:"model"`
    Temperature  float64         `yaml:"temperature"`
    MaxTokens    int             `yaml:"max_tokens"`
    Capabilities map[string]bool `yaml:"capabilities"`
    // 没有 TopP！
}
```

**后果**: `top_p` 不会被解析，模型调用时丢失参数。

---

## 修复清单

| 优先级 | 缺陷     | 修复方案                                                                |
| ------ | -------- | ----------------------------------------------------------------------- |
| **P0** | Defect 1 | 在 `types.go` 添加 `NodeTypeMemoryRetrieval`，在 `factory.go` 添加 case |
| **P0** | Defect 2 | 使用 UUID 格式 ID，或改 Schema 为 VARCHAR                               |
| **P0** | Defect 3 | 统一参数名为 `exit_on_score` 或 `exit_condition`                        |
| **P0** | Defect 4 | 添加 Join 节点或修复 Engine 并发聚合逻辑                                |
| **P0** | Defect 5 | 给 `workflow_templates` 添加 `updated_at` 列                            |
| **P1** | Defect 6 | 在 SPEC 中明确数据流传递机制                                            |
| **P1** | Defect 7 | 在 `AgentConfig` 添加 `TopP` 字段                                       |

---

## 结论

**❌ 方案不可实施**

之前的审计只看了文档层面，没有深入到代码与 Schema 的交叉验证。  
以上 5 个 P0 缺陷中任何一个都会导致系统在启动或运行时崩溃。

**建议**: 立即修复 P0 缺陷后再进行下一轮实施。
