# SPEC-005: WebSocket 消息处理优化

> **优先级**: P1 | **预估工时**: 2h  
> **关联 PRD**: F.4.1 流程监控, F.4.2 结构化对话流  
> **关联 TDD**: 03_communication.md

---

## 1. 设计背景

### 1.1 当前问题

```typescript
// 当前 useConnectStore.ts - 简单的消息存储
ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  set({ lastMessage: data }); // 仅存储最后一条
};
```

**问题**:
1. 无消息路由机制，所有组件需要自行过滤
2. 流式消息处理重复且容易出错
3. 无重连/心跳机制
4. 无消息类型校验

### 1.2 目标

1. 统一的消息分发中心
2. 类型安全的事件处理
3. 自动重连和心跳
4. 分离连接管理与业务逻辑

---

## 2. 架构设计

```
┌─────────────────┐      ┌────────────────────┐      ┌─────────────────────┐
│ useConnectStore │ ───▶ │ useWebSocketRouter │ ───▶ │ useSessionStore     │
│   (连接管理)    │      │   (消息路由/解析)  │      │ useWorkflowRunStore │
└─────────────────┘      └────────────────────┘      └─────────────────────┘
```

---

## 3. 接口规格

### 3.1 WebSocket 消息类型

```typescript
// types/websocket.ts

// 下行事件 (Server -> Client)
export type WSEventType =
  | 'token_stream'        // Token 流
  | 'node_state_change'   // 节点状态变化
  | 'node:parallel_start' // 并行节点开始
  | 'token_usage'         // Token 使用统计
  | 'execution:paused'    // 执行已暂停
  | 'execution:completed' // 执行完成
  | 'error';              // 错误

export interface WSMessage<T = unknown> {
  event: WSEventType;
  data: T;
  timestamp?: string;
}

// 具体事件数据类型
export interface TokenStreamData {
  node_id: string;
  agent_id: string;
  chunk: string;
  is_thinking?: boolean;
}

export interface NodeStateChangeData {
  node_id: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
}

export interface TokenUsageData {
  node_id: string;
  agent_id: string;
  input_tokens: number;
  output_tokens: number;
  estimated_cost_usd: number;
}

export interface ParallelStartData {
  node_id: string;
  branches: string[];
}

// 上行命令 (Client -> Server)
export type WSCommandType = 'start_session' | 'pause_session' | 'resume_session' | 'user_input';

export interface WSCommand<T = unknown> {
  cmd: WSCommandType;
  data?: T;
}
```

### 3.2 增强版 useConnectStore

```typescript
// stores/useConnectStore.ts

interface ConnectState {
  socket: WebSocket | null;
  status: 'disconnected' | 'connecting' | 'connected' | 'reconnecting';
  lastError: string | null;
  reconnectAttempts: number;
  
  // Actions
  connect: (url: string) => void;
  disconnect: () => void;
  send: <T>(command: WSCommand<T>) => void;
  
  // Internal
  _onMessage: (msg: WSMessage) => void;
  _scheduleReconnect: () => void;
}

const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY = 3000;
const HEARTBEAT_INTERVAL = 30000;

export const useConnectStore = create<ConnectState>((set, get) => ({
  socket: null,
  status: 'disconnected',
  lastError: null,
  reconnectAttempts: 0,

  connect: (url: string) => {
    if (get().socket) return;
    
    set({ status: 'connecting' });
    const ws = new WebSocket(url);

    ws.onopen = () => {
      set({ status: 'connected', reconnectAttempts: 0, lastError: null });
      // 启动心跳
      get()._startHeartbeat();
    };

    ws.onclose = (e) => {
      set({ status: 'disconnected', socket: null });
      if (!e.wasClean) {
        get()._scheduleReconnect();
      }
    };

    ws.onerror = () => {
      set({ lastError: 'WebSocket connection error' });
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data) as WSMessage;
        get()._onMessage(msg);
      } catch (e) {
        console.error('Failed to parse WS message:', e);
      }
    };

    set({ socket: ws });
  },

  disconnect: () => {
    const { socket } = get();
    socket?.close(1000, 'Client disconnect');
    set({ socket: null, status: 'disconnected' });
  },

  send: (command) => {
    const { socket, status } = get();
    if (socket && status === 'connected') {
      socket.send(JSON.stringify(command));
    } else {
      console.warn('Cannot send: not connected');
    }
  },

  _onMessage: (msg) => {
    // 触发订阅者 (通过 Zustand subscribeWithSelector)
    // 详见 useWebSocketRouter hook
  },

  _scheduleReconnect: () => {
    const { reconnectAttempts } = get();
    if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
      set({ lastError: 'Max reconnection attempts reached' });
      return;
    }
    
    set({ status: 'reconnecting', reconnectAttempts: reconnectAttempts + 1 });
    setTimeout(() => {
      get().connect(/* stored url */);
    }, RECONNECT_DELAY * (reconnectAttempts + 1));
  },
}));
```

---

## 4. 消息路由 Hook

