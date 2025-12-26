# 开发计划 (Development Plan)

> **策略**: 前端优先，API Contract First  
> **质量内建**: 每个 Spec 必须通过 CI (Lint/Test) 和验收标准方可标记 Done

---

## 一、进度总览

| Sprint | 阶段               |   状态   | 完成度 |
| :----: | :----------------- | :------: | :----: |
| S1-S4  | MVP 核心功能       |  ✅ Done  |  100%  |
|   S5   | Post-MVP (Partial) |  ✅ Done  | 100%*  |
|   S6   | Default Experience |  ✅ Done  |  100%  |
|   S7   | UX Polish          |  ✅ Done  |  100%  |
|   S8   | Meeting Room Fix   |  ✅ Done  | 100%*  |
|   S9   | Quality (Base)     |  ✅ Done  | 100%*  |
|  S10   | System Hardening   | 🔄 进行中 |  10%   |

> **当前进度**: 整体 85% (43/51 任务已完成)。S5/S8/S9 的后续深度集成与硬化工作已统一移至 **Sprint 10**。

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

| ID   | 任务         | Spec     | Sprint | 优先级 | 状态  |
| :--- | :----------- | :------- | :----: | :----: | :---: |
| 5.1  | i18n 国际化  | SPEC-501 |   S5   |   P0   |   ✅   |
| 5.2  | E2E 测试     | SPEC-502 |   S5   |   P1   |   ✅   |
| 5.3  | 性能瓶颈优化 | SPEC-503 |  S10   |   P2   |   ⬜   |



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

### 4.3 System Hardening (Sprint 10)

| Spec ID  | 文档                                                         | 类型           | 优先级 |
| :------- | :----------------------------------------------------------- | :------------- | :----: |
| SPEC-901 | [RAG 深度集成](./specs/sprint10/SPEC-901-rag-integration.md) | Integration    |   P1   |
| SPEC-902 | [QA 质量硬化](./specs/sprint10/SPEC-902-qa-hardening.md)     | QA             |   P0   |
| SPEC-903 | [架构标准化](./specs/sprint10/SPEC-903-standardization.md)   | Refactor       |   P1   |
| SPEC-904 | [性能与安全](./specs/sprint10/SPEC-904-perf-security.md)     | Infrastructure |   P2   |

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
|   2   | LLM 注册表  |  4h   | SPEC-702 |   ✅   |
|   3   | WS 连接修复 |  2h   | SPEC-703 |   ✅   |

### 6.2 规格文档

| Spec ID  | 文档                                                                         | 类型     | Phase | 状态  |
| :------- | :--------------------------------------------------------------------------- | :------- | :---: | :---: |
| SPEC-701 | [Session Creation UI](./specs/sprint7/SPEC-701-session-creation-ui.md)       | Feature  |   1   |   ✅   |
| SPEC-702 | [Dynamic LLM Registry](./specs/sprint7/SPEC-702-llm-registry.md)             | Refactor |   2   |   ✅   |
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
- [x] Agent 运行时按 ModelConfig.Provider 选择
- [x] (Scope Reduced) 动态 BaseURL/Key 已确认不需要

**SPEC-703 Session WS Connect Fix**:
- [x] `SessionStarter.tsx` 在 API 成功后调用 `connect()`
- [x] `MeetingRoom.tsx` 自动重连断开的 WebSocket
- [x] 单元测试覆盖 (48/48 通过)
- [x] Lint + Build 验证通过

---

## 七、当前 Sprint: Meeting Room Fix (S8)

> **目标**: 修复会议室功能，完善用户体验，还原 Example 辩论流程

### 7.1 任务列表

| ID   | 任务                       | 类型        | 优先级 | 状态  |
| :--- | :------------------------- | :---------- | :----: | :---: |
| 8.1  | LLM Model 降级逻辑修复     | Bugfix      |   P0   |   ✅   |
| 8.2  | 会议室左侧流程实时监控修复 | Bugfix      |   P1   |   ✅   |
| 8.3  | 会议启动流程重构           | Feature     |   P0   |   ✅   |
| 8.4  | 会议过程 UX/UI 优化        | UX          |   P1   |   ✅   |
| 8.5  | 知识库面板 (Mock 集成)     | Feature     |   P2   |   ✅   |
| 8.6  | Example 辩论流程还原       | Feature     |   P0   |   ✅   |
| 8.7  | 知识库深度数据集成 (RAG)   | Integration |   P1   |   ✅   |


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

