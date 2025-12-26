# Sprint 9 执行计划: Quality & Stability

> **Sprint 目标**: 偿还技术债务，建立自动化质量防线  
> **优先级**: P0 (阻塞生产发布)  
> **预计工时**: 51 小时 (3 周)  
> **开始日期**: 2025-12-30  
> **结束日期**: 2026-01-17

---

## 一、目标与背景

### 1.1 Sprint 目标

建立自动化质量防线，消除技术债务，为 v1.0 生产发布做好准备。

**核心目标**:
1. 前后端 WebSocket 消息类型自动同步
2. 核心层代码覆盖率达到 80%
3. 所有默认值可配置化
4. Core 层解耦，不依赖具体 Agent 实现
5. 建立端到端消息格式验证机制

### 1.2 问题背景

在 Sprint 6-8 的开发过程中，发现以下系统性隐患:

**技术债务**:
- WebSocket 消息类型未自动生成，前后端类型不一致风险高
- 核心层测试覆盖率不足 (~65%)，重构风险大
- 硬编码默认值散落各处，配置不灵活
- Core 层依赖具体 Agent 实现 (Adjudicator)，耦合度高
- ID 参数命名语义不明确，跨层传递易出错

**预防需求**:
- 建立前后端类型契约自动校验机制
- 提升测试覆盖率，建立质量红线
- 统一配置管理，支持环境变量覆盖
- 清理架构依赖，提升可维护性

---

## 二、任务分解与优先级

### 2.1 任务总览

| 任务 ID | 任务名称 | 优先级 | 工时 | 依赖 | 负责模块 |
|---------|----------|--------|------|------|----------|
| 9.1 | WebSocket 类型自动生成 (tygo) | P0 | 8h | 无 | 基础设施 |
| 9.2 | 端到端消息格式测试 | P0 | 10h | 9.1 | QA |
| 9.3 | 核心层测试覆盖率提升 | P1 | 12h | 无 | 后端 |
| 9.4 | 硬编码默认值清理 & 配置化 | P1 | 8h | 无 | 全栈 |
| 9.7 | Core 层解耦 (移除 Adjudicator) | P1 | 6h | 无 | 后端 |
| 9.5 | ID 命名规范审查与重构 | P2 | 4h | 无 | 全栈 |
| 9.6 | Store 类型 UI 字段完整性检查 | P2 | 3h | 无 | 前端 |

**总计**: 51 小时

---

## 三、详细任务规格

### 3.1 任务 9.1: WebSocket 类型自动生成 (tygo)

**优先级**: P0  
**工时**: 8 小时  
**依赖**: 无

#### 目标

使用 tygo 工具自动从 Go 结构体生成 TypeScript 类型定义，确保前后端 WebSocket 消息类型一致。

#### 实现方案

**Step 1: 安装 tygo**
```bash
go install github.com/gzuidhof/tygo@latest
```

**Step 2: 配置 tygo.yaml**
```yaml
# tygo.yaml
packages:
  - path: "github.com/yourusername/council/internal/api/ws"
    type_mappings:
      time.Time: "string"
      json.RawMessage: "any"
    output_path: "frontend/src/types/websocket.generated.ts"
    
  - path: "github.com/yourusername/council/internal/core/workflow"
    output_path: "frontend/src/types/workflow.generated.ts"
```

**Step 3: 标记需要生成的类型**
```go
// internal/api/ws/message.go

//go:generate tygo

// WebSocket 下行事件
type TokenChunk struct {
    Type    string `json:"type" tygo:"required"`
    NodeID  string `json:"node_id" tygo:"required"`
    Content string `json:"content" tygo:"required"`
}

type NodeStatusUpdate struct {
    Type   string `json:"type" tygo:"required"`
    NodeID string `json:"node_id" tygo:"required"`
    Status string `json:"status" tygo:"required"` // running, completed, failed
}

// ... 其他消息类型
```

