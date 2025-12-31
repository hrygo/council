# SPEC-1303: NodeFactory 接口

> **Sprint**: 13  
> **状态**: Draft  
> **作者**: AI Assistant  
> **创建日期**: 2025-12-31  
> **优先级**: P1 (架构优化)  
> **预估工作量**: 2 天  
> **依赖**: SPEC-1302  
> **后续**: SPEC-1304

---

## 1. 目标

定义 `NodeFactory` 接口，实现应用层可插拔。未来可支持多种工作流类型（Council、其他业务）。

### 1.1 当前问题

```go
// ❌ 问题：Engine.NodeFactory 是函数类型，与应用层耦合
type Engine struct {
    NodeFactory func(node *Node) (NodeProcessor, error) // 闭包，难以替换
}
```

### 1.2 目标状态

```go
// ✅ 目标：NodeFactory 是接口，可注入不同实现
type Engine struct {
    NodeFactory NodeFactory // 接口类型
}

// Council 应用层实现
type CouncilNodeFactory struct { ... }
func (f *CouncilNodeFactory) Create(node *Node, deps FactoryDeps) (NodeProcessor, error)

// 未来：其他应用层
type OtherAppNodeFactory struct { ... }
```

---

## 2. 详细设计

### 2.1 骨架层接口

新建 `internal/core/workflow/factory.go`：

```go
package workflow

// NodeFactory 定义节点处理器工厂接口。
// 应用层实现此接口以提供业务专属的 Processor。
type NodeFactory interface {
	// Create 根据节点定义创建处理器。
	// deps 包含运行时依赖，由 Engine 在执行时提供。
	Create(node *Node, deps FactoryDeps) (NodeProcessor, error)
}

// FactoryDeps 定义工厂创建 Processor 所需的依赖。
// 使用 interface{} 避免骨架层依赖具体实现。
type FactoryDeps struct {
	// LLMRegistry 提供 LLM Provider 访问
	LLMRegistry interface{}
	
	// AgentRepository 提供 Agent 数据访问
	AgentRepository interface{}
	
	// Session 当前工作流会话
	Session *Session
	
	// MemoryManager 提供记忆系统访问
	MemoryManager interface{}
	
	// Config 应用层配置 (透传字段等)
	Config interface{}
}

// DefaultNodeFactory 提供骨架层默认工厂。
// 只处理通用节点类型 (start, end, parallel)。
// 应用层应组合此工厂处理业务节点。
type DefaultNodeFactory struct{}

func (f *DefaultNodeFactory) Create(node *Node, deps FactoryDeps) (NodeProcessor, error) {
	switch node.Type {
	case NodeTypeStart:
		return &GenericPassthroughProcessor{}, nil
	case NodeTypeEnd:
		return &GenericEndProcessor{}, nil
	default:
		return nil, fmt.Errorf("unknown node type: %s (application factory required)", node.Type)
	}
}

// GenericPassthroughProcessor 是一个通用透传处理器。
// 它将所有输入字段原样传递到输出。
type GenericPassthroughProcessor struct{}

func (p *GenericPassthroughProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
	stream <- StreamEvent{Type: "node_state_change", Data: map[string]interface{}{"status": "running"}}
	
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}
	
	stream <- StreamEvent{Type: "node_state_change", Data: map[string]interface{}{"status": "completed"}}
	return output, nil
}

// GenericEndProcessor 是一个通用结束处理器。
// 它只输出完成事件，不做 LLM 调用。
type GenericEndProcessor struct{}

func (p *GenericEndProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
	stream <- StreamEvent{Type: "node_state_change", Data: map[string]interface{}{"status": "running"}}
	stream <- StreamEvent{Type: "execution:completed"}
	stream <- StreamEvent{Type: "node_state_change", Data: map[string]interface{}{"status": "completed"}}
	return input, nil
}
```

### 2.2 修改 Engine

`internal/core/workflow/engine.go`：

```go
type Engine struct {
	// ... existing fields ...
	
	// NodeFactory 创建节点处理器。
	// 必须在 Run() 之前设置，否则使用 DefaultNodeFactory。
	NodeFactory NodeFactory
	
	// factoryDeps 工厂依赖，在 NewEngine 时初始化。
	factoryDeps FactoryDeps
}

func NewEngine(session *Session, deps FactoryDeps) *Engine {
	e := &Engine{
		// ... existing init ...
		NodeFactory: &DefaultNodeFactory{},
		factoryDeps: deps,
	}
	e.computeInDegrees()
	return e
}

func (e *Engine) executeNode(ctx context.Context, nodeID string, input map[string]interface{}) {
	// ... existing code ...
	
	// 使用接口调用替代函数调用
	processor, err := e.NodeFactory.Create(node, e.factoryDeps)
	if err != nil {
		e.emitError(nodeID, err)
		return
	}
	
	// ... rest of execution ...
}
```

### 2.3 应用层实现

新建 `internal/council/factory.go`：