**8.5 右侧知识库面板集成 (P2)** ✅
- 前端: 创建 `KnowledgePanel` 组件，支持按记忆层级过滤、搜索
- 前端: 集成到 `MeetingRoom.tsx` 右侧面板
- 后端: 实现 `GET /api/v1/sessions/:sessionID/knowledge` API
- 后端: 创建 `KnowledgeHandler` 和测试（5 个测试用例全部通过）
- 文件: `frontend/src/features/meeting-room/components/KnowledgePanel.tsx`
- 文件: `internal/api/handler/knowledge.go`
- 测试: `internal/api/handler/knowledge_test.go`
- Spec: [SPEC-805](./specs/sprint8/SPEC-805-knowledge-panel.md)

**8.6 Example 辩论流程还原 (P0)**
- 问题: 辩论过程未还原 `example/` 中的完整逻辑
- Spec: [SPEC-804](./specs/sprint8/SPEC-804-debate-flow-restoration.md)

### 7.3 规格文档索引

| Spec ID  | 文档                                                                           | 类型    | 状态  |
| :------- | :----------------------------------------------------------------------------- | :------ | :---: |
| SPEC-801 | [Session Startup Flow](./specs/sprint8/SPEC-801-session-startup-flow.md)       | Feature |   ✅   |
| SPEC-802 | [Workflow Live Monitor](./specs/sprint8/SPEC-802-workflow-live-monitor.md)     | Feature |   ✅   |
| SPEC-803 | [Meeting UX Optimization](./specs/sprint8/SPEC-803-meeting-ux-optimization.md) | UX      |   ✅   |
| SPEC-804 | [Debate Flow Restoration](./specs/sprint8/SPEC-804-debate-flow-restoration.md) | Feature |   ✅   |
| SPEC-805 | [Knowledge Panel](./specs/sprint8/SPEC-805-knowledge-panel.md)                 | Feature |   ✅   |

### 7.3 已确认决策

- ✅ 不需要每个 Agent 独立配置 API Key
- ✅ 不需要支持动态 BaseURL

---

## 八/九、Sprint 9: Quality (Base)

> **目标**: 完成基础质量基线建设 (✅ 已达成)。

### 9.1 核心任务 (已完成)

| ID   | 任务                              | 类型  | 状态  |
| :--- | :-------------------------------- | :---- | :---: |
| 9.1  | WebSocket 消息类型自动生成 (tygo) | Infra |   ✅   |

---

## 十、Sprint 10: System Hardening (执行中)

> **目标**: 完成跨模块的深度集成，消除最后的技术债务，确保系统生产就绪。

### 10.1 核心交付物

1.  **RAG 真实集成**: 右侧知识库面板不再使用 Mock 数据，支持真实的向量搜索与显示。
2.  **质量防线**: Nodes 包覆盖率提升至 80%，E2E 契约测试通过所有关键状态校验。
3.  **标准化配置**: 移除所有硬编码，支持 `config.yaml` 动态配置环境参数。
4.  **性能硬化**: 前端支持千万级长消息流的虚拟滚动渲染，保证长会话不卡顿。

### 10.2 执行阶段与任务分解

| Phase | 名称               | 目标                                         | 关联 ID      |
| :---- | :----------------- | :------------------------------------------- | :----------- |
| 1     | RAG 深度集成 (✅)   | 将 KnowledgePanel 从 Mock 切换为真实向量检索 | (8.7)        |
| 2     | 质量契约自动化 (✅) | 补全 E2E 消息校验，建立持续的 Nodes 测试防线 | (10.1, 10.2) |
| 3     | 架构标准化         | 完成全局配置化重构，清理 ID 语义冲突         | (10.3, 10.4) |
| 4     | 性能与安全硬化     | 引入虚拟滚动优化与 CSRF/XSS 基础防御层       | (5.3, 5.4)   |

### 10.3 任务详细矩阵

