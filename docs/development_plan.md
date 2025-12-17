# 开发计划 (Development Plan)

> **基于**: [审计报告](docs/reports/audit_report.md)  
> **策略**: 前端优先，API Contract First  
> **质量内建**: 每个 Spec 必须通过 CI (Lint/Test) 和验收标准 (Acceptance Criteria) 方可标记 Done  
> **预估周期**: 4 周

---


## 开发进度跟踪矩阵 (Progress Matrix)

| ID                  | 任务名称 (Task)               | Spec ID       | Sprint | 优先级 | 状态     | 负责人    |
| ------------------- | ----------------------------- | ------------- | ------ | ------ | -------- | --------- |
| **前端 (Frontend)** |
| 1.1                 | 重写 `useSessionStore.ts`     | SPEC-001      | S1     | P0     | [x] Done | @Frontend |
| 1.2                 | 实现 `useWorkflowRunStore.ts` | SPEC-002      | S1     | P0     | [x] Done | @Frontend |
| 1.3                 | ChatPanel 分组消息            | SPEC-003      | S1     | P1     | [x] Done | @Frontend |
| 1.4                 | ChatPanel 并行消息            | SPEC-004      | S1     | P1     | [x] Done | @Frontend |
| 1.5                 | WebSocket 优化                | SPEC-005      | S1     | P1     | [x] Done | @Frontend |
| 2.1                 | 群组管理页面 + 列表           | SPEC-101, 102 | S2     | P0     | [x] Done | @Frontend |
| 2.2                 | Agent 管理页面 + 列表         | SPEC-103, 104 | S2     | P0     | [x] Done | @Frontend |
| 2.3                 | Agent 模型配置                | SPEC-105      | S2     | P1     | [x] Done | @Frontend |
| 3.1                 | 节点属性面板 (PropertyPanel)  | SPEC-201      | S3     | P0     | [x] Done | @Frontend |
| 3.2                 | Vote/Loop 节点 UI             | SPEC-202      | S3     | P1     | [x] Done | @Frontend |
| 3.3                 | FactCheck/HumanReview UI      | SPEC-203      | S3     | P1     | [x] Done | @Frontend |
| 3.4                 | 模版库侧边栏                  | SPEC-204      | S3     | P2     | [x] Done | @Frontend |
| 3.5                 | 保存为模版                    | SPEC-205      | S3     | P2     | [x] Done | @Frontend |
| 3.6                 | 向导模式 (`WizardMode`)       | SPEC-206      | S3     | P0     | [x] Done | @Frontend |
| 4.1                 | HumanReviewModal              | SPEC-301      | S4     | P0     | [x] Done | @Frontend |
| 4.2                 | CostEstimator                 | SPEC-302      | S4     | P1     | [ ] Todo | @Frontend |
| 4.3                 | 文档引用跳转                  | SPEC-303      | S4     | P1     | [ ] Todo | @Frontend |
| 4.4                 | 快捷键支持                    | SPEC-304      | S4     | P2     | [ ] Todo | @Frontend |
| 4.5                 | KaTeX 公式渲染                | SPEC-305      | S4     | P2     | [ ] Todo | @Frontend |
| **后端 (Backend)**  |
| B.1                 | SequenceProcessor             | SPEC-401      | S1-2   | P1     | [x] Done | @Backend  |
| B.2                 | VoteProcessor                 | SPEC-402      | S1-2   | P1     | [x] Done | @Backend  |
| B.3                 | LoopProcessor                 | SPEC-403      | S1-2   | P2     | [x] Done | @Backend  |
| B.4                 | FactCheckProcessor            | SPEC-404      | S3-4   | P1     | [x] Done | @Backend  |
| B.5                 | HumanReviewProcessor          | SPEC-405      | S3-4   | P0     | [x] Done | @Backend  |
| B.6                 | Templates API                 | SPEC-406      | S3-4   | P1     | [x] Done | @Backend  |
| B.7                 | Cost Estimation API           | SPEC-407      | S3-4   | P1     | [ ] Todo | @Backend  |
| B.8                 | 三层记忆协议                  | SPEC-408      | S3-4   | P0     | [ ] Todo | @Backend  |
| B.9                 | 逻辑熔断 (Circuit Breaker)    | SPEC-409      | S3-4   | P0     | [x] Done | @Backend  |
| B.10                | 防幻觉传播                    | SPEC-410      | S3-4   | P1     | [x] Done | @Backend  |
| B.11                | 联网搜索集成                  | SPEC-411      | S3-4   | P1     | [ ] Todo | @Backend  |

