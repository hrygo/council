# 0. PRD-TDD 追溯矩阵 (Traceability Matrix)

本章节建立 PRD 需求与 TDD 技术方案的映射关系，确保需求覆盖可追踪。

### 功能需求追溯

| PRD 编号  | 需求描述                     | TDD 章节                                                                                                                                | 状态 |
| --------- | ---------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- | ---- |
| **F.1.1** | 创建/管理群组                | [2.4 群组管理模块](../tdd/02_core/04_group_management.md)                                                                               | ✅    |
| **F.1.2** | 群记忆隔离                   | [2.3 双层记忆系统](../tdd/02_core/03_rag.md), [4.3 数据库 Schema](../tdd/04_storage.md#43-数据库-schema-er-diagram)                     | ✅    |
| **F.2.1** | 角色定义                     | [2.5 Agent 工厂模块](../tdd/02_core/05_agent_factory.md)                                                                                | ✅    |
| **F.2.2** | 模型配置 (Model Agnostic)    | [2.2 AI 网关](../tdd/02_core/02_ai_gateway.md), [2.5 Agent 工厂模块](../tdd/02_core/05_agent_factory.md)                                | ✅    |
| **F.2.3** | 能力开关                     | [2.5 Agent 工厂模块](../tdd/02_core/05_agent_factory.md)                                                                                | ✅    |
| **F.3.1** | 节点编辑器                   | [2.1 流程编排引擎](../tdd/02_core/01_workflow_engine.md), [2.6 节点处理器](../tdd/02_core/06_node_processors.md)                        | ✅    |
| **F.3.2** | 模板库                       | [2.8 工作流模板模块](../tdd/02_core/08_template_library.md)                                                                             | ✅    |
| **F.3.3** | 向导模式 (Wizard Mode)       | [2.13 NL2Workflow 模块](../tdd/02_core/13_nl2workflow.md)                                                                               | ✅    |
| **F.4.0** | 弹性布局框架                 | [5.1 弹性布局实现](../tdd/05_frontend.md#51-弹性布局实现), [5.6 布局持久化](../tdd/05_frontend.md#56-布局状态持久化-layout-persistence) | ✅    |
| **F.4.1** | 流程监控 (含 Token 预估)     | [3.1 WebSocket 事件](../tdd/03_communication.md#31-websocket-事件定义)                                                                  | ✅    |
| **F.4.2** | 结构化对话流                 | [5.4 并行消息渲染](../tdd/05_frontend.md#54-并行消息渲染-parallel-message-ui)                                                           | ✅    |
| **F.4.3** | 双向索引                     | [5.5 双向文档索引](../tdd/05_frontend.md#55-双向文档索引-bidirectional-document-reference)                                              | ✅    |
| **F.4.4** | 成本预估模块                 | [3.1 WebSocket 事件](../tdd/03_communication.md#31-websocket-事件定义) (token_usage 事件)                                               | ✅    |
| **F.5.1** | 记忆写入 (萃取)              | [2.7 会议萃取引擎](../tdd/02_core/07_extraction_engine.md)                                                                              | ✅    |
| **F.5.2** | 记忆检索 (Context Injection) | [2.3 双层记忆系统](../tdd/02_core/03_rag.md)                                                                                            | ✅    |

### 非功能需求追溯

| PRD 章节     | 需求描述                       | TDD 章节                                                           | 状态 |
| ------------ | ------------------------------ | ------------------------------------------------------------------ | ---- |
| **4.隐私性** | 数据安全, Key 加密             | [7.1 安全与隐私](../tdd/07_nfr.md#71-安全与隐私-security--privacy) | ✅    |
| **4.性能**   | 3 Agent 并发, 首字延迟 < 500ms | [7.2 性能指标](../tdd/07_nfr.md#72-性能指标-performance-metrics)   | ✅    |
| **4.扩展性** | OpenAI Function Calling 标准   | [7.3 扩展性设计](../tdd/07_nfr.md#73-扩展性设计-extensibility)     | ✅    |
