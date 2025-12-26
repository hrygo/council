# SPEC-1003: 架构标准化与配置重构 (Standardization)

## 1. 需求背景
系统目前存在多处硬编码（如 Agent 的默认模型、数据库超时时间、WebSocket 重连频率）。这使得不同环境下部署调试非常困难。此外，ID 命名在不同层级存在语义不一致的风险。

## 2. 目标定义
- 建立统一的配置加载机制。
- 移除核心逻辑中的硬编码字面量。
- 标准化前后端跨层传递的 ID 语义。

## 3. 设计方案

### 3.1 配置文件化 (`internal/pkg/config/`)
1.  **扩展 Config 结构**: 增加 `LLM_DEFAULT_MODEL`, `DB_TIMEOUT`, `WS_WINDOW_MS` 等字段。
2.  **默认值注入**: 所有的 `NewXXX` 工厂方法必须接收 `config` 对象，禁止在内部定义常量。
3.  **Agent 处理器重构**: 修改 `AgentProcessor.Process`，当 `ag.ModelConfig.Model` 为空时，从全局配置读取，而非代码内定义。

### 3.2 ID 语义规范
- 数据库主键: 统一使用 `_uuid` 后缀（如 `agent_uuid`, `session_uuid`）。
- 逻辑标识符: 统一使用 `_id` 后缀（如 `node_id`, `workflow_id`）。
- **核查清单**: 遍历所有 `v1` 路由，确保参数命名遵循此规则。

## 4. 验收标准
- [ ] 搜索代码，除配置加载类外，无 `gpt-4`, `8080` 等特定业务/环境常量的硬编码。
- [ ] 修改 `config.yaml` 后的配置能实时反映到系统运行行为中。
- [ ] 前后端 API 请求参数命名符合语义规范。
