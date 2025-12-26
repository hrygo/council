# SPEC-1001: RAG 深度数据集成 (Knowledge Panel Real Integration)

## 1. 需求背景
目前会议室右侧的 `KnowledgePanel` (见 `SPEC-805`) 使用的是 `KnowledgeHandler` 中硬编码的 Mock 数据。为了实现真正的智能助手价值，系统需要能够根据当前会话上下文，从向量数据库（PGVector）和 LTM（长期记忆）中实时检索相关知识。

## 2. 目标定义
- 替换 Mock 实现，建立从 UI 到持久层的数据通路。
- 支持按知识的相关度（Relevance Score）排序。
- 支持按记忆层级（Sandboxed/Working/Long-term）过滤。

## 3. 设计方案

### 3.1 后端重构 (`internal/api/handler/knowledge.go`)
1.  **依赖注入**: 在 `NewKnowledgeHandler` 中注入 `memory.Service`。
2.  **方法实现**: 修改 `GetSessionKnowledge`，解析 `sessionID` 和 `layer`。
3.  **检索逻辑**:
    - 调用 `h.memoryService.Retrieve(ctx, query, sessionID)`。
    - 将返回的 `memory.ContextItem` 转换为 `KnowledgeItem` DTO。
    - 处理 `RelevanceScore` 的格式化（将 0-1 的相似度映射为 1-5 级的显示）。

### 3.2 存储层适配 (`internal/core/memory/service.go`)
- 确保 `Retrieve` 方法能正确处理 `sessionID`（当前 Service 接口使用的是 `groupID`，需兼容 `sessionID` 或确保映射关系正确）。
- 优化向量检索 SQL，增加对 `layer` 标签的过滤支持。

### 3.3 数据闭环流程
1.  **节点触发**: 当 `memory_retrieval` 节点执行时，后端主动向 WebSocket 推送 `knowledge:updated` 事件。
2.  **前端响应**: 前端监听到事件后，调用 `GET /api/v1/sessions/:sessionID/knowledge` 刷新列表。

## 4. 验收标准
- [ ] 移除 `knowledge.go` 中的 `mockKnowledgeItems` 函数及相关引用。
- [ ] 在 `internal/resources/seeds` 中添加测试用记忆数据。
- [ ] 手动测试：输入特定的关键词，右侧面板应显示包含该关键词的历史片段。
- [ ] 单元测试：`GetSessionKnowledge` 能够成功返回向量数据库中的真实记录。
