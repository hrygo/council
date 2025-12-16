# The Council - 技术架构方案 (Technical Architecture Document)

这是一个基于 **The Council PRD v1.3.0** 的详细技术架构方案。

> [!NOTE]
> 本文档已拆分为模块化子文档。请点击以下链接查看详细内容。

---

## 目录 (Table of Contents)

### 0. 需求追踪
- [0. PRD-TDD 追溯矩阵 (Traceability Matrix)](./tdd/00_traceability.md)

### 1. 总体架构
- [1. 总体架构设计 (High-Level Architecture)](./tdd/01_architecture.md)

### 2. 核心模块详细设计 (Core Modules)
- [2.1 流程编排引擎 (Workflow Engine)](./tdd/02_core/01_workflow_engine.md)
- [2.2 AI 网关与模型路由 (AI Model Router)](./tdd/02_core/02_ai_gateway.md)
- [2.3 三层记忆协议 (Three-Tier Memory Protocol)](./tdd/02_core/03_rag.md)
- [2.4 群组管理模块 (Group Management)](./tdd/02_core/04_group_management.md)
- [2.5 Agent 工厂模块 (Agent Factory)](./tdd/02_core/05_agent_factory.md)
- [2.6 节点处理器详解 (Node Processors)](./tdd/02_core/06_node_processors.md)
- [2.7 会议萃取引擎 (Session Extraction Engine)](./tdd/02_core/07_extraction_engine.md)
- [2.8 工作流模板模块 (Template Library)](./tdd/02_core/08_template_library.md)
- [2.9 Context Injection 优先级构建器](./tdd/02_core/09_context_builder.md)
- [2.10 搜索工具集成模块 (Search Tool Integration)](./tdd/02_core/10_search_tool.md)
- [2.11 并发执行配置 (Parallel Execution)](./tdd/02_core/11_parallel_execution.md)
- [2.12 FactCheck 节点处理器 (Fact Verification Node)](./tdd/02_core/12_fact_check.md)
- [2.13 自然语言转工作流模块 (NL2Workflow)](./tdd/02_core/13_nl2workflow.md)
- [2.14 防御性中间件架构 (Defense Middleware Architecture)](./tdd/02_core/14_defense_mechanisms.md)

### 3. 通信与存储
- [3. 前后端通信协议 (Communication Protocol)](./tdd/03_communication.md)
- [4. 数据存储架构 (Storage)](./tdd/04_storage.md)

### 4. 前端与实施
- [5. 前端架构关键点 (Frontend Specifics)](./tdd/05_frontend.md)
- [6. 开发实施路线 (Implementation Phases)](./tdd/06_implementation_plan.md)

### 5. 其他
- [7. 非功能需求实现 (Non-Functional Requirements)](./tdd/07_nfr.md)
- [8. 技术决策记录 (Architecture Decision Records)](./tdd/08_adr.md)
- [9. 部署与打包架构 (Deployment Architecture)](./tdd/09_deployment.md)

---

| 属性         | 内容                                                 |
| ------------ | ---------------------------------------------------- |
| **项目名称** | The Council (理事会)                                 |
| **架构模式** | **Cloud-Native WebApp** (React SPA + Go Backend API) |
| **核心语言** | TypeScript (Frontend), Go (Backend)                  |
| **数据存储** | PostgreSQL + pgvector (Docker 部署)                  |
