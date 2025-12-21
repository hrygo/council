# 开发计划 (Development Plan)

> **策略**: 前端优先，API Contract First  
> **质量内建**: 每个 Spec 必须通过 CI (Lint/Test) 和验收标准方可标记 Done

---

## 一、进度总览

| Sprint | 阶段               |   状态   | 完成度 |
| :----: | :----------------- | :------: | :----: |
| S1-S4  | MVP 核心功能       |  ✅ Done  |  100%  |
|   S5   | Post-MVP 优化      | 🔄 进行中 |  50%   |
|   S6   | Default Experience |  ✅ Done  |  100%  |
|   S7   | UX Polish          |  ✅ Done  |  100%  |
|   S8   | Meeting Room Fix   | 🔄 进行中 |   0%   |

---

## 二、里程碑

| 时间   | 里程碑         | 验收标准                     | 状态  |
| :----- | :------------- | :--------------------------- | :---: |
| Week 1 | M1: Run Mode   | 运行简单工作流，消息正确显示 |   ✅   |
| Week 2 | M2: 管理页面   | 群组和 Agent CRUD 完成       |   ✅   |
| Week 3 | M3: Builder    | 所有节点类型，模版可用       |   ✅   |
| Week 4 | M4: MVP        | 人类裁决、成本预估、文档引用 |   ✅   |
| Week 5 | M5: 国际化     | i18n + E2E 测试              |   ✅   |
| Week 6 | M6: Out-of-Box | The Council 开箱即用         |   ✅   |

---

## 三、任务跟踪矩阵

### 3.1 前端任务

| ID   | 任务                          | Spec         | Sprint | 优先级 | 状态  |
| :--- | :---------------------------- | :----------- | :----: | :----: | :---: |
| 1.1  | 重写 `useSessionStore.ts`     | SPEC-001     |   S1   |   P0   |   ✅   |
| 1.2  | 实现 `useWorkflowRunStore.ts` | SPEC-002     |   S1   |   P0   |   ✅   |
| 1.3  | ChatPanel 分组消息            | SPEC-003     |   S1   |   P1   |   ✅   |
| 1.4  | ChatPanel 并行消息            | SPEC-004     |   S1   |   P1   |   ✅   |
| 1.5  | WebSocket 优化                | SPEC-005     |   S1   |   P1   |   ✅   |
| 2.1  | 群组管理页面                  | SPEC-101/102 |   S2   |   P0   |   ✅   |
| 2.2  | Agent 管理页面                | SPEC-103/104 |   S2   |   P0   |   ✅   |
| 2.3  | Agent 模型配置                | SPEC-105     |   S2   |   P1   |   ✅   |
| 3.1  | 节点属性面板                  | SPEC-201     |   S3   |   P0   |   ✅   |
| 3.2  | Vote/Loop 节点 UI             | SPEC-202     |   S3   |   P1   |   ✅   |
| 3.3  | FactCheck/HumanReview UI      | SPEC-203     |   S3   |   P1   |   ✅   |
| 3.4  | 模版库侧边栏                  | SPEC-204     |   S3   |   P2   |   ✅   |
| 3.5  | 保存为模版                    | SPEC-205     |   S3   |   P2   |   ✅   |
| 3.6  | 向导模式                      | SPEC-206     |   S3   |   P0   |   ✅   |
| 4.1  | HumanReviewModal              | SPEC-301     |   S4   |   P0   |   ✅   |
| 4.2  | CostEstimator                 | SPEC-302     |   S4   |   P1   |   ✅   |
| 4.3  | 文档引用跳转                  | SPEC-303     |   S4   |   P1   |   ✅   |
| 4.4  | 快捷键支持                    | SPEC-304     |   S4   |   P2   |   ✅   |
| 4.5  | KaTeX 公式渲染                | SPEC-305     |   S4   |   P2   |   ✅   |

### 3.2 后端任务

