# SPEC-408: 三层记忆协议 (Memory Purification Protocol)

> **优先级**: P0 | **预估工时**: 6h  
> **关联 PRD**: F.5.1-F.5.3 | **关联 TDD**: 02_core/03_rag.md

---

## 1. 架构概述

```
┌─────────────────────────────────────────────────────┐
│                  用户交互层                          │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│ Layer 1: Quarantine (隔离区)                        │
│ - 会议原始产出，物理隔离                            │
│ - TTL: 永久 (直到晋升或删除)                        │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│ Layer 2: Working Memory (工作记忆)                  │
│ - 热缓存，入口过滤                                  │
│ - TTL: 24 小时                                      │
│ - Scope: Project ID 隔离                            │
└─────────────────────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────┐
│ Layer 3: Long-Term DB (长期记忆)                    │
│ - 经验证的高价值知识                                │
│ - 向量化存储 (pgvector)                             │
└─────────────────────────────────────────────────────┘
```

---

## 2. 数据模型

```go
type MemoryEntry struct {
    ID          string    `json:"id"`
    ProjectID   string    `json:"project_id"`
    SessionID   string    `json:"session_id"`
    Content     string    `json:"content"`
    Embedding   []float32 `json:"embedding"`
    Layer       string    `json:"layer"` // quarantine, working, longterm
    CreatedAt   time.Time `json:"created_at"`
    ExpiresAt   *time.Time `json:"expires_at"` // 仅 working 层有效
    Verified    bool      `json:"verified"`
    PromotedAt  *time.Time `json:"promoted_at"`
}
```

---

## 3. Layer 1: Quarantine (隔离区)

```go
type QuarantineService struct {
    Store MemoryStore
}

func (q *QuarantineService) Save(sessionID string, content string) error {
    entry := &MemoryEntry{
        ID:        uuid.New().String(),
        SessionID: sessionID,
        Content:   content,
        Layer:     "quarantine",
        Verified:  false,
    }
    return q.Store.Insert(entry)
}
```

**特点**:
- 所有会议产出默认进入此层
- 不参与 RAG 检索
- 防止低质量内容污染核心库

---

## 4. Layer 2: Working Memory (工作记忆)

```go
type WorkingMemoryService struct {
    Store       MemoryStore
    Filter      IngressFilter
    TTL         time.Duration // 24h
}

func (w *WorkingMemoryService) Ingest(entry *MemoryEntry) error {
    // 入口过滤：自洽性检查
    if !w.Filter.CheckConsistency(entry.Content) {
        return ErrFailedConsistencyCheck
    }
    
    entry.Layer = "working"
    entry.ExpiresAt = time.Now().Add(w.TTL)
    return w.Store.Insert(entry)
}

// 定时清理过期条目
func (w *WorkingMemoryService) Cleanup() {
    w.Store.DeleteExpired("working")
}
```

**入口过滤 (Ingress Filter)**:
```go
type IngressFilter interface {
    CheckConsistency(content string) bool
}

// 基于 LLM 的自洽性检查
func (f *LLMFilter) CheckConsistency(content string) bool {
    prompt := `判断以下内容是否自相矛盾或包含明显幻觉：
    %s
    回答 YES 或 NO`
    response := f.LLM.Call(prompt)
    return response == "NO"
}
```

---

## 5. Layer 3: Knowledge Promotion (知识晋升)

```go
type PromotionService struct {
    Store     MemoryStore
    Embedder  EmbeddingService
}

// 一键晋升
func (p *PromotionService) Promote(entryID string) error {
    entry, _ := p.Store.Get(entryID)
    
    // 生成向量
    embedding := p.Embedder.Embed(entry.Content)
    entry.Embedding = embedding
    entry.Layer = "longterm"
    entry.PromotedAt = time.Now()
    entry.Verified = true
    
    return p.Store.Update(entry)
}

// 智能简报生成
func (p *PromotionService) GenerateDigest(projectID string) *Digest {
    entries := p.Store.ListQuarantine(projectID)
    // 聚类 + 摘要
    return p.cluster(entries)
}
```

---

## 6. 前端 UI

### 简报面板

```tsx
const KnowledgeDigest: FC<{ projectId: string }> = ({ projectId }) => {
  const { data: digest } = useDigest(projectId);
  const { mutate: promote } = usePromote();

  return (
    <Card>
      <CardHeader>
        <h3>📚 知识简报 (本周)</h3>
      </CardHeader>
      <CardContent>
        {digest?.insights.map((insight, i) => (
          <div key={i} className="flex items-center justify-between">
            <p>{insight.summary}</p>
            <Button size="sm" onClick={() => promote(insight.id)}>
              ⬆️ 晋升
            </Button>
          </div>
        ))}
      </CardContent>
    </Card>
  );
};
```

### 临时上下文标识

```tsx
// 引用 Working Memory 时显示
<Badge variant="outline" className="text-amber-500">
  ⚡️ 临时上下文
</Badge>
```

---

## 7. 测试要点

- [ ] Quarantine 隔离正确
- [ ] Working Memory 24h 过期
- [ ] Ingress Filter 过滤幻觉
- [ ] 晋升后可 RAG 检索
