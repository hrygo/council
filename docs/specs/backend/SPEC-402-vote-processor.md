# SPEC-402: VoteProcessor

> **优先级**: P1 | **关联 PRD**: F.3.1 Vote 节点 | **关联 TDD**: 06_node_processors.md

---

## 1. 数据结构

```go
type VoteProcessor struct {
    Threshold    float64  // 0.5-1.0
    VoteType     string   // "yes_no" | "score_1_10"
    AgentIDs     []string
    AgentResults map[string]VoteResult
}

type VoteResult struct {
    AgentID   string
    Decision  string  // "yes", "no", "abstain"
    Score     int     // 1-10 (仅 score 模式)
    Reasoning string
}
```

---

## 2. 实现逻辑

```go
func (v *VoteProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    prompt := v.buildVotePrompt(input)
    
    // 并发调用所有 Agent
    var wg sync.WaitGroup
    results := make(chan VoteResult, len(v.AgentIDs))
    
    for _, agentID := range v.AgentIDs {
        wg.Add(1)
        go func(aid string) {
            defer wg.Done()
            result := v.callAgentForVote(ctx, aid, prompt, stream)
            results <- result
        }(agentID)
    }
    
    wg.Wait()
    close(results)
    
    // 收集结果
    for r := range results {
        v.AgentResults[r.AgentID] = r
    }
    
    // 计算是否通过
    passed := v.calculatePassed()
    
    return map[string]interface{}{
        "passed":  passed,
        "votes":   v.AgentResults,
        "summary": v.generateSummary(),
    }, nil
}

func (v *VoteProcessor) calculatePassed() bool {
    yesCount := 0
    total := 0
    for _, r := range v.AgentResults {
        if r.Decision != "abstain" {
            total++
            if r.Decision == "yes" {
                yesCount++
            }
        }
    }
    if total == 0 {
        return false
    }
    return float64(yesCount)/float64(total) >= v.Threshold
}
```

---

## 3. 强制结构化输出

```go
const votePromptTemplate = `基于以下内容进行投票：
%s

你必须严格按照 JSON 格式输出：
{"decision": "yes/no/abstain", "score": 1-10, "reasoning": "你的理由"}`
```

---

## 4. 测试用例

- 全部 yes → 通过
- 投票不足阈值 → 不通过
- 有 abstain → 正确排除
