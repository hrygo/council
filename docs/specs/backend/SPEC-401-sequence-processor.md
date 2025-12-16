# SPEC-401: SequenceProcessor

> **优先级**: P1 | **关联 PRD**: F.3.1 Sequence 节点

---

## 1. 职责

按顺序依次执行子节点，前一个节点输出作为后一个节点输入。

---

## 2. 接口

```go
type SequenceProcessor struct {
    ChildNodeIDs []string
}

func (s *SequenceProcessor) Process(
    ctx context.Context, 
    input map[string]interface{}, 
    stream chan<- StreamEvent,
) (map[string]interface{}, error)
```

---

## 3. 实现逻辑

```go
func (s *SequenceProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    currentInput := input
    
    for i, childID := range s.ChildNodeIDs {
        stream <- StreamEvent{
            Event: "sequence:step",
            Data: map[string]interface{}{
                "step":    i + 1,
                "total":   len(s.ChildNodeIDs),
                "node_id": childID,
            },
        }
        
        processor, err := s.NodeFactory(childID)
        if err != nil {
            return nil, fmt.Errorf("failed to create processor for %s: %w", childID, err)
        }
        
        output, err := processor.Process(ctx, currentInput, stream)
        if err != nil {
            return nil, fmt.Errorf("node %s failed: %w", childID, err)
        }
        
        currentInput = output
    }
    
    return currentInput, nil
}
```

---

## 4. 测试用例

```go
func TestSequenceProcessor(t *testing.T) {
    // 测试顺序执行
    // 测试错误传播
    // 测试上下文取消
}
```
