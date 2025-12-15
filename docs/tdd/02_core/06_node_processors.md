# 2.6 节点处理器详解 (Node Processors)

补充 PRD F.3.1 中 Vote 和 Loop 节点的技术实现。

#### 2.6.1 Vote 节点 (表决)

强制 Agent 输出结构化投票结果。

```go
type VoteProcessor struct {
    Threshold    float64        // 通过阈值 (如 0.67 = 2/3 多数)
    VoteType     string         // "yes_no" | "score_1_10"
    AgentResults map[string]VoteResult
}

type VoteResult struct {
    AgentID   string
    Decision  string  // "yes", "no", "abstain"
    Score     int     // 1-10 (仅 score 模式)
    Reasoning string  // 投票理由
}

func (v *VoteProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    // 1. 构造强制结构化输出的 Prompt
    prompt := fmt.Sprintf(`基于以下内容进行投票，你必须严格按照 JSON 格式输出：
%s

请输出：
{"decision": "yes/no", "score": 1-10, "reasoning": "你的理由"}`, input["context"])
    
    // 2. 并发调用所有 Agent 获取投票
    // 3. 聚合结果，判断是否通过阈值
    passed := v.calculateResult()
    
    return map[string]interface{}{
        "passed":  passed,
        "votes":   v.AgentResults,
        "summary": v.generateSummary(),
    }, nil
}
```

#### 2.6.2 Loop 节点 (循环辩论)

实现多轮对话交锋。

```go
type LoopProcessor struct {
    MaxRounds   int              // 最大轮次
    AgentPairs  [][2]string      // 辩论配对 [[A, B], [B, A]]
    ExitCondition string         // "max_rounds" | "consensus"
}

func (l *LoopProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    history := []Message{}
    
    for round := 1; round <= l.MaxRounds; round++ {
        for _, pair := range l.AgentPairs {
            attacker, defender := pair[0], pair[1]
            
            // 构造带历史的 Prompt
            prompt := l.buildDebatePrompt(attacker, defender, history, round)
            
            // 调用 Agent 生成回复
            response := l.callAgent(ctx, attacker, prompt, stream)
            history = append(history, Message{Agent: attacker, Content: response, Round: round})
            
            // 推送事件
            stream <- StreamEvent{
                Event: "debate_round",
                Data: map[string]interface{}{
                    "round":   round,
                    "agent":   attacker,
                    "content": response,
                },
            }
        }
        
        // 检查是否达成共识 (可选)
        if l.ExitCondition == "consensus" && l.checkConsensus(history) {
            break
        }
    }
    
    return map[string]interface{}{
        "debate_history": history,
        "total_rounds":   len(history) / len(l.AgentPairs),
    }, nil
}

#### 2.6.3 HumanReview 节点 (人类裁决)

对应 PRD F.3.1 强制性的人机回环节点。

```go
type HumanReviewProcessor struct {
    ReviewType string // "approve_reject" | "edit_content"
}

func (h *HumanReviewProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    // 1. 生成决策草案
    draft := input["summary"]
    
    // 2. 暂停工作流，等待用户信号
    stream <- StreamEvent{
        Event: "human_review_required",
        Data: map[string]interface{}{
            "node_id": "current_node_id",
            "draft":   draft,
        },
    }
    
    // 3. 阻塞等待用户操作 (实际实现可能通过 channel 或数据库状态轮询)
    // 这里简化为等待一个 channel 信号
    userAction := <-h.UserSignal
    
    if userAction.Action == "reject" {
        return nil, fmt.Errorf("user rejected the draft: %s", userAction.Reason)
    }
    
    // 4. 用户可能修改了草案
    finalContent := userAction.Content
    if finalContent == "" {
        finalContent = draft.(string)
    }
    
    return map[string]interface{}{
        "reviewed_content": finalContent,
        "reviewer":         userAction.UserID,
        "approved_at":      time.Now(),
    }, nil
}
```
```
