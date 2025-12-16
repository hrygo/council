# 0. PRD-TDD 追溯矩阵 (Traceability Matrix)

本章节建立 PRD 需求与 TDD 技术方案的映射关系，确保需求覆盖可追踪。

### 功能需求追溯

| PRD 编号  | 需求描述                    | TDD 章节                                                                                                             | 状态 |
| --------- | --------------------------- | -------------------------------------------------------------------------------------------------------------------- | ---- |
| **F.1.1** | 创建/管理群组               | [2.4 群组管理模块](../tdd/02_core/04_group_management.md)                                                            | ✅    |
| **F.1.2** | 群记忆隔离                  | [2.3 三层记忆协议](../tdd/02_core/03_rag.md)                                                                         | ✅    |
| **F.2.1** | 角色定义                    | [2.5 Agent 工厂模块](../tdd/02_core/05_agent_factory.md)                                                             | ✅    |
| **F.2.2** | 模型配置 (Model Agnostic)   | [2.2 AI 网关](../tdd/02_core/02_ai_gateway.md)                                                                       | ✅    |
| **F.3.1** | 节点编辑器 (Build Mode)     | [2.1 流程编排引擎](../tdd/02_core/01_workflow_engine.md), [5. 前端架构](../tdd/05_frontend.md)                       | ✅    |
| **F.3.3** | 向导模式 (Wizard Mode)      | [2.13 NL2Workflow 模块](../tdd/02_core/13_nl2workflow.md)                                                            | ✅    |
| **F.3.5** | 记忆净化协议 (Memory Logic) | [2.3 三层记忆协议](../tdd/02_core/03_rag.md)                                                                         | ✅    |
| **F.4.x** | 沉浸式会议 (Run Mode)       | [5. 前端架构](../tdd/05_frontend.md)                                                                                 | ✅    |
| **F.4.4** | 防御性 UX                   | [2.14 防御性中间件](../tdd/02_core/14_defense_mechanisms.md)                                                         | ✅    |
| **F.6.1** | 逻辑熔断 (Hard Stop)        | [2.14 防御性中间件](../tdd/02_core/14_defense_mechanisms.md)                                                         | ✅    |
| **F.6.2** | 防幻觉传播 (Fact Check)     | [2.14 防御性中间件](../tdd/02_core/14_defense_mechanisms.md), [2.12 FactCheck Node](../tdd/02_core/12_fact_check.md) | ✅    |

### 非功能需求追溯

| PRD 章节     | 需求描述                    | TDD 章节                                                           | 状态 |
| ------------ | --------------------------- | ------------------------------------------------------------------ | ---- |
| **5.隐私性** | 数据安全, Key 加密          | [7.1 安全与隐私](../tdd/07_nfr.md#71-安全与隐私-security--privacy) | ✅    |
| **5.性能**   | 并发处理, Run Mode 轻量渲染 | [7.2 性能指标](../tdd/07_nfr.md#72-性能指标-performance-metrics)   | ✅    |
