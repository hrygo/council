# 2.1 流程编排引擎 (Workflow Engine)

这是系统的核心。我们需要实现一个支持**并发**、**流式输出**的 DAG（有向无环图）调度器。

* **设计模式**：生产者-消费者模型 + 递归遍历。
* **并发控制**：`sync.WaitGroup` + `Channels`。

**Go 核心结构体定义：**

```go
// NodeProcessor 定义了所有节点（Start, Agent, Vote）必须实现的接口
type NodeProcessor interface {
    // Process 执行节点逻辑
    // ctx: 上下文
    // input: 上游节点的输出
    // stream: 用于向 WebSocket 推送流式数据的 Channel
    Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (output map[string]interface{}, err error)
}

// Engine 负责解析 Graph 并调度执行
type Engine struct {
    Graph    *GraphDefinition
    Status   map[string]NodeStatus
    mu       sync.RWMutex
}

// Run 启动流程
func (e *Engine) Run(ctx context.Context, startNodeID string) {
    // 1. 拓扑排序检查是否有环 (MVP 可跳过，假设 React Flow 保证无环)
    // 2. 从 Start 节点开始递归执行
    e.executeNode(ctx, startNodeID, nil)
}

func (e *Engine) executeNode(ctx context.Context, nodeID string, input interface{}) {
    node := e.Graph.Nodes[nodeID]
    
    // 如果是并行节点 (Parallel)，查找所有子节点
    if node.Type == "parallel" {
        var wg sync.WaitGroup
        for _, childID := range node.NextIDs {
            wg.Add(1)
            go func(cid string) {
                defer wg.Done()
                e.executeNode(ctx, cid, input) // 递归并发调用
            }(childID)
        }
        wg.Wait()
        return
    }
    
    // 普通节点执行
    processor := NodeFactory(node)
    output, _ := processor.Process(ctx, input, e.StreamChannel)
    
    // 继续执行下游
    for _, nextID := range node.NextIDs {
        e.executeNode(ctx, nextID, output)
    }
}
```

#### 2.1.1 Start 节点处理器 (提案入口)

对应 PRD F.3.1 🟢 Start 节点，处理用户输入和附件。

```go
type StartProcessor struct {
    ProposalText string
    Attachments  []Attachment // PDF/MD 文件
}

type Attachment struct {
    FileName    string
    FilePath    string
    ContentType string // "application/pdf", "text/markdown"
    ParsedText  string // 解析后的文本内容
}

func (s *StartProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    // 1. 推送开始事件
    stream <- StreamEvent{
        Event: "node_state_change",
        Data: map[string]interface{}{"node_id": "start", "status": "running"},
    }
    
    // 2. 解析所有附件
    var parsedContents []string
    for _, att := range s.Attachments {
        content, err := s.parseAttachment(att)
        if err != nil {
            return nil, fmt.Errorf("failed to parse %s: %w", att.FileName, err)
        }
        att.ParsedText = content
        parsedContents = append(parsedContents, content)
    }
    
    // 3. 构造初始上下文
    output := map[string]interface{}{
        "proposal":           s.ProposalText,
        "attachments":        s.Attachments,
        "combined_context":   strings.Join(parsedContents, "\n\n---\n\n"),
        "metadata": map[string]interface{}{
            "started_at":       time.Now(),
            "attachment_count": len(s.Attachments),
        },
    }
    
    stream <- StreamEvent{
        Event: "node_state_change",
        Data: map[string]interface{}{"node_id": "start", "status": "completed"},
    }
    
    return output, nil
}

func (s *StartProcessor) parseAttachment(att Attachment) (string, error) {
    switch att.ContentType {
    case "application/pdf":
        return pdf.ExtractText(att.FilePath) // 使用 pdfcpu 或 unidoc
    case "text/markdown":
        content, err := os.ReadFile(att.FilePath)
        return string(content), err
    default:
        return "", fmt.Errorf("unsupported content type: %s", att.ContentType)
    }
}
```

#### 2.1.2 End 节点处理器 (总结输出)

对应 PRD F.3.1 🔴 End 节点，生成最终报告并触发萃取。

```go
type EndProcessor struct {
    SummaryPrompt string   // 自定义总结提示词
    SessionID     string
    GroupID       string
}

func (e *EndProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    stream <- StreamEvent{
        Event: "node_state_change",
        Data: map[string]interface{}{"node_id": "end", "status": "running"},
    }
    
    // 1. 汇总所有上游输出
    allContent := e.aggregateUpstreamContent(input)
    
    // 2. 调用 LLM 生成结构化总结
    summaryPrompt := e.SummaryPrompt
    if summaryPrompt == "" {
        summaryPrompt = `请对以下会议讨论进行总结，输出格式：
## 核心结论
## 主要分歧点
## 建议行动项`
    }
    
    summary, err := llmClient.Chat(ctx, ChatRequest{
        Messages: []Message{
            {Role: "system", Content: summaryPrompt},
            {Role: "user", Content: allContent},
        },
    })
    if err != nil {
        return nil, fmt.Errorf("failed to generate summary: %w", err)
    }
    
    // 3. 流式输出总结
    stream <- StreamEvent{
        Event: "token_stream",
        Data: map[string]interface{}{
            "node_id": "end",
            "chunk":   summary,
        },
    }
    
    // 4. 异步触发萃取引擎
    go func() {
        extractionEngine.Extract(context.Background(), e.SessionID, e.GroupID)
    }()
    
    stream <- StreamEvent{
        Event: "node_state_change",
        Data: map[string]interface{}{"node_id": "end", "status": "completed"},
    }
    
    return map[string]interface{}{
        "final_report": summary,
        "ended_at":     time.Now(),
    }, nil
}
```

#### 2.1.3 Sequence 节点处理器 (串行执行)

对应 PRD F.3.1 🔶 Sequence 逻辑节点，确保子节点按顺序执行。

```go
type SequenceProcessor struct {
    ChildNodeIDs []string
    engine       *Engine
}

func (s *SequenceProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    var lastOutput = input
    
    for i, nodeID := range s.ChildNodeIDs {
        // 推送当前执行进度
        stream <- StreamEvent{
            Event: "sequence_progress",
            Data: map[string]interface{}{
                "current_step": i + 1,
                "total_steps":  len(s.ChildNodeIDs),
                "node_id":      nodeID,
            },
        }
        
        // 获取节点处理器并执行
        node := s.engine.Graph.Nodes[nodeID]
        processor := NodeFactory(node)
        
        output, err := processor.Process(ctx, lastOutput, stream)
        if err != nil {
            return nil, fmt.Errorf("sequence step %d (%s) failed: %w", i+1, nodeID, err)
        }
        
        // 将当前输出作为下一个节点的输入
        lastOutput = output
    }
    
    return lastOutput, nil
}
```