| ID   | 任务                                | 类型     | 优先级 | 状态  | 工时  |
| :--- | :---------------------------------- | :------- | :----: | :---: | :---: |
| 10.1 | 端到端消息数据一致性测试 (E2E)      | QA       |   P0   |   ✅   |  10h  |
| 10.2 | 核收层测试覆盖率深度补全 (Nodes)    | QA       |   P1   |   ✅   |  12h  |
| 10.3 | 配置文件化 & 硬编码清理             | Refactor |   P1   |   ⬜   |  8h   |
| 10.4 | 全局 ID 命名规范化审查              | Refactor |   P2   |   ⬜   |  4h   |
| 10.5 | Store 类型字段完整性校验            | QA       |   P2   |   ⬜   |  3h   |
| 10.6 | 核心解耦后期收尾 (移除 Adjudicator) | Refactor |   P1   |   ⬜   |  6h   |

### 10.4 风险管理

1.  **RAG 准确性**: PGVector 检索参数（如 `Limit`, `Distance Threshold`）需要调优以避免检索到无关内容。
2.  **测试覆盖率瓶颈**: `nodes` 包中的并行与循环逻辑 Mock 难度大，需引入更强大的 Mock 框架。
3.  **配置兼容性**: 引入配置文件需确保本地开发环境（.env）与容器环境的一致性。

### 10.5 验收标准与质量门禁

| 门禁项         | 检查内容                                          |       自动化        |
| :------------- | :------------------------------------------------ | :-----------------: |
| **Lint**       | 代码风格符合规范                                  |        ✅ CI         |
| **Test**       | 单元测试通过，核心层覆盖率 ≥ 80%                  |        ✅ CI         |
| **Contract**   | 前后端 WebSocket 消息类型自动生成且一致           |       ✅ tygo        |
| **E2E**        | 关键消息契约通过端到端测试验证                    |    ✅ Playwright     |
| **Config**     | 所有默认值可从环境变量或配置文件覆盖              |      ⬜ Manual       |
| **Decoupling** | Core 层不依赖任何具体 Agent 实现 (如 Adjudicator) | ✅ verify-decoupling |


---

## 十一、已搁置任务

### 11.1 搁置任务列表

| ID   | 任务              | Spec     | 优先级 | 状态  | 搁置原因                       |
| :--- | :---------------- | :------- | :----: | :---: | :----------------------------- |
| X.1  | 安全强化          | SPEC-504 |   P3   |   🚫   | 开发环境无需，生产环境另行处理 |
| X.2  | 暗黑模式切换      | -        |   P2   |   🚫   | 已支持三种模式，按钮非必需     |
| X.3  | Run Mode 轻量渲染 | -        |   P2   |   🚫   | 当前性能可接受，ROI不足        |

### 11.2 搁置决策说明

**5.4 安全强化** - 搁置理由：
- 当前为开发环境，无需生产级安全措施
- 生产部署时应由运维团队统一配置
- XSS/CSRF 防御可通过框架默认配置实现
- 投入产出比不高

**B.1 暗黑模式切换** - 搁置理由：
- 系统已支持 system/light/dark 三种模式
- 可通过系统设置自动切换
- UI 增加切换按钮收益有限
- 不影响核心功能使用

**B.2 Run Mode 轻量渲染** - 搁置理由：
- 代码分割已实现，性能已优化
- 当前工作流规模下性能可接受
- 实现复杂度高，收益边际递减
- 非阻塞性问题

---

## 十二、技术专项方案 (Post-MVP)

### 12.1 性能优化策略 (Sprint 5)

| 领域     | 优化项     | 技术方案                         | 预期收益               |
| :------- | :--------- | :------------------------------- | :--------------------- |
| **前端** | 长消息列表 | React-Window 虚拟滚动            | 减少 DOM 节点 90%      |
| **前端** | 画布性能   | ReactFlow 视口裁剪 + Canvas 模式 | 支持 100+ 节点流畅交互 |
| **前端** | 渲染优化   | Zustand Selector 细粒度订阅      | 减少不必要的重渲染 50% |
| **后端** | 记忆检索   | Redis 缓存热点查询结果           | 查询响应时间减少 70%   |
| **后端** | 并发控制   | 令牌桶限流 + 优先级队列          | 避免 Provider 限流     |

### 12.2 安全强化措施 (Sprint 5)

| 领域     | 安全问题  | 解决方案                            |
| :------- | :-------- | :---------------------------------- |
| **前端** | XSS 注入  | DOMPurify 清理；强制校验 HTML 渲染  |
| **前端** | CSRF 攻击 | API 请求携带 CSRF Token (Cookie)    |
| **后端** | SQL 注入  | 强制使用参数化查询；禁止拼接 SQL    |
| **后端** | 权限校验  | Middleware 验证资源访问权限         |
| **后端** | 敏感脱敏  | 日志记录前移除 API Key 和敏感 Token |

