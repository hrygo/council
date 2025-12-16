# SPEC-405: HumanReviewProcessor

> **优先级**: P0 | **关联 PRD**: F.3.1 HumanReview | **API**: human_review.md

---

## 1. 数据结构

```go
type HumanReviewProcessor struct {
    SessionID     string
    NodeID        string
    ReviewStore   ReviewStore
    Timeout       time.Duration
    AllowSkip     bool
}
```

---

## 2. 实现逻辑

```go
func (h *HumanReviewProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    // 1. 生成决策草案
    draft := h.generateDraft(input)
    
    // 2. 创建裁决请求
    review := &HumanReviewRequest{
        ID:        uuid.New().String(),
        SessionID: h.SessionID,
        NodeID:    h.NodeID,
        Status:    "pending",
        Draft:     draft,
        ExpiresAt: time.Now().Add(h.Timeout),
    }
    h.ReviewStore.Save(review)
    
    // 3. 通知前端
    stream <- StreamEvent{
        Event: "human_review:requested",
        Data:  review,
    }
    
    // 4. 等待决策 (带超时)
    signalCh := h.Session.GetSignalChannel(h.NodeID)
    
    select {
    case decision := <-signalCh:
        return h.handleDecision(decision, draft)
    case <-time.After(h.Timeout):
        if h.AllowSkip {
            return map[string]interface{}{"auto_approved": true, "content": draft}, nil
        }
        return nil, fmt.Errorf("human review timeout")
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

func (h *HumanReviewProcessor) handleDecision(decision interface{}, draft interface{}) (map[string]interface{}, error) {
    d := decision.(*HumanDecision)
    switch d.Action {
    case "approve":
        return map[string]interface{}{"approved": true, "content": draft}, nil
    case "modify":
        return map[string]interface{}{"approved": true, "content": d.ModifiedContent}, nil
    case "reject":
        return nil, fmt.Errorf("rejected: %s", d.RejectionReason)
    default:
        return nil, fmt.Errorf("unknown action: %s", d.Action)
    }
}
```

---

## 3. API 集成

当收到 `POST /api/v1/sessions/:id/review` 时：

```go
func (h *WorkflowHandler) SubmitReview(c *gin.Context) {
    var req ReviewSubmission
    c.BindJSON(&req)
    
    session := h.SessionManager.Get(c.Param("id"))
    session.SendSignal(req.NodeID, &HumanDecision{
        Action:          req.Action,
        ModifiedContent: req.ModifiedContent,
        RejectionReason: req.RejectionReason,
    })
    
    c.JSON(200, gin.H{"status": "accepted"})
}
```

---

## 4. 测试用例

- 审批通过
- 修改后通过
- 驳回
- 超时处理
