# 2.3 双层记忆系统 (RAG Implementation)

利用 PostgreSQL 的 `pgvector` 实现。

**数据库查询逻辑 (SQL)**：

当一个 `Agent` 在 `Group` 中发言时，我们需要构建混合查询：

```sql
SELECT content, 1 - (embedding <=> $query_embedding) AS similarity
FROM memories
WHERE 
    -- 1. 必须属于当前群组 (Project Context)
    group_id = $current_group_id 
    AND (
        -- 2. 要么是这个群里的通用记忆
        agent_origin_id IS NULL 
        OR 
        -- 3. 要么是这个 Agent 类型跨项目的经验 (Persona Context)
        -- 注意：这里需要业务逻辑层的配合，agent_origin_id 应该指向具体的 Agent 原型 ID
        agent_origin_id = $current_agent_id
    )
ORDER BY similarity DESC
LIMIT 5;
```