**Step 4: 集成到 Makefile**
```makefile
.PHONY: generate-types
generate-types:
	@echo "Generating TypeScript types from Go..."
	tygo generate
	@echo "✅ Type generation complete"

.PHONY: verify-types
verify-types: generate-types
	@echo "Verifying type consistency..."
	git diff --exit-code frontend/src/types/*.generated.ts || \
		(echo "❌ Generated types are out of sync! Run 'make generate-types'" && exit 1)
	@echo "✅ Type consistency verified"
```

**Step 5: CI 集成**
```yaml
# .github/workflows/ci.yml
- name: Verify type consistency
  run: make verify-types
```

#### 验收标准

- [x] tygo 配置文件创建 (`tygo.yaml`)
- [x] 所有 WebSocket 消息类型自动生成
- [x] Makefile 增加 `generate-types` 命令
- [x] CI 中强制检查类型一致性
- [x] 文档更新 (如何添加新消息类型)

#### 产出物

- `tygo.yaml` - tygo 配置文件
- `frontend/src/types/websocket.generated.ts` - 生成的 WebSocket 类型
- `frontend/src/types/workflow.generated.ts` - 生成的工作流类型
- `Makefile` 更新
- `.github/workflows/ci.yml` 更新

---

### 3.2 任务 9.2: 端到端消息格式测试

**优先级**: P0  
**工时**: 10 小时  
**依赖**: 任务 9.1

#### 目标

建立 E2E 测试套件，验证前后端 WebSocket 消息格式完全一致，防止运行时类型错误。

#### 实现方案

**Step 1: 创建消息契约测试**
```typescript
// e2e/tests/websocket-contract.spec.ts

import { test, expect } from '@playwright/test';

test.describe('WebSocket Message Contract', () => {
  test('TokenChunk message format', async ({ page }) => {
    await page.goto('/chat');
    
    // 启动一个简单的工作流
    await page.click('[data-testid="start-session"]');
    
    // 监听 WebSocket 消息
    const messages: any[] = [];
    page.on('websocket', ws => {
      ws.on('framereceived', frame => {
        const msg = JSON.parse(frame.payload as string);
        messages.push(msg);
      });
    });
    
    // 等待收到 TokenChunk 消息
    await page.waitForTimeout(2000);
    
    const tokenChunk = messages.find(m => m.type === 'token_chunk');
    expect(tokenChunk).toBeDefined();
    expect(tokenChunk).toMatchObject({
      type: expect.any(String),
      node_id: expect.any(String),
      content: expect.any(String),
    });
  });

  test('NodeStatusUpdate message format', async ({ page }) => {
    // 类似测试...
  });
  
  // 为每种消息类型添加测试
});
```

**Step 2: 创建消息模拟器 (Mock Server)**
```typescript
// e2e/utils/ws-mock-server.ts

export class WebSocketMockServer {
  sendTokenChunk(nodeId: string, content: string) {
    this.send({
      type: 'token_chunk',
      node_id: nodeId,
      content,
    });
  }

  sendNodeStatusUpdate(nodeId: string, status: 'running' | 'completed' | 'failed') {
    this.send({
      type: 'node_status_update',
      node_id: nodeId,
      status,
    });
  }

  private send(message: any) {
    // 发送到客户端
  }
}
```

**Step 3: 集成到 CI**
```yaml
# .github/workflows/e2e.yml
- name: Run WebSocket contract tests
  run: |
    cd e2e
    npm run test:contract
```

#### 验收标准

- [x] 为所有 WebSocket 消息类型添加 E2E 测试
- [x] 测试覆盖以下场景:
  - [x] TokenChunk (流式输出)
  - [x] NodeStatusUpdate (节点状态)
  - [x] ExecutionPaused (暂停)
  - [x] ExecutionCompleted (完成)
  - [x] HumanInteractionRequired (人工审核)
  - [x] Error (错误)
- [x] Mock Server 支持所有消息类型
- [x] CI 中自动运行契约测试

#### 产出物

- `e2e/tests/websocket-contract.spec.ts` - 消息契约测试
- `e2e/utils/ws-mock-server.ts` - WebSocket Mock Server
- `.github/workflows/e2e.yml` 更新

---

### 3.3 任务 9.3: 核心层测试覆盖率提升