---

## Sprint 1: 运行时状态重构 (Week 1)

### 目标
修复 **前端 Run Mode 核心阻断项**，使会议室可以正常运行。

### 前端任务

| 任务                                                  | 优先级 | 预估 |
| ----------------------------------------------------- | ------ | ---- |
| **重写 `useSessionStore.ts`** - 实现完整 Session 状态 | P0     | 4h   |
| 实现 `useWorkflowRunStore.ts` - 工作流运行时状态      | P0     | 3h   |
| `ChatPanel` 支持按节点分组消息                        | P1     | 2h   |
| `ChatPanel` 支持并行消息并排显示                      | P1     | 2h   |
| WebSocket 消息处理优化                                | P1     | 2h   |

### API 需求 (已有，需验证)

- `POST /api/v1/workflows/execute` - 启动工作流
- `POST /api/v1/sessions/:id/control` - 暂停/恢复/停止
- `WebSocket /ws` - 实时事件流

### 交付物
- [ ] 会议室可以运行一个简单的工作流
- [ ] 消息按节点分组显示

### 📚 详细规格文档

> **[Sprint 1 Specs 目录](./specs/sprint1/README.md)**

| Spec ID  | 文档                                                                  | 类型        | 优先级 |
| -------- | --------------------------------------------------------------------- | ----------- | ------ |
| SPEC-001 | [useSessionStore 重写](./specs/sprint1/SPEC-001-session-store.md)     | Store       | P0     |
| SPEC-002 | [useWorkflowRunStore](./specs/sprint1/SPEC-002-workflow-run-store.md) | Store       | P0     |
| SPEC-003 | [ChatPanel 分组消息](./specs/sprint1/SPEC-003-chat-panel-grouped.md)  | Component   | P1     |
| SPEC-004 | [并行消息 UI](./specs/sprint1/SPEC-004-parallel-message-ui.md)        | Component   | P1     |
| SPEC-005 | [WebSocket 优化](./specs/sprint1/SPEC-005-websocket-optimization.md)  | Integration | P1     |

---

## Sprint 2: 管理页面 (Week 2)

### 目标
完成 **群组管理** 和 **Agent 管理** 的 CRUD 页面。

### 前端任务

| 任务                                | 优先级 | 预估 |
| ----------------------------------- | ------ | ---- |
| 创建 `/groups` 路由和页面           | P0     | 3h   |
| 群组列表 + 创建/编辑/删除           | P0     | 4h   |
| 创建 `/agents` 路由和页面           | P0     | 3h   |
| Agent 列表 + 创建/编辑/删除         | P0     | 4h   |
| Agent 模型配置面板 (Model Selector) | P1     | 3h   |

### API 需求 (已有)

- `GET/POST /api/v1/groups`
- `GET/PUT/DELETE /api/v1/groups/:id`
- `GET/POST /api/v1/agents`
- `GET/PUT/DELETE /api/v1/agents/:id`

### 交付物
- [x] 可创建/管理群组
- [x] 可创建/管理 Agent，并配置模型

### 📚 详细规格文档

> **[Sprint 2 Specs 目录](./specs/sprint2/README.md)**

