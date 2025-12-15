# 2.5 Agent 工厂模块 (Agent Factory)

对应 PRD F.2.x，实现模型无关的角色配置。

**核心数据结构：**

```go
type Agent struct {
    ID           uuid.UUID      `gorm:"type:uuid;primary_key"`
    Name         string         `gorm:"size:64;not null"`
    Avatar       string         `gorm:"size:256"`
    Description  string         `gorm:"size:512"`            // 职业描述
    PersonaPrompt string        `gorm:"type:text;not null"`  // 人设提示词
    
    // 模型配置 (Model Agnostic)
    ModelConfig  ModelConfig    `gorm:"serializer:json"`
    
    // 能力开关 (MVP 预留)
    Capabilities Capabilities   `gorm:"serializer:json"`
    
    CreatedAt    time.Time
    UpdatedAt    time.Time
}

type ModelConfig struct {
    Provider    string  `json:"provider"`    // openai, anthropic, ollama, deepseek, google
    Model       string  `json:"model"`       // gpt-4o, claude-3.5-sonnet, llama3, etc.
    Temperature float64 `json:"temperature"` // 0.0 - 2.0
    TopP        float64 `json:"top_p"`       // 0.0 - 1.0
    MaxTokens   int     `json:"max_tokens"`  // 输出上限
}

type Capabilities struct {
    WebSearch      bool   `json:"web_search"`       // MVP: 核心必选，默认 true
    SearchProvider string `json:"search_provider"`  // "tavily" | "serper" | "duckduckgo"
    CodeExecution  bool   `json:"code_execution"`   // Phase 2 启用，默认 false
}
```

**REST API 设计：**

| Method   | Endpoint                       | 描述                |
| -------- | ------------------------------ | ------------------- |
| `GET`    | `/api/v1/agents`               | 获取 Agent 库列表   |
| `POST`   | `/api/v1/agents`               | 创建 Agent          |
| `GET`    | `/api/v1/agents/:id`           | 获取 Agent 详情     |
| `PUT`    | `/api/v1/agents/:id`           | 更新 Agent 配置     |
| `DELETE` | `/api/v1/agents/:id`           | 删除 Agent          |
| `POST`   | `/api/v1/agents/:id/duplicate` | 复制 Agent 为新角色 |
