# SPEC-411: 联网搜索集成 (Search Integration)

> **优先级**: P1 | **预估工时**: 3h  
> **关联 PRD**: F.2.3 联网搜索 | **关联 TDD**: 02_core/12_fact_check.md

---

## 1. 概述

集成 Tavily/Serper API 进行事实核查和信息检索。

---

## 2. 接口抽象

```go
type SearchClient interface {
    Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error)
}

type SearchOptions struct {
    MaxResults int
    SearchType string // "search" | "answer"
    Domains    []string
}

type SearchResult struct {
    Query   string
    Results []SearchItem
    Answer  string // Tavily 的直接回答
}

type SearchItem struct {
    Title   string
    URL     string
    Content string
    Score   float64
}
```

---

## 3. Tavily 实现

```go
type TavilyClient struct {
    APIKey  string
    BaseURL string
}

func (t *TavilyClient) Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
    req := map[string]interface{}{
        "api_key":           t.APIKey,
        "query":             query,
        "search_depth":      "advanced",
        "include_answer":    true,
        "max_results":       opts.MaxResults,
    }
    
    resp, err := http.Post(t.BaseURL+"/search", "application/json", toJSON(req))
    if err != nil {
        return nil, fmt.Errorf("tavily search failed: %w", err)
    }
    
    var result TavilyResponse
    json.NewDecoder(resp.Body).Decode(&result)
    
    return &SearchResult{
        Query:   query,
        Answer:  result.Answer,
        Results: mapResults(result.Results),
    }, nil
}
```

---

## 4. Serper 实现

```go
type SerperClient struct {
    APIKey string
}

func (s *SerperClient) Search(ctx context.Context, query string, opts SearchOptions) (*SearchResult, error) {
    req := map[string]interface{}{
        "q":   query,
        "num": opts.MaxResults,
    }
    
    httpReq, _ := http.NewRequest("POST", "https://google.serper.dev/search", toJSON(req))
    httpReq.Header.Set("X-API-KEY", s.APIKey)
    
    resp, err := http.DefaultClient.Do(httpReq)
    // ...
}
```

---

## 5. 工厂模式

```go
func NewSearchClient(provider string) (SearchClient, error) {
    switch provider {
    case "tavily":
        return &TavilyClient{
            APIKey:  os.Getenv("TAVILY_API_KEY"),
            BaseURL: "https://api.tavily.com",
        }, nil
    case "serper":
        return &SerperClient{
            APIKey: os.Getenv("SERPER_API_KEY"),
        }, nil
    default:
        return nil, fmt.Errorf("unknown search provider: %s", provider)
    }
}
```

---

## 6. FactCheck 节点集成

```go
type FactCheckProcessor struct {
    SearchClient SearchClient
    // ...
}

func (f *FactCheckProcessor) verify(claim string) *VerifyResult {
    result, err := f.SearchClient.Search(ctx, claim, SearchOptions{
        MaxResults: 3,
        SearchType: "answer",
    })
    if err != nil {
        return &VerifyResult{Passed: false, Error: err.Error()}
    }
    
    // 分析搜索结果判断声明是否属实
    return f.analyzeResults(claim, result)
}
```

---

## 7. 配置

```yaml
# config.yaml
search:
  provider: tavily  # tavily | serper
  max_queries_per_session: 10
  timeout_seconds: 5
```

---

## 8. 测试要点

- [ ] Tavily API 调用
- [ ] Serper API 调用
- [ ] 错误处理 (API 限流)
- [ ] 结果解析
