# 2.2 AI 网关与模型路由 (AI Model Router)

为了实现"模型无关性"，我们需要一个适配器层。

* **设计模式**：策略模式 (Strategy Pattern)。
* **接口定义**：

```go
type LLMProvider interface {
    // StreamChat 生成流式对话
    StreamChat(ctx context.Context, req ChatRequest) (<-chan string, error)
    // Embed 生成向量
    Embed(ctx context.Context, text string) ([]float32, error)
}

// Router 根据配置分发请求 (纯云服务)
func (r *Router) GetProvider(config JSONB) LLMProvider {
    switch config.Provider {
    case "openai":
        return NewOpenAIProvider(config.APIKey, config.Model)
    case "anthropic":
        return NewClaudeProvider(config.APIKey, config.Model)
    case "google":
        return NewGeminiProvider(config.APIKey, config.Model)
    case "deepseek":
        return NewDeepSeekProvider(config.APIKey, config.Model)
    }
    return nil
}
```
