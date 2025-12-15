# 3. 前后端通信协议 (Communication Protocol)

由于涉及 3 个 Agent 并发输出，HTTP 请求无法满足需求，必须使用 **WebSocket** 进行全双工通信。

### 3.1 WebSocket 事件定义

**Server -> Client (Downstream):**

```json
{
  "event": "token_stream",
  "data": {
    "node_id": "node_abc123",    // 哪个节点在说话
    "agent_id": "agent_finance", // 哪个 Agent
    "chunk": "成本",             // 文本片段
    "is_thinking": false         // 是否在推理阶段 (针对 o1/DeepSeek-R1)
  }
}
```

```json
{
  "event": "node_state_change",
  "data": {
    "node_id": "node_abc123",
    "status": "running" // pending, running, completed, error
  }
}
```

```json
{
  "event": "token_usage",
  "data": {
    "node_id": "node_abc123",
    "agent_id": "agent_finance",
    "input_tokens": 1250,
    "output_tokens": 350,
    "total_tokens": 1600,
    "estimated_cost_usd": 0.0048,
    "model": "gpt-4o",
    "cumulative_session_tokens": 4500,
    "cumulative_session_cost_usd": 0.0135
  }
}
```

> [!NOTE]
> `token_usage` 事件在每个 Agent 完成发言后触发，用于前端实时显示 Token 消耗预估（对应 PRD F.4.1）。

#### 3.1.1 会议前成本预估 API (PRD F.4.4)

在点击"开始会议"前，前端调用此 API 获取预估成本：

**Request:**
```http
POST /api/v1/sessions/estimate
Content-Type: application/json

{
  "workflow_id": "uuid",
  "proposal": {
    "text": "用户输入的提案内容",
    "files": ["file1.pdf"]
  }
}
```

**Response:**
```json
{
  "estimated_tokens": 12500,
  "estimated_cost_usd": 0.35,
  "estimated_duration_seconds": 120,
  "breakdown": [
    {"agent_id": "uuid", "agent_name": "财务顾问", "model": "gpt-4o", "tokens": 4000, "cost": 0.12},
    {"agent_id": "uuid", "agent_name": "技术专家", "model": "claude-3.5-sonnet", "tokens": 5000, "cost": 0.15},
    {"agent_id": "uuid", "agent_name": "风控分析师", "model": "gpt-4o-mini", "tokens": 3500, "cost": 0.08}
  ],
  "warning": null
}
```

> 前端根据此预估结果显示："本次会议预计消耗 ~$0.35，耗时 ~2 分钟"，用户可选择继续或更换更经济的模型配置。

**Client -> Server (Upstream):**

* `cmd: start_session`: 启动会议。
* `cmd: pause_session`: 暂停。
* `cmd: user_input`: 用户在会议中途插话或回复 Vote 节点。