| Spec ID  | 文档                                                        | 类型      | 优先级 |
| -------- | ----------------------------------------------------------- | --------- | ------ |
| SPEC-101 | [Groups 页面](./specs/sprint2/SPEC-101-groups-page.md)      | Page      | P0     |
| SPEC-102 | [GroupList 组件](./specs/sprint2/SPEC-102-group-list.md)    | Component | P0     |
| SPEC-103 | [Agents 页面](./specs/sprint2/SPEC-103-agents-page.md)      | Page      | P0     |
| SPEC-104 | [AgentList 组件](./specs/sprint2/SPEC-104-agent-list.md)    | Component | P0     |
| SPEC-105 | [ModelSelector](./specs/sprint2/SPEC-105-model-selector.md) | Component | P1     |

---

## Sprint 3: Builder 增强 (Week 3)

### 目标
完善 **工作流编辑器**，支持更多节点类型和模版库。

### 前端任务

| 任务                                   | 优先级 | 预估 |
| -------------------------------------- | ------ | ---- |
| 节点属性面板 (PropertyPanel)           | P0     | 4h   |
| 新增节点类型 UI: Vote/Loop             | P1     | 3h   |
| 新增节点类型 UI: FactCheck/HumanReview | P1     | 3h   |
| 模版库侧边栏                           | P2     | 3h   |
| 保存为模版功能                         | P2     | 2h   |

### API 需求 (需新增)

> 详见 [API 设计文档](./api/templates.md)

| 端点                    | 方法   | 说明         |
| ----------------------- | ------ | ------------ |
| `/api/v1/templates`     | GET    | 获取模版列表 |
| `/api/v1/templates`     | POST   | 创建模版     |
| `/api/v1/templates/:id` | GET    | 获取模版详情 |
| `/api/v1/templates/:id` | DELETE | 删除模版     |

### 交付物
- [ ] 节点可配置属性
- [ ] 支持所有 PRD 定义的节点类型 UI
- [ ] 模版保存和加载

### 📚 详细规格文档

> **[Sprint 3 Specs 目录](./specs/sprint3/README.md)**

| Spec ID  | 文档                                                                             | 类型      | 优先级 |
| -------- | -------------------------------------------------------------------------------- | --------- | ------ |
| SPEC-201 | [PropertyPanel](./specs/sprint3/SPEC-201-property-panel.md)                      | Component | P0     |
| SPEC-202 | [Vote/Loop 节点](./specs/sprint3/SPEC-202-vote-loop-nodes.md)                    | Component | P1     |
| SPEC-203 | [FactCheck/HumanReview](./specs/sprint3/SPEC-203-factcheck-humanreview-nodes.md) | Component | P1     |
| SPEC-204 | [模版库侧边栏](./specs/sprint3/SPEC-204-template-sidebar.md)                     | Component | P2     |
| SPEC-205 | [保存为模版](./specs/sprint3/SPEC-205-save-template.md)                          | Feature   | P2     |
| SPEC-206 | [向导模式](./specs/sprint3/SPEC-206-wizard-mode.md)                              | Feature   | P0     |

---

## Sprint 4: 高级功能 (Week 4)

### 目标
实现 **人类裁决**、**成本预估**、**文档引用** 等高级功能。

### 前端任务

| 任务                                  | 优先级 | 预估 |
| ------------------------------------- | ------ | ---- |
| `HumanReviewModal` 组件               | P0     | 3h   |
| 成本预估面板 (`CostEstimator`)        | P1     | 3h   |
| 文档双向索引 (`[Ref: P3]` 解析)       | P1     | 4h   |
| 键盘快捷键 (`useFullscreenShortcuts`) | P2     | 2h   |
| 公式渲染 (KaTeX 集成)                 | P2     | 2h   |

### API 需求 (需新增)

> 详见 [API 设计文档](./api/cost_estimation.md)

| 端点                             | 方法 | 说明             |
| -------------------------------- | ---- | ---------------- |
| `/api/v1/workflows/:id/estimate` | POST | 预估成本         |
| `/api/v1/sessions/:id/review`    | POST | 提交人类裁决结果 |

### 交付物
- [ ] 人类裁决节点可暂停并等待用户输入
- [ ] 启动前显示成本预估
- [ ] 文档引用可点击跳转

### 📚 详细规格文档

