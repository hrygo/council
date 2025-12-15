# 2.9 Context Injection 优先级构建器 (Context Builder)

对应 PRD F.5.2，实现记忆检索时的优先级拼接逻辑。

```go
// ContextBuilder 按 PRD F.5.2 定义的 5 层优先级构建上下文
type ContextBuilder struct {
    Group       *Group
    Agent       *Agent
    Session     *Session
    MemoryStore MemoryRepository
    Embedder    LLMProvider
}

// ContextPriority 定义注入优先级
type ContextPriority int

const (
    PriorityGroupSystem    ContextPriority = 1 // 群定位 (最高)
    PriorityAgentPersona   ContextPriority = 2 // 人设
    PriorityGroupMemory    ContextPriority = 3 // 群相关记忆
    PriorityAgentMemory    ContextPriority = 4 // 角色经验记忆
    PriorityCurrentProposal ContextPriority = 5 // 当前提案内容
)

// Build 按优先级构建完整上下文
func (c *ContextBuilder) Build(ctx context.Context, query string) (string, int, error) {
    var parts []string
    var totalTokens int
    
    // ─────────────────────────────────────────────
    // Priority 1: 群定位 (System Prompt) - 最高优先级
    // ─────────────────────────────────────────────
    if c.Group.SystemPrompt != "" {
        section := fmt.Sprintf("## 📋 群组定位\n%s", c.Group.SystemPrompt)
        parts = append(parts, section)
        totalTokens += estimateTokens(section)
    }
    
    // ─────────────────────────────────────────────
    // Priority 2: 人设 (Persona Prompt)
    // ─────────────────────────────────────────────
    personaSection := fmt.Sprintf("## 🎭 你的角色\n%s", c.Agent.PersonaPrompt)
    parts = append(parts, personaSection)
    totalTokens += estimateTokens(personaSection)
    
    // ─────────────────────────────────────────────
    // Priority 3: 群相关记忆 (RAG with group_id)
    // ─────────────────────────────────────────────
    queryEmbedding, _ := c.Embedder.Embed(ctx, query)
    groupMemories, _ := c.MemoryStore.Search(ctx, MemoryQuery{
        GroupID:   &c.Group.ID,
        Embedding: queryEmbedding,
        Limit:     5,
    })
    if len(groupMemories) > 0 {
        memSection := "## 📚 项目历史记忆\n" + formatMemories(groupMemories)
        parts = append(parts, memSection)
        totalTokens += estimateTokens(memSection)
    }
    
    // ─────────────────────────────────────────────
    // Priority 4: 角色经验记忆 (RAG with agent_id)
    // ─────────────────────────────────────────────
    agentMemories, _ := c.MemoryStore.Search(ctx, MemoryQuery{
        AgentID:   &c.Agent.ID,
        Embedding: queryEmbedding,
        Limit:     3,
    })
    if len(agentMemories) > 0 {
        expSection := "## 💡 你的历史经验\n" + formatMemories(agentMemories)
        parts = append(parts, expSection)
        totalTokens += estimateTokens(expSection)
    }
    
    // ─────────────────────────────────────────────
    // Priority 5: 当前提案内容
    // ─────────────────────────────────────────────
    proposalSection := fmt.Sprintf("## 📝 当前提案\n%s", c.Session.Proposal.Text)
    if len(c.Session.Proposal.Files) > 0 {
        proposalSection += "\n\n### 附件内容\n" + c.Session.Proposal.ParsedContent
    }
    parts = append(parts, proposalSection)
    totalTokens += estimateTokens(proposalSection)
    
    return strings.Join(parts, "\n\n---\n\n"), totalTokens, nil
}

func formatMemories(memories []Memory) string {
    var lines []string
    for i, m := range memories {
        lines = append(lines, fmt.Sprintf("%d. %s (相似度: %.2f)", i+1, m.Content, m.Similarity))
    }
    return strings.Join(lines, "\n")
}

func estimateTokens(text string) int {
    // 粗略估算: 中文约 1.5 字符/token, 英文约 4 字符/token
    return len([]rune(text)) / 2
}
```
