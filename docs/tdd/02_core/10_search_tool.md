# 2.10 搜索工具集成模块 (Search Tool Integration)

对应 PRD F.2.3 联网搜索核心必选，集成外部搜索 API。

```go
type SearchTool interface {
    Search(ctx context.Context, query string, opts SearchOptions) ([]SearchResult, error)
}

type SearchOptions struct {
    Depth      string // "basic" | "advanced"
    MaxResults int
    TimeRange  string // "day" | "week" | "month" | "year"
}

type SearchResult struct {
    Title   string `json:"title"`
    URL     string `json:"url"`
    Snippet string `json:"snippet"`
    Score   float64 `json:"score"`
}

// TavilySearchTool 实现 Tavily API 集成
type TavilySearchTool struct {
    APIKey  string
    BaseURL string
    client  *http.Client
}

func NewTavilySearchTool(apiKey string) *TavilySearchTool {
    return &TavilySearchTool{
        APIKey:  apiKey,
        BaseURL: "https://api.tavily.com/v1",
        client:  &http.Client{Timeout: 10 * time.Second},
    }
}

func (t *TavilySearchTool) Search(ctx context.Context, query string, opts SearchOptions) ([]SearchResult, error) {
    req := map[string]interface{}{
        "api_key":        t.APIKey,
        "query":          query,
        "search_depth":   opts.Depth,
        "include_answer": true,
        "max_results":    opts.MaxResults,
    }
    
    body, _ := json.Marshal(req)
    httpReq, _ := http.NewRequestWithContext(ctx, "POST", t.BaseURL+"/search", bytes.NewReader(body))
    httpReq.Header.Set("Content-Type", "application/json")
    
    resp, err := t.client.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("tavily search failed: %w", err)
    }
    defer resp.Body.Close()
    
    var result struct {
        Results []SearchResult `json:"results"`
    }
    json.NewDecoder(resp.Body).Decode(&result)
    
    return result.Results, nil
}

// 搜索工具工厂
func NewSearchTool(provider, apiKey string) SearchTool {
    switch provider {
    case "tavily":
        return NewTavilySearchTool(apiKey)
    case "serper":
        return NewSerperSearchTool(apiKey)
    default:
        return NewDuckDuckGoSearchTool() // 免费备选
    }
}
```
