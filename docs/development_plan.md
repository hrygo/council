# 开发计划 (Development Plan)

> **策略**: 前端优先，API Contract First  
> **质量内建**: 每个 Spec 必须通过 CI (Lint/Test) 和验收标准方可标记 Done  
> **文档版本**: v2.0 (2025-12-30 重构)

---

## 一、进度总览 (Dashboard)

### 1.1 Sprint 状态

| Sprint | 名称                   |  状态  | 完成度 |
| :----: | :--------------------- | :----: | :----: |
| S1-S4  | MVP 核心功能           | ✅ Done |  100%  |
|   S5   | Post-MVP (i18n, E2E)   | ✅ Done |  100%  |
|   S6   | Default Experience     | ✅ Done |  100%  |
|   S7   | UX Polish              | ✅ Done |  100%  |
|   S8   | Meeting Room Fix       | ✅ Done |  100%  |
|   S9   | Quality Base           | ✅ Done |  100%  |
|  S10   | System Hardening       | ✅ Done |  100%  |
|  S11   | Dialecta 2.0 Evolution | ✅ Done |  100%  |
|  S12   | Integrated Visibility  | ✅ Done |  100%  |

### 1.2 里程碑达成

| 里程碑 | 名称             | 验收标准                     | 状态  |
| :----: | :--------------- | :--------------------------- | :---: |
|   M1   | Run Mode         | 运行简单工作流，消息正确显示 |   ✅   |
|   M2   | 管理页面         | 群组和 Agent CRUD 完成       |   ✅   |
|   M3   | Builder          | 所有节点类型，模版可用       |   ✅   |
|   M4   | MVP              | 人类裁决、成本预估、文档引用 |   ✅   |
|   M5   | 国际化           | i18n + E2E 测试              |   ✅   |
|   M6   | Out-of-Box       | The Council 开箱即用         |   ✅   |
|   M7   | Meeting Room Fix | 会议室功能完善               |   ✅   |
|   M8   | Dialecta 2.0     | Tool Use + Logic Loop        |   ✅   |
|   M9   | Visibility       | VFS + Analytics              |   ✅   |

### 1.3 完成度统计

| 维度     | 已完成 | 总计  | 完成率 |
| :------- | :----: | :---: | :----: |
| Sprint   |   11   |  12   |  92%   |
| 前端任务 |   22   |  22   |  100%  |
| 后端任务 |   17   |  17   |  100%  |
| 总任务数 |   56   |  58   |  97%   |

---

## 二、当前进行中 (Active Sprints)

### 2.1 Sprint 12: Integrated Visibility

> **目标**: 将 Backend 的自主进化引擎可视化，提供 VFS 文件浏览器、Diff 审查和实时循环分析图表。

#### 任务矩阵

| ID   | 任务                                | Spec      | 优先级 | 状态  |
| :--- | :---------------------------------- | :-------- | :----: | :---: |
| 12.1 | VFS Explorer UI (Codebase Tab)      | SPEC-1201 |   P0   |   ✅   |
| 12.2 | Advanced Human Review (Diff Editor) | SPEC-1202 |   P0   |   ✅   |
| 12.3 | Loop Analytics (Score Chart)        | SPEC-1203 |   P1   |   ✅   |
| 12.4 | E2E Integration (Run Optimize Flow) | -         |   P0   |   ✅   |
| 12.5 | Naming Standardization (Refactor)   | SPEC-1205 |   P1   |   ✅   |

#### 规格文档

| SPEC ID   | 文档                                                          | 类型     | 状态  |
| :-------- | :------------------------------------------------------------ | :------- | :---: |
| SPEC-1201 | [VFS Frontend](./specs/sprint12/SPEC-1201-vfs-frontend.md)    | Frontend |   ✅   |
| SPEC-1202 | [Diff Review](./specs/sprint12/SPEC-1202-diff-review.md)      | Frontend |   ✅   |
| SPEC-1203 | [Loop Charts](./specs/sprint12/SPEC-1203-loop-charts.md)      | Frontend |   ✅   |
| SPEC-1205 | [Tech Debt Refactor](./specs/sprint12/SPEC-1205-tech-debt.md) | Refactor |   ✅   |