> **[Sprint 4 Specs 目录](./specs/sprint4/README.md)**

| Spec ID  | 文档                                                                | 类型      | 优先级 |
| -------- | ------------------------------------------------------------------- | --------- | ------ |
| SPEC-301 | [HumanReviewModal](./specs/sprint4/SPEC-301-human-review-modal.md)  | Component | P0     |
| SPEC-302 | [CostEstimator](./specs/sprint4/SPEC-302-cost-estimator.md)         | Component | P1     |
| SPEC-303 | [DocumentReference](./specs/sprint4/SPEC-303-document-reference.md) | Feature   | P1     |
| SPEC-304 | [快捷键](./specs/sprint4/SPEC-304-fullscreen-shortcuts.md)          | Hook      | P2     |
| SPEC-305 | [KaTeX 渲染](./specs/sprint4/SPEC-305-katex-rendering.md)           | Feature   | P2     |

---

## 后端补充任务 (并行进行)

### Sprint 1-2 并行

| 任务                     | 优先级 |
| ------------------------ | ------ |
| 实现 `SequenceProcessor` | P1     |
| 实现 `VoteProcessor`     | P1     |
| 实现 `LoopProcessor`     | P2     |

### Sprint 3-4 并行

| 任务                        | 优先级 |
| --------------------------- | ------ |
| 实现 `FactCheckProcessor`   | P1     |
| 实现 `HumanReviewProcessor` | P0     |
| 实现 Templates CRUD API     | P1     |
| 实现 Cost Estimation API    | P1     |
| 完善三层记忆协议            | P2     |

### 📚 后端详细规格文档

> **[Backend Specs 目录](./specs/backend/README.md)**

| Spec ID  | 文档                                                                       | 类型        | 优先级 |
| -------- | -------------------------------------------------------------------------- | ----------- | ------ |
| SPEC-401 | [SequenceProcessor](./specs/backend/SPEC-401-sequence-processor.md)        | Processor   | P1     |
| SPEC-402 | [VoteProcessor](./specs/backend/SPEC-402-vote-processor.md)                | Processor   | P1     |
| SPEC-403 | [LoopProcessor](./specs/backend/SPEC-403-loop-processor.md)                | Processor   | P2     |
| SPEC-404 | [FactCheckProcessor](./specs/backend/SPEC-404-factcheck-processor.md)      | Processor   | P1     |
| SPEC-405 | [HumanReviewProcessor](./specs/backend/SPEC-405-human-review-processor.md) | Processor   | P0     |
| SPEC-406 | [Templates API](./specs/backend/SPEC-406-templates-api.md)                 | API         | P1     |
| SPEC-407 | [Cost Estimation API](./specs/backend/SPEC-407-cost-estimation-api.md)     | API         | P1     |
| SPEC-408 | [三层记忆协议](./specs/backend/SPEC-408-memory-protocol.md)                | Core        | P0     |
| SPEC-409 | [逻辑熔断](./specs/backend/SPEC-409-circuit-breaker.md)                    | Core        | P0     |
| SPEC-410 | [防幻觉传播](./specs/backend/SPEC-410-anti-hallucination.md)               | Core        | P1     |
| SPEC-411 | [联网搜索集成](./specs/backend/SPEC-411-search-integration.md)             | Integration | P1     |

---

## 里程碑

| 时间      | 里程碑                | 验收标准                         |
| --------- | --------------------- | -------------------------------- |
| Week 1 末 | **M1: Run Mode 可用** | 能运行简单工作流，消息正确显示   |
| Week 2 末 | **M2: 管理页面完整**  | 能管理群组和 Agent               |
| Week 3 末 | **M3: Builder 完整**  | 支持所有节点类型，模版可用       |
| Week 4 末 | **M4: MVP 完成**      | 人类裁决、成本预估、文档引用可用 |

---

## 技术债务清理 (持续进行)

- [ ] 测试覆盖率提升至 80%
- [ ] i18n 翻译完善
- [ ] 暗黑模式切换
- [ ] 性能优化 (Run Mode 轻量渲染)
