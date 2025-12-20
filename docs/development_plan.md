# 开发计划 (Development Plan)

> **策略**: 前端优先，API Contract First  
> **质量内建**: 每个 Spec 必须通过 CI (Lint/Test) 和验收标准方可标记 Done

---

## 一、进度总览

| Sprint | 阶段               |   状态   | 完成度 |
| :----: | :----------------- | :------: | :----: |
| S1-S4  | MVP 核心功能       |  ✅ Done  |  100%  |
|   S5   | Post-MVP 优化      | 🔄 进行中 |  50%   |
|   S6   | Default Experience | 🔄 进行中 |  95%   |

---

## 二、里程碑

| 时间   | 里程碑         | 验收标准                     |   状态   |
| :----- | :------------- | :--------------------------- | :------: |
| Week 1 | M1: Run Mode   | 运行简单工作流，消息正确显示 |    ✅     |
| Week 2 | M2: 管理页面   | 群组和 Agent CRUD 完成       |    ✅     |
| Week 3 | M3: Builder    | 所有节点类型，模版可用       |    ✅     |
| Week 4 | M4: MVP        | 人类裁决、成本预估、文档引用 |    ✅     |
| Week 5 | M5: 国际化     | i18n + E2E 测试              |    ✅     |
| Week 6 | M6: Out-of-Box | The Council 开箱即用         | 🔄 进行中 |

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
- [ ] 3 个系统 Agent 存在
- [ ] "The Council" 群组存在
- [ ] Debate + Optimize 流程存在
- [ ] `memory_retrieval` 节点可用
- [ ] HumanReview 前自动备份
- [ ] 完整 Optimize 循环可运行

**skill.md 覆盖**:
- [ ] Step 1: Memory Retrieval
- [ ] Step 2: Parallel + Agent
- [ ] Step 3: Scoring Matrix
- [ ] Step 4: Versioning
- [ ] Step 5: HumanReview
- [ ] Step 6: Loop

**解耦验证**:
- [ ] `make verify-decoupling` 通过
- [ ] 删除 `example/` 后系统正常

---

## 六、技术债务

| 任务                 | 优先级 | 状态  |
| :------------------- | :----: | :---: |
| 测试覆盖率提升至 80% |   P1   |   ⬜   |
| 暗黑模式切换         |   P2   |   ⬜   |
| Run Mode 轻量渲染    |   P2   |   ⬜   |
