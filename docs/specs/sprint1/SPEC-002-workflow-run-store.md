# SPEC-002: useWorkflowRunStore 工作流运行时状态

> **优先级**: P0 (阻断项)  
> **预估工时**: 3h  
> **关联 PRD**: F.4.1 流程监控, F.4.4 成本预估  
> **关联 TDD**: 05_frontend.md §5.2, 03_communication.md

---

## 1. 设计背景

### 1.1 职责划分

| Store                 | 职责                                            |
| --------------------- | ----------------------------------------------- |
| `useSessionStore`     | 管理会话元数据、消息流 (聊天视图)               |
| `useWorkflowRunStore` | 管理工作流图状态、节点高亮、执行控制 (画布视图) |

### 1.2 分离原因

1. **渲染优化**: 画布更新和消息更新可独立触发渲染
2. **单一职责**: 图操作与消息操作逻辑分离
3. **可测试性**: 更细粒度的测试边界

---

## 2. 接口规格

### 2.1 核心类型

```typescript
// types/workflow-run.ts

/**
 * 运行时节点数据 (覆盖 React Flow Node)
 */
export interface RuntimeNode {
  id: string;
  type: string;
  label: string;
  status: NodeStatus;
  progress?: number;           // 0-100, 用于长时间节点
  error?: string;              // 错误信息
  tokenUsage?: {
    input: number;
    output: number;
    cost: number;
  };
}

/**
 * 控制命令
 */
export type ControlAction = 'pause' | 'resume' | 'stop';

/**
 * 运行控制状态
 */
export interface RunControlState {
  canPause: boolean;
  canResume: boolean;
  canStop: boolean;
}
```

### 2.2 Store 接口

```typescript
// stores/useWorkflowRunStore.ts
import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import type { Node, Edge } from 'reactflow';

interface WorkflowRunState {
  // === State ===
  
  /**
   * React Flow 节点 (运行时增强)
   */
  nodes: Node<RuntimeNode>[];
  
  /**
   * React Flow 边
   */
  edges: Edge[];
  
  /**
   * 当前高亮的节点 ID 列表
   */
  activeNodeIds: Set<string>;
  
  /**
   * 执行状态
   */
  executionStatus: 'idle' | 'running' | 'paused' | 'completed' | 'failed';
  
  /**
   * 累计统计
   */
  stats: {
    totalNodes: number;
    completedNodes: number;
    failedNodes: number;
    totalTokens: number;
    totalCostUsd: number;
    elapsedTimeMs: number;
  };
  
  /**
   * 控制状态 (派生)
   */
  readonly controlState: RunControlState;
  
  // === Actions ===
  
  /**
   * 加载工作流图
   */
  loadWorkflow: (nodes: Node[], edges: Edge[]) => void;
  
  /**
   * 清除工作流
   */
  clearWorkflow: () => void;
  
  /**
   * 更新单个节点状态
   */
  updateNodeStatus: (nodeId: string, status: NodeStatus, error?: string) => void;
  
  /**
   * 批量设置活跃节点
   */
  setActiveNodes: (nodeIds: string[]) => void;
  
  /**
   * 追加活跃节点 (并行场景)
   */
  addActiveNode: (nodeId: string) => void;
  
  /**
   * 移除活跃节点
   */
  removeActiveNode: (nodeId: string) => void;
  
  /**
   * 更新节点 Token 消耗
   */
  updateNodeTokenUsage: (nodeId: string, usage: RuntimeNode['tokenUsage']) => void;
  
  /**
   * 设置执行状态
   */
  setExecutionStatus: (status: WorkflowRunState['executionStatus']) => void;
  
  /**
   * 发送控制命令 (通过 API)
   */
  sendControl: (sessionId: string, action: ControlAction) => Promise<void>;
  
  /**
   * 启动计时器
   */
  startTimer: () => void;
  
  /**
   * 停止计时器
   */
  stopTimer: () => void;
}
```

---

## 3. 实现规格

### 3.1 节点状态样式映射

```typescript
// utils/nodeStyles.ts

export const getNodeStatusStyles = (status: NodeStatus): React.CSSProperties => {
  switch (status) {
    case 'pending':
      return { opacity: 0.6 };
    case 'running':
      return {
        boxShadow: '0 0 0 2px #3B82F6',
        animation: 'pulse 1.5s ease-in-out infinite',
      };
    case 'completed':
      return {
        borderColor: '#10B981',
        boxShadow: '0 0 8px rgba(16, 185, 129, 0.3)',
      };
    case 'failed':
      return {
        borderColor: '#EF4444',
        boxShadow: '0 0 8px rgba(239, 68, 68, 0.3)',
      };
    default:
      return {};
  }
};

// 节点状态图标
export const getNodeStatusIcon = (status: NodeStatus): string => {
  switch (status) {
    case 'pending': return '⏳';
    case 'running': return '🔄';
    case 'completed': return '✅';
    case 'failed': return '❌';
    default: return '';
  }
};
```

