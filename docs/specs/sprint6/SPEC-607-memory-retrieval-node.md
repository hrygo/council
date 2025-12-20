# SPEC-607: Memory Retrieval Node

> **优先级**: P0  
> **类型**: Feature (Workflow Node)  
> **预估工时**: 4h

## 1. 概述

实现 `NodeTypeMemoryRetrieval` 节点，用于在 Workflow 执行过程中自动检索相关历史上下文。

## 2. 目标

- 在 Agent 节点执行前，自动从 Memory 系统检索相关历史
- 将检索结果注入到后续节点的 Context 中
- 支持 `skill.md` Step 1 的 "Compress History Context" 功能

## 3. 节点规格

### 3.1 节点定义

| 属性           | 值                                                 |
| :------------- | :------------------------------------------------- |
| **Type**       | `memory_retrieval`                                 |
| **Input**      | Session Context (当前话题/文档)                    |
| **Output**     | 历史摘要注入到 Context                             |
| **Properties** | `max_results`, `time_range`, `relevance_threshold` |

### 3.2 配置示例

```json
{
  "id": "memory_retrieval",
  "type": "memory_retrieval",
  "name": "Load History Context",
  "properties": {
    "max_results": 5,
    "time_range_days": 7,
    "relevance_threshold": 0.7,
    "include_verdicts": true
  },
  "next_ids": ["parallel_analysis"]
}
```

## 4. 技术实现

### 4.1 文件结构

```
internal/core/workflow/nodes/
  memory_retrieval.go      # 节点实现
  memory_retrieval_test.go # 测试
```

### 4.2 接口定义

```go
// internal/core/workflow/nodes/memory_retrieval.go
package nodes

import (
    "context"
    
    "github.com/hrygo/council/internal/core/memory"
    "github.com/hrygo/council/internal/core/workflow"
)

type MemoryRetrievalProcessor struct {
    memoryService memory.Service
}

func NewMemoryRetrievalProcessor(ms memory.Service) *MemoryRetrievalProcessor {
    return &MemoryRetrievalProcessor{memoryService: ms}
}

func (p *MemoryRetrievalProcessor) Process(
    ctx context.Context,
    node *workflow.Node,
    input map[string]interface{},
    emitter workflow.EventEmitter,
) (map[string]interface{}, error) {
    // 1. 提取当前话题/文档标识
    topic := input["topic"].(string)
    
    // 2. 从 Memory 系统检索相关历史
    props := node.Properties
    results, err := p.memoryService.Retrieve(ctx, memory.Query{
        Topic:              topic,
        MaxResults:         props.GetInt("max_results", 5),
        TimeRangeDays:      props.GetInt("time_range_days", 7),
        RelevanceThreshold: props.GetFloat("relevance_threshold", 0.7),
    })
    if err != nil {
        return nil, err
    }
    
    // 3. 组装历史摘要
    historySummary := p.formatHistorySummary(results)
    
    // 4. 注入到输出 Context 中
    output := make(map[string]interface{})
    for k, v := range input {
        output[k] = v
    }
    output["history_context"] = historySummary
    
    // 5. 发送事件
    emitter.Emit(workflow.StreamEvent{
        Type:    "memory_retrieved",
        Payload: map[string]interface{}{"count": len(results)},
    })
    
    return output, nil
}

func (p *MemoryRetrievalProcessor) formatHistorySummary(results []memory.Entry) string {
    // 格式化为 Markdown 摘要
    // ...
}
```

### 4.3 注册到 Factory

```go
// internal/core/workflow/nodes/factory.go
case workflow.NodeTypeMemoryRetrieval:
    return NewMemoryRetrievalProcessor(deps.MemoryService)
```

### 4.4 新增 NodeType 常量

```go
// internal/core/workflow/types.go
const (
    // ...existing types...
    NodeTypeMemoryRetrieval NodeType = "memory_retrieval"
)
```

## 5. 与 Memory 系统集成

### 5.1 Memory Service 接口

```go
// internal/core/memory/service.go
type Service interface {
    Retrieve(ctx context.Context, query Query) ([]Entry, error)
    Store(ctx context.Context, entry Entry) error
}

type Query struct {
    Topic              string
    MaxResults         int
    TimeRangeDays      int
    RelevanceThreshold float64
}

type Entry struct {
    ID        string
    Content   string
    Timestamp time.Time
    Score     float64
    Metadata  map[string]interface{}
}
```