| ID   | 任务                 | Spec     | Sprint | 优先级 | 状态  |
| :--- | :------------------- | :------- | :----: | :----: | :---: |
| B.1  | SequenceProcessor    | SPEC-401 |  S1-2  |   P1   |   ✅   |
| B.2  | VoteProcessor        | SPEC-402 |  S1-2  |   P1   |   ✅   |
| B.3  | LoopProcessor        | SPEC-403 |  S1-2  |   P2   |   ✅   |
| B.4  | FactCheckProcessor   | SPEC-404 |  S3-4  |   P1   |   ✅   |
| B.5  | HumanReviewProcessor | SPEC-405 |  S3-4  |   P0   |   ✅   |
| B.6  | Templates API        | SPEC-406 |  S3-4  |   P1   |   ✅   |
| B.7  | Cost Estimation API  | SPEC-407 |  S3-4  |   P1   |   ✅   |
| B.8  | 三层记忆协议         | SPEC-408 |  S3-4  |   P0   |   ✅   |
| B.9  | 逻辑熔断             | SPEC-409 |  S3-4  |   P0   |   ✅   |
| B.10 | 防幻觉传播           | SPEC-410 |  S3-4  |   P1   |   ✅   |
| B.11 | 联网搜索集成         | SPEC-411 |  S3-4  |   P1   |   ✅   |

### 3.3 Post-MVP 任务

| ID   | 任务        | Spec     | Sprint | 优先级 | 状态  |
| :--- | :---------- | :------- | :----: | :----: | :---: |
| 5.1  | i18n 国际化 | SPEC-501 |   S5   |   P0   |   ✅   |
| 5.2  | E2E 测试    | SPEC-502 |   S5   |   P1   |   ✅   |
| 5.3  | 性能优化    | SPEC-503 |   S5   |   P2   |   ⬜   |
| 5.4  | 安全强化    | SPEC-504 |   S5   |   P3   |   ⬜   |

---

## 四、规格文档索引

### 4.1 MVP 核心 (Sprint 1-4)

| Sprint  | 目录                                        | 说明             |
| :------ | :------------------------------------------ | :--------------- |
| S1      | [specs/sprint1/](./specs/sprint1/README.md) | 运行时状态重构   |
| S2      | [specs/sprint2/](./specs/sprint2/README.md) | 管理页面         |
| S3      | [specs/sprint3/](./specs/sprint3/README.md) | Builder 增强     |
| S4      | [specs/sprint4/](./specs/sprint4/README.md) | 高级功能         |
| Backend | [specs/backend/](./specs/backend/README.md) | 后端处理器和 API |

### 4.2 Post-MVP (Sprint 5)

| Spec ID  | 文档                                                             | 类型     | 优先级 |
| :------- | :--------------------------------------------------------------- | :------- | :----: |
| SPEC-501 | [i18n 国际化](./specs/sprint5/SPEC-501-i18n.md)                  | Feature  |   P0   |
| SPEC-502 | [E2E 测试](./specs/sprint5/SPEC-502-e2e-testing.md)              | QA       |   P1   |
| SPEC-503 | [性能优化](./specs/sprint5/SPEC-503-performance-optimization.md) | Refactor |   P2   |
| SPEC-504 | [安全强化](./specs/sprint5/SPEC-504-security-hardening.md)       | Security |   P3   |

---

## 五、当前 Sprint: Default Experience (S6)

> **目标**: 将 `example/` 迁移为系统默认开箱即用体验，100% 覆盖 `skill.md`

### 5.1 执行阶段

| Phase | 名称     |  工时   | Specs                   |
| :---: | :------- | :-----: | :---------------------- |
|   1   | 基础设施 |   10h   | SPEC-609, SPEC-608      |
|   2   | 数据注入 |   12h   | SPEC-607, 601, 602, 603 |
|   3   | 功能增强 |   7h    | SPEC-605, SPEC-606      |
|   4   | 集成验证 |   4h    | E2E 测试                |
|       | **总计** | **33h** |                         |

### 5.2 依赖关系

```
SPEC-609 ─► SPEC-608 ─► SPEC-601 ─┬─► SPEC-602 ─► SPEC-606
                                  │
SPEC-609 ─► SPEC-607 ─────────────┴─► SPEC-603 ─► Integration
                                              │
SPEC-605 ─────────────────────────────────────┘
```

### 5.3 规格文档