```typescript
// hooks/useWebSocketRouter.ts

import { useEffect, useRef } from 'react';
import { useConnectStore } from '@/stores/useConnectStore';
import { useSessionStore } from '@/stores/useSessionStore';
import { useWorkflowRunStore } from '@/stores/useWorkflowRunStore';
import type { WSMessage, TokenStreamData, NodeStateChangeData } from '@/types/websocket';

export const useWebSocketRouter = () => {
  const sessionStore = useSessionStore();
  const workflowStore = useWorkflowRunStore();
  const processedRef = useRef<Set<string>>(new Set());

  useEffect(() => {
    // 订阅 WebSocket 消息
    const unsubscribe = useConnectStore.subscribe(
      (state) => state._lastMessage,
      (message) => {
        if (!message) return;
        
        // 防重复处理
        const msgId = `${message.event}-${message.timestamp}`;
        if (processedRef.current.has(msgId)) return;
        processedRef.current.add(msgId);
        
        // 路由到对应处理器
        routeMessage(message, sessionStore, workflowStore);
      }
    );

    return unsubscribe;
  }, [sessionStore, workflowStore]);
};

function routeMessage(
  msg: WSMessage,
  sessionStore: ReturnType<typeof useSessionStore.getState>,
  workflowStore: ReturnType<typeof useWorkflowRunStore.getState>
) {
  switch (msg.event) {
    case 'token_stream': {
      const data = msg.data as TokenStreamData;
      sessionStore.appendMessage({
        nodeId: data.node_id,
        agentId: data.agent_id,
        role: 'agent',
        content: data.chunk,
        isStreaming: true,
        isChunk: true,
      });
      break;
    }

    case 'node_state_change': {
      const data = msg.data as NodeStateChangeData;
      workflowStore.updateNodeStatus(data.node_id, data.status);
      
      if (data.status === 'running') {
        workflowStore.addActiveNode(data.node_id);
      } else if (data.status === 'completed' || data.status === 'failed') {
        workflowStore.removeActiveNode(data.node_id);
        sessionStore.finalizeMessage(data.node_id);
      }
      break;
    }

    case 'node:parallel_start': {
      const data = msg.data as ParallelStartData;
      workflowStore.setActiveNodes(data.branches);
      sessionStore.handleParallelNodeStart(data.node_id, data.branches);
      break;
    }

    case 'token_usage': {
      const data = msg.data as TokenUsageData;
      sessionStore.updateTokenUsage(data.node_id, data.agent_id, {
        inputTokens: data.input_tokens,
        outputTokens: data.output_tokens,
        estimatedCostUsd: data.estimated_cost_usd,
      });
      workflowStore.updateNodeTokenUsage(data.node_id, {
        input: data.input_tokens,
        output: data.output_tokens,
        cost: data.estimated_cost_usd,
      });
      break;
    }

    case 'execution:paused':
      workflowStore.setExecutionStatus('paused');
      sessionStore.updateSessionStatus('paused');
      break;

    case 'execution:completed':
      workflowStore.setExecutionStatus('completed');
      sessionStore.updateSessionStatus('completed');
      break;

    case 'error': {
      const data = msg.data as { node_id?: string; error: string };
      if (data.node_id) {
        workflowStore.updateNodeStatus(data.node_id, 'failed', data.error);
      }
      console.error('WS Error:', data.error);
      break;
    }
  }
}
```

---

## 5. 使用示例

```tsx
// pages/MeetingRoom.tsx

export const MeetingRoom: FC = () => {
  // 初始化消息路由
  useWebSocketRouter();
  
  const { connect, status } = useConnectStore();
  
  useEffect(() => {
    connect('ws://localhost:8080/ws');
  }, [connect]);
  
  if (status !== 'connected') {
    return <ConnectionStatus status={status} />;
  }
  
  return (
    <PanelGroup direction="horizontal">
      <Panel><WorkflowCanvas /></Panel>
      <PanelResizeHandle />
      <Panel><ChatPanel /></Panel>
      <PanelResizeHandle />
      <Panel><DocumentReader /></Panel>
    </PanelGroup>
  );
};
```

---

## 6. 测试规格

```typescript
describe('useWebSocketRouter', () => {
  it('should route token_stream to sessionStore', () => {
    const { result } = renderHook(() => useWebSocketRouter());
    
    // Simulate message
    act(() => {
      useConnectStore.setState({
        _lastMessage: {
          event: 'token_stream',
          data: { node_id: 'n1', agent_id: 'a1', chunk: 'Hello' },
        },
      });
    });
    
    expect(useSessionStore.getState().messageGroups[0].messages[0].content)
      .toContain('Hello');
  });

  it('should update both stores on node_state_change', () => {
    // ...
  });

  it('should handle reconnection', async () => {
    // ...
  });
});
```

---

## 7. 检查清单

- [ ] 创建 `types/websocket.ts`
- [ ] 增强 `useConnectStore` (重连/心跳)
- [ ] 创建 `useWebSocketRouter` hook
- [ ] 集成到 `MeetingRoom` 页面
- [ ] 编写测试

---

## 8. 变更日志

| 日期       | 版本 | 变更内容 |
| ---------- | ---- | -------- |
| 2024-12-16 | v1.0 | 初始规格 |