#### 验收标准

- [x] 系统界面右侧栏包含 "Codebase" 标签页，可浏览 VFS 文件树 (SPEC-1201)
- [x] 点击文件可查看内容，支持 Diff View 对比版本 (SPEC-1201)
- [x] Human Review 弹窗内嵌入 Diff Editor (SPEC-1202)
- [x] 顶部或侧边显示 "Optimization Score" 趋势图 (SPEC-1203)
- [x] ID 命名规范化完成 (`workflow_uuid`, `context_data` 等) (SPEC-1205)
- [x] 完整跑通 `council_optimize` 流程

---

### 2.2 Sprint 10: System Hardening (并行)

> **目标**: 完成跨模块的深度集成，消除技术债务，确保系统生产就绪。

#### 任务矩阵

| ID   | 任务                           | 类型     | 优先级 | 状态  |
| :--- | :----------------------------- | :------- | :----: | :---: |
| 10.1 | E2E 消息数据一致性测试         | QA       |   P0   |   ✅   |
| 10.2 | Nodes 测试覆盖率补全           | QA       |   P1   |   ✅   |
| 10.3 | 配置文件化 & 硬编码清理        | Refactor |   P1   |   ✅   |
| 10.4 | 全局 ID 命名规范化审查         | Refactor |   P2   |   ✅   |
| 10.5 | Store 类型字段完整性校验       | QA       |   P2   |   ✅   |
| 10.6 | Core 解耦收尾 (移除特定 Agent) | Refactor |   P1   |   ✅   |

#### 规格文档

| Spec ID  | 文档                                                         | 类型           |
| :------- | :----------------------------------------------------------- | :------------- |
| SPEC-901 | [RAG 深度集成](./specs/sprint10/SPEC-901-rag-integration.md) | Integration    |
| SPEC-902 | [QA 质量硬化](./specs/sprint10/SPEC-902-qa-hardening.md)     | QA             |
| SPEC-903 | [架构标准化](./specs/sprint10/SPEC-903-standardization.md)   | Refactor       |
| SPEC-904 | [性能与安全](./specs/sprint10/SPEC-904-perf-security.md)     | Infrastructure |

---

## 三、已完成 Sprints (归档)

> 以下为历史 Sprint 的精简记录，详细规格请参阅对应目录。

### 3.1 Sprint 11: Dialecta 2.0 Evolution ✅

**目标**: 升级核心引擎，实现 Tool Use 与 Logic Loop 自主进化闭环。

| ID   | 任务                      | Spec      | 状态  |
| :--- | :------------------------ | :-------- | :---: |
| 11.0 | Relational VFS Foundation | SPEC-1100 |   ✅   |
| 11.1 | Core Tool Infrastructure  | SPEC-1101 |   ✅   |
| 11.2 | System Surgeon Agent      | SPEC-1102 |   ✅   |
| 11.3 | Logic Loop Processor      | SPEC-1103 |   ✅   |
| 11.4 | Context Synthesizer       | SPEC-1104 |   ✅   |

**规格目录**: [specs/sprint11/](./specs/sprint11/)

---

### 3.2 Sprint 6-9: Experience & Quality ✅

| Sprint | 名称               | 核心交付                               | 规格目录                           |
| :----: | :----------------- | :------------------------------------- | :--------------------------------- |
|   S6   | Default Experience | 系统默认 Agents/Workflows, 开箱即用    | [specs/sprint6/](./specs/sprint6/) |
|   S7   | UX Polish          | Session Creation UI, LLM Registry      | [specs/sprint7/](./specs/sprint7/) |
|   S8   | Meeting Room Fix   | 会议启动流程, 知识库面板, 辩论流程还原 | [specs/sprint8/](./specs/sprint8/) |
|   S9   | Quality Base       | WebSocket 类型自动生成 (tygo)          | -                                  |

---

