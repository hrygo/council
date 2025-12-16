# Human Review API 设计文档

> **状态**: 待实现  
> **优先级**: Sprint 4 (P0 关键路径)  
> **前端依赖**: HumanReviewModal 组件

---

## 概述

人类裁决 API 用于支持 **PRD F.3.1 HumanReview 节点**。当工作流执行到人类裁决节点时：

1. 后端暂停执行，生成决策草案
2. 通过 WebSocket 通知前端显示裁决弹窗
3. 用户审核并提交决定（通过/修改/驳回）
4. 后端恢复执行或终止流程

---

## 数据模型

### HumanReviewRequest

```typescript
interface HumanReviewRequest {
  id: string;                      // 裁决请求 ID
  session_id: string;              // 会话 ID
  node_id: string;                 // 触发裁决的节点 ID
  created_at: string;              // 请求创建时间
  expires_at?: string;             // 可选超时时间
  status: 'pending' | 'approved' | 'rejected' | 'expired';
  
  // 决策草案
  draft: {
    title: string;                 // 草案标题
    content: string;               // Markdown 格式的草案内容
    summary: {                     // 结构化摘要
      core_conclusions: string[];  // 核心结论
      disagreements: string[];     // 主要分歧
      action_items: string[];      // 建议行动项
    };
  };
  
  // 上下文
  context: {
    preceding_messages: Message[];  // 前序对话
    agent_votes?: AgentVote[];      // 如果有投票，显示投票结果
  };
  
  // 用户决定
  decision?: HumanDecision;
}

interface HumanDecision {
  action: 'approve' | 'modify' | 'reject';
  modified_content?: string;       // 如果是 modify，用户修改后的内容
  rejection_reason?: string;       // 如果是 reject，驳回理由
  decided_by: string;              // 决策人 ID
  decided_at: string;              // 决策时间
}

interface AgentVote {
  agent_id: string;
  agent_name: string;
  vote: 'yes' | 'no' | 'abstain';
  reason: string;
}
```

---

## WebSocket 事件

### 1. 裁决请求通知 (Server → Client)

当工作流到达 HumanReview 节点时，后端发送此事件：

```json
{
  "type": "human_review:requested",
  "timestamp": "2024-12-16T10:00:00Z",
  "data": {
    "review_id": "review-uuid",
    "session_id": "session-uuid",
    "node_id": "node-human-review",
    "draft": {
      "title": "商业计划评审结论",
      "content": "## 核心结论\n\n1. 市场机会真实存在...\n\n## 主要分歧\n\n...",
      "summary": {
        "core_conclusions": ["市场机会真实", "技术可行"],
        "disagreements": ["融资规模存在分歧"],
        "action_items": ["进一步调研竞品"]
      }
    },
    "context": {
      "preceding_messages": [/* ... */],
      "agent_votes": [
        {"agent_id": "ceo", "agent_name": "CEO", "vote": "yes", "reason": "..."},
        {"agent_id": "cfo", "agent_name": "CFO", "vote": "no", "reason": "..."}
      ]
    },
    "expires_at": "2024-12-16T10:30:00Z"  // 30分钟超时
  }
}
```

### 2. 裁决提醒 (Server → Client)

超时前 5 分钟提醒：

```json
{
  "type": "human_review:reminder",
  "timestamp": "2024-12-16T10:25:00Z",
  "data": {
    "review_id": "review-uuid",
    "expires_in_seconds": 300,
    "message": "人类裁决将在 5 分钟后超时"
  }
}
```

### 3. 裁决超时 (Server → Client)

```json
{
  "type": "human_review:expired",
  "timestamp": "2024-12-16T10:30:00Z",
  "data": {
    "review_id": "review-uuid",
    "action_taken": "session_paused",  // 或 "default_rejected"
    "message": "人类裁决已超时，会话已暂停"
  }
}
```

---

## REST API 端点

### 1. 提交人类裁决

```http
POST /api/v1/sessions/:sessionId/review
```

**Request Body:**

```json
{
  "review_id": "review-uuid",
  "action": "approve",  // "approve" | "modify" | "reject"
  "modified_content": null,
  "rejection_reason": null
}
```

**Response 200 (approve/modify):**

```json
{
  "status": "accepted",
  "message": "裁决已接受，工作流继续执行",
  "next_node_id": "node-end"
}
```

**Response 200 (reject):**

```json
{
  "status": "rejected",
  "message": "裁决已驳回，工作流已终止",
  "session_status": "terminated"
}
```

**Response 404:**

```json
{
  "error": "Review request not found or expired"
}
```

**Response 409:**

```json
{
  "error": "Review already submitted"
}
```

---

### 2. 获取待处理裁决列表

```http
GET /api/v1/reviews/pending
```

用于显示用户所有待处理的裁决请求（多会话场景）。

**Response 200:**

