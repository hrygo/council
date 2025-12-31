# SPEC-1302: 配置驱动透传

> **Sprint**: 13  
> **状态**: Draft  
> **作者**: AI Assistant  
> **创建日期**: 2025-12-31  
> **优先级**: P1 (架构优化)  
> **预估工作量**: 1.5 天  
> **依赖**: SPEC-1301  
> **后续**: SPEC-1303

---

## 1. 目标

将 Processor 中硬编码的透传字段列表改为**配置注入**，实现骨架层零业务感知。

### 1.1 当前问题

```go
// ❌ 违规：骨架层硬编码业务字段
// internal/core/workflow/nodes/agent.go
var agentPassthroughKeys = []string{
    "document_content",      // ← Council 专属
    "proposal",              // ← Council 专属
}
```

### 1.2 目标状态

```go
// ✅ 正确：Processor 接收配置
// internal/council/processors/agent.go
type CouncilAgentProcessor struct {
    PassthroughKeys []string // 从构造函数注入
}
```

---

## 2. 详细设计

### 2.1 骨架层：通用透传工具

新建 `internal/core/workflow/passthrough.go`：

```go
package workflow

// PassthroughConfig 定义透传配置。
// 这是一个通用结构，不包含任何业务字段名。
type PassthroughConfig struct {
	// Keys 需要从 input 透传到 output 的字段名列表
	Keys []string
}

// ApplyPassthrough 将配置中指定的字段从 input 复制到 output。
// 这是一个纯工具函数，不包含业务逻辑。
func ApplyPassthrough(input, output map[string]interface{}, config PassthroughConfig) {
	for _, key := range config.Keys {
		if val, ok := input[key]; ok {
			output[key] = val
		}
	}
}
```

### 2.2 应用层：Processor 参数化

#### 2.2.1 修改 AgentProcessor

**Before** (`internal/core/workflow/nodes/agent.go`)：
```go
var agentPassthroughKeys = []string{"document_content", ...}

func (a *AgentProcessor) Process(...) {
    for _, key := range agentPassthroughKeys { // ❌ 硬编码
        if val, ok := input[key]; ok {
            output[key] = val
        }
    }
}
```

**After** (`internal/core/workflow/nodes/agent.go`)：
```go
type AgentProcessor struct {
    // ...existing fields...
    PassthroughKeys []string // ← 新增：从外部注入
}

func (a *AgentProcessor) Process(...) {
    workflow.ApplyPassthrough(input, output, workflow.PassthroughConfig{
        Keys: a.PassthroughKeys, // ✅ 使用注入的配置
    })
}
```

#### 2.2.2 修改 LoopProcessor

**After** (`internal/core/workflow/nodes/loop.go`)：
```go
type LoopProcessor struct {
    MaxRounds       int
    ExitOnScore     int
    Session         *workflow.Session
    PassthroughKeys []string // ← 新增
}

func (l *LoopProcessor) Process(...) {
    // ... existing logic ...
    
    workflow.ApplyPassthrough(input, output, workflow.PassthroughConfig{
        Keys: l.PassthroughKeys,
    })
}
```

#### 2.2.3 修改 StartProcessor

**After** (`internal/core/workflow/nodes/start.go`)：
```go
type StartProcessor struct {
    OutputKeys []string // ← 新增：需要输出的字段
}

func (s *StartProcessor) Process(...) {
    output := make(map[string]interface{})
    
    for _, key := range s.OutputKeys {
        if val, ok := input[key]; ok {
            output[key] = val
        }
    }
    
    // ... rest of logic ...
}
```

#### 2.2.4 修改 EndProcessor

