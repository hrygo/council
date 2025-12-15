# 2.11 并发执行配置 (Parallel Execution)

纯云服务架构，默认使用并发模式以获得最佳性能。

```go
type EngineConfig struct {
    MaxConcurrent int // 最大并发 Agent 数，默认 3
}

// 引擎并发执行逻辑
func (e *Engine) executeParallelNode(ctx context.Context, node *Node, input interface{}) {
    // 使用信号量限制并发数
    sem := make(chan struct{}, e.Config.MaxConcurrent)
    var wg sync.WaitGroup
    
    for _, childID := range node.NextIDs {
        wg.Add(1)
        sem <- struct{}{} // 获取信号量
        
        go func(cid string) {
            defer wg.Done()
            defer func() { <-sem }() // 释放信号量
            e.executeNode(ctx, cid, input)
        }(childID)
    }
    
    wg.Wait()
}
```