### 5.2 3-Tier 映射

| Memory Tier         | 用途               | 本节点访问   |
| ------------------- | ------------------ | ------------ |
| Tier 1 (Quarantine) | 原始对话日志       | ❌ 不直接访问 |
| Tier 2 (Working)    | 工作记忆/向量检索  | ✅ 主要数据源 |
| Tier 3 (Long-term)  | 用户标记的重要知识 | ✅ 辅助数据源 |

## 6. 验收标准

- [ ] `NodeTypeMemoryRetrieval` 常量已定义
- [ ] `MemoryRetrievalProcessor` 实现完成
- [ ] Factory 正确注册该节点类型
- [ ] 节点可从 Memory 系统检索历史
- [ ] 检索结果正确注入到输出 Context
- [ ] Workflow Canvas UI 可渲染该节点

## 7. 测试

```go
func TestMemoryRetrievalProcessor_Process(t *testing.T) {
    mockMemory := &MockMemoryService{
        Results: []memory.Entry{
            {ID: "1", Content: "Previous verdict: Approved", Score: 0.9},
        },
    }
    
    processor := NewMemoryRetrievalProcessor(mockMemory)
    
    output, err := processor.Process(ctx, node, input, emitter)
    
    assert.NoError(t, err)
    assert.Contains(t, output["history_context"], "Previous verdict")
}
```

## 8. 前端 UI 设计 (Issue 4 Remediation)

### 8.1 节点外观

| 属性     | 值               |
| :------- | :--------------- |
| **图标** | 📚 (书籍/历史)    |
| **颜色** | #6366F1 (Indigo) |
| **形状** | 圆角矩形         |
| **标签** | "Load History"   |

### 8.2 节点配置面板

当用户在 Workflow Canvas 中选中 `memory_retrieval` 节点时，右侧面板显示：

```
┌────────────────────────────────────────┐
│ 📚 Memory Retrieval                    │
├────────────────────────────────────────┤
│ Max Results:      [5]  ▼               │
│ Time Range (days): [7]  ▼              │
│ Relevance Threshold: [0.7] ──●───      │
│ Include Verdicts:  ☑                   │
├────────────────────────────────────────┤
│ Preview:                               │
│ ┌────────────────────────────────────┐ │
│ │ Will retrieve up to 5 historical   │ │
│ │ items from the last 7 days with    │ │
│ │ relevance score >= 0.7             │ │
│ └────────────────────────────────────┘ │
└────────────────────────────────────────┘
```

### 8.3 React 组件定义

```typescript
// frontend/src/components/workflow/nodes/MemoryRetrievalNode.tsx
export interface MemoryRetrievalNodeProps {
  id: string;
  data: {
    label: string;
    properties: {
      max_results: number;
      time_range_days: number;
      relevance_threshold: number;
      include_verdicts: boolean;
    };
  };
}

export const MemoryRetrievalNode: React.FC<NodeProps<MemoryRetrievalNodeProps>> = ({ data }) => {
  return (
    <div className="memory-retrieval-node">
      <div className="node-icon">📚</div>
      <div className="node-label">{data.label}</div>
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
};
```

### 8.4 CSS 样式

```css
/* frontend/src/styles/nodes.css */
.memory-retrieval-node {
  background: linear-gradient(135deg, #6366F1 0%, #4F46E5 100%);
  border-radius: 8px;
  padding: 12px 16px;
  color: white;
  min-width: 120px;
  text-align: center;
  box-shadow: 0 4px 6px rgba(99, 102, 241, 0.3);
}

.memory-retrieval-node .node-icon {
  font-size: 24px;
  margin-bottom: 4px;
}
```

### 8.5 节点注册

```typescript
// frontend/src/components/workflow/nodeTypes.ts
import { MemoryRetrievalNode } from './nodes/MemoryRetrievalNode';

export const nodeTypes = {
  // ...existing types
  memory_retrieval: MemoryRetrievalNode,
};
```

## 9. 依赖

- **Memory System**: `internal/core/memory` 必须已实现 `Retrieve` 方法

