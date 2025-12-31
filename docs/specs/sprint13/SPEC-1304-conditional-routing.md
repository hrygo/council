# SPEC-1304: 条件路由抽象

> **Sprint**: 13  
> **状态**: Draft  
> **作者**: AI Assistant  
> **创建日期**: 2025-12-31  
> **优先级**: P2 (架构优化)  
> **预估工作量**: 1.5 天  
> **依赖**: SPEC-1303  
> **后续**: 无

---

## 1. 目标

将 Loop 节点的条件路由逻辑从 Engine 中抽离，实现骨架层零业务感知。

### 1.1 当前问题

```go
// ❌ 违规：Engine 硬编码 Loop 的路由逻辑和字段名
// internal/core/workflow/engine.go
func (e *Engine) deliverToDownstream(...) {
    if node.Type == NodeTypeLoop && len(node.NextIDs) >= 2 {
        shouldExit, _ := output["should_exit"].(bool)  // ← 业务字段硬编码
        if shouldExit {
            targetNextIDs = []string{node.NextIDs[1]}
        } else {
            targetNextIDs = []string{node.NextIDs[0]}
        }
    }
}
```

### 1.2 目标状态

```go
// ✅ 目标：Processor 实现路由接口，Engine 只调用接口
type ConditionalRouter interface {
    GetNextNodes(output map[string]interface{}) []string
}

// LoopProcessor 实现 ConditionalRouter
func (l *LoopProcessor) GetNextNodes(output map[string]interface{}) []string {
    if output["should_exit"].(bool) {
        return []string{l.ExitNodeID}
    }
    return []string{l.ContinueNodeID}
}
```

---

## 2. 详细设计

### 2.1 骨架层接口

在 `internal/core/workflow/processor.go` 中添加：

```go
// ConditionalRouter 定义条件路由能力。
// Processor 可选实现此接口以自定义下游节点选择。
// 如果 Processor 未实现此接口，Engine 使用节点定义的 NextIDs。
type ConditionalRouter interface {
	// GetNextNodes 根据处理输出决定下游节点。
	// 参数:
	//   - output: 当前节点的处理输出
	//   - allNextIDs: 节点定义的所有可能下游节点
	// 返回:
	//   - 实际应触发的下游节点 ID 列表
	GetNextNodes(output map[string]interface{}, allNextIDs []string) []string
}
```

### 2.2 修改 Engine

`internal/core/workflow/engine.go` 中的 `deliverToDownstream`：

```go
func (e *Engine) deliverToDownstream(ctx context.Context, nodeID string, output map[string]interface{}) {
	node := e.Graph.Nodes[nodeID]
	if node == nil {
		return
	}

	// 确定下游节点
	targetNextIDs := node.NextIDs

	// 检查 Processor 是否实现了 ConditionalRouter
	processor, err := e.NodeFactory.Create(node, e.factoryDeps)
	if err == nil {
		if router, ok := processor.(ConditionalRouter); ok {
			// Processor 自行决定路由
			targetNextIDs = router.GetNextNodes(output, node.NextIDs)
		}
	}

	// 检查是否为 Loop-back 投递 (用于绕过入度检查)
	isLoopBack := e.isLoopBackDelivery(node, targetNextIDs, output)

	for _, nextID := range targetNextIDs {
		// ... existing join logic ...
	}
}

// isLoopBackDelivery 检测是否为循环回路投递。
// 循环回路需要绕过入度等待以避免死锁。
func (e *Engine) isLoopBackDelivery(node *Node, targetNextIDs []string, output map[string]interface{}) bool {
	// 如果 Processor 实现了 ConditionalRouter 且选择了第一个 NextID
	// 且节点类型支持循环 (Loop, Vote 等)，则认为是循环回路
	if len(targetNextIDs) == 1 && len(node.NextIDs) >= 2 {
		if targetNextIDs[0] == node.NextIDs[0] {
			// 第一个 NextID 通常是"继续"路径
			return true
		}
	}
	return false
}
```

### 2.3 LoopProcessor 实现

修改 `internal/core/workflow/nodes/loop.go`：

```go
// Compile-time check
var _ workflow.ConditionalRouter = (*LoopProcessor)(nil)

type LoopProcessor struct {
	MaxRounds       int
	ExitOnScore     int
	Session         *workflow.Session
	PassthroughKeys []string
	
	// 新增：路由目标 (由 Factory 注入)
	ContinueNodeID string // 循环继续的目标节点
	ExitNodeID     string // 循环退出的目标节点
}

// GetNextNodes 实现 ConditionalRouter 接口。
func (l *LoopProcessor) GetNextNodes(output map[string]interface{}, allNextIDs []string) []string {
	shouldExit, _ := output["should_exit"].(bool)
	
	if shouldExit {
		// 退出路径：第二个 NextID
		if len(allNextIDs) >= 2 {
			return []string{allNextIDs[1]}
		}
		return allNextIDs
	}
	
	// 继续路径：第一个 NextID
	if len(allNextIDs) >= 1 {
		return []string{allNextIDs[0]}
	}
	return allNextIDs
}
```