### 3.2 控制状态派生

```typescript
// 在 Store 中作为 getter 实现
get controlState(): RunControlState {
  const status = get().executionStatus;
  return {
    canPause: status === 'running',
    canResume: status === 'paused',
    canStop: status === 'running' || status === 'paused',
  };
}
```

### 3.3 发送控制命令

```typescript
sendControl: async (sessionId: string, action: ControlAction) => {
  const { setExecutionStatus } = get();
  
  try {
    const response = await fetch(`/api/v1/sessions/${sessionId}/control`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ action }),
    });
    
    if (!response.ok) {
      throw new Error(`Control action failed: ${response.statusText}`);
    }
    
    // 乐观更新
    switch (action) {
      case 'pause':
        setExecutionStatus('paused');
        break;
      case 'resume':
        setExecutionStatus('running');
        break;
      case 'stop':
        setExecutionStatus('failed'); // 或 'cancelled'
        break;
    }
  } catch (error) {
    console.error('Failed to send control command:', error);
    throw error;
  }
},
```

### 3.4 计时器实现

```typescript
private timerRef: ReturnType<typeof setInterval> | null = null;

startTimer: () => {
  const { stopTimer } = get();
  stopTimer(); // 防止重复启动
  
  const startTime = Date.now();
  set({ stats: { ...get().stats, elapsedTimeMs: 0 } });
  
  timerRef = setInterval(() => {
    set(state => ({
      stats: { ...state.stats, elapsedTimeMs: Date.now() - startTime }
    }));
  }, 100);
},

stopTimer: () => {
  if (timerRef) {
    clearInterval(timerRef);
    timerRef = null;
  }
},
```

---

## 4. 组件集成

### 4.1 WorkflowCanvas 只读模式

```tsx
// components/workflow/WorkflowCanvas.tsx

interface WorkflowCanvasProps {
  readOnly?: boolean;  // Run Mode 下为 true
}

export const WorkflowCanvas: FC<WorkflowCanvasProps> = ({ readOnly }) => {
  const { nodes, edges, activeNodeIds } = useWorkflowRunStore();
  
  // 为活跃节点添加动画类
  const enhancedNodes = useMemo(() => 
    nodes.map(node => ({
      ...node,
      className: cn(
        node.className,
        activeNodeIds.has(node.id) && 'node-active-pulse'
      ),
      data: {
        ...node.data,
        isActive: activeNodeIds.has(node.id),
      },
    })),
    [nodes, activeNodeIds]
  );
  
  return (
    <ReactFlow
      nodes={enhancedNodes}
      edges={edges}
      nodesDraggable={!readOnly}
      nodesConnectable={!readOnly}
      elementsSelectable={!readOnly}
      fitView
    >
      <Background />
      <Controls showInteractive={!readOnly} />
    </ReactFlow>
  );
};
```

### 4.2 执行控制栏

```tsx
// components/meeting/ExecutionControlBar.tsx

export const ExecutionControlBar: FC<{ sessionId: string }> = ({ sessionId }) => {
  const { executionStatus, controlState, sendControl, stats } = useWorkflowRunStore();
  
  const formatTime = (ms: number) => {
    const seconds = Math.floor(ms / 1000);
    const minutes = Math.floor(seconds / 60);
    return `${minutes}:${(seconds % 60).toString().padStart(2, '0')}`;
  };
  
  return (
    <div className="flex items-center gap-4 p-2 bg-gray-50 rounded-lg">
      {/* 状态指示器 */}
      <div className="flex items-center gap-2">
        <StatusDot status={executionStatus} />
        <span className="text-sm font-medium capitalize">{executionStatus}</span>
      </div>
      
      {/* 控制按钮 */}
      <div className="flex gap-2">
        {controlState.canPause && (
          <Button 
            variant="outline" 
            size="sm"
            onClick={() => sendControl(sessionId, 'pause')}
          >
            <Pause size={16} className="mr-1" /> 暂停
          </Button>
        )}
        
        {controlState.canResume && (
          <Button 
            variant="outline" 
            size="sm"
            onClick={() => sendControl(sessionId, 'resume')}
          >
            <Play size={16} className="mr-1" /> 继续
          </Button>
        )}
        
        {controlState.canStop && (
          <Button 
            variant="destructive" 
            size="sm"
            onClick={() => sendControl(sessionId, 'stop')}
          >
            <Square size={16} className="mr-1" /> 停止
          </Button>
        )}
      </div>
      
      {/* 统计信息 */}
      <div className="ml-auto flex items-center gap-4 text-sm text-gray-500">
        <span>⏱️ {formatTime(stats.elapsedTimeMs)}</span>
        <span>📊 {stats.completedNodes}/{stats.totalNodes} 节点</span>
        <span>💰 ${stats.totalCostUsd.toFixed(4)}</span>
      </div>
    </div>
  );
};
```

