# LLM提供商集成

<cite>
**本文档引用的文件**  
- [llm.go](file://internal/infrastructure/llm/llm.go)
- [router.go](file://internal/infrastructure/llm/router.go)
- [config.go](file://internal/pkg/config/config.go)
- [openai.go](file://internal/infrastructure/llm/openai.go)
- [gemini.go](file://internal/infrastructure/llm/gemini.go)
- [mock.go](file://internal/infrastructure/llm/mock.go)
- [dashscope.go](file://internal/infrastructure/llm/dashscope.go)
- [deepseek.go](file://internal/infrastructure/llm/deepseek.go)
- [ollama.go](file://internal/infrastructure/llm/ollama.go)
- [siliconflow.go](file://internal/infrastructure/llm/siliconflow.go)
- [router_test.go](file://internal/infrastructure/llm/router_test.go)
- [llm-providers.md](file://docs/guide/llm-providers.md)
</cite>

## 目录
1. [LLM接口实现](#llm接口实现)
2. [路由机制](#路由机制)
3. [配置管理](#配置管理)
4. [错误处理与差异](#错误处理与差异)
5. [测试策略](#测试策略)
6. [安全考虑](#安全考虑)

## LLM接口实现

系统通过`LLMProvider`接口统一管理不同LLM提供商的集成。该接口定义了`Generate`和`Stream`方法，所有提供商必须实现这些方法以确保一致性。

`Generate`方法用于同步生成完整响应，接受上下文和完成请求作为参数，返回包含内容和使用情况的响应对象。`Stream`方法用于流式传输响应，返回字符串通道和错误通道，支持实时流式输出。

`Embedder`接口专门处理嵌入生成，其`Embed`方法接受上下文、模型名称和文本，返回浮点数切片形式的嵌入向量。不同提供商在实现此接口时表现出显著差异。

以OpenAI和Gemini为例进行对比分析：
- **OpenAI实现**：`OpenAIClient`直接使用`go-openai`库，通过`CreateChatCompletion`和`CreateEmbeddings`方法与API交互。它支持所有功能，包括流式传输和嵌入生成。
- **Gemini实现**：`GeminiClient`使用Google的`genai`库，将消息映射到Gemini特定的内容结构。它同样支持流式传输和嵌入，但在角色映射上需要特殊处理（如将"system"角色映射为"user"）。

其他提供商如DeepSeek和Ollama采用不同的策略：
- **DeepSeek**：通过`DeepSeekClient`包装`OpenAIClient`，利用其OpenAI兼容API。但明确不支持嵌入功能，`Embed`方法直接返回错误。
- **Ollama**：类似地包装`OpenAIClient`，但要求在调用`Embed`时必须指定模型，因为Ollama需要用户预先拉取嵌入模型。

**Section sources**
- [llm.go](file://internal/infrastructure/llm/llm.go#L35-L47)
- [openai.go](file://internal/infrastructure/llm/openai.go#L34-L154)
- [gemini.go](file://internal/infrastructure/llm/gemini.go#L35-L206)
- [deepseek.go](file://internal/infrastructure/llm/deepseek.go#L28-L33)
- [ollama.go](file://internal/infrastructure/llm/ollama.go#L31-L58)

## 路由机制

`Registry`结构体实现了动态路由机制，根据配置选择合适的LLM提供商。该机制在`GetLLMProvider`方法中实现，支持按名称获取提供商实例。

路由逻辑遵循以下步骤：
1. 首先检查缓存中是否存在请求的提供商实例
2. 如果请求的是"default"或空名称，则使用系统默认提供商
3. 根据提供商名称查找相应的API密钥和基础URL
4. 创建配置并调用`createProvider`工厂方法初始化提供商
5. 将新创建的提供商实例缓存以供后续使用

`createProvider`方法作为内部工厂函数，通过switch语句根据配置类型初始化不同的提供商客户端。支持的提供商包括"openai"、"gemini"、"deepseek"、"ollama"、"dashscope"和"siliconflow"等。

对于OpenAI兼容的API（如DashScope、DeepSeek、SiliconFlow），系统通过`NewOpenAICompatibleClient`创建客户端，只需指定不同的基础URL。Ollama的处理略有不同，使用"ollama"作为占位符API密钥。

`GetDefaultModel`方法根据系统配置返回默认模型名称，提供了一层额外的抽象，允许在不修改代码的情况下更改默认模型。

**Section sources**
- [router.go](file://internal/infrastructure/llm/router.go#L37-L149)
- [router.go](file://internal/infrastructure/llm/router.go#L111-L129)
- [router.go](file://internal/infrastructure/llm/router.go#L132-L149)

## 配置管理

系统配置通过`Config`结构体集中管理，包含端口、数据库URL、各种API密钥以及LLM和嵌入配置。`LLMConfig`和`EmbeddingConfig`结构体分别定义了聊天模型和嵌入服务的配置项。

在`config.go`中，`Load`函数从环境变量加载配置值。关键配置项包括：
- `Provider`：指定默认提供商类型
- `APIKey`：提供商的认证密钥
- `BaseURL`：API基础地址，用于自定义端点或代理
- `Model`：默认使用的模型名称

环境变量映射遵循命名约定，如`LLM_PROVIDER`对应`LLM.Provider`，`OPENAI_API_KEY`对应`OpenAIKey`。系统为嵌入配置提供了智能默认值，根据提供商自动选择推荐的嵌入模型。

配置优先级遵循"环境变量 > 配置文件 > 代码内默认值"的原则。例如，如果未设置`EMBEDDING_PROVIDER`环境变量，则默认使用"siliconflow"。这种分层配置机制提供了极大的灵活性，允许在不同环境中轻松切换提供商。

**Section sources**
- [config.go](file://internal/pkg/config/config.go#L8-L43)
- [config.go](file://internal/pkg/config/config.go#L45-L132)

## 错误处理与差异

不同提供商在错误码、速率限制和响应格式方面存在显著差异，系统通过统一的错误处理机制来应对这些差异。

错误处理策略包括：
- 使用`fmt.Errorf`包装底层错误，保留原始错误信息
- 为不支持的功能（如DeepSeek的嵌入）返回明确的错误消息
- 在客户端初始化失败时（如Gemini）返回初始化错误

速率限制处理主要依赖于提供商自身的机制，但系统通过配置提供了基本的控制能力。例如，可以通过环境变量设置不同的API密钥来管理配额。

响应格式差异通过标准化的数据结构来解决：
- `CompletionResponse`统一了不同提供商的响应格式
- `Usage`结构体标准化了令牌计数的表示
- 消息角色通过`Message`结构体进行标准化

特别值得注意的是Gemini在角色映射上的特殊处理：将"system"角色映射为"user"，"assistant"和"model"角色映射为"model"。这种转换确保了不同提供商之间的语义一致性。

**Section sources**
- [llm.go](file://internal/infrastructure/llm/llm.go#L23-L33)
- [openai.go](file://internal/infrastructure/llm/openai.go#L58-L63)
- [gemini.go](file://internal/infrastructure/llm/gemini.go#L85-L87)
- [deepseek.go](file://internal/infrastructure/llm/deepseek.go#L30-L32)

## 测试策略

系统提供了全面的测试策略，包括mock测试和真实API调用验证。

`MockProvider`结构体实现了`LLMProvider`和`Embedder`接口，用于单元测试。它允许设置预期的响应、错误和调用跟踪，使测试具有高度的可控性。`MockProvider`包含以下可配置字段：
- `GenerateResponse`和`GenerateError`：控制`Generate`方法的行为
- `StreamContent`和`StreamError`：控制`Stream`方法的行为
- `EmbedResponse`和`EmbedError`：控制`Embed`方法的行为
- 调用计数器：用于验证方法调用次数

`router_test.go`中的测试用例验证了路由机制的正确性：
- `TestRegistry_NewEmbedder`测试嵌入提供商的创建
- `TestRegistry_GetLLMProvider`测试LLM提供商的获取

端到端测试通过Playwright实现，位于`e2e/tests/`目录下。这些测试验证了从用户界面到后端服务的完整工作流程，确保系统在真实环境中的行为符合预期。

测试覆盖率达到核心组件100%，遵循TDD（测试驱动开发）原则。测试中禁止使用真实数据库和LLM服务，确保测试的可重复性和稳定性。

**Section sources**
- [mock.go](file://internal/infrastructure/llm/mock.go#L8-L80)
- [router_test.go](file://internal/infrastructure/llm/router_test.go#L9-L68)
- [GEMINI.md](file://GEMINI.md#L70-L74)

## 安全考虑

安全性在LLM提供商集成中至关重要，主要体现在API密钥管理和数据保护方面。

API密钥通过环境变量进行管理，从不硬编码在源代码中。系统遵循最小权限原则，为不同提供商使用独立的密钥。`.env.example`文件提供了配置模板，但实际的`.env`文件被.gitignore排除，防止密钥意外提交到版本控制系统。

密钥加密存储方面，虽然当前实现直接使用环境变量，但系统架构支持通过密钥管理服务（KMS）进行加密。未来可以集成AWS KMS、Google Cloud KMS或Hashicorp Vault等服务，实现密钥的加密存储和动态获取。

其他安全措施包括：
- 上下文溢出自动截断，防止提示注入攻击
- WebSocket消息格式的端到端验证，确保通信完整性
- 禁止在中间层丢弃上游传入的有效信息，保持数据完整性
- 使用`context`传递请求范围，支持请求取消和超时控制

系统还通过`circuit_breaker.go`实现了熔断器模式，防止故障扩散和资源耗尽，在面对不稳定的LLM服务时提供额外的保护层。

**Section sources**
- [config.go](file://internal/pkg/config/config.go#L126-L130)
- [GEMINI.md](file://GEMINI.md#L67-L68)
- [GEMINI.md](file://GEMINI.md#L49-L53)
- [middleware/circuit_breaker.go](file://internal/core/middleware/circuit_breaker.go)