### 2.4 VoteProcessor 实现 (可选)

如果 Vote 节点也有条件路由需求：

```go
var _ workflow.ConditionalRouter = (*VoteProcessor)(nil)

func (v *VoteProcessor) GetNextNodes(output map[string]interface{}, allNextIDs []string) []string {
	approved, _ := output["approved"].(bool)
	
	if approved {
		// 通过路径
		if len(allNextIDs) >= 1 {
			return []string{allNextIDs[0]}
		}
	} else {
		// 拒绝路径
		if len(allNextIDs) >= 2 {
			return []string{allNextIDs[1]}
		}
	}
	return allNextIDs
}
```

### 2.5 删除 Engine 中的硬编码

删除以下代码：

```go
// 删除这段硬编码逻辑
if node.Type == NodeTypeLoop && len(node.NextIDs) >= 2 {
    shouldExit, _ := output["should_exit"].(bool)
    if shouldExit {
        targetNextIDs = []string{node.NextIDs[1]}
    } else {
        targetNextIDs = []string{node.NextIDs[0]}
    }
}
```

---

## 3. 实施步骤

### Step 1: 定义接口
在 `processor.go` 中添加 `ConditionalRouter` 接口。

### Step 2: 修改 Engine
更新 `deliverToDownstream` 使用接口检测。

### Step 3: 实现 LoopProcessor
让 LoopProcessor 实现 `ConditionalRouter`。

### Step 4: 实现 VoteProcessor (可选)
如需要，让 VoteProcessor 也实现接口。

### Step 5: 删除硬编码
从 Engine 中删除 NodeTypeLoop 特殊处理。

### Step 6: 验证
```bash
grep -r "should_exit" internal/core/workflow/engine.go
# 应返回空

go build ./...
go test ./internal/...
```

---

## 4. 验收标准

- [ ] `ConditionalRouter` 接口定义在骨架层
- [ ] Engine 不包含 `should_exit` 字段名
- [ ] Engine 不包含 `NodeTypeLoop` 特殊路由逻辑
- [ ] LoopProcessor 实现 `ConditionalRouter`
- [ ] 循环工作流 (council_optimize) 正常执行
- [ ] 编译通过，测试通过，E2E 正常

---

## 5. 风险

| 风险               | 影响 | 缓解                         |
| :----------------- | :--- | :--------------------------- |
| 接口检测性能开销   | 低   | 每次执行只检测一次，可忽略   |
| Processor 创建两次 | 中   | 可缓存第一次创建的 Processor |

### 5.1 性能优化 (可选)

如果 Processor 创建开销大，可以缓存：

```go
func (e *Engine) executeNode(...) {
    processor, err := e.NodeFactory.Create(node, e.factoryDeps)
    // ... execution ...
    
    // 缓存 processor 供 deliverToDownstream 使用
    e.lastProcessor = processor
}

func (e *Engine) deliverToDownstream(...) {
    if router, ok := e.lastProcessor.(ConditionalRouter); ok {
        targetNextIDs = router.GetNextNodes(output, node.NextIDs)
    }
}
```

---

## 6. 修改文件清单

| 操作     | 文件                                   |
| :------- | :------------------------------------- |
| 修改     | `internal/core/workflow/processor.go`  |
| 修改     | `internal/core/workflow/engine.go`     |
| 修改     | `internal/core/workflow/nodes/loop.go` |
| 可选修改 | `internal/core/workflow/nodes/vote.go` |

---

## 7. 完成后架构状态

完成 SPEC-1301 ~ 1304 后，架构达到目标状态：

```
骨架层 (internal/core/workflow/):
├── engine.go           ✅ 无业务字段，只调用接口
├── processor.go        ✅ 定义 NodeProcessor, ConditionalRouter
├── factory.go          ✅ 定义 NodeFactory 接口
├── merge.go            ✅ 定义 MergeStrategy 接口
├── passthrough.go      ✅ 通用工具函数
└── nodes/              ✅ 只包含通用处理器 (可选)

应用层 (internal/council/):
├── context_keys.go     ✅ Council 专属字段定义
├── merge_strategy.go   ✅ CouncilMergeStrategy
├── factory.go          ✅ CouncilNodeFactory
└── processors/         ✅ Council 专属处理器 (可选)

验证命令:
$ grep -r "document_content\|should_exit\|aggregated_outputs" internal/core/workflow/
(无输出 = 成功)
```
