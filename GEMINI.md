# 🛠️ The Council - 工程开发规约

> **Version**: v0.6.1 | **Architecture**: WebApp (React SPA + Go API + Docker PostgreSQL)

---

## 1. 目录结构

```text
council/
├── cmd/council/              # Go 入口 (main.go)
├── internal/
│   ├── api/                 # HTTP/WS Handler
│   ├── core/                # 领域逻辑 (Workflow, Agent, Memory)
│   └── infrastructure/      # DB, LLM Client, SearchTool
├── pkg/                     # 公共库
├── frontend/src/
│   ├── components/          # React 组件
│   ├── stores/              # Zustand (Session/Config/Layout)
│   └── i18n/locales/        # 中英文案 (zh-CN, en-US)
├── prompts/                 # AI 提示词版本管理
├── docker-compose.yml       # PostgreSQL + pgvector
└── docs/                    # PRD.md, TDD.md
```

---

## 2. 后端规约 (Go)

| 规则         | 说明                                               |
| ------------ | -------------------------------------------------- |
| **格式化**   | `gofmt` / `goimports` 保存自动执行                 |
| **Linter**   | CI 集成 `golangci-lint`                            |
| **错误处理** | 必须 `%w` 包装，禁止 `_ = func()`                  |
| **并发**     | 必须传递 `ctx context.Context`，禁止裸 `go func()` |
| **接口**     | Accept Interfaces, Return Structs                  |

```go
// 错误包装示例
return fmt.Errorf("failed to init agent %s: %w", id, err)
```

---

## 3. 前端规约 (React/TS)

| 规则     | 说明                                                                    |
| -------- | ----------------------------------------------------------------------- |
| **组件** | FC + Hooks，禁止 Class Component                                        |
| **命名** | 组件 `PascalCase.tsx`，Props 用 `interface`                             |
| **状态** | Zustand 分 Store：`useSessionStore`, `useConfigStore`, `useLayoutStore` |
| **样式** | TailwindCSS 优先，`clsx` 处理动态类                                     |
| **i18n** | `react-i18next`，禁止 Hardcode 文案                                     |

```tsx
const { t } = useTranslation('chat');
<input placeholder={t('input_placeholder')} />
```

---

## 4. API 规约

### REST (CRUD)
```
GET/POST   /api/v1/groups
GET/PUT    /api/v1/groups/:id
GET/POST   /api/v1/agents
POST       /api/v1/workflows/generate
```

### WebSocket (实时流)
```json
{"event": "agent:speaking", "data": {...}}
{"event": "node:completed", "data": {...}}
{"event": "token_usage", "data": {...}}
```

### 类型同步
- 后端 Struct → 前端 TS 类型（`tygo` 或手动 `types/api.d.ts`）

---

## 5. 数据库规约 (PostgreSQL + pgvector)

| 规则     | 说明                                                 |
| -------- | ---------------------------------------------------- |
| **命名** | `snake_case`，主键 `id` (UUID)，外键 `xxx_id`        |
| **迁移** | `golang-migrate`，文件名 `YYYYMMDDHHMMSS_xxx.up.sql` |
| **向量** | 字段名 `embedding`，注释维度 `-- 1536 dim`           |

---

## 6. AI/Prompt 规约

| 规则              | 说明                                      |
| ----------------- | ----------------------------------------- |
| **禁止 Hardcode** | Prompt 存放 `/prompts/*.md`               |
| **模版占位符**    | 必须预留 `{{.Context}}`, `{{.UserQuery}}` |
| **降级策略**      | 处理 Token Overflow，自动截断历史         |
| **版本控制**      | Prompt 变更需 Code Review                 |

---

## 7. 核心技术选型

| 模块       | 技术                             | 理由               |
| ---------- | -------------------------------- | ------------------ |
| 前端框架   | React + Vite                     | 快速 HMR           |
| 状态管理   | Zustand                          | 极简 API + persist |
| 工作流编辑 | React Flow                       | 自定义节点体验佳   |
| 后端框架   | Gin                              | 高性能 + WebSocket |
| 数据库     | PostgreSQL + pgvector            | 向量与关系统一存储 |
| 搜索工具   | Tavily API                       | 事实核查           |
| LLM        | OpenAI/Anthropic/Google/DeepSeek | 纯云服务           |

---

## 8. 测试规约 (Testing Standards)

| 规则         | 说明                                                                                                      |
| :----------- | :-------------------------------------------------------------------------------------------------------- |
| **命令**     | 使用 `make test` 执行测试，严禁直接运行 `go test`                                                         |
| **覆盖率**   | 核心逻辑必须 100% 覆盖。基础设施 (infrastructure) 代码可通过 `make test` 自动排除统计。                   |
| **Mock**     | 核心业务测试禁止依赖真实 DB/LLM，必须使用 `MockProvider` 或 `Interface Stub`。                            |
| **开发模式** | **TDD (Test-Driven Development)**。先写测试（红），再写实现（绿），最后重构（蓝）。禁止先写代码后补测试。 |

```bash
# 执行测试并查看覆盖率报告 (自动过滤 infrastructure 噪音)
make test
```
```
