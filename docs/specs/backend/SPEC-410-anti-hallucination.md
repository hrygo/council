# SPEC-410: 防幻觉传播 (Anti-Hallucination)

> **优先级**: P1 | **预估工时**: 3h  
> **关联 PRD**: F.6.2 | **关联 TDD**: 02_core/14_defense_mechanisms.md

---

## 1. 机制概述

在 Agent 信息传递链路中增加"事实校验层"，对未经验证的关键数据自动标记警示。

---

## 2. 校验流程

```
Agent A 输出 ──▶ Fact Verification Layer ──▶ Agent B 输入
                         │
                         ▼
               标记 "Verify Pending" (如需)
```

---

## 3. 实现

```go
type FactVerificationLayer struct {
    Patterns []regexp.Regexp // 需要验证的模式
}

var defaultPatterns = []string{
    `\d+%`,                     // 百分比
    `\$[\d,.]+`,                // 金额
    `\d{4}年`,                  // 年份
    `根据.*?(报告|研究|数据)`,   // 引用声明
}

type VerificationResult struct {
    Content      string
    Claims       []Claim
    NeedsVerify  bool
}

type Claim struct {
    Text       string
    Type       string  // percentage, currency, citation
    Verified   bool
    VerifyNote string
}

func (f *FactVerificationLayer) Analyze(content string) *VerificationResult {
    claims := []Claim{}
    
    for _, pattern := range f.Patterns {
        matches := pattern.FindAllString(content, -1)
        for _, match := range matches {
            claims = append(claims, Claim{
                Text:     match,
                Verified: false,
            })
        }
    }
    
    return &VerificationResult{
        Content:     content,
        Claims:      claims,
        NeedsVerify: len(claims) > 0,
    }
}
```

---

## 4. 消息标记

```go
func (e *Engine) executeNode(ctx context.Context, nodeID string) error {
    output := agent.Generate(ctx, input)
    
    // 事实校验层
    verifyResult := e.FactVerifier.Analyze(output)
    
    // 添加元数据
    e.StreamChannel <- StreamEvent{
        Event: "message",
        Data: map[string]interface{}{
            "content":      output,
            "needs_verify": verifyResult.NeedsVerify,
            "claims":       verifyResult.Claims,
        },
    }
    
    return nil
}
```

---

## 5. 前端 UI

```tsx
const MessageBubble: FC<{ message: Message }> = ({ message }) => {
  return (
    <div className="relative">
      <ReactMarkdown>{message.content}</ReactMarkdown>
      
      {message.needsVerify && (
        <Badge 
          className="absolute top-2 right-2 bg-yellow-100 text-yellow-800"
          title="此消息包含未验证的事实声明"
        >
          ⚠️ Verify Pending
        </Badge>
      )}
      
      {/* 高亮需验证的声明 */}
      {message.claims?.map((claim, i) => (
        <Tooltip key={i} content={`未验证: ${claim.text}`}>
          <span className="underline decoration-wavy decoration-yellow-500">
            {claim.text}
          </span>
        </Tooltip>
      ))}
    </div>
  );
};
```

---

## 6. 与 FactCheck 节点集成

当消息标记 `needsVerify` 后，下游 FactCheck 节点会优先验证这些声明：

```go
func (f *FactCheckProcessor) Process(ctx context.Context, input map[string]interface{}) {
    claims := input["claims"].([]Claim)
    
    for _, claim := range claims {
        if !claim.Verified {
            result := f.verify(claim.Text)
            claim.Verified = result.Passed
            claim.VerifyNote = result.Source
        }
    }
}
```

---

## 7. 测试要点

- [ ] 数字/金额/引用模式检测
- [ ] Verify Pending 标记显示
- [ ] 与 FactCheck 节点联动