| Spec ID  | 文档                                                                       | 类型           | Phase |
| :------- | :------------------------------------------------------------------------- | :------------- | :---: |
| SPEC-609 | [Architecture Fixes](./specs/sprint6/SPEC-609-architecture-fixes.md)       | Bug Fix        |   1   |
| SPEC-608 | [Prompt Embed](./specs/sprint6/SPEC-608-prompt-embed.md)                   | Infrastructure |   1   |
| SPEC-607 | [Memory Retrieval Node](./specs/sprint6/SPEC-607-memory-retrieval-node.md) | Go Node        |   2   |
| SPEC-601 | [Default Agents](./specs/sprint6/SPEC-601-default-agents.md)               | Go Seeder      |   2   |
| SPEC-602 | [Default Group](./specs/sprint6/SPEC-602-default-group.md)                 | Go Seeder      |   2   |
| SPEC-603 | [Default Workflows](./specs/sprint6/SPEC-603-default-workflows.md)         | Go Seeder      |   2   |
| SPEC-605 | [Versioning Middleware](./specs/sprint6/SPEC-605-versioning-middleware.md) | Middleware     |   3   |
| SPEC-606 | [Documentation](./specs/sprint6/SPEC-606-documentation.md)                 | Docs           |   3   |

### 5.4 验收标准

**功能验收**:
- [x] 3 个系统 Agent 存在 (`seeder.go`: Affirmative, Negative, Adjudicator)
- [x] "The Council" 群组存在 (`seeder.go`: SeedGroups)
- [x] Debate + Optimize 流程存在 (`seeder.go`: debateWorkflowGraph, optimizeWorkflowGraph)
- [x] `memory_retrieval` 节点可用 (`internal/core/workflow/nodes/memory_retrieval.go`)
- [x] HumanReview 前自动备份 (`internal/core/middleware/versioning.go`)
- [x] 完整 Optimize 循环可运行

**skill.md 覆盖**:
- [x] Step 1: Memory Retrieval (`nodes/memory_retrieval.go`)
- [x] Step 2: Parallel + Agent (`nodes/parallel.go`, `nodes/agent.go`)
- [x] Step 3: Scoring Matrix (`nodes/vote.go`)
- [x] Step 4: Versioning (`middleware/versioning.go`)
- [x] Step 5: HumanReview (`nodes/human_review.go`)
- [x] Step 6: Loop (`nodes/loop.go`)

**解耦验证**:
- [x] `make verify-decoupling` 通过
- [x] 删除 `example/` 后系统正常

---

## 六、当前 Sprint: UX Polish (S7)

> **目标**: 优化用户体验，填补 UI/UX 缺口，使核心功能（会话创建、管理）闭环。

### 6.1 执行阶段

| Phase | 名称        | 工时  | Specs    | 状态  |
| :---: | :---------- | :---: | :------- | :---: |
|   1   | UX 闭环     |  8h   | SPEC-701 |   ✅   |
|   2   | LLM 注册表  |  4h   | SPEC-702 |   ⬜   |
|   3   | WS 连接修复 |  2h   | SPEC-703 |   ✅   |

### 6.2 规格文档

| Spec ID  | 文档                                                                         | 类型     | Phase | 状态  |
| :------- | :--------------------------------------------------------------------------- | :------- | :---: | :---: |
| SPEC-701 | [Session Creation UI](./specs/sprint7/SPEC-701-session-creation-ui.md)       | Feature  |   1   |   ✅   |
| SPEC-702 | [Dynamic LLM Registry](./specs/sprint7/SPEC-702-llm-registry.md)             | Refactor |   2   |   ⬜   |
| SPEC-703 | [Session WS Connect Fix](./specs/sprint7/SPEC-703-session-ws-connect-fix.md) | Bugfix   |   3   |   ✅   |

### 6.3 验收标准

**SPEC-701 Session Creation UI**:
- [x] `/chat` 页面在无会话时显示 "Start Session" 界面 (`SessionStarter.tsx`)
- [x] 支持选择 "Council Debate" 模板并启动
- [x] 用户可输入讨论主题
- [x] 点击 Launch 成功启动后端流程
- [x] 聊天界面立即反映新会话 (WebSocket 连接)

