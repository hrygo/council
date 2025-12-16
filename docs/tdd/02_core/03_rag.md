# 2.3 三层记忆协议 (Three-Tier Memory Protocol) 🆕

为解决传统 RAG 的 "Memory Pollution" (记忆污染) 与 "Archival Latency" (归档时滞) 问题，系统实施严格的三层数据流转协议。

### 1. 架构分层

| 层级       | 名称                               | 存储介质                              | 生命周期            | 可见性                         | 写入时机                   |
| :--------- | :--------------------------------- | :------------------------------------ | :------------------ | :----------------------------- | :------------------------- |
| **Tier 1** | **Cortex Quarantine** (记忆隔离区) | PostgreSQL (Table: `quarantine_logs`) | 永久 (直至手动处理) | **不可见** (仅供后台审计/简报) | 会议产生的任何 Token       |
| **Tier 2** | **Working Memory** (工作记忆)      | Redis (Key: `wm:{project_id}`)        | **24 Hours** (TTL)  | **可见** (标记为 Ephemeral)    | 通过 Ingress Filter 后写入 |
| **Tier 3** | **Long-Term Memory** (核心智库)    | pgvector (Table: `vectors`)           | 永久                | **可见** (Verified)            | 用户 One-Click Promote 后  |

### 2. 核心逻辑详解

#### 2.1 Tier 2: Ingress Filter (入口过滤)
只有满足以下条件的 Quarantine 数据才能进入 Hot Cache：
*   **自洽性检查**: Agent 自身对输出的 Confidence Score > 0.8。
*   **非重复**: 与当前上下文相似度 < 0.9。

#### 2.2 混合检索策略 (Hybrid Retrieval)

当 `Agent` 需要 Query 时，System 同时并行查询 Tier 2 和 Tier 3：

```go
// 伪代码逻辑
func FetchContext(query string) []Context {
    // 1. 并行查询
    chanHot := Go(Redis.Search(query)) // 速度极快
    chanCold := Go(PGVector.Search(query)) // 速度较慢
    
    // 2. 结果合并与降权
    hotResults := <- chanHot // 标记为 [⚡️ Unverified]
    coldResults := <- chanCold // 标记为 [✅ Verified]
    
    // 3. 排序策略：Hot数据如果相关性极高，优先展示，但给予醒目UI提示
    return MergeAndRank(hotResults, coldResults)
}
```

#### 2.3 Tier 3: Knowledge Promotion (晋升机制)
*   **Smart Digest**: 每周 CronJob 扫描 `quarantine_logs`，通过 LLM 生成摘要。
*   **Commit**: 用户点击 "Promote" 按钮，触发向量化服务 (Embedding Service)，将摘要写入 `vectors` 表。

### 3. 数据表设计更新

**新增 `quarantine_logs`**:
```sql
CREATE TABLE quarantine_logs (
    id UUID PRIMARY KEY,
    session_id UUID,
    content TEXT,
    raw_metadata JSONB, -- 包含 agent_id, timestamp, confidence_score
    created_at TIMESTAMP DEFAULT NOW()
);
```

**核心 `vectors` 表 (Tier 3)**:
保持原有设计，但数据来源仅限于 `Promote` 动作。