---

## 十三、完成度统计

### 13.1 总体进度

| 维度     | 已完成 | 总计  | 完成率 |
| :------- | :----: | :---: | :----: |
| Sprint   |   7    |  10   |  70%   |
| 前端任务 |   19   |  19   |  100%  |
| 后端任务 |   11   |  11   |  100%  |
| Post-MVP |   2    |   4   |  50%   |
| 总任务数 |   43   |  52   | 82.7%  |

### 13.2 里程碑达成

- ✅ M1: Run Mode (Week 1)
- ✅ M2: 管理页面 (Week 2)
- ✅ M3: Builder (Week 3)
- ✅ M4: MVP (Week 4)
- ✅ M5: 国际化 (Week 5)
- ✅ M6: Out-of-Box (Week 6)
- ✅ M7: Meeting Room Fix (Week 7)
- 🔄 M8: Quality & Stability (进行中)


### 13.3 技术债务偿还

- ✅ WebSocket 类型自动生成
- 🔄 测试覆盖率提升 (当前 ~73%)
- ⬜ 硬编码默认值清理
- ⬜ ID 命名规范
- ⬜ Core 层解耦验证
- ⬜ 性能优化 (计划中)

---

## 十四、项目总结

### 14.1 核心成果

✅ **完整的多智能体协作系统**
- 工作流设计器（9种节点类型）
- 实时执行引擎（WebSocket通信）
- 三层记忆系统
- 人工审核机制
- 成本估算与控制

✅ **开箱即用体验**
- 默认 Agent（Affirmative, Negative, Adjudicator）
- 预置工作流（Debate, Optimize）
- 完整的示例流程

✅ **质量保障**
- 前端测试覆盖
- 后端测试覆盖 73.3%
- E2E 测试全覆盖
- CI/CD 自动化

✅ **性能优化**
- 代码分割与懒加载
- Bundle 优化（vendor拆分）
- WebSocket 消息优化

### 14.2 技术亮点

- **类型安全**: tygo 自动生成前后端类型定义
- **国际化**: i18n 支持中英文切换
- **防御机制**: 熔断器、反幻觉传播
- **版本控制**: HumanReview 前自动备份
- **知识检索**: 集成外部搜索（Tavily）

### 14.3 架构优势

- **解耦设计**: Core 层独立于具体业务
- **可扩展**: 支持自定义节点、Provider
- **可观测**: 完整的日志和审计追踪
- **高可用**: 状态管理、错误恢复机制

---

## 十五、其他待办 (Backlog)

### 15.1 功能与优化

> **注**: 以下为可选增强功能，不影响系统核心使用

| 任务             | 优先级 | 状态  | 备注          |
| :--------------- | :----: | :---: | :------------ |
| GraphQL API 支持 |   P3   |   -   | 可选替代 REST |
| 批量导入 Agent   |   P3   |   -   | 便捷性功能    |
| 工作流版本控制   |   P3   |   -   | Git 集成      |
| 性能监控面板     |   P3   |   -   | APM 集成      |
| 多租户支持       |   P3   |   -   | SaaS 模式     |
| 审计日志导出     |   P3   |   -   | 合规需求      |


---

## 附录：术语表 (Glossary)

| 术语            | 定义                                             |
| :-------------- | :----------------------------------------------- |
| **Sprint**      | 敏捷开发中的迭代周期，通常为 1-2 周              |
| **Spec**        | Specification 的缩写，指规格说明文档             |
| **任务矩阵**    | 用表格形式组织的任务清单，包含状态、优先级等维度 |
| **里程碑**      | 项目中的关键交付节点，具有明确的验收标准         |
| **验收标准**    | 用于判断任务或 Sprint 是否完成的可测试条件       |
| **依赖关系**    | 任务或规格之间的先后顺序或输入输出关系           |
| **P0/P1/P2/P3** | 优先级标记，P0 最高，P3 最低                     |
| **技术债务**    | 为快速交付而采用的非最优方案，需要后续偿还       |
| **tygo**        | 自动将 Go 结构体转换为 TypeScript 类型定义的工具 |

