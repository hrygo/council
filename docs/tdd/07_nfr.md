# 7. 非功能需求实现 (Non-Functional Requirements)

对应 PRD 第4章非功能需求。

### 7.1 安全与隐私 (Security & Privacy)

#### 7.1.1 API Key 安全存储

**后端存储策略**：

```go
import (
    "crypto/aes"
    "crypto/cipher"
    "encoding/base64"
    "os"
)

// 使用环境变量中的主密钥加密 API Key
func EncryptAPIKey(plainKey string) (string, error) {
    masterKey := os.Getenv("COUNCIL_MASTER_KEY") // 32 bytes for AES-256
    block, err := aes.NewCipher([]byte(masterKey))
    if err != nil {
        return "", err
    }
    
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    ciphertext := gcm.Seal(nonce, nonce, []byte(plainKey), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// API Key 存储在 PostgreSQL 中，已加密
type APIKeyRecord struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key"`
    Provider     string    `gorm:"size:32;not null"`  // openai, anthropic, google, deepseek
    EncryptedKey string    `gorm:"type:text;not null"`
    CreatedAt    time.Time
}
```

**前端存储策略**：
- API Key 仅在内存中临时持有，不存储在 localStorage
- 通过 HTTPS 传输至后端加密存储

#### 7.1.2 数据安全隔离

- 所有数据存储在用户自有 Docker PostgreSQL 中，不上传至第三方服务器
- 网络请求仅发生在调用云端 LLM API 时
- 可配置 Proxy 支持企业内网环境

### 7.2 性能指标 (Performance Metrics)

| 指标             | 目标值  | 实现策略                                |
| ---------------- | ------- | --------------------------------------- |
| **首字延迟**     | < 500ms | 使用 WebSocket 流式推送，避免 HTTP 轮询 |
| **3 Agent 并发** | 支持    | Go Goroutines + Channel 调度            |
| **向量检索**     | < 100ms | pgvector IVFFlat 索引，lists=100        |
| **UI 渲染**      | 60fps   | React 虚拟列表 + 增量 DOM 更新          |
| **内存占用**     | < 500MB | 流式处理，避免全量加载对话历史          |

**并发控制实现：**

```go
type ConcurrencyLimiter struct {
    sem chan struct{}
}

func NewLimiter(maxConcurrent int) *ConcurrencyLimiter {
    return &ConcurrencyLimiter{
        sem: make(chan struct{}, maxConcurrent),
    }
}

func (l *ConcurrencyLimiter) Acquire() { l.sem <- struct{}{} }
func (l *ConcurrencyLimiter) Release() { <-l.sem }

// 使用示例：限制最多 3 个 Agent 同时调用 LLM
var agentLimiter = NewLimiter(3)
```

### 7.3 扩展性设计 (Extensibility)

#### 7.3.1 插件预留架构

遵循 OpenAI Function Calling 标准，为未来能力扩展做准备。

```go
type ToolDefinition struct {
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Parameters  json.RawMessage `json:"parameters"` // JSON Schema
}

type ToolCall struct {
    ID       string `json:"id"`
    Type     string `json:"type"` // "function"
    Function struct {
        Name      string `json:"name"`
        Arguments string `json:"arguments"`
    } `json:"function"`
}

// Agent 能力注册表
type CapabilityRegistry struct {
    tools map[string]ToolExecutor
}

type ToolExecutor interface {
    Execute(ctx context.Context, args json.RawMessage) (json.RawMessage, error)
}

// 内置能力（MVP 预留）
func init() {
    registry.Register("web_search", &WebSearchTool{})       // 联网搜索
    registry.Register("code_execution", &CodeSandbox{})     // 代码执行
    registry.Register("file_read", &FileReaderTool{})       // 文件读取
}
```

#### 7.3.2 模型适配器扩展

新增 LLM Provider 只需实现接口：

```go
type LLMProvider interface {
    StreamChat(ctx context.Context, req ChatRequest) (<-chan StreamChunk, error)
    Embed(ctx context.Context, text string) ([]float32, error)
    ListModels(ctx context.Context) ([]ModelInfo, error)
}

// 注册新 Provider 示例
func init() {
    llm.RegisterProvider("deepseek", NewDeepSeekProvider)
    llm.RegisterProvider("google", NewGeminiProvider)
    llm.RegisterProvider("mistral", NewMistralProvider)
}
```