---

## 5. WebSocket 事件处理

### 5.1 事件映射表

| WebSocket Event       | Store Action                                |
| --------------------- | ------------------------------------------- |
| `node_state_change`   | `updateNodeStatus()`                        |
| `node:parallel_start` | `setActiveNodes(branchIds)`                 |
| `token_usage`         | `updateNodeTokenUsage()`                    |
| `execution:paused`    | `setExecutionStatus('paused')`              |
| `execution:completed` | `setExecutionStatus('completed')`           |
| `error`               | `updateNodeStatus(nodeId, 'failed', error)` |

### 5.2 事件处理 Hook

```typescript
// hooks/useWorkflowEvents.ts

export const useWorkflowEvents = () => {
  const { lastMessage } = useConnectStore();
  const workflowStore = useWorkflowRunStore();
  
  useEffect(() => {
    if (!lastMessage) return;
    
    const { event, data } = lastMessage as { event: string; data: any };
    
    switch (event) {
      case 'node_state_change':
        workflowStore.updateNodeStatus(data.node_id, data.status);
        if (data.status === 'running') {
          workflowStore.addActiveNode(data.node_id);
        } else if (data.status === 'completed' || data.status === 'failed') {
          workflowStore.removeActiveNode(data.node_id);
        }
        break;
        
      case 'node:parallel_start':
        workflowStore.setActiveNodes(data.branches);
        break;
        
      case 'token_usage':
        workflowStore.updateNodeTokenUsage(data.node_id, {
          input: data.input_tokens,
          output: data.output_tokens,
          cost: data.estimated_cost_usd,
        });
        break;
        
      case 'execution:paused':
        workflowStore.setExecutionStatus('paused');
        workflowStore.stopTimer();
        break;
        
      case 'execution:completed':
        workflowStore.setExecutionStatus('completed');
        workflowStore.stopTimer();
        break;
        
      case 'error':
        workflowStore.updateNodeStatus(data.node_id, 'failed', data.error);
        break;
    }
  }, [lastMessage, workflowStore]);
};
```

---

## 6. 测试规格

### 6.1 单元测试

```typescript
describe('useWorkflowRunStore', () => {
  describe('updateNodeStatus', () => {
    it('should update node status and trigger style change', () => {
      const { loadWorkflow, updateNodeStatus, nodes } = useWorkflowRunStore.getState();
      
      loadWorkflow([
        { id: 'node-1', type: 'start', position: { x: 0, y: 0 }, data: {} as RuntimeNode },
      ], []);
      
      updateNodeStatus('node-1', 'running');
      
      const updatedNodes = useWorkflowRunStore.getState().nodes;
      expect(updatedNodes[0].data.status).toBe('running');
    });
  });

  describe('controlState', () => {
    it('should derive correct control states', () => {
      const store = useWorkflowRunStore.getState();
      
      store.setExecutionStatus('running');
      expect(store.controlState.canPause).toBe(true);
      expect(store.controlState.canResume).toBe(false);
      
      store.setExecutionStatus('paused');
      expect(store.controlState.canPause).toBe(false);
      expect(store.controlState.canResume).toBe(true);
    });
  });

  describe('activeNodeIds', () => {
    it('should manage active nodes for parallel execution', () => {
      const store = useWorkflowRunStore.getState();
      
      store.setActiveNodes(['node-a', 'node-b']);
      expect(store.activeNodeIds.size).toBe(2);
      
      store.removeActiveNode('node-a');
      expect(store.activeNodeIds.has('node-a')).toBe(false);
      expect(store.activeNodeIds.has('node-b')).toBe(true);
    });
  });
});
```

---

## 7. CSS 动画

```css
/* index.css */

@keyframes node-pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.5);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(59, 130, 246, 0);
  }
}

.node-active-pulse {
  animation: node-pulse 1.5s ease-in-out infinite;
}

.react-flow__node.node-running {
  border: 2px solid #3B82F6;
}

.react-flow__node.node-completed {
  border: 2px solid #10B981;
}

.react-flow__node.node-failed {
  border: 2px solid #EF4444;
}
```

---

## 8. 检查清单

- [ ] 创建类型定义 `types/workflow-run.ts`
- [ ] 实现 `stores/useWorkflowRunStore.ts`
- [ ] 实现节点样式工具函数
- [ ] 更新 `WorkflowCanvas` 支持只读模式
- [ ] 创建 `ExecutionControlBar` 组件
- [ ] 实现 `useWorkflowEvents` hook
- [ ] 添加 CSS 动画
- [ ] 编写单元测试

---

## 9. 变更日志

| 日期       | 版本 | 作者 | 变更内容     |
| ---------- | ---- | ---- | ------------ |
| 2024-12-16 | v1.0 | -    | 初始规格创建 |
