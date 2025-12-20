# SPEC-603: Default Workflow Templates

> **优先级**: P0  
> **类型**: SQL Migration + Go Seeder  
> **预估工时**: 6h  
> **依赖**: SPEC-601 (Agents), SPEC-607 (Memory Node)

## 1. 概述

创建系统默认的两个工作流模板："Debate" 和 "Optimize"。

## 2. Workflow 定义

### 2.1 Debate Workflow (辩论流程)

**ID**: `council_debate`  
**简化流程，用于单次辩论**

```
Start → Parallel(Aff, Neg) → Adjudicator → End
```

### 2.2 Optimize Workflow (优化流程) - **已增强**

**ID**: `council_optimize`  
**完整流程，对应 skill.md 6 步**

```
Start
  │
  ▼
[Memory Retrieval] ──► 读取历史上下文
  │
  ▼
[Parallel]
  ├─► Affirmative
  └─► Negative
  │
  ▼
[Adjudicator] ──► 输出结构化评分 + 行动建议
  │
  ▼
[Human Review] ──► 用户决定: 继续/应用/退出
  │
  ├─► (继续) → Loop 回到 Memory Retrieval
  └─► (退出) → End
```

## 3. 技术实现

### 3.1 GraphDefinition (council_optimize)

```json
{
  "id": "council_optimize",
  "name": "Council Optimize",
  "description": "迭代优化循环，含历史上下文检索",
  "start_node_id": "start",
  "nodes": {
    "start": {
      "id": "start",
      "type": "start",
      "name": "Input Document",
      "next_ids": ["memory_retrieval"]
    },
    "memory_retrieval": {
      "id": "memory_retrieval",
      "type": "memory_retrieval",
      "name": "Load History",
      "properties": {
        "max_results": 5,
        "time_range_days": 7,
        "include_verdicts": true
      },
      "next_ids": ["parallel_debate"]
    },
    "parallel_debate": {
      "id": "parallel_debate",
      "type": "parallel",
      "name": "Debate Round",
      "next_ids": ["agent_affirmative", "agent_negative"]
    },
    "agent_affirmative": {
      "id": "agent_affirmative",
      "type": "agent",
      "name": "Affirmative",
      "properties": {"agent_id": "system_affirmative"},
      "next_ids": ["agent_adjudicator"]
    },
    "agent_negative": {
      "id": "agent_negative",
      "type": "agent",
      "name": "Negative",
      "properties": {"agent_id": "system_negative"},
      "next_ids": ["agent_adjudicator"]
    },
    "agent_adjudicator": {
      "id": "agent_adjudicator",
      "type": "agent",
      "name": "Adjudicator",
      "properties": {
        "agent_id": "system_adjudicator",
        "output_format": "structured_verdict"
      },
      "next_ids": ["human_review"]
    },
    "human_review": {
      "id": "human_review",
      "type": "human_review",
      "name": "Review & Apply",
      "properties": {
        "show_score": true,
        "actions": ["continue", "apply", "exit", "rollback"]
      },
      "next_ids": ["loop_decision"]
    },
    "loop_decision": {
      "id": "loop_decision",
      "type": "loop",
      "name": "Continue?",
      "properties": {
        "max_rounds": 5,
        "exit_on_score": 90
      },
      "next_ids": ["memory_retrieval", "end"]
    },
    "end": {
      "id": "end",
      "type": "end",
      "name": "Final Report"
    }
  }
}
```

### 3.2 实现方式

使用 Go Seeder 注入 Workflow 模板（与 Agent 一致）。

```go
// internal/resources/seeder.go
func (s *Seeder) SeedWorkflows(ctx context.Context) error {
    workflows := []struct {
        ID          string
        Name        string
        Description string
        Graph       string // JSON
    }{
        {
            ID:          "council_debate",
            Name:        "Council Debate",
            Description: "三方辩论，生成综合裁决报告",
            Graph:       debateGraphJSON,
        },
        {
            ID:          "council_optimize",
            Name:        "Council Optimize",
            Description: "迭代优化循环，含历史上下文检索",
            Graph:       optimizeGraphJSON,
        },
    }
    
    for _, wf := range workflows {
        _, err := s.db.ExecContext(ctx, `
            INSERT INTO workflow_templates (id, name, description, graph_definition, created_at, updated_at)
            VALUES ($1, $2, $3, $4::jsonb, NOW(), NOW())
            ON CONFLICT (id) DO NOTHING
        `, wf.ID, wf.Name, wf.Description, wf.Graph)
        if err != nil {
            return err
        }
    }
    return nil
}
```

## 4. skill.md 步骤映射

| skill.md 步骤              | Workflow 节点              | 实现方式        |
| -------------------------- | -------------------------- | --------------- |
| Step 1: Compress History   | `memory_retrieval`         | SPEC-607        |
| Step 2: Convene Council    | `parallel_debate` + Agents | 现有            |
| Step 3: Verify Consistency | `agent_adjudicator`        | Prompt 内含评分 |
| Step 4: Snapshot Backup    | `human_review` 前触发      | SPEC-605        |
| Step 5: The Surgeon        | `human_review`             | 用户在 UI 应用  |
| Step 6: Loop Decision      | `loop_decision`            | 现有 Loop 节点  |

## 5. 验收标准

- [ ] 数据库存在 `council_debate` 和 `council_optimize` 两个模板
- [ ] `council_optimize` 包含 `memory_retrieval` 节点
- [ ] Workflow Canvas 可正确渲染两个模板
- [ ] 执行 `council_optimize` 时，Memory 节点被调用
- [ ] HumanReview 节点显示评分和行动按钮

## 6. 依赖

- **SPEC-601**: Default Agents
- **SPEC-607**: Memory Retrieval Node
