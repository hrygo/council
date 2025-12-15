# 2.12 FactCheck 节点处理器 (Fact Verification Node)

对应 PRD F.3.1 FactCheck 节点，阻断集体幻觉。

```go
type FactCheckProcessor struct {
    WebSearchTool   SearchTool
    LocalSearchTool SearchTool // 🆕 本地知识库搜索
    LLM             LLMProvider
    Threshold       float64
}

type FactCheckResult struct {
    Claim      string   `json:"claim"`
    Verified   bool     `json:"verified"`
    Confidence float64  `json:"confidence"`
    Sources    []string `json:"sources"`
    Correction string   `json:"correction,omitempty"`
}

func (f *FactCheckProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    content := input["content"].(string)
    
    // 1. 提取可验证的事实断言
    claims := f.extractClaims(ctx, content)
    
    stream <- StreamEvent{
        Event: "fact_check_start",
        Data: map[string]interface{}{"total_claims": len(claims)},
    }
    
    // 2. 逐个验证
    var results []FactCheckResult
    for i, claim := range claims {
        // 混合搜索验证 (Web + Local)
        webResults, _ := f.WebSearchTool.Search(ctx, claim, SearchOptions{MaxResults: 3})
        localResults, _ := f.LocalSearchTool.Search(ctx, claim, SearchOptions{MaxResults: 2})
        searchResults := append(webResults, localResults...)
        
        // LLM 判断
        result := f.verifyClaim(ctx, claim, searchResults)
        results = append(results, result)
        
        stream <- StreamEvent{
            Event: "fact_check_progress",
            Data: map[string]interface{}{
                "current":  i + 1,
                "claim":    claim,
                "verified": result.Verified,
            },
        }
    }
    
    // 3. 计算通过率
    passRate := f.calculatePassRate(results)
    
    return map[string]interface{}{
        "fact_check_results": results,
        "pass_rate":          passRate,
        "overall_passed":     passRate >= f.Threshold,
    }, nil
}

func (f *FactCheckProcessor) extractClaims(ctx context.Context, content string) []string {
    prompt := `从以下文本中提取所有可验证的事实断言（数字、日期、事件、声明）。
每行一个断言，不要编号，只输出断言本身：

` + content

    resp, _ := f.LLM.Chat(ctx, ChatRequest{
        Messages: []Message{{Role: "user", Content: prompt}},
    })
    
    lines := strings.Split(strings.TrimSpace(resp), "\n")
    var claims []string
    for _, line := range lines {
        if line = strings.TrimSpace(line); line != "" {
            claims = append(claims, line)
        }
    }
    return claims
}
```
