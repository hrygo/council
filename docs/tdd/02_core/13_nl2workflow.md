# 2.13 自然语言转工作流模块 (NL2Workflow)

对应 PRD F.3.3 Wizard Mode，降低交互门槛。

```go
type NL2WorkflowService struct {
    LLM          LLMProvider
    TemplateRepo TemplateRepository
}

type WorkflowGenerationRequest struct {
    Description string  // 用户自然语言描述
    GroupID     string  // 目标群组
    AgentPool   []Agent // 可用 Agent 列表
}

func (s *NL2WorkflowService) Generate(ctx context.Context, req WorkflowGenerationRequest) (*GraphDefinition, error) {
    // 1. 先尝试匹配现有模版
    similarTemplates := s.TemplateRepo.FindSimilar(ctx, req.Description, 3)
    
    prompt := fmt.Sprintf(`你是工作流设计专家。根据用户需求生成 JSON 格式的工作流定义。

用户需求: %s

可参考的相似模版:
%s

可用 Agent:
%s

生成包含以下字段的 JSON（只输出 JSON，不要解释）:
{
  "nodes": [
    {"id": "string", "type": "start|agent|parallel|sequence|vote|loop|fact_check|end", "config": {...}}
  ],
  "edges": [{"source": "string", "target": "string"}],
  "recommended_template_id": "string 或 null (如果与某模版高度相似)"
}`,
        req.Description,
        formatTemplates(similarTemplates),
        formatAgents(req.AgentPool))

    resp, err := s.LLM.Chat(ctx, ChatRequest{
        Messages:       []Message{{Role: "user", Content: prompt}},
        ResponseFormat: &ResponseFormat{Type: "json_object"},
    })
    if err != nil {
        return nil, err
    }
    
    var result struct {
        Nodes                 []NodeDefinition `json:"nodes"`
        Edges                 []EdgeDefinition `json:"edges"`
        RecommendedTemplateID *string          `json:"recommended_template_id"`
    }
    if err := json.Unmarshal([]byte(resp), &result); err != nil {
        return nil, fmt.Errorf("failed to parse workflow JSON: %w", err)
    }
    
    return &GraphDefinition{
        Nodes: result.Nodes,
        Edges: result.Edges,
    }, nil
}
```

**REST API：**

| Method | Endpoint                     | 描述               |
| ------ | ---------------------------- | ------------------ |
| `POST` | `/api/v1/workflows/generate` | 自然语言生成工作流 |
