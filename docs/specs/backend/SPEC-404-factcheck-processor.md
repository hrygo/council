# SPEC-404: FactCheckProcessor

> **优先级**: P1 | **关联 PRD**: F.3.1 FactCheck 节点

---

## 1. 数据结构

```go
type FactCheckProcessor struct {
    SearchSources   []string  // ["tavily", "serper", "local_kb"]
    MaxQueries      int
    VerifyThreshold float64
    SearchClient    SearchClient
}

type FactCheckResult struct {
    Claim       string
    Verified    bool
    Confidence  float64
    Sources     []Source
    Explanation string
}

type Source struct {
    Title string
    URL   string
    Quote string
}
```

---

## 2. 实现逻辑

```go
func (f *FactCheckProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    // 1. 从输入中提取需要核查的声明
    claims := f.extractClaims(input)
    
    results := []FactCheckResult{}
    
    for _, claim := range claims {
        stream <- StreamEvent{
            Event: "factcheck:checking",
            Data:  map[string]interface{}{"claim": claim},
        }
        
        // 2. 搜索验证
        searchResults := f.search(ctx, claim)
        
        // 3. 分析置信度
        result := f.analyze(claim, searchResults)
        results = append(results, result)
        
        stream <- StreamEvent{
            Event: "factcheck:result",
            Data:  result,
        }
    }
    
    // 4. 汇总
    allVerified := true
    for _, r := range results {
        if !r.Verified {
            allVerified = false
            break
        }
    }
    
    return map[string]interface{}{
        "verified": allVerified,
        "results":  results,
    }, nil
}
```

---

## 3. 声明提取

```go
func (f *FactCheckProcessor) extractClaims(input map[string]interface{}) []string {
    // 使用 LLM 从输入中提取可验证的事实性声明
    prompt := `从以下内容中提取所有可验证的事实性声明，每行一条：
    %s`
    // ...
}
```

---

## 4. 测试用例

- 正确提取声明
- 搜索 API 调用
- 置信度计算
- 低置信度标记
