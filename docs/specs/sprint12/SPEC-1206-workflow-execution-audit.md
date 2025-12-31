# SPEC-1206: Workflow Execution Engine 严格审计与修复

> **优先级**: P0 (Critical)  
> **类型**: Core Engine Refactor  
> **预估工时**: 8h  
> **依赖**: SPEC-603 (Default Workflows)

## 1. 审计背景 (Audit Context)

用户观察到 Council Debate 工作流执行时，Agent 的发言"完全不知所云"。经过对 `internal/core/workflow/` 目录下核心代码的深入审计，发现了导致此问题的 **两个致命架构缺陷**。

---

## 2. 发现的致命缺陷 (Critical Defects Found)

### Defect-1: Fan-in 逻辑缺失 (Missing Join/Barrier)

**严重程度**: 🔴 Critical

**现象描述**:  
在 `council_debate` 工作流中，`Affirmative` (正方) 和 `Negative` (反方) 节点并行执行后，都指向 `Adjudicator` (裁决官)。

**期望行为**:  
裁决官应该 **等待** 正反双方都完成后，**聚合** 双方输出，然后 **执行一次**。

**实际行为** (Bug):  
当前 `engine.go` 的 `executeNode` 逻辑在每个节点完成后，直接遍历 `NextIDs` 并递归触发下游节点。

```go
// engine.go L144-L152 (Current Flawed Logic)
for _, nextID := range node.NextIDs {
    wg.Add(1)
    go func(nid string) {
        defer wg.Done()
        e.executeNode(ctx, nid, output) // ← 每个上游完成都触发一次
    }(nextID)
}
```

**后果**:  
1. 裁决官被触发 **两次** (正方完成触发一次，反方完成触发一次)。
2. 每次触发时，裁决官只收到 **单方** 的输出，无法进行真正的辩论裁决。
3. 两次执行的输出相互覆盖或混乱，导致最终结果"不知所云"。

**根因分析**:  
引擎缺少对节点 **入度 (In-degree)** 的跟踪，以及在多入度节点上的 **屏障/聚合 (Barrier/Join)** 机制。

---

### Defect-2: 上下文透传断裂 (Context Propagation Broken)

**严重程度**: 🔴 Critical

**现象描述**:  
用户提交的原始文档 (`document_content`) 在经过第一个 Agent 节点后丢失，后续节点无法获取。

**期望行为**:  
原始文档应贯穿整个工作流，每个节点都能访问。

**实际行为** (Bug):  
`AgentProcessor.Process()` 的输出仅包含该 Agent 产生的新内容，未透传上游输入。

```go
// nodes/agent.go L184-L188 (Current Flawed Logic)
output := map[string]interface{}{
    "agent_output": finalResponse,
    "agent_id":     a.AgentID,
    "timestamp":    time.Now(),
}
// ← 未包含 input 中的原始数据！
```

**后果**:  
1. `Start` 节点提供文档。
2. `Affirmative` 节点收到文档，产生分析，但输出只有分析，**文档丢失**。
3. `Adjudicator` 收到的输入中没有原始文档，只有正方的分析片段。
4. 裁决官在完全没有原始材料的情况下盲目输出，导致"不知所云"。

**根因分析**:  
缺少 **上下文合并 (Context Merge)** 策略。每个 Processor 应将核心上下文字段透传到输出。

---

## 3. 架构原则 (Architecture Principles)

> ⚠️ **核心约束**: 骨架 (Engine) 与应用 (Council Debate) 必须解耦

### 3.1 分层职责

| 层级                         | 职责                             | 示例                                 |
| :--------------------------- | :------------------------------- | :----------------------------------- |
| **Engine (骨架)**            | 通用的图执行、入度计算、数据路由 | `engine.go`                          |
| **Node Processor (节点)**    | 节点级逻辑、上下文消费方式       | `agent.go`, `vote.go`                |
| **Workflow Template (应用)** | 业务编排、字段命名约定           | `council_debate`, `council_optimize` |

### 3.2 设计准则

1. **Engine 无业务感知**: Engine 不应硬编码任何业务字段名 (如 `document_content`, `agent_output`)
2. **策略可配置**: 聚合策略 (如何 merge 多个上游输出) 应由节点属性或全局配置决定
3. **透传由节点自治**: 每个 Processor 自行决定透传哪些字段，Engine 只负责路由

---

## 4. 修复方案 (Remediation Plan)

### Fix-1: 通用 Join 机制 (Engine 层)

**目标**: 在 Engine 层实现通用的入度跟踪和数据聚合，不引入任何业务字段硬编码。

**设计要点**:
- Engine **仅负责**: 入度计算、数据收集、触发时机控制
- Engine **不负责**: 具体的聚合逻辑 (由可插拔的 `MergeStrategy` 接口处理)

**实现方案**:

1. **定义 MergeStrategy 接口** (骨架层):
   ```go
   // internal/core/workflow/merge.go (新增)
   
   // MergeStrategy 定义多入度节点如何聚合上游输出
   // 骨架层只定义接口，不实现具体策略
   type MergeStrategy interface {
       // Merge 接收多个上游输出，返回聚合后的输入
       Merge(inputs []map[string]interface{}) map[string]interface{}
   }
   
   // DefaultMergeStrategy 默认策略：简单合并，保留所有字段
   // 同名字段按索引区分 (branch_0, branch_1, ...)
   type DefaultMergeStrategy struct{}
   
   func (s *DefaultMergeStrategy) Merge(inputs []map[string]interface{}) map[string]interface{} {
       merged := make(map[string]interface{})
       
       for i, inp := range inputs {
           // 保留每个分支的完整输出
           merged[fmt.Sprintf("branch_%d", i)] = inp
           
           // 透传首次出现的非 branch 字段 (优先级: 先到先得)
           for k, v := range inp {
               if _, exists := merged[k]; !exists {
                   merged[k] = v
               }
           }
       }
       
       return merged
   }
   ```

2. **Engine 新增入度跟踪**:
   ```go
   // engine.go (修改)
   type Engine struct {
       // ... existing fields
       inDegree      map[string]int                    // 节点入度
       pendingInputs map[string][]map[string]interface{} // 待聚合数据
       joinMu        sync.Mutex
       MergeStrategy MergeStrategy  // 可注入的聚合策略
   }
   
   func NewEngine(session *Session) *Engine {
       e := &Engine{
           // ... existing init
           pendingInputs: make(map[string][]map[string]interface{}),
           MergeStrategy: &DefaultMergeStrategy{}, // 默认策略
       }
       e.computeInDegrees()
       return e
   }
   
   func (e *Engine) computeInDegrees() {
       e.inDegree = make(map[string]int)
       for _, node := range e.Graph.Nodes {
           for _, nextID := range node.NextIDs {
               e.inDegree[nextID]++
           }
       }
   }
   ```

3. **替换直接触发为"投递-聚合"模式**:
   ```go
   // engine.go (替换原 L144-L152)
   func (e *Engine) deliverToDownstream(ctx context.Context, nodeID string, output map[string]interface{}) {
       node := e.Graph.Nodes[nodeID]
       
       for _, nextID := range node.NextIDs {
           e.joinMu.Lock()
           e.pendingInputs[nextID] = append(e.pendingInputs[nextID], output)
           ready := len(e.pendingInputs[nextID]) >= e.inDegree[nextID]
           
           var mergedInput map[string]interface{}
           if ready {
               mergedInput = e.MergeStrategy.Merge(e.pendingInputs[nextID])
               e.pendingInputs[nextID] = nil // 清空
           }
           e.joinMu.Unlock()
           
           if ready {
               go e.executeNode(ctx, nextID, mergedInput)
           }
       }
   }
   ```

---

### Fix-2: 上下文透传 (Processor 层自治)

**目标**: 让每个 Processor 自主决定透传哪些字段，Engine 不干预。

**设计要点**:
- **无全局白名单**: 不在 Engine 或公共包中硬编码字段列表
- **Processor 自治**: AgentProcessor 可在内部定义需要透传的字段
- **可扩展**: 不同 Processor 可有不同的透传逻辑

**实现方案**:

1. **AgentProcessor 内部定义透传字段** (应用层):
   ```go
   // nodes/agent.go (修改)
   
   // agentPassthroughKeys 定义 Agent 节点需要透传的字段
   // 这是 AgentProcessor 的私有配置，不污染 Engine
   var agentPassthroughKeys = []string{
       "document_content",
       "proposal",
       "optimization_objective",
       "attachments",
       "combined_context",
       "session_id",
       "aggregated_outputs", // 用于接收 Join 后的聚合数据
   }
   
   // 在 Process() 的输出构建中
   output := map[string]interface{}{
       "agent_output": finalResponse,
       "agent_id":     a.AgentID,
       "timestamp":    time.Now(),
   }
   
   // 透传上下文 (Processor 自治)
   for _, key := range agentPassthroughKeys {
       if val, ok := input[key]; ok {
           output[key] = val
       }
   }
   ```

2. **修复 Prompt 构建逻辑**:
   ```go
   // nodes/agent.go (修改 constructHistory)
   func constructHistory(systemPrompt string, input map[string]interface{}) []llm.Message {
       var contextBuilder strings.Builder
       
       // 构建结构化上下文 (顺序: 文档 -> 上游分析 -> 目标)
       sections := []struct {
           key   string
           label string
       }{
           {"document_content", "document_content"},
           {"proposal", "proposal"},
           {"aggregated_outputs", "previous_analyses"},
           {"optimization_objective", "optimization_objective"},
       }
       
       for _, sec := range sections {
           if val, ok := input[sec.key].(string); ok && val != "" {
               contextBuilder.WriteString(fmt.Sprintf("<%s>\n%s\n</%s>\n\n", sec.label, val, sec.label))
           }
       }
       
       userContent := contextBuilder.String()
       if userContent == "" {
           userContent = "Begin task."
       }
       
       return []llm.Message{
           {Role: "system", Content: systemPrompt},
           {Role: "user", Content: userContent},
       }
   }
   ```