**优先级**: P1  
**工时**: 12 小时  
**依赖**: 无

#### 目标

将 `internal/core` 层代码覆盖率从 ~65% 提升至 ≥ 80%。

#### 实现方案

**Step 1: 审查当前覆盖率**
```bash
cd internal/core
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

**Step 2: 识别低覆盖率模块**

优先补充以下模块的测试:
- `workflow/executor.go` (工作流执行器)
- `workflow/nodes/vote.go` (投票节点)
- `workflow/nodes/loop.go` (循环节点)
- `middleware/versioning.go` (版本控制)
- `memory/retriever.go` (记忆检索)

**Step 3: 补充单元测试**
```go
// internal/core/workflow/executor_test.go

func TestExecutor_Execute_WithVoteNode(t *testing.T) {
    // 测试投票节点执行逻辑
    executor := NewExecutor(...)
    
    workflow := &Workflow{
        Nodes: []Node{
            {Type: "start"},
            {Type: "vote", Config: VoteConfig{...}},
            {Type: "end"},
        },
    }
    
    err := executor.Execute(context.Background(), workflow)
    assert.NoError(t, err)
    
    // 验证投票结果
    result := executor.GetResult()
    assert.Equal(t, "expected_winner", result.WinnerID)
}

// 为边界情况添加测试
func TestExecutor_Execute_WithInvalidVoteScore(t *testing.T) {
    // 测试无效分数处理
}

func TestExecutor_Execute_WithTieBreaking(t *testing.T) {
    // 测试平局处理
}
```

**Step 4: 集成测试**
```go
// internal/core/workflow/integration_test.go

func TestWorkflow_EndToEnd_DebateFlow(t *testing.T) {
    // 测试完整的辩论流程
    // Memory Retrieval -> Parallel Agents -> Vote -> HumanReview -> Loop
}
```

**Step 5: CI 门禁**
```yaml
# .github/workflows/ci.yml
- name: Test coverage
  run: |
    cd internal/core
    go test -coverprofile=coverage.out ./...
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$COVERAGE < 80" | bc -l) )); then
      echo "❌ Coverage is $COVERAGE%, must be >= 80%"
      exit 1
    fi
    echo "✅ Coverage is $COVERAGE%"
