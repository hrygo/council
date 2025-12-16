# SPEC-001: useSessionStore 完整重写

> **优先级**: P0 (阻断项)  
> **预估工时**: 4h  
> **关联 PRD**: F.4.0 弹性布局, F.4.1 流程监控, F.4.2 结构化对话流  
> **关联 TDD**: 05_frontend.md §5.2

---

## 1. 问题背景

### 当前实现
```typescript
// 当前 useSessionStore.ts (严重空虚)
interface SessionState {
    user: { id: string; name: string } | null;
    setUser: (user: { id: string; name: string }) => void;
    logout: () => void;
}
```

### 问题分析
1. **职责混乱**: 当前 Store 仅管理用户信息，与 "Session" (工作流运行实例) 概念完全不符
2. **核心缺失**: 无法支撑 Run Mode 的任何功能
3. **审计结论**: "前端 `useSessionStore` 严重空虚，无法支撑 Run Mode" (审计报告 §3.2)

---

## 2. 设计目标

### 2.1 职责定义

`useSessionStore` 负责管理 **当前活跃的工作流运行会话**，包括：
- 会话元数据 (ID, 状态, 时间)
- 工作流图结构 (nodes, edges)
- 节点执行状态
- 消息流 (按节点分组)

### 2.2 设计原则

1. **单一职责**: 仅管理运行时会话状态，用户认证移至独立 Store
2. **规范化存储**: 使用 `Map<nodeId, T>` 结构便于增量更新
3. **派生状态**: 通过 Zustand selector 派生计算属性
4. **持久化**: 不持久化运行时状态 (刷新后需重新连接)

---

## 3. 接口规格

### 3.1 核心类型定义

```typescript
// types/session.ts

/**
 * 节点执行状态枚举
 * 对应后端 NodeStatus
 */
export type NodeStatus = 'pending' | 'running' | 'completed' | 'failed';

/**
 * 会话整体状态
 * 对应后端 SessionStatus
 */
export type SessionStatus = 
  | 'idle'       // 未开始
  | 'running'    // 执行中
  | 'paused'     // 已暂停
  | 'completed'  // 已完成
  | 'failed'     // 执行失败
  | 'cancelled'; // 用户取消

/**
 * 消息角色
 */
export type MessageRole = 'user' | 'agent' | 'system';

/**
 * 单条消息
 */
export interface Message {
  id: string;
  nodeId: string;          // 所属节点 ID
  agentId?: string;        // Agent ID (如有)
  agentName?: string;      // Agent 显示名
  agentAvatar?: string;    // Agent 头像
  role: MessageRole;
  content: string;
  isStreaming: boolean;    // 是否正在流式输出
  timestamp: Date;
  tokenUsage?: {
    inputTokens: number;
    outputTokens: number;
    estimatedCostUsd: number;
  };
}

/**
 * 消息组 (按节点分组)
 */
export interface MessageGroup {
  nodeId: string;
  nodeName: string;
  nodeType: 'start' | 'agent' | 'parallel' | 'sequence' | 'vote' | 'loop' | 'fact_check' | 'human_review' | 'end';
  isParallel: boolean;     // 是否为并行组
  messages: Message[];
  status: NodeStatus;
}

/**
 * 节点状态快照
 */
export interface NodeStateSnapshot {
  id: string;
  status: NodeStatus;
  startedAt?: Date;
  completedAt?: Date;
  tokenUsage?: {
    inputTokens: number;
    outputTokens: number;
  };
}

/**
 * 会话状态
 */
export interface WorkflowSession {
  id: string;
  workflowId: string;
  groupId: string;
  status: SessionStatus;
  startedAt?: Date;
  completedAt?: Date;
  
  // 工作流图
  nodes: Map<string, NodeStateSnapshot>;
  
  // 当前高亮节点 (可能多个，如并行执行)
  activeNodeIds: string[];
  
  // 累计统计
  totalTokens: number;
  totalCostUsd: number;
}
```

### 3.2 Store 接口