**After** (`internal/core/workflow/nodes/end.go`)：
```go
type EndProcessor struct {
    LLM       llm.LLMProvider
    Model     string
    Prompt    string
    InputKeys []string // ← 新增：需要读取的字段
}

func (e *EndProcessor) Process(...) {
    var contentBuilder strings.Builder
    
    for _, key := range e.InputKeys {
        if val, ok := input[key].(string); ok && val != "" {
            contentBuilder.WriteString(fmt.Sprintf("## %s\n%s\n\n", key, val))
        }
    }
    
    // ... rest of logic ...
}
```

### 2.3 修改 Factory

`internal/core/workflow/nodes/factory.go` 需要注入配置：

```go
func NewNodeFactory(deps FactoryDeps) func(*Node) (NodeProcessor, error) {
    return func(node *Node) (NodeProcessor, error) {
        switch node.Type {
        case NodeTypeAgent:
            return &AgentProcessor{
                // ...existing fields...
                PassthroughKeys: deps.Config.AgentPassthroughKeys,
            }, nil
            
        case NodeTypeLoop:
            return &LoopProcessor{
                // ...existing fields...
                PassthroughKeys: deps.Config.LoopPassthroughKeys,
            }, nil
            
        case NodeTypeStart:
            return &StartProcessor{
                OutputKeys: deps.Config.StartOutputKeys,
            }, nil
            
        case NodeTypeEnd:
            return &EndProcessor{
                // ...existing fields...
                InputKeys: deps.Config.EndInputKeys,
            }, nil
        }
    }
}
```

### 2.4 API 层注入配置

`internal/api/handler/workflow.go`：

```go
import "github.com/hrygo/council/internal/council"

func (h *WorkflowHandler) Execute(c *gin.Context) {
    // ...
    
    engine.NodeFactory = nodes.NewNodeFactory(nodes.FactoryDeps{
        Registry:  h.registry,
        AgentRepo: h.agentRepo,
        Session:   session,
        MemoryMgr: h.memoryMgr,
        Config: nodes.ProcessorConfig{
            AgentPassthroughKeys: council.AgentPassthroughKeys,
            LoopPassthroughKeys:  council.LoopPassthroughKeys,
            StartOutputKeys:      council.StartOutputKeys,
            EndInputKeys:         council.EndInputKeys,
        },
    })
}
```

---

## 3. 实施步骤

### Step 1: 创建 passthrough.go
创建骨架层通用透传工具。

### Step 2: 修改 Processor 结构体
为 AgentProcessor, LoopProcessor, StartProcessor, EndProcessor 添加配置字段。

### Step 3: 修改 Processor 逻辑
使用 `ApplyPassthrough` 替代硬编码循环。

### Step 4: 修改 Factory
注入配置到 Processor。

### Step 5: 修改 API 层
从 council 包读取配置并注入。

### Step 6: 删除硬编码
删除 `agentPassthroughKeys`, `loopPassthroughKeys` 等硬编码变量。

### Step 7: 验证
```bash
grep -r "document_content" internal/core/workflow/nodes/
# 应返回空（除注释外）

go build ./...
go test ./internal/...
```

---

## 4. 验收标准

- [ ] `internal/core/workflow/passthrough.go` 存在
- [ ] AgentProcessor, LoopProcessor, StartProcessor, EndProcessor 均接收配置参数
- [ ] `grep -r "document_content" internal/core/workflow/nodes/*.go` 返回空（不计注释）
- [ ] 配置从 `internal/council/context_keys.go` 读取
- [ ] 编译通过，测试通过，E2E 正常

---

## 5. 修改文件清单

| 操作 | 文件                                      |
| :--- | :---------------------------------------- |
| 新建 | `internal/core/workflow/passthrough.go`   |
| 修改 | `internal/core/workflow/nodes/agent.go`   |
| 修改 | `internal/core/workflow/nodes/loop.go`    |
| 修改 | `internal/core/workflow/nodes/start.go`   |
| 修改 | `internal/core/workflow/nodes/end.go`     |
| 修改 | `internal/core/workflow/nodes/factory.go` |
| 修改 | `internal/api/handler/workflow.go`        |