---

### Fix-3: 应用层聚合策略 (可选注入)

**目标**: 允许特定工作流覆盖默认聚合策略。

**实现方案**:

```go
// internal/core/workflow/merge_council.go (新增，应用层)

// CouncilMergeStrategy 专为 Council Debate 设计的聚合策略
// 会特别处理 agent_output 字段的聚合
type CouncilMergeStrategy struct{}

func (s *CouncilMergeStrategy) Merge(inputs []map[string]interface{}) map[string]interface{} {
    merged := make(map[string]interface{})
    var agentOutputs []string
    
    for i, inp := range inputs {
        // 收集所有 agent_output
        if out, ok := inp["agent_output"].(string); ok {
            agentOutputs = append(agentOutputs, out)
        }
        
        // 透传首次出现的上下文字段
        for k, v := range inp {
            if k == "agent_output" {
                continue // 已特殊处理
            }
            if _, exists := merged[k]; !exists {
                merged[k] = v
            }
        }
        
        // 保留分支数据供调试
        merged[fmt.Sprintf("branch_%d", i)] = inp
    }
    
    // 聚合所有 Agent 输出
    if len(agentOutputs) > 0 {
        merged["aggregated_outputs"] = strings.Join(agentOutputs, "\n\n---\n\n")
    }
    
    return merged
}
```

**注入时机** (在 Session/Handler 创建 Engine 时):
```go
// api/handler/workflow_run.go (示例)
engine := workflow.NewEngine(session)
engine.MergeStrategy = &workflow.CouncilMergeStrategy{} // 覆盖默认策略
```

---

## 5. 修改文件清单 (Files to Modify)

| 文件路径                                  | 修改类型       | 层级 | 说明                                      |
| :---------------------------------------- | :------------- | :--- | :---------------------------------------- |
| `internal/core/workflow/merge.go`         | New File       | 骨架 | MergeStrategy 接口 + DefaultMergeStrategy |
| `internal/core/workflow/engine.go`        | Major Refactor | 骨架 | 入度计算、deliverToDownstream、策略注入   |
| `internal/core/workflow/merge_council.go` | New File       | 应用 | CouncilMergeStrategy 实现                 |
| `internal/core/workflow/nodes/agent.go`   | Modify         | 应用 | 透传字段列表、Prompt 构建修复             |
| `internal/core/workflow/engine_test.go`   | Add Tests      | 测试 | Join 机制单元测试                         |
| `internal/api/handler/workflow_run.go`    | Modify         | 应用 | 注入 CouncilMergeStrategy                 |

---

## 5. 验收标准 (Acceptance Criteria)

- [ ] **TC-1**: 在 `council_debate` 流程中，`Adjudicator` 仅执行 **一次**
- [ ] **TC-2**: `Adjudicator` 收到的输入包含 `aggregated_outputs` 字段，内容为正反双方的完整分析
- [ ] **TC-3**: `Adjudicator` 收到的输入包含原始 `document_content`
- [ ] **TC-4**: 所有现有单元测试通过 (`make test`)
- [ ] **TC-5**: 新增 Join 机制的单元测试覆盖

---

## 6. 风险评估 (Risk Assessment)

| 风险                      | 影响                   | 缓解措施                             |
| :------------------------ | :--------------------- | :----------------------------------- |
| 并发竞态条件              | 数据丢失或重复执行     | 使用 `sync.Mutex` 保护 `pendingData` |
| 入度计算错误              | 节点永不触发或提前触发 | 增加 Graph 验证步骤                  |
| 上下文过大导致 Token 超限 | LLM 调用失败           | 可在后续迭代中引入截断策略           |

---

## 7. 依赖与阻塞 (Dependencies)

- **无外部依赖**: 此修复为核心引擎改进，不依赖新的外部服务
- **向后兼容**: 对于入度=1的节点，行为与当前一致

---

## 8. 附录：工作流图示

### 当前错误流程

```
Start → Parallel → [Aff]  →→→ Adjudicator (触发1，只有正方)
                 ↘ [Neg] →→→ Adjudicator (触发2，只有反方)
```

### 修复后正确流程

```
Start → Parallel → [Aff]  ─┐
                 ↘ [Neg] ─┴→ [Join/Merge] → Adjudicator (触发1，正反双方聚合)
```