### 3.3 Sprint 1-5: MVP Core ✅

| Sprint | 名称           | 核心交付                                  | 规格目录                           |
| :----: | :------------- | :---------------------------------------- | :--------------------------------- |
|   S1   | 运行时状态重构 | SessionStore, WorkflowRunStore, ChatPanel | [specs/sprint1/](./specs/sprint1/) |
|   S2   | 管理页面       | 群组/Agent CRUD, 模型配置                 | [specs/sprint2/](./specs/sprint2/) |
|   S3   | Builder 增强   | 节点属性面板, 模版库, 向导模式            | [specs/sprint3/](./specs/sprint3/) |
|   S4   | 高级功能       | HumanReview, CostEstimator, 文档引用      | [specs/sprint4/](./specs/sprint4/) |
|   S5   | Post-MVP       | i18n 国际化, E2E 测试                     | [specs/sprint5/](./specs/sprint5/) |

**后端规格**: [specs/backend/](./specs/backend/) (Processors, API, 三层记忆, 熔断, 防幻觉)

---

## 四、规格文档总索引

### 4.1 按类型分类

| 类型           | Spec 范围        | 说明                        |
| :------------- | :--------------- | :-------------------------- |
| **Frontend**   | SPEC-001 ~ 305   | 前端组件、Store、UI 功能    |
| **Backend**    | SPEC-401 ~ 411   | 后端 Processor、API、中间件 |
| **Post-MVP**   | SPEC-501 ~ 504   | i18n, E2E, 性能, 安全       |
| **Default**    | SPEC-601 ~ 609   | 开箱即用体验                |
| **UX/Bugfix**  | SPEC-701 ~ 703   | UX 优化与问题修复           |
| **Meeting**    | SPEC-801 ~ 805   | 会议室功能增强              |
| **Hardening**  | SPEC-901 ~ 904   | 系统硬化与标准化            |
| **Evolution**  | SPEC-1100 ~ 1104 | Dialecta 2.0 进化引擎       |
| **Visibility** | SPEC-1201 ~ 1205 | 可视化与控制                |

### 4.2 快速导航

| 目录                                        | 说明             |
| :------------------------------------------ | :--------------- |
| [specs/sprint1/](./specs/sprint1/README.md) | S1: 运行时状态   |
| [specs/sprint2/](./specs/sprint2/README.md) | S2: 管理页面     |
| [specs/sprint3/](./specs/sprint3/README.md) | S3: Builder      |
| [specs/sprint4/](./specs/sprint4/README.md) | S4: 高级功能     |
| [specs/sprint5/](./specs/sprint5/)          | S5: Post-MVP     |
| [specs/sprint6/](./specs/sprint6/)          | S6: Default      |
| [specs/sprint7/](./specs/sprint7/)          | S7: UX Polish    |
| [specs/sprint8/](./specs/sprint8/)          | S8: Meeting Room |
| [specs/sprint10/](./specs/sprint10/)        | S10: Hardening   |
| [specs/sprint11/](./specs/sprint11/)        | S11: Evolution   |
| [specs/sprint12/](./specs/sprint12/)        | S12: Visibility  |
| [specs/backend/](./specs/backend/README.md) | 后端 Processors  |

---

## 五、技术债务与待办

### 5.1 技术债务跟踪

| 状态  | 债务项                 | 关联 Spec | 说明                                 |
| :---: | :--------------------- | :-------- | :----------------------------------- |
|   ✅   | WebSocket 类型自动生成 | S9        | tygo 集成完成                        |
|   ✅   | ID 命名规范化          | SPEC-1205 | `workflow_uuid`, `context_data` 统一 |
|   ⬜   | 硬编码默认值清理       | SPEC-903  | 移入 config.yaml                     |
|   ⬜   | Core 层解耦验证        | 10.6      | 移除特定 Agent 依赖                  |
|   ⬜   | 性能优化 (虚拟滚动)    | SPEC-904  | 长消息列表优化                       |

### 5.2 已搁置任务

