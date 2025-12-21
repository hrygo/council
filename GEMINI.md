# 🛠️ The Council - 开发与交付规约 (v0.15.0)

> **原则**: 务求实效 (Pragmatism) | 前端驱动 (Contract First) | 模拟优先 (Mock First) | 渐进交付 (Atomic Delivery) | TDD 开发

## 1. 核心架构与目录 (Architecture)

| 领域     | 技术栈                    | 关键目录                          |
| :------- | :------------------------ | :-------------------------------- |
| **App**  | **Web App (React SPA)**   | `frontend/src/` (Vite, Tailwind)  |
| **API**  | **Go (Gin, WebSocket)**   | `cmd/council/`, `internal/api/`   |
| **Core** | **Workflow Engine**       | `internal/core/` (Agents, Memory) |
| **Data** | **PostgreSQL + pgvector** | `internal/infrastructure/db/`     |
| **Docs** | **PRD / Specs**           | `docs/`                           |

**交付规约**:
*   **Atomic Delivery**: 每次 PR 必须是完整可运行单元，禁止 Broken Build。
*   **Atomic Commits**: 单一逻辑变更/commit，规范 message (`feat:`, `fix:`)。
*   **Strict Quality Gates**: 每个 SPEC 完成后必须通过所有验收标准 (Acceptance Criteria) 及 CI 检查 (Lint + Test)。
*   **Track Progress**: 任务完成后必须更新 `docs/development_plan.md` 进度矩阵。

## 2. 统一编码规约 (Coding Standards)

| 维度     | Go (Backend)                                   | React/TS (Frontend)                            |
| :------- | :--------------------------------------------- | :--------------------------------------------- |
| **风格** | `gofmt` + `goimports` (Auto-save)              | `Prettier` + `ESLint`                          |
| **Lint** | `golangci-lint` (CI 强制)                      | No `any`, Strict Mode                          |
| **命名** | `snake_case` (DB/JSON), `PascalCase` (Structs) | `PascalCase` (Components), `camelCase` (Props) |
| **错误** | 必须 wrap: `fmt.Errorf("...: %w", err)`        | Error Boundary + Toast 通知                    |
| **状态** | 接受 Interface，返回 Struct                    | Zustand Stores (`useSessionStore`)             |
| **并发** | 必须传递 `ctx`, 禁止裸 `go func`               | `useEffect` cleanups, RQ/SWR                   |
| **UI**   | N/A                                            | TailwindCSS 优先, `clsx` 动态类                |
| **i18n** | N/A                                            | `react-i18next`, 禁止 Hardcode                 |

## 3. 接口与数据 (API & Data)

**RESTful / WebSocket** (`/api/v1`)
*   **Sync**: 后端 Struct ↔ 前端 TS 类型 (`tygo`)。
*   **Vector DB**: PGvector (`embedding`), `uuid` 主键。
*   **Migrations**: `golang-migrate` (`YYYYMMDDHHMMSS_name.up.sql`)。

## 4. 前后端协作规约 (Frontend-Backend Contract) 🆕

> 基于 2025-12-21 WebSocket 调试经验总结

### 4.1 消息格式一致性

| 规则              | 说明                                                     |
| :---------------- | :------------------------------------------------------- |
| **JSON 字段命名** | 前后端必须使用**完全相同**的字段名，优先采用前端命名惯例 |
| **类型同步**      | 所有共享类型必须通过 `tygo` 自动生成，**禁止手动维护**   |
| **测试覆盖**      | WebSocket 消息需有**端到端格式验证测试**                 |

### 4.2 ID 语义规范

| ID 类型    | 命名规范        | 示例                         |
| :--------- | :-------------- | :--------------------------- |
| 数据库主键 | `{entity}_uuid` | `agent_uuid`, `session_uuid` |
| 逻辑标识符 | `{context}_id`  | `node_id`, `workflow_id`     |

**禁止混用**: 在事件/消息中，必须使用逻辑标识符，而非数据库 UUID。

### 4.3 数据传递完整性

*   **类型设计**: 必须包含 UI 渲染所需的全部字段（如 `name`, `type`）
*   **信息传递**: 禁止在中间层丢弃上游传入的有效信息

## 5. 默认值与配置规约 (Defaults & Config) 🆕

| 规则           | 说明                                     |
| :------------- | :--------------------------------------- |
| **禁止硬编码** | 所有默认值必须通过配置文件或工厂方法获取 |
| **降级一致性** | 降级逻辑必须与当前 Provider/环境兼容     |
| **配置优先级** | 环境变量 > 配置文件 > 代码内默认值       |

**示例**:
```go
// ❌ 禁止
if model == "" { model = "gpt-4" }

// ✅ 正确
if model == "" { model = registry.GetDefaultModel() }
```

## 6. AI & Prompt Engineering

*   **Prompt Management**: 存放在 `/prompts/*.md`，禁止 Hardcode。
*   **Template**: `{{.Context}}`, `{{.UserQuery}}` 占位符。
*   **Safety**: 处理 Context Overflow (自动截断)，版本变更需 Review。

## 7. 测试与质量 (QA & Testing)

*   **Command**: 使用 `make test` (自动过滤 infra 噪音)，严禁 `go test`。
*   **Coverage**: 核心业务逻辑 (Core) **100%** 覆盖。
*   **Mock Strategy**: 业务测试禁止连真实 DB/LLM，使用 `MockProvider`。
*   **TDD**: 红 (Test) -> 绿 (Impl) -> 蓝 (Refactor)。
*   **端到端测试** 🆕: WebSocket 消息格式需跨前后端验证。