```json
{
  "reviews": [
    {
      "id": "review-uuid-1",
      "session_id": "session-uuid-1",
      "session_name": "商业计划评审",
      "created_at": "2024-12-16T10:00:00Z",
      "expires_at": "2024-12-16T10:30:00Z"
    }
  ],
  "total": 1
}
```

---

### 3. 获取裁决详情

```http
GET /api/v1/reviews/:reviewId
```

**Response 200:** 返回完整的 `HumanReviewRequest` 对象

---

## 前端实现示例

```tsx
// HumanReviewModal.tsx
const HumanReviewModal: FC = () => {
  const { pendingReview, submitReview } = useSessionStore();
  const [content, setContent] = useState(pendingReview?.draft.content || '');
  const [rejectReason, setRejectReason] = useState('');
  const [mode, setMode] = useState<'view' | 'edit' | 'reject'>('view');
  
  if (!pendingReview) return null;
  
  const handleApprove = () => {
    submitReview({
      review_id: pendingReview.id,
      action: 'approve',
    });
  };
  
  const handleModify = () => {
    submitReview({
      review_id: pendingReview.id,
      action: 'modify',
      modified_content: content,
    });
  };
  
  const handleReject = () => {
    submitReview({
      review_id: pendingReview.id,
      action: 'reject',
      rejection_reason: rejectReason,
    });
  };
  
  return (
    <Dialog open={true}>
      <DialogContent className="max-w-4xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Shield className="text-amber-500" />
            🛡️ 需要人类裁决
          </DialogTitle>
          <DialogDescription>
            AI 已生成决策草案，请仔细审查。您具有最终决定权。
          </DialogDescription>
        </DialogHeader>
        
        {/* 投票结果展示 */}
        {pendingReview.context.agent_votes && (
          <div className="flex gap-4 my-4">
            {pendingReview.context.agent_votes.map(vote => (
              <VoteBadge key={vote.agent_id} vote={vote} />
            ))}
          </div>
        )}
        
        {/* 草案内容 */}
        {mode === 'view' && (
          <div className="prose prose-sm max-w-none p-4 bg-gray-50 rounded-lg">
            <ReactMarkdown>{pendingReview.draft.content}</ReactMarkdown>
          </div>
        )}
        
        {mode === 'edit' && (
          <Textarea
            className="min-h-[300px] font-mono"
            value={content}
            onChange={e => setContent(e.target.value)}
          />
        )}
        
        {mode === 'reject' && (
          <Textarea
            placeholder="请输入驳回理由..."
            value={rejectReason}
            onChange={e => setRejectReason(e.target.value)}
          />
        )}
        
        {/* 操作按钮 */}
        <DialogFooter className="gap-2">
          {mode === 'view' && (
            <>
              <Button variant="outline" onClick={() => setMode('reject')}>
                驳回
              </Button>
              <Button variant="outline" onClick={() => setMode('edit')}>
                修改
              </Button>
              <Button onClick={handleApprove}>
                签署并通过
              </Button>
            </>
          )}
          {mode === 'edit' && (
            <>
              <Button variant="ghost" onClick={() => setMode('view')}>取消</Button>
              <Button onClick={handleModify}>提交修改</Button>
            </>
          )}
          {mode === 'reject' && (
            <>
              <Button variant="ghost" onClick={() => setMode('view')}>返回</Button>
              <Button variant="destructive" onClick={handleReject}>确认驳回</Button>
            </>
          )}
        </DialogFooter>
        
        {/* 超时提示 */}
        <TimeoutIndicator expiresAt={pendingReview.expires_at} />
      </DialogContent>
    </Dialog>
  );
};
```

---

## 后端实现要点

### HumanReviewProcessor

```go
type HumanReviewProcessor struct {
    SessionID   string
    NodeID      string
    ReviewStore ReviewStore
    Timeout     time.Duration // 默认 30 分钟
}

func (h *HumanReviewProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
    // 1. 生成决策草案
    draft := h.generateDraft(input)
    
    // 2. 创建裁决请求
    review := &HumanReviewRequest{
        ID:        uuid.New().String(),
        SessionID: h.SessionID,
        NodeID:    h.NodeID,
        Status:    "pending",
        Draft:     draft,
        ExpiresAt: time.Now().Add(h.Timeout),
    }
    h.ReviewStore.Save(review)
    
    // 3. 通知前端
    stream <- StreamEvent{
        Type: "human_review:requested",
        Data: review,
    }
    
    // 4. 阻塞等待决策
    decision, err := h.waitForDecision(ctx, review.ID)
    if err != nil {
        return nil, fmt.Errorf("human review failed: %w", err)
    }
    
    // 5. 处理决策
    switch decision.Action {
    case "approve":
        return input, nil  // 继续执行
    case "modify":
        return map[string]interface{}{
            "human_modified": decision.ModifiedContent,
        }, nil
    case "reject":
        return nil, fmt.Errorf("rejected by human: %s", decision.RejectionReason)
    default:
        return nil, fmt.Errorf("unknown action: %s", decision.Action)
    }
}
```