```go
package council

import (
	"fmt"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

// CouncilNodeFactory 是 Council 工作流专用的节点工厂。
type CouncilNodeFactory struct {
	// BaseFactory 用于处理通用节点
	BaseFactory workflow.NodeFactory
}

// NewCouncilNodeFactory 创建 Council 节点工厂。
func NewCouncilNodeFactory() *CouncilNodeFactory {
	return &CouncilNodeFactory{
		BaseFactory: &workflow.DefaultNodeFactory{},
	}
}

// Create 创建 Council 专属处理器。
func (f *CouncilNodeFactory) Create(node *workflow.Node, deps workflow.FactoryDeps) (workflow.NodeProcessor, error) {
	// 类型断言获取具体依赖
	registry, _ := deps.LLMRegistry.(*llm.Registry)
	agentRepo, _ := deps.AgentRepository.(agent.Repository)
	memMgr, _ := deps.MemoryManager.(memory.MemoryManager)

	switch node.Type {
	case workflow.NodeTypeStart:
		return &nodes.StartProcessor{
			OutputKeys: StartOutputKeys,
		}, nil

	case workflow.NodeTypeEnd:
		provider, _ := registry.GetLLMProvider("default")
		return &nodes.EndProcessor{
			LLM:       provider,
			Model:     registry.GetDefaultModel(),
			InputKeys: EndInputKeys,
		}, nil

	case workflow.NodeTypeAgent:
		agentID, _ := node.Properties["agent_uuid"].(string)
		return &nodes.AgentProcessor{
			NodeID:          node.ID,
			AgentID:         agentID,
			AgentRepo:       agentRepo,
			Registry:        registry,
			Session:         deps.Session,
			PassthroughKeys: AgentPassthroughKeys,
		}, nil

	case workflow.NodeTypeLoop:
		maxRounds, _ := node.Properties["max_rounds"].(float64)
		exitScore, _ := node.Properties["exit_on_score"].(float64)
		return &nodes.LoopProcessor{
			MaxRounds:       int(maxRounds),
			ExitOnScore:     int(exitScore),
			Session:         deps.Session,
			PassthroughKeys: LoopPassthroughKeys,
		}, nil

	case workflow.NodeTypeHumanReview:
		timeout, _ := node.Properties["timeout_minutes"].(float64)
		return &nodes.HumanReviewProcessor{
			TimeoutMinutes: int(timeout),
		}, nil

	case workflow.NodeTypeVote:
		threshold, _ := node.Properties["threshold"].(float64)
		return &nodes.VoteProcessor{
			Threshold: threshold,
		}, nil

	case workflow.NodeTypeMemoryRetrieval:
		return nodes.NewMemoryRetrievalProcessor(memMgr), nil

	case workflow.NodeTypeContextSynth:
		return &nodes.ContextSynthesizerProcessor{}, nil

	case workflow.NodeTypeFactCheck:
		provider, _ := registry.GetLLMProvider("default")
		return &nodes.FactCheckProcessor{
			LLM:   provider,
			Model: registry.GetDefaultModel(),
		}, nil

	default:
		// Fallback to base factory
		return f.BaseFactory.Create(node, deps)
	}
}
```

### 2.4 API 层注入

`internal/api/handler/workflow.go`：

```go
import "github.com/hrygo/council/internal/council"

func (h *WorkflowHandler) Execute(c *gin.Context) {
	// ...
	
	// 创建 Engine，注入 Council 工厂
	deps := workflow.FactoryDeps{
		LLMRegistry:     h.registry,
		AgentRepository: h.agentRepo,
		Session:         session,
		MemoryManager:   h.memoryMgr,
	}
	engine := workflow.NewEngine(session, deps)
	engine.NodeFactory = council.NewCouncilNodeFactory()
	engine.MergeStrategy = &council.CouncilMergeStrategy{}
	
	// ...
}
```

---

## 3. 实施步骤

### Step 1: 创建骨架层接口
创建 `internal/core/workflow/factory.go`。

### Step 2: 修改 Engine
更新 Engine 结构体和 NewEngine 签名。

### Step 3: 创建通用处理器
实现 GenericPassthroughProcessor, GenericEndProcessor。

### Step 4: 创建 Council 工厂
创建 `internal/council/factory.go`。

### Step 5: 修改 API 层
注入 CouncilNodeFactory。

### Step 6: 迁移现有 Factory
保留 `nodes/factory.go` 作为辅助，或逐步删除。

### Step 7: 验证
```bash
go build ./...
go test ./internal/...
```

---

## 4. 验收标准

- [ ] `workflow.NodeFactory` 是接口类型
- [ ] `workflow.FactoryDeps` 使用 `interface{}` 避免依赖
- [ ] `CouncilNodeFactory` 实现 `NodeFactory` 接口
- [ ] Engine 可接受不同的 NodeFactory 实现
- [ ] 所有节点类型在 CouncilNodeFactory 中有处理
- [ ] 编译通过，测试通过，E2E 正常

---

## 5. 修改文件清单

| 操作     | 文件                                                     |
| :------- | :------------------------------------------------------- |
| 新建     | `internal/core/workflow/factory.go`                      |
| 新建     | `internal/council/factory.go`                            |
| 修改     | `internal/core/workflow/engine.go`                       |
| 修改     | `internal/api/handler/workflow.go`                       |
| 可选删除 | `internal/core/workflow/nodes/factory.go` (或保留为辅助) |