```typescript
// stores/useSessionStore.ts
import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';

interface SessionState {
  // === State ===
  
  /**
   * 当前活跃会话
   */
  currentSession: WorkflowSession | null;
  
  /**
   * 消息列表 (按节点分组)
   */
  messageGroups: MessageGroup[];
  
  /**
   * 连接状态
   */
  connectionStatus: 'disconnected' | 'connecting' | 'connected' | 'error';
  
  // === Actions ===
  
  /**
   * 初始化新会话
   */
  initSession: (params: {
    sessionId: string;
    workflowId: string;
    groupId: string;
    nodes: Array<{ id: string; name: string; type: string }>;
  }) => void;
  
  /**
   * 更新会话状态
   */
  updateSessionStatus: (status: SessionStatus) => void;
  
  /**
   * 更新节点状态
   */
  updateNodeStatus: (nodeId: string, status: NodeStatus) => void;
  
  /**
   * 设置当前活跃节点
   */
  setActiveNodes: (nodeIds: string[]) => void;
  
  /**
   * 追加/更新消息 (处理流式输出)
   * @param isStreaming - 如果为 true，则追加到最后一条消息；否则创建新消息
   */
  appendMessage: (message: Omit<Message, 'id' | 'timestamp'> & { isChunk?: boolean }) => void;
  
  /**
   * 标记消息流式完成
   */
  finalizeMessage: (nodeId: string, agentId?: string) => void;
  
  /**
   * 更新 Token 使用量
   */
  updateTokenUsage: (nodeId: string, agentId: string, usage: {
    inputTokens: number;
    outputTokens: number;
    estimatedCostUsd: number;
  }) => void;
  
  /**
   * 清理当前会话
   */
  clearSession: () => void;
  
  /**
   * 设置连接状态
   */
  setConnectionStatus: (status: SessionState['connectionStatus']) => void;
}

// === Selectors ===

/**
 * 获取当前活跃节点的消息组
 */
export const selectActiveMessageGroups = (state: SessionState): MessageGroup[] => {
  if (!state.currentSession) return [];
  return state.messageGroups.filter(
    g => state.currentSession!.activeNodeIds.includes(g.nodeId)
  );
};

/**
 * 获取指定节点的状态
 */
export const selectNodeStatus = (nodeId: string) => (state: SessionState): NodeStatus | null => {
  return state.currentSession?.nodes.get(nodeId)?.status ?? null;
};

/**
 * 获取累计成本
 */
export const selectTotalCost = (state: SessionState): number => {
  return state.currentSession?.totalCostUsd ?? 0;
};
```

---

## 4. 实现规格

### 4.1 消息流式处理逻辑

```typescript
appendMessage: (msg) => {
  set(state => {
    // 1. 查找对应的消息组
    let group = state.messageGroups.find(g => g.nodeId === msg.nodeId);
    
    // 2. 如果不存在，创建新组
    if (!group) {
      group = {
        nodeId: msg.nodeId,
        nodeName: '节点', // 从 nodes 获取
        nodeType: 'agent',
        isParallel: false,
        messages: [],
        status: 'running',
      };
      state.messageGroups.push(group);
    }
    
    // 3. 处理流式消息
    if (msg.isChunk && msg.isStreaming) {
      // 查找同一 Agent 的最后一条流式消息
      const existingMsg = group.messages.find(
        m => m.agentId === msg.agentId && m.isStreaming
      );
      
      if (existingMsg) {
        // 追加内容
        existingMsg.content += msg.content;
        return;
      }
    }
    
    // 4. 创建新消息
    group.messages.push({
      id: crypto.randomUUID(),
      nodeId: msg.nodeId,
      agentId: msg.agentId,
      agentName: msg.agentName,
      agentAvatar: msg.agentAvatar,
      role: msg.role,
      content: msg.content,
      isStreaming: msg.isStreaming,
      timestamp: new Date(),
    });
  });
};
```

### 4.2 并行消息识别

```typescript
// 当收到 node:parallel_start 事件时
handleParallelStart: (nodeId: string, branchIds: string[]) => {
  set(state => {
    // 创建并行消息组
    state.messageGroups.push({
      nodeId,
      nodeName: 'Parallel Execution',
      nodeType: 'parallel',
      isParallel: true,
      messages: [], // 并行消息将被收集到这里
      status: 'running',
    });
    
    // 标记分支节点为活跃
    if (state.currentSession) {
      state.currentSession.activeNodeIds = branchIds;
    }
  });
};
```

---

## 5. 测试规格

### 5.1 单元测试

