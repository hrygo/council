# 单元测试覆盖率提升计划 (Unit Test Coverage Improvement Plan)

## 1. 当前现状 (Current Status)

| 领域             | 当前覆盖率 (Statements) | 目标覆盖率 | 差距   |
| :--------------- | :---------------------- | :--------- | :----- |
| **后端 (Go)**    | 28.5%                   | 80%        | -51.5% |
| **前端 (React)** | 62.4%                   | 80%        | -17.6% |
| **总体平均**     | 35.2% (加权)            | 80%        | -44.8% |

### 1.1 后端薄弱环节 (Backend Weak Points)
- `internal/api/ws`: 0.0% (关键的 WebSocket 处理逻辑)
- `internal/infrastructure/*`: 0.0% (LLM, DB, Cache, Search 基建抽象)
- `internal/core/memory`: 28.9% (三层记忆协议实现)
- `internal/core/agent`: 0.0% (Agent 核心逻辑，目前无测试文件)

### 1.2 前端薄弱环节 (Frontend Weak Points)
- `src/hooks/useWebSocketRouter.ts`: 40.0% (核心消息分发逻辑)
- `src/components/chat/ChatInput.tsx`: 46.15% (消息发送逻辑)
- `src/stores/useLayoutStore.ts`: 42.85% (UI 布局状态)
- `src/stores/useWorkflowRunStore.ts`: 52.7% (工作流运行核心状态)

---

## 2. 提升目标与阶段 (Phases & Objectives)

### 第一阶段 (Sprint 1): 基建与核心 (Core & Infra)
- [x] **实现后端 Mock 基建**: 完善 `internal/infrastructure/mocks`，支持 LLM、DB、Search 的模拟。
- [x] **提升 Core 覆盖率**: 将 `internal/core` 下所有包的覆盖率提升至 50% 以上 (nodes 已达 56%)。
- [x] **前端 Store 补全**: 将 `useWorkflowRunStore` 和 `useSessionStore` 提升至 75% 以上。

### 第二阶段 (Sprint 2): API 与 交互 (API & UI Logic)
- [ ] **WebSocket 测试**: 实现 `internal/api/ws` 的单元/集成测试，覆盖率为 60% 以上。
- [ ] **Handler 强化**: 将 `internal/api/handler` 提升至 70% 以上。
- [ ] **前端 Hook 覆盖**: 将 `useWebSocketRouter` 提升至 80% 以上。

---

## 3. 具体任务清单 (Action Items)

### 3.1 后端 (Backend)
| 任务 ID  | 包/文件                         | 描述                                                    | 优先级 | 状态     |
| :------- | :------------------------------ | :------------------------------------------------------ | :----- | :------- |
| BE-QA-01 | `internal/infrastructure/mocks` | 为 `LLMProvider` 和 `DB` 增加通用 Mock 实现             | P0     | [x] Done |
| BE-QA-02 | `internal/core/agent`           | 新增 `agent_test.go`，覆盖 Agent 思考与响应流程         | P0     | [x] Done |
| BE-QA-03 | `internal/core/memory`          | 增加 LRU 缓存及向量检索部分的边界测试                   | P1     | [ ] Todo |
| BE-QA-04 | `internal/api/ws`               | 使用 `httptest.NewServer` 模拟 WebSocket 握手及消息循环 | P1     | [ ] Todo |

### 3.2 前端 (Frontend)
| 任务 ID  | 组件/Store                    | 描述                                                     | 优先级 | 状态     |
| :------- | :---------------------------- | :------------------------------------------------------- | :----- | :------- |
| FE-QA-01 | `useWebSocketRouter.test.ts`  | 模拟各种业务消息 (NodeUpdate, ChatMessage)，验证分发逻辑 | P0     | [x] Done |
| FE-QA-02 | `ChatInput.test.tsx`          | 测试发送消息、清空输入、加载状态、快捷键响应             | P1     | [ ] Todo |
| FE-QA-03 | `useWorkflowRunStore.test.ts` | 深度测试状态机分支、节点拓扑更新逻辑                     | P0     | [x] Done |

---

## 4. 推进计划 (Execution Roadmap)

2. **[Today/Tomorrow]** 完成核心 Store 的补齐 (FE-QA-03)。
3. **[Next]** 逐一推进 P0 级任务。
