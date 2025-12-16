# SPEC-409: 逻辑熔断 (Circuit Breaker)

> **优先级**: P0 | **预估工时**: 4h  
> **关联 PRD**: F.6.1 逻辑熔断 | **关联 TDD**: 02_core/14_defense_mechanisms.md

---

## 1. 触发条件

| 条件           | 阈值                    | 说明     |
| -------------- | ----------------------- | -------- |
| Token 消耗激增 | > 3x 预估               | 成本失控 |
| 死循环检测     | 连续 3 轮对话熵值 < 0.1 | 无效重复 |
| 执行超时       | > 10 分钟               | 异常阻塞 |

---

## 2. 状态机

```
┌─────────┐     触发熔断      ┌─────────────────┐
│ RUNNING │ ─────────────────▶│ SUSPENDED_LOCKED│
└─────────┘                   └─────────────────┘
                                       │
                                       │ 用户解锁
                                       ▼
                              ┌─────────────────┐
                              │    RESUMED      │
                              └─────────────────┘
```

---

## 3. 实现

```go
type CircuitBreaker struct {
    Session        *Session
    TokenThreshold float64
    EntropyWindow  int
    Timeout        time.Duration
}

type BreakerStatus string

const (
    StatusOpen     BreakerStatus = "SUSPENDED_LOCKED"
    StatusClosed   BreakerStatus = "RUNNING"
    StatusRecovery BreakerStatus = "PENDING_UNLOCK"
)

func (cb *CircuitBreaker) Monitor(ctx context.Context, events <-chan StreamEvent) {
    entropyHistory := []float64{}
    startTime := time.Now()
    estimatedTokens := cb.Session.EstimatedTokens
    actualTokens := 0

    for event := range events {
        switch event.Event {
        case "token_usage":
            data := event.Data.(TokenUsageData)
            actualTokens += data.InputTokens + data.OutputTokens
            
            // 检查 Token 激增
            if float64(actualTokens) > float64(estimatedTokens)*cb.TokenThreshold {
                cb.Trip("TOKEN_SURGE", "Token 消耗超过预估 3 倍")
                return
            }

        case "message":
            entropy := cb.calculateEntropy(event.Data.(string))
            entropyHistory = append(entropyHistory, entropy)
            
            // 检查死循环
            if len(entropyHistory) >= cb.EntropyWindow {
                recent := entropyHistory[len(entropyHistory)-cb.EntropyWindow:]
                if cb.avgEntropy(recent) < 0.1 {
                    cb.Trip("LOOP_DETECTED", "检测到无效循环对话")
                    return
                }
            }
        }

        // 检查超时
        if time.Since(startTime) > cb.Timeout {
            cb.Trip("TIMEOUT", "执行超时")
            return
        }
    }
}

func (cb *CircuitBreaker) Trip(reason, message string) {
    cb.Session.Status = StatusOpen
    cb.Session.LockReason = reason
    cb.Session.LockMessage = message
    
    // 通知前端
    cb.Session.Stream <- StreamEvent{
        Event: "circuit_breaker:tripped",
        Data: map[string]string{
            "reason":  reason,
            "message": message,
        },
    }
}
```

---

## 4. 恢复流程

```go
type UnlockRequest struct {
    SessionID         string
    RiskJustification string   // 风险陈述
    SafetyChecks      []bool   // 三项安全自查
}

func (cb *CircuitBreaker) Unlock(req UnlockRequest) error {
    // 验证解锁条件
    if req.RiskJustification == "" && !allTrue(req.SafetyChecks) {
        return ErrUnlockConditionNotMet
    }
    
    cb.Session.Status = StatusClosed
    cb.Session.LockReason = ""
    return nil
}
```

---

## 5. 前端 UI

```tsx
const CircuitBreakerModal: FC = () => {
  const { lockedSession } = useSessionStore();
  const [justification, setJustification] = useState('');
  const [checks, setChecks] = useState([false, false, false]);

  if (!lockedSession) return null;

  return (
    <Dialog open className="bg-gray-900">
      {/* 灰阶模式 */}
      <DialogContent className="grayscale-[50%]">
        <DialogHeader>
          <DialogTitle className="text-red-500 flex items-center gap-2">
            🚨 系统已锁定
          </DialogTitle>
        </DialogHeader>
        
        <Alert variant="destructive">
          <p>触发原因: {lockedSession.lockReason}</p>
          <p>{lockedSession.lockMessage}</p>
        </Alert>
        
        <div className="space-y-4">
          <Textarea
            label="风险陈述 (可选)"
            placeholder="请说明您理解的风险并确认继续..."
            value={justification}
            onChange={e => setJustification(e.target.value)}
          />
          
          <div className="space-y-2">
            <p className="font-medium">或完成安全自查：</p>
            <Checkbox
              checked={checks[0]}
              onChange={v => setChecks([v, checks[1], checks[2]])}
            >
              我已确认当前对话内容符合预期
            </Checkbox>
            <Checkbox
              checked={checks[1]}
              onChange={v => setChecks([checks[0], v, checks[2]])}
            >
              我理解继续可能产生额外费用
            </Checkbox>
            <Checkbox
              checked={checks[2]}
              onChange={v => setChecks([checks[0], checks[1], v])}
            >
              我接受后续结果的风险
            </Checkbox>
          </div>
        </div>
        
        <DialogFooter>
          <Button variant="destructive" onClick={handleUnlock}>
            解锁并继续
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
```

---

## 6. 测试要点

- [ ] Token 激增触发
- [ ] 死循环检测
- [ ] 超时触发
- [ ] 解锁流程
