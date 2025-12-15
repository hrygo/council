# 2.7 会议萃取引擎 (Session Extraction Engine)

对应 PRD F.5.1，会议结束后自动萃取关键信息。

```go
type ExtractionEngine struct {
    Embedder  LLMProvider  // 用于生成向量
    Splitter  TextSplitter // 文本切分器
}

type ExtractionResult struct {
    Conclusions []Conclusion  // 关键结论
    Decisions   []Decision    // 决策点
    ActionItems []ActionItem  // 待办事项
}

func (e *ExtractionEngine) Extract(ctx context.Context, session *Session) (*ExtractionResult, error) {
    // 1. 汇总所有节点输出
    fullContent := e.aggregateContent(session)
    
    // 2. 调用 LLM 萃取结构化信息
    extractPrompt := `请从以下会议记录中提取：
1. 关键结论 (conclusions)
2. 决策点 (decisions)  
3. 待办事项 (action_items)

会议内容：
%s

请以 JSON 格式输出。`
    
    result := e.callLLM(ctx, fmt.Sprintf(extractPrompt, fullContent))
    
    // 3. 切片向量化存入记忆库
    chunks := e.Splitter.Split(fullContent, 512) // 512 tokens per chunk
    for _, chunk := range chunks {
        embedding, _ := e.Embedder.Embed(ctx, chunk.Text)
        memory := &Memory{
            GroupID:   session.GroupID,
            AgentID:   nil, // 群记忆，不绑定特定 Agent
            Content:   chunk.Text,
            Embedding: embedding,
            SessionID: session.ID,
            CreatedAt: time.Now(),
        }
        e.saveMemory(ctx, memory)
    }
    
    return result, nil
}

// MemorySanitation 给定的会议结论打标或无效化
func (e *ExtractionEngine) SanitizeMemory(ctx context.Context, sessionID uuid.UUID, userFeedback []UserCorrection) error {
    for _, correction := range userFeedback {
        if correction.Action == "mark_invalid" {
            // 软删除或标记为负面样本
            e.MemoryStore.MarkInvalid(ctx, correction.MemoryID)
        }
    }
    return nil
}
```
