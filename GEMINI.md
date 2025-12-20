# 🛠️ The Council - 开发与交付规约 (v0.14.0)

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
*   **Atomic Commits**: 单一逻辑变更/commit，规范 message (`feat:`, `fix:`).
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
*   **Sync**: 后端 Struct ↔ 前端 TS 类型 (`tygo`).
*   **Vector DB**: PGvector (`embedding`), `uuid` 主键.
*   **Migrations**: `golang-migrate` (`YYYYMMDDHHMMSS_name.up.sql`).

## 4. AI & Prompt Engineering

*   **Prompt Management**: 存放在 `/prompts/*.md`，禁止 Hardcode。
*   **Template**: `{{.Context}}`, `{{.UserQuery}}` 占位符。
*   **Safety**: 处理 Context Overflow (自动截断)，版本变更需 Review。

## 5. 测试与质量 (QA & Testing)

*   **Command**: 使用 `make test` (自动过滤 infra 噪音)，严禁 `go test`.
*   **Coverage**: 核心业务逻辑 (Core) **100%** 覆盖。
*   **Mock Strategy**: 业务测试禁止连真实 DB/LLM，使用 `MockProvider`.
*   **TDD**: 红 (Test) -> 绿 (Impl) -> 蓝 (Refactor).