```

#### 验收标准

- [x] Core 层覆盖率 ≥ 80%
- [x] 所有关键路径有单元测试
- [x] 边界情况有专门测试
- [x] CI 中强制检查覆盖率
- [x] 覆盖率报告可视化 (HTML)

#### 产出物

- `internal/core/*/test.go` - 新增单元测试
- `internal/core/workflow/integration_test.go` - 集成测试
- `.github/workflows/ci.yml` 更新

---

### 3.4 任务 9.4: 硬编码默认值清理 & 配置化

**优先级**: P1  
**工时**: 8 小时  
**依赖**: 无

#### 目标

将散落在代码中的硬编码默认值统一迁移到配置文件或环境变量。

#### 实现方案

**Step 1: 审查硬编码值**

扫描以下模式:
```bash
grep -r "gpt-4\|gemini-\|http://\|localhost:8080" internal/
```

**Step 2: 创建默认值注册表**
```go
// internal/pkg/config/defaults.go

package config

type Defaults struct {
    LLM      LLMDefaults
    Server   ServerDefaults
    Database DatabaseDefaults
}

type LLMDefaults struct {
    DefaultProvider string // 默认 Provider
    DefaultModel    string // 默认模型
    Timeout         int    // 超时时间 (秒)
    MaxRetries      int    // 最大重试次数
}

type ServerDefaults struct {
    Host string
    Port int
    Env  string // development, production
}

type DatabaseDefaults struct {
    MaxOpenConns int
    MaxIdleConns int
}

// 从环境变量或配置文件加载
func LoadDefaults() *Defaults {
    return &Defaults{
        LLM: LLMDefaults{
            DefaultProvider: getEnv("LLM_DEFAULT_PROVIDER", "gemini"),
            DefaultModel:    getEnv("LLM_DEFAULT_MODEL", "gemini-2.0-flash-exp"),
            Timeout:         getEnvInt("LLM_TIMEOUT", 60),
            MaxRetries:      getEnvInt("LLM_MAX_RETRIES", 3),
        },
        Server: ServerDefaults{
            Host: getEnv("SERVER_HOST", "localhost"),
            Port: getEnvInt("SERVER_PORT", 8080),
            Env:  getEnv("ENV", "development"),
        },
        Database: DatabaseDefaults{
            MaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 25),
            MaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),
        },
    }
}
```

**Step 3: 更新代码使用配置**
```go
// internal/core/workflow/nodes/agent.go

// Before:
defaultModel := "gpt-4" // ❌ 硬编码

// After:
defaultModel := a.Config.Defaults.LLM.DefaultModel // ✅ 配置化
```

**Step 4: 更新 .env.example**
```bash
# .env.example

# LLM 配置
LLM_DEFAULT_PROVIDER=gemini
LLM_DEFAULT_MODEL=gemini-2.0-flash-exp
LLM_TIMEOUT=60
LLM_MAX_RETRIES=3

# 服务器配置
SERVER_HOST=localhost
SERVER_PORT=8080
ENV=development

# 数据库配置
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
```

**Step 5: 验证脚本**
```bash
# scripts/verify-no-hardcodes.sh

#!/bin/bash

HARDCODES=$(grep -r "gpt-4\|gemini-2\|localhost:8080" internal/ | grep -v "_test.go" | grep -v "// OK")

if [ -n "$HARDCODES" ]; then
    echo "❌ Found hardcoded values:"
    echo "$HARDCODES"
    exit 1
fi

echo "✅ No hardcoded values found"
```

#### 验收标准

- [x] 所有硬编码默认值迁移到配置
- [x] 支持环境变量覆盖
- [x] `.env.example` 包含所有配置项
- [x] CI 中验证无硬编码值
- [x] 配置文档更新

#### 产出物

- `internal/pkg/config/defaults.go` - 默认值注册表
- `.env.example` 更新
- `scripts/verify-no-hardcodes.sh` - 验证脚本
- `docs/guide/configuration.md` - 配置文档

---

### 3.5 任务 9.7: Core 层解耦 (移除 Adjudicator)

**优先级**: P1  
**工时**: 6 小时  
**依赖**: 无

#### 目标

移除 Core 层对具体 Agent 实现 (Adjudicator) 的依赖，通过接口抽象提升可维护性。

#### 实现方案

**Step 1: 识别耦合点**
```bash
grep -r "Adjudicator" internal/core/
```

**Step 2: 定义抽象接口**
```go
// internal/core/agent/interface.go

package agent

type DecisionMaker interface {
    MakeDecision(ctx context.Context, options []Option) (Decision, error)
}

type Option struct {
    ID      string
    Content string
    Score   float64
}

type Decision struct {
    SelectedID string
    Reason     string
}
```

**Step 3: 重构 Vote 节点**
```go
// internal/core/workflow/nodes/vote.go

// Before:
func (v *VoteNode) Execute(ctx context.Context) error {
    // 硬编码 Adjudicator
    adjudicator := agent.NewAdjudicator()
    decision := adjudicator.Decide(scores)
}

// After:
func (v *VoteNode) Execute(ctx context.Context) error {
    // 通过配置注入 DecisionMaker
    decisionMaker := v.Config.DecisionMaker
    decision := decisionMaker.MakeDecision(ctx, options)
}
```

**Step 4: 更新测试**
```go
// internal/core/workflow/nodes/vote_test.go

type MockDecisionMaker struct{}

func (m *MockDecisionMaker) MakeDecision(ctx context.Context, options []agent.Option) (agent.Decision, error) {
    // Mock 实现
    return agent.Decision{SelectedID: "option1"}, nil
}

func TestVoteNode_WithMockDecisionMaker(t *testing.T) {
    node := &VoteNode{
        Config: VoteConfig{
            DecisionMaker: &MockDecisionMaker{},
        },
    }
    
    err := node.Execute(context.Background())
    assert.NoError(t, err)
}
```

#### 验收标准

- [x] Core 层不再导入 Adjudicator 实现
- [x] DecisionMaker 接口定义清晰
- [x] Vote 节点支持注入 DecisionMaker
- [x] 测试使用 Mock 实现
- [x] 架构文档更新

#### 产出物

- `internal/core/agent/interface.go` - DecisionMaker 接口
- `internal/core/workflow/nodes/vote.go` 重构
- `internal/core/workflow/nodes/vote_test.go` 更新
- `docs/tdd/01_architecture.md` 更新

---

### 3.6 任务 9.5: ID 命名规范审查与重构

**优先级**: P2  
**工时**: 4 小时  
**依赖**: 无

#### 目标

统一 ID 参数命名，消除语义歧义，提升代码可读性。

#### 实现方案

**Step 1: 建立命名规范**
```markdown
# ID 命名规范

| 实体 | ID 字段名 | 示例 |
|------|-----------|------|
| 会话 | sessionID | "sess_abc123" |
| 工作流 | workflowID | "wf_xyz789" |
| 节点 | nodeID | "node_001" |
| Agent | agentID | "agent_affirmative" |
| 群组 | groupID | "group_council" |
| 运行 | runID | "run_20250101_001" |

**命名原则**:
- 使用驼峰命名法 (camelCase)
- 包含实体类型 (如 sessionID 而非 id)
- 避免缩写 (如 wfID -> workflowID)
```

**Step 2: 扫描不规范命名**
```bash
grep -r "\"id\":" internal/ frontend/src/
```

**Step 3: 批量重构**
```bash
# 示例: 重命名 id -> sessionID
find internal/ -name "*.go" -exec sed -i 's/SessionID string `json:"id"`/SessionID string `json:"session_id"`/g' {} \;
```

**Step 4: 更新前端类型**
```typescript
// frontend/src/types/session.ts

// Before:
interface Session {
  id: string; // ❌ 不明确
}

// After:
interface Session {
  sessionID: string; // ✅ 明确
}
```

#### 验收标准

- [x] ID 命名规范文档完成
- [x] 所有 ID 字段符合规范
- [x] 前后端类型一致
- [x] 代码审查通过

#### 产出物

- `docs/guide/naming-conventions.md` - 命名规范文档
- 代码重构 (后端 + 前端)

---

### 3.7 任务 9.6: Store 类型 UI 字段完整性检查

**优先级**: P2  
**工时**: 3 小时  
**依赖**: 无

#### 目标

确保所有 Zustand Store 类型定义包含 UI 所需的全部字段。

#### 实现方案

**Step 1: 审查 Store 定义**
```typescript
// frontend/src/stores/useSessionStore.ts

interface SessionStore {
  sessions: Session[];
  currentSessionID: string | null;
  
  // ❌ 缺少 UI 需要的字段
  // isLoading: boolean;
  // error: string | null;
}
```

**Step 2: 补全字段**
```typescript
interface SessionStore {
  // 数据
  sessions: Session[];
  currentSessionID: string | null;
  
  // UI 状态
  isLoading: boolean;
  error: string | null;
  
  // 操作
  fetchSessions: () => Promise<void>;
  selectSession: (sessionID: string) => void;
  clearError: () => void;
}
```

**Step 3: 创建检查清单**
```markdown
# Store 字段完整性检查清单

对于每个 Store，确保包含:
- [x] 数据字段 (如 sessions, agents)
- [x] UI 状态字段 (isLoading, error)
- [x] 操作方法 (fetch, create, update, delete)
- [x] 选择器 (getById, getByName)
```

#### 验收标准

- [x] 所有 Store 包含必需字段
- [x] UI 组件无类型错误
- [x] 检查清单文档完成

#### 产出物

- `frontend/src/stores/*.ts` 更新
- `docs/guide/store-checklist.md` - Store 检查清单

---

## 四、执行计划

### 4.1 Week 1: P0 任务 (2025-12-30 - 2026-01-03)

| 日期 | 任务 | 工时 | 负责人 |
|------|------|------|--------|
| 12/30 | 任务 9.1: tygo 配置 | 4h | Backend Lead |
| 12/31 | 任务 9.1: Makefile + CI | 4h | DevOps |
| 01/02 | 任务 9.3: 覆盖率审查 | 4h | Backend Team |
| 01/03 | 任务 9.3: 补充测试 (1/3) | 8h | Backend Team |

**里程碑**: tygo 配置完成，类型自动生成可用

---

### 4.2 Week 2: P0-P1 任务 (2026-01-06 - 2026-01-10)

| 日期 | 任务 | 工时 | 负责人 |
|------|------|------|--------|
| 01/06 | 任务 9.2: E2E 契约测试 | 10h | QA Team |
| 01/07 | 任务 9.4: 默认值清理 | 4h | Full Stack |
| 01/08 | 任务 9.4: 配置化重构 | 4h | Full Stack |
| 01/09 | 任务 9.3: 补充测试 (2/3) | 8h | Backend Team |
| 01/10 | 任务 9.7: Core 解耦 | 6h | Backend Lead |

**里程碑**: E2E 消息测试完成，覆盖率达标

---

### 4.3 Week 3: P1-P2 任务 (2026-01-13 - 2026-01-17)

| 日期 | 任务 | 工时 | 负责人 |
|------|------|------|--------|
| 01/13 | 任务 9.5: ID 命名审查 | 2h | Tech Lead |
| 01/14 | 任务 9.5: ID 重构 | 2h | Full Stack |
| 01/15 | 任务 9.6: Store 字段检查 | 3h | Frontend Team |
| 01/16 | 集成测试 | 4h | QA Team |
| 01/17 | Sprint 回顾 | 2h | 全员 |

**里程碑**: Sprint 9 全部完成，质量门禁建立

---

## 五、验收标准

### 5.1 功能验收

- [x] 前后端 WebSocket 消息类型自动同步
- [x] E2E 测试覆盖所有消息类型
- [x] 核心层代码覆盖率 ≥ 80%
- [x] 所有默认值可配置化
- [x] Core 层不依赖具体 Agent 实现
- [x] ID 命名符合规范
- [x] Store 类型包含 UI 必需字段

### 5.2 质量验收

- [x] CI 中所有检查通过
- [x] 无 Lint 错误
- [x] 无已知 Bug
- [x] 文档完整更新

### 5.3 非功能验收

- [x] 性能无退化
- [x] 向后兼容 (现有功能不受影响)
- [x] 代码审查通过

---

## 六、风险与应对

### 6.1 风险识别

| 风险 | 影响 | 概率 | 应对措施 |
|------|------|------|----------|
| tygo 工具兼容性问题 | 高 | 低 | 提前验证，准备手动生成方案 |
| 测试覆盖率提升困难 | 中 | 中 | 聚焦核心路径，接受非关键代码低覆盖率 |
| 配置化重构引入新 Bug | 中 | 中 | 小步快跑，每次重构后回归测试 |
| 工作量被低估 | 低 | 低 | 优先 P0 任务，P2 任务可延期 |

### 6.2 应急计划

如果 Week 2 结束时进度落后:
- 优先完成任务 9.1, 9.2 (P0)
- 任务 9.5, 9.6 (P2) 延期到 Sprint 10

---

## 七、成功指标

| 指标 | 目标值 | 测量方式 |
|------|--------|----------|
| 任务完成率 | 100% | 所有任务状态为 Done |
| 代码覆盖率 (Core) | ≥ 80% | `go test -cover` |
| E2E 测试通过率 | 100% | Playwright 报告 |
| CI 通过率 | 100% | GitHub Actions |
| 配置化率 | 100% | 无硬编码值检查 |

---

## 八、总结

Sprint 9 是项目质量防线建设的关键 Sprint，完成后将:
- ✅ 消除前后端类型不一致风险
- ✅ 建立测试覆盖率红线
- ✅ 实现配置灵活管理
- ✅ 提升架构可维护性
- ✅ 为 v1.0 生产发布扫清障碍

**下一步**: 完成 Sprint 5 性能优化和安全强化，准备生产发布。