**SPEC-702 Dynamic LLM Registry**:
- [x] `Registry` 结构替代单一 Provider (`router.go` 已实现)
- [x] 支持多 Provider 动态切换 (gemini, deepseek, openai 等)
- [ ] Agent 运行时按 ModelConfig.Provider 选择
- [ ] 完整集成测试覆盖

**SPEC-703 Session WS Connect Fix**:
- [x] `SessionStarter.tsx` 在 API 成功后调用 `connect()`
- [x] `MeetingRoom.tsx` 自动重连断开的 WebSocket
- [x] 单元测试覆盖 (48/48 通过)
- [x] Lint + Build 验证通过

---

## 七、当前 Sprint: Meeting Room Fix (S8)

> **目标**: 修复会议室功能，完善用户体验，还原 Example 辩论流程

### 7.1 任务列表

| ID   | 任务                       | 类型    | 优先级 | 状态  |
| :--- | :------------------------- | :------ | :----: | :---: |
| 8.1  | LLM Model 降级逻辑修复     | Bugfix  |   P0   |   ✅   |
| 8.2  | 会议室左侧流程实时监控修复 | Bugfix  |   P1   |   ⬜   |
| 8.3  | 会议启动流程重构           | Feature |   P0   |   ⬜   |
| 8.4  | 会议过程 UX/UI 优化        | UX      |   P1   |   ⬜   |
| 8.5  | 右侧知识库面板集成         | Feature |   P2   |   ⬜   |
| 8.6  | Example 辩论流程还原       | Feature |   P0   |   ⬜   |

### 7.2 任务详情

**8.1 LLM Model 降级逻辑修复 (P0)** ✅
- 问题: `agent.go` 第 66-68 行硬编码 `gpt-4` 作为默认模型
- 修复: 改为 `a.Registry.GetDefaultModel()`
- 文件: `internal/core/workflow/nodes/agent.go`

**8.2 会议室左侧流程实时监控修复 (P1)**
- 问题: ReactFlow 画布为空白，节点状态未同步
- 预期: 显示工作流图并实时高亮当前执行节点
- Spec: [SPEC-802](./specs/sprint8/SPEC-802-workflow-live-monitor.md)

**8.3 会议启动流程重构 (P0)**
- 问题: 选择模板后会议自动运行，无用户参与
- 预期: 用户可上传文件、输入目标、确认后启动
- Spec: [SPEC-801](./specs/sprint8/SPEC-801-session-startup-flow.md)

**8.4 会议过程 UX/UI 优化 (P1)**
- 改进消息展示、状态指示、Agent 头像等
- Spec: [SPEC-803](./specs/sprint8/SPEC-803-meeting-ux-optimization.md)

**8.5 右侧知识库面板集成 (P2)**
- 问题: 右侧知识面板未被使用
- 预期: 显示会议相关知识、上下文、引用文档

**8.6 Example 辩论流程还原 (P0)**
- 问题: 辩论过程未还原 `example/` 中的完整逻辑
- Spec: [SPEC-804](./specs/sprint8/SPEC-804-debate-flow-restoration.md)

### 7.3 规格文档索引

| Spec ID  | 文档                                                                           | 类型    | 状态  |
| :------- | :----------------------------------------------------------------------------- | :------ | :---: |
| SPEC-801 | [Session Startup Flow](./specs/sprint8/SPEC-801-session-startup-flow.md)       | Feature |   ⬜   |
| SPEC-802 | [Workflow Live Monitor](./specs/sprint8/SPEC-802-workflow-live-monitor.md)     | Feature |   ⬜   |
| SPEC-803 | [Meeting UX Optimization](./specs/sprint8/SPEC-803-meeting-ux-optimization.md) | UX      |   ⬜   |
| SPEC-804 | [Debate Flow Restoration](./specs/sprint8/SPEC-804-debate-flow-restoration.md) | Feature |   ⬜   |

### 7.3 已确认决策

- ✅ 不需要每个 Agent 独立配置 API Key
- ✅ 不需要支持动态 BaseURL

---

## 八、技术债务

| 任务                 | 优先级 | 状态  |
| :------------------- | :----: | :---: |
| 测试覆盖率提升至 80% |   P1   |   ⬜   |
| 暗黑模式切换         |   P2   |   ⬜   |
| Run Mode 轻量渲染    |   P2   |   ⬜   |
