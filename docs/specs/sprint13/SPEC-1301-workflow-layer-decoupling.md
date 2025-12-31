# SPEC-1301: 应用层目录结构

> **Sprint**: 13  
> **状态**: Ready  
> **作者**: AI Assistant  
> **创建日期**: 2025-12-31  
> **优先级**: P1 (架构优化)  
> **预估工作量**: 1 天  
> **依赖**: 无  
> **后续**: SPEC-1302, SPEC-1303, SPEC-1304

---

## 1. 背景

本 SPEC 是 Sprint 13 架构解耦重构的**第一阶段**，专注于创建应用层目录结构，为后续重构奠定基础。

### 1.1 系列 SPEC 概览

| SPEC     | 名称             | 范围                                 | 状态   |
| :------- | :--------------- | :----------------------------------- | :----- |
| **1301** | 应用层目录结构   | 创建目录 + 移动 MergeStrategy        | 本文档 |
| 1302     | 配置驱动透传     | PassthroughConfig + Processor 参数化 | Draft  |
| 1303     | NodeFactory 接口 | 工厂接口 + CouncilNodeFactory        | Draft  |
| 1304     | 条件路由抽象     | ConditionalRouter + 动态路由         | Draft  |

### 1.2 本阶段目标

1. 创建 `internal/council/` 目录作为应用层容器
2. 将 Council 专属的 `merge_council.go` 移动到应用层
3. 定义 Council 上下文字段常量，为后续 Processor 参数化做准备

---

## 2. 详细设计

### 2.1 目录结构变更

```
internal/
├── core/
│   └── workflow/
│       ├── engine.go           # 不变
│       ├── merge.go            # 不变 (接口定义)
│       ├── merge_council.go    # ← 删除 (移动到 council)
│       └── ...
│
└── council/                    # ← 新建
    ├── context_keys.go         # Council 上下文字段常量
    └── merge_strategy.go       # CouncilMergeStrategy (从 core 移动)
```

### 2.2 新建文件

#### 2.2.1 `internal/council/context_keys.go`

```go
// Package council 包含 Council Debate/Optimize 工作流的应用层实现。
// 此包依赖 internal/core/workflow (骨架层)，但反之不可。
package council

// CouncilContextKeys 定义 Council 工作流使用的所有上下文字段。
// 这些字段名是业务概念，不应出现在骨架层代码中。
var CouncilContextKeys = []string{
	"document_content",       // 原始文档内容
	"proposal",               // 方案摘要
	"optimization_objective", // 优化目标
	"attachments",            // 附件列表
	"combined_context",       // 合并后的附件内容
	"session_id",             // 会话 ID
	"aggregated_outputs",     // 聚合的 Agent 输出
	"agent_output",           // 单个 Agent 输出
	"history_context",        // 历史上下文
}

// AgentPassthroughKeys 定义 Agent 节点需要透传的字段。
var AgentPassthroughKeys = []string{
	"document_content",
	"proposal",
	"optimization_objective",
	"attachments",
	"combined_context",
	"session_id",
	"aggregated_outputs",
}

// LoopPassthroughKeys 定义 Loop 节点需要透传的字段。
var LoopPassthroughKeys = []string{
	"document_content",
	"proposal",
	"optimization_objective",
	"combined_context",
	"session_id",
}

// StartOutputKeys 定义 Start 节点需要输出的字段。
var StartOutputKeys = []string{
	"document_content",
	"proposal",
	"optimization_objective",
	"attachments",
	"combined_context",
}

// EndInputKeys 定义 End 节点需要读取的字段。
var EndInputKeys = []string{
	"document_content",
	"combined_context",
	"proposal",
	"aggregated_outputs",
	"agent_output",
}
```

#### 2.2.2 `internal/council/merge_strategy.go`

```go
// Package council 包含 Council Debate/Optimize 工作流的应用层实现。
package council

import (
	"strings"

	"github.com/hrygo/council/internal/core/workflow"
)

// Compile-time check: CouncilMergeStrategy implements MergeStrategy
var _ workflow.MergeStrategy = (*CouncilMergeStrategy)(nil)

// CouncilMergeStrategy 是 Council 工作流专用的聚合策略。
// 它将多个 Agent 的 agent_output 聚合为 aggregated_outputs，
// 并透传其他上下文字段。
type CouncilMergeStrategy struct{}

// Merge 聚合多个上游节点的输出。
func (s *CouncilMergeStrategy) Merge(inputs []map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	var agentOutputs []string

	for _, input := range inputs {
		// 聚合 agent_output
		if output, ok := input["agent_output"].(string); ok && output != "" {
			agentOutputs = append(agentOutputs, output)
		}

		// 透传上下文字段 (first-win 策略)
		for _, key := range CouncilContextKeys {
			if _, exists := merged[key]; !exists {
				if val, ok := input[key]; ok {
					merged[key] = val
				}
			}
		}
	}

	// 将多个 agent_output 聚合为 aggregated_outputs
	if len(agentOutputs) > 0 {
		merged["aggregated_outputs"] = strings.Join(agentOutputs, "\n\n---\n\n")
	}

	return merged
}
```

### 2.3 修改文件

#### 2.3.1 `internal/api/handler/workflow.go`

更新 import 路径：

```go
// Before
engine.MergeStrategy = &workflow.CouncilMergeStrategy{}

// After
import "github.com/hrygo/council/internal/council"
engine.MergeStrategy = &council.CouncilMergeStrategy{}
```

### 2.4 删除文件

| 文件                                      | 原因                                        |
| :---------------------------------------- | :------------------------------------------ |
| `internal/core/workflow/merge_council.go` | 移动到 `internal/council/merge_strategy.go` |

---

## 3. 实施步骤

### Step 1: 创建目录和文件
```bash
mkdir -p internal/council
touch internal/council/context_keys.go
touch internal/council/merge_strategy.go
```

### Step 2: 编写代码
按照 2.2 节的代码创建文件。

### Step 3: 更新 workflow.go
修改 import 路径和 MergeStrategy 初始化。

### Step 4: 删除旧文件
```bash
rm internal/core/workflow/merge_council.go
```

### Step 5: 验证
```bash
go build ./...
go test ./internal/...
```

---

## 4. 验收标准

- [ ] `internal/council/` 目录存在
- [ ] `internal/council/context_keys.go` 定义所有 Council 字段常量
- [ ] `internal/council/merge_strategy.go` 包含 CouncilMergeStrategy
- [ ] `internal/core/workflow/merge_council.go` 已删除
- [ ] `go build ./...` 编译通过
- [ ] `go test ./internal/...` 全部通过
- [ ] 工作流执行正常 (E2E)

---

## 5. 风险

| 风险                | 影响     | 缓解                     |
| :------------------ | :------- | :----------------------- |
| Import 路径更新遗漏 | 编译失败 | IDE 自动修复 + grep 检查 |

---

## 6. 修改文件清单

| 操作 | 文件                                      |
| :--- | :---------------------------------------- |
| 新建 | `internal/council/context_keys.go`        |
| 新建 | `internal/council/merge_strategy.go`      |
| 修改 | `internal/api/handler/workflow.go`        |
| 删除 | `internal/core/workflow/merge_council.go` |