```typescript
// stores/__tests__/useSessionStore.test.ts

describe('useSessionStore', () => {
  beforeEach(() => {
    useSessionStore.getState().clearSession();
  });

  describe('initSession', () => {
    it('should initialize session with correct structure', () => {
      const { initSession, currentSession } = useSessionStore.getState();
      
      initSession({
        sessionId: 'sess-123',
        workflowId: 'wf-456',
        groupId: 'group-789',
        nodes: [
          { id: 'node-1', name: 'Start', type: 'start' },
          { id: 'node-2', name: 'Analyst', type: 'agent' },
        ],
      });
      
      const session = useSessionStore.getState().currentSession;
      expect(session).not.toBeNull();
      expect(session?.id).toBe('sess-123');
      expect(session?.status).toBe('idle');
      expect(session?.nodes.size).toBe(2);
    });
  });

  describe('appendMessage (streaming)', () => {
    it('should append chunks to existing streaming message', () => {
      const store = useSessionStore.getState();
      store.initSession({ /* ... */ });
      
      // First chunk
      store.appendMessage({
        nodeId: 'node-2',
        agentId: 'agent-1',
        role: 'agent',
        content: 'Hello ',
        isStreaming: true,
        isChunk: true,
      });
      
      // Second chunk
      store.appendMessage({
        nodeId: 'node-2',
        agentId: 'agent-1',
        role: 'agent',
        content: 'World!',
        isStreaming: true,
        isChunk: true,
      });
      
      const groups = useSessionStore.getState().messageGroups;
      expect(groups[0].messages.length).toBe(1);
      expect(groups[0].messages[0].content).toBe('Hello World!');
    });
  });

  describe('updateTokenUsage', () => {
    it('should accumulate total cost', () => {
      const store = useSessionStore.getState();
      store.initSession({ /* ... */ });
      
      store.updateTokenUsage('node-2', 'agent-1', {
        inputTokens: 100,
        outputTokens: 50,
        estimatedCostUsd: 0.01,
      });
      
      store.updateTokenUsage('node-3', 'agent-2', {
        inputTokens: 200,
        outputTokens: 100,
        estimatedCostUsd: 0.02,
      });
      
      const session = useSessionStore.getState().currentSession;
      expect(session?.totalCostUsd).toBe(0.03);
    });
  });
});
```

### 5.2 集成测试场景

| ID     | 场景         | 预期结果                                  |
| ------ | ------------ | ----------------------------------------- |
| IT-001 | 启动工作流   | Session 初始化，状态为 `idle` → `running` |
| IT-002 | 接收流式消息 | 消息实时追加，UI 同步更新                 |
| IT-003 | 并行节点执行 | 多个 Agent 消息按组并排显示               |
| IT-004 | 暂停/恢复    | 状态正确切换，消息流暂停/继续             |
| IT-005 | Token 统计   | 累计成本实时更新                          |

---

## 6. 注意事项

### 6.1 用户认证分离

当前 `useSessionStore` 中的 `user` 相关逻辑需要迁移到独立的 `useAuthStore`：

```typescript
// stores/useAuthStore.ts (新增)
interface AuthState {
  user: { id: string; name: string } | null;
  setUser: (user: { id: string; name: string }) => void;
  logout: () => void;
}
```

### 6.2 与 useConnectStore 的协作

`useSessionStore` 不直接管理 WebSocket 连接，而是：
1. `useConnectStore` 负责 WebSocket 生命周期
2. 消息分发逻辑在 `useWebSocketHandler` hook 中处理
3. `useSessionStore` 仅接收处理后的状态更新

```
┌─────────────────┐      ┌─────────────────────┐      ┌─────────────────┐
│ useConnectStore │ ───▶ │ useWebSocketHandler │ ───▶ │ useSessionStore │
│   (WS 连接)     │      │   (消息路由/解析)   │      │   (状态更新)    │
└─────────────────┘      └─────────────────────┘      └─────────────────┘
```

---

## 7. 检查清单

- [ ] 创建新的类型定义文件 `types/session.ts`
- [ ] 重写 `stores/useSessionStore.ts`
- [ ] 创建 `stores/useAuthStore.ts` (迁移用户逻辑)
- [ ] 编写单元测试
- [ ] 更新相关组件引用
- [ ] 与 WebSocket 消息处理集成

---

## 8. 变更日志

| 日期       | 版本 | 作者 | 变更内容     |
| ---------- | ---- | ---- | ------------ |
| 2024-12-16 | v1.0 | -    | 初始规格创建 |