| 任务              | 原优先级 | 搁置原因                       |
| :---------------- | :------: | :----------------------------- |
| 安全强化          |    P3    | 开发环境无需，生产环境另行处理 |
| 暗黑模式切换按钮  |    P2    | 已支持系统自动切换，ROI 不足   |
| Run Mode 轻量渲染 |    P2    | 当前性能可接受                 |

### 5.3 Backlog (低优先级)

| 任务             | 备注          |
| :--------------- | :------------ |
| GraphQL API 支持 | 可选替代 REST |
| 批量导入 Agent   | 便捷性功能    |
| 工作流版本控制   | Git 集成      |
| 性能监控面板     | APM 集成      |
| 多租户支持       | SaaS 模式     |

---

## 六、质量门禁

| 门禁项         | 检查内容                            |       自动化        |
| :------------- | :---------------------------------- | :-----------------: |
| **Lint**       | 代码风格符合规范                    |        ✅ CI         |
| **Test**       | 单元测试通过，核心层覆盖率 ≥ 80%    |        ✅ CI         |
| **Contract**   | 前后端 WebSocket 消息类型一致       |       ✅ tygo        |
| **E2E**        | 关键流程端到端测试通过              |    ✅ Playwright     |
| **Decoupling** | Core 层不依赖具体 Agent 实现        | ✅ verify-decoupling |
| **Config**     | 所有默认值可从环境变量/配置文件覆盖 |      ⬜ Manual       |

---

## 七、核心成果总结

### 7.1 系统能力

✅ **完整的多智能体协作系统**
- 工作流设计器（9 种节点类型）
- 实时执行引擎（WebSocket 通信）
- 三层记忆系统（Short/Long/External）
- 人工审核机制 + VFS 版本控制
- 成本估算与控制

✅ **开箱即用体验**
- 默认 Agent（Affirmative, Negative, Adjudicator, Surgeon）
- 预置工作流（Debate, Optimize）
- Tool Use 支持（write_file, read_file）

✅ **质量保障**
- 前端单元测试覆盖
- 后端测试覆盖 80%+
- E2E 关键流程测试
- CI/CD 自动化

### 7.2 技术亮点

- **类型安全**: tygo 自动生成前后端类型定义
- **国际化**: i18n 支持中英文切换
- **防御机制**: 熔断器、反幻觉传播
- **版本控制**: HumanReview 前自动备份 (VFS)
- **知识检索**: 集成外部搜索（Tavily）
- **Tool Use**: Agent 可直接修改 VFS 文件

### 7.3 架构优势

- **解耦设计**: Core 层独立于具体业务
- **可扩展**: 支持自定义节点、Provider、Tool
- **可观测**: 完整的日志和审计追踪
- **高可用**: 状态管理、错误恢复机制

---

## 附录

### A. 术语表

| 术语            | 定义                                             |
| :-------------- | :----------------------------------------------- |
| **Sprint**      | 敏捷开发中的迭代周期，通常为 1-2 周              |
| **Spec**        | Specification 的缩写，指规格说明文档             |
| **P0/P1/P2/P3** | 优先级标记，P0 最高，P3 最低                     |
| **技术债务**    | 为快速交付而采用的非最优方案，需要后续偿还       |
| **tygo**        | 自动将 Go 结构体转换为 TypeScript 类型定义的工具 |
| **VFS**         | Virtual File System，会话内的虚拟文件系统        |
| **Tool Use**    | Agent 调用外部工具（如文件写入）的能力           |

### B. 变更日志

| 日期       | 版本 | 变更说明                                    |
| :--------- | :--- | :------------------------------------------ |
| 2025-12-30 | v2.0 | 重构文档结构，统一格式，更新 Sprint 12 状态 |
| 2025-12-29 | v1.5 | 添加 Sprint 11/12，完成 SPEC-1205           |
| 2025-12-26 | v1.4 | Sprint 11 完成，Dialecta 2.0 发布           |

---

*文档维护: 请在每个 Sprint 结束时更新进度总览和归档对应 Sprint。*
