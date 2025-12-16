# Cost Estimation API 设计文档

> **状态**: 待实现  
> **优先级**: Sprint 4  
> **前端依赖**: CostEstimator 成本预估面板

---

## 概述

成本预估 API 用于在用户启动会议前，根据选定的模型配置和工作流复杂度预估本次会议的 Token 消耗和费用。

对应 PRD **F.4.4 成本预估模块**：
> 在点击"开始会议"前，根据选定模型和流程复杂度显示预估成本：
> 示例："本次会议预计消耗 ~$0.35，耗时 ~2 分钟"

---

## 数据模型

### CostEstimate

```typescript
interface CostEstimate {
  workflow_id: string;
  estimated_tokens: {
    input: number;    // 预估输入 Token
    output: number;   // 预估输出 Token
    total: number;
  };
  estimated_cost: {
    amount: number;   // 预估费用 (美元)
    currency: 'USD';
  };
  estimated_duration: {
    seconds: number;  // 预估耗时 (秒)
    formatted: string; // e.g. "~2 分钟"
  };
  breakdown: CostBreakdownItem[];  // 分项明细
  warnings: CostWarning[];         // 警告信息
}

interface CostBreakdownItem {
  node_id: string;
  node_name: string;
  agent_name?: string;
  model: string;
  estimated_tokens: number;
  estimated_cost: number;
}

interface CostWarning {
  type: 'high_cost' | 'long_duration' | 'complex_workflow';
  message: string;
  suggestion?: string;
}
```

---

## API 端点

### 1. 预估工作流成本

```http
POST /api/v1/workflows/:id/estimate
```

**Request Body:**

```json
{
  "proposal": "请分析这份商业计划书的可行性",  // 可选，提案内容
  "attachments": [                             // 可选，附件信息
    {
      "filename": "business_plan.pdf",
      "size_bytes": 102400
    }
  ],
  "agent_overrides": {                         // 可选，临时覆盖 Agent 模型配置
    "agent-cfo": {
      "model": "gpt-4o-mini"                   // 使用更便宜的模型
    }
  }
}
```

**Response 200:**

```json
{
  "workflow_id": "wf-001",
  "estimated_tokens": {
    "input": 12500,
    "output": 8000,
    "total": 20500
  },
  "estimated_cost": {
    "amount": 0.35,
    "currency": "USD"
  },
  "estimated_duration": {
    "seconds": 120,
    "formatted": "~2 分钟"
  },
  "breakdown": [
    {
      "node_id": "node-ceo",
      "node_name": "CEO 视角分析",
      "agent_name": "CEO",
      "model": "gpt-4-turbo",
      "estimated_tokens": 5000,
      "estimated_cost": 0.15
    },
    {
      "node_id": "node-cfo",
      "node_name": "CFO 财务审查",
      "agent_name": "CFO",
      "model": "claude-3.5-sonnet",
      "estimated_tokens": 4500,
      "estimated_cost": 0.12
    },
    {
      "node_id": "node-factcheck",
      "node_name": "事实核查",
      "agent_name": null,
      "model": "tavily-search",
      "estimated_tokens": 1000,
      "estimated_cost": 0.02
    }
  ],
  "warnings": [
    {
      "type": "high_cost",
      "message": "本次会议预估成本超过 $0.30",
      "suggestion": "可考虑将部分 Agent 切换为更经济的模型 (如 gpt-4o-mini)"
    }
  ]
}
```

---

### 2. 获取模型定价信息

```http
GET /api/v1/models/pricing
```

**Response 200:**

```json
{
  "models": [
    {
      "provider": "openai",
      "model": "gpt-4-turbo",
      "input_price_per_1m": 10.00,  // $10 / 1M tokens
      "output_price_per_1m": 30.00,
      "context_window": 128000
    },
    {
      "provider": "openai",
      "model": "gpt-4o-mini",
      "input_price_per_1m": 0.15,
      "output_price_per_1m": 0.60,
      "context_window": 128000
    },
    {
      "provider": "anthropic",
      "model": "claude-3.5-sonnet",
      "input_price_per_1m": 3.00,
      "output_price_per_1m": 15.00,
      "context_window": 200000
    },
    {
      "provider": "google",
      "model": "gemini-1.5-pro",
      "input_price_per_1m": 1.25,
      "output_price_per_1m": 5.00,
      "context_window": 2000000
    },
    {
      "provider": "deepseek",
      "model": "deepseek-chat",
      "input_price_per_1m": 0.14,
      "output_price_per_1m": 0.28,
      "context_window": 64000
    }
  ],
  "last_updated": "2024-12-01T00:00:00Z"
}
```

---

## 预估算法

### Token 预估逻辑

```go
func EstimateTokens(workflow *GraphDefinition, proposal string, attachments []Attachment) int {
    baseTokens := 0
    
    // 1. 提案内容
    baseTokens += EstimateTokenCount(proposal)
    
    // 2. 附件内容 (假设 PDF 每页约 500 tokens)
    for _, att := range attachments {
        baseTokens += EstimateAttachmentTokens(att)
    }
    
    // 3. 遍历节点估算
    nodeMultiplier := 1.0
    for _, node := range workflow.Nodes {
        switch node.Type {
        case "agent":
            nodeMultiplier += 0.8  // 每个 Agent 节点约增加 80% 输入
        case "parallel":
            nodeMultiplier += 0.5 * float64(len(node.NextIDs))  // 并行分支
        case "loop":
            rounds := node.Properties["max_rounds"].(int)
            nodeMultiplier *= float64(rounds)  // 循环倍增
        }
    }
    
    return int(float64(baseTokens) * nodeMultiplier)
}
```

### 费用计算

```go
func CalculateCost(tokens int, model string) float64 {
    pricing := GetModelPricing(model)
    inputTokens := tokens * 0.6   // 假设 60% 输入
    outputTokens := tokens * 0.4  // 假设 40% 输出
    
    inputCost := float64(inputTokens) / 1_000_000 * pricing.InputPricePer1M
    outputCost := float64(outputTokens) / 1_000_000 * pricing.OutputPricePer1M
    
    return inputCost + outputCost
}
```

---

## 前端集成示例

```tsx
// CostEstimator.tsx
const CostEstimator: FC<{ workflowId: string }> = ({ workflowId }) => {
  const [estimate, setEstimate] = useState<CostEstimate | null>(null);
  const [loading, setLoading] = useState(false);
  
  const fetchEstimate = async () => {
    setLoading(true);
    const res = await fetch(`/api/v1/workflows/${workflowId}/estimate`, {
      method: 'POST',
      body: JSON.stringify({ proposal: currentProposal }),
    });
    setEstimate(await res.json());
    setLoading(false);
  };
  
  useEffect(() => { fetchEstimate(); }, [workflowId]);
  
  if (loading) return <Spinner />;
  if (!estimate) return null;
  
  return (
    <div className="p-4 bg-gray-50 rounded-lg">
      <div className="flex items-center gap-2 text-lg font-medium">
        <DollarSign size={20} />
        <span>预估成本: ${estimate.estimated_cost.amount.toFixed(2)}</span>
      </div>
      <div className="flex items-center gap-2 text-gray-600">
        <Clock size={16} />
        <span>预估耗时: {estimate.estimated_duration.formatted}</span>
      </div>
      {estimate.warnings.map((w, i) => (
        <Alert key={i} type="warning">{w.message}</Alert>
      ))}
    </div>
  );
};
```
