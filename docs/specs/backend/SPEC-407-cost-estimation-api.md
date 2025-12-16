# SPEC-407: Cost Estimation API

> **优先级**: P1 | **详细设计**: docs/api/cost_estimation.md

---

## 1. 端点

```http
POST /api/v1/workflows/:id/estimate
```

---

## 2. Handler 实现

```go
func (h *WorkflowHandler) Estimate(c *gin.Context) {
    var req EstimateRequest
    c.BindJSON(&req)
    
    workflow := h.WorkflowStore.Get(c.Param("id"))
    
    estimate := h.Estimator.Calculate(workflow, req.Proposal, req.Attachments)
    
    c.JSON(200, estimate)
}
```

---

## 3. 预估算法

```go
type Estimator struct {
    PricingTable map[string]ModelPricing
}

func (e *Estimator) Calculate(workflow *Workflow, proposal string, attachments []Attachment) *CostEstimate {
    baseTokens := e.estimateInputTokens(proposal, attachments)
    
    breakdown := []CostBreakdownItem{}
    totalTokens := 0
    totalCost := 0.0
    
    for _, node := range workflow.Graph.Nodes {
        if node.Type != "agent" {
            continue
        }
        
        agent := workflow.GetAgent(node.AgentID)
        pricing := e.PricingTable[agent.ModelConfig.Model]
        
        nodeTokens := baseTokens + e.estimateOutputTokens(node)
        nodeCost := e.calculateCost(nodeTokens, pricing)
        
        breakdown = append(breakdown, CostBreakdownItem{
            NodeID:          node.ID,
            NodeName:        node.Name,
            AgentName:       agent.Name,
            Model:           agent.ModelConfig.Model,
            EstimatedTokens: nodeTokens,
            EstimatedCost:   nodeCost,
        })
        
        totalTokens += nodeTokens
        totalCost += nodeCost
    }
    
    return &CostEstimate{
        EstimatedTokens: totalTokens,
        EstimatedCost:   totalCost,
        Breakdown:       breakdown,
        Warnings:        e.generateWarnings(totalCost),
    }
}
```

---

## 4. 定价表

```go
var DefaultPricing = map[string]ModelPricing{
    "gpt-4-turbo":       {InputPer1M: 10.00, OutputPer1M: 30.00},
    "gpt-4o":            {InputPer1M: 2.50, OutputPer1M: 10.00},
    "gpt-4o-mini":       {InputPer1M: 0.15, OutputPer1M: 0.60},
    "claude-3.5-sonnet": {InputPer1M: 3.00, OutputPer1M: 15.00},
    "gemini-1.5-pro":    {InputPer1M: 1.25, OutputPer1M: 5.00},
    "deepseek-chat":     {InputPer1M: 0.14, OutputPer1M: 0.28},
}
```

---

## 5. 测试用例

- 正确计算成本
- 警告阈值触发
- 多节点累加
