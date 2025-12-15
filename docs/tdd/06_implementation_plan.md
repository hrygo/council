# 6. 开发实施路线 (Implementation Phases)

### Phase 1: 骨架搭建 (Week 1)

1. **Go Init**: 搭建 Gin 服务，集成 GORM/pgx 和 pgvector。
2. **DB Schema**: 完成 SQL 迁移脚本编写。
3. **LLM Driver**: 写通 OpenAI 和 Ollama 的基础调用接口。

### Phase 2: 引擎核心 (Week 2-3)

1. **DAG Scheduler**: 实现 NodeProcessor 接口，跑通简单的 "Start -> Agent -> End" 流程。
2. **WebSocket Hub**: 实现多路流式推送。

### Phase 3: 前端交互 (Week 4-5)

1. **Canvas**: 集成 React Flow，实现 JSON 数据到图的互转。
2. **Chat UI**: 实现分栏渲染 markdown 流。

### Phase 4: 记忆与集成 (Week 6)

1. **RAG**: 实现文件解析 (Text Splitter) 和向量入库。
2. **联调**: 串联全流程，测试并发性能。
