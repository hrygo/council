# 2.14 防御性中间件架构 (Defense Middleware Architecture)

> 对应 PRD: [4. 安全与风控]

## 1. 逻辑熔断器 (Logic Circuit Breaker)

为了防止 LLM 在自主运行过程中产生“死循环”、“Token 爆炸”或“无效重复对话”，系统在 **Workflow Engine** 调度层引入熔断机制。

### 1.1 熔断策略 (Breaker Strategy)

熔断器作为 Middleware 运行在每次 `Node Execution` 之前。

| 监控指标               | 阈值 (Threshold)  | 说明                                                  |
| :--------------------- | :---------------- | :---------------------------------------------------- |
| **Recursion Depth**    | > 10              | 单个 Loop 节点的最大循环次数，防止死循环。            |
| **Token Velocity**     | > 1000 tokens/sec | 检测生成的 Token 速率是否异常（通常意味着乱码喷涌）。 |
| **Pattern Repetition** | Similarity > 0.95 | 连续 3 次输出几乎完全相同的内容。                     |
| **Entropy Check**      | < 0.5             | 信息熵过低，表示 Agent 陷入复读机模式。               |

### 1.2 熔断响应 (Hard Stop)

当触发阈值时：

1.  **System Lock**: 立即调用 `CancelFunc` 终止所有相关的 Go Context，停止底层 LLM 请求。
2.  **State Upate**: 将 Session 状态置为 `LOCKED (SUSPENDED)`.
3.  **Persist Log**: 记录 `breaker_events` 表，包含触发原因和现场快照。

```go
type CircuitBreaker struct {
    MaxDepth    int
    MaxVelocity float64
}

func (cb *CircuitBreaker) Check(ctx context.Context, session *Session) error {
    if session.Depth > cb.MaxDepth {
        return fmt.Errorf("recursion depth exceeded: %d", session.Depth)
    }
    // ... 其他检查
    return nil
}
```

## 2. 防幻觉传播 (Anti-Hallucination)

### 2.1 事实校验层 (Fact Verification Layer)

在 Agent 传递信息（Message Passing）的间隙，增加一个轻量级的校验层。

*   **Trigger**: 正则匹配检测到 Agent 输出了 `[Specific Metric]` 或 `[External Citation]`.
*   **Action**: 
    1.  自动调用 `FactCheck` 工具（如 Tavily）进行验证。
    2.  如果验证失败，自动在消息体中注入 `Wait! This claim seems unverified.` 的 System Note。
    3.  UI 前端渲染时，将此类消息标记为 **黄色警示**。

## 3. 实现位置

代码结构上，这些逻辑应实现在 `internal/core/middleware/` 包中，并通过 `Chain of Responsibility` 模式挂载到 Workflow Engine。
