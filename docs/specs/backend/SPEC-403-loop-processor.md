# SPEC-403: LoopProcessor

> **优先级**: P2 | **关联 PRD**: F.3.1 Loop 节点 | **关联 TDD**: 06_node_processors.md

---

## 1. 数据结构

```go
type LoopProcessor struct {
    MaxRounds     int
    AgentPairs    [][2]string  // [[A, B], [B, A]]
    ExitCondition string       // "max_rounds" | "consensus"
}

type DebateMessage struct {
    Agent   string
    Content string
    Round   int
}
```

---

## 2. 实现逻辑

```go
func (l *LoopProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    history := []DebateMessage{}
    
    for round := 1; round <= l.MaxRounds; round++ {
        for _, pair := range l.AgentPairs {
            attacker, defender := pair[0], pair[1]
            
            prompt := l.buildDebatePrompt(attacker, defender, history, round)
            response := l.callAgent(ctx, attacker, prompt, stream)
            
            history = append(history, DebateMessage{
                Agent:   attacker,
                Content: response,
                Round:   round,
            })
            
            stream <- StreamEvent{
                Event: "debate:round",
                Data: map[string]interface{}{
                    "round":   round,
                    "agent":   attacker,
                    "content": response,
                },
            }
        }
        
        // 检查共识
        if l.ExitCondition == "consensus" && l.checkConsensus(history) {
            break
        }
    }
    
    return map[string]interface{}{
        "debate_history": history,
        "total_rounds":   len(history) / len(l.AgentPairs),
    }, nil
}
```

---

## 3. 共识检测

```go
func (l *LoopProcessor) checkConsensus(history []DebateMessage) bool {
    // 简单实现：检查最后一轮是否有明确 agree 关键词
    // 更复杂可用 LLM 判断
    if len(history) < 2 {
        return false
    }
    lastMsg := history[len(history)-1].Content
    return strings.Contains(strings.ToLower(lastMsg), "i agree") ||
           strings.Contains(lastMsg, "达成共识")
}
```

---

## 4. 测试用例

- 达到最大轮数退出
- 提前共识退出
- 历史正确传递
