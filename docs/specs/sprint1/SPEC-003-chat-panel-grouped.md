# SPEC-003: ChatPanel 按节点分组消息显示

> **优先级**: P1  
> **预估工时**: 2h  
> **关联 PRD**: F.4.2 结构化对话流  
> **关联 TDD**: 05_frontend.md §5.4

---

## 1. 设计背景

### 1.1 当前问题

```tsx
// 当前 ChatPanel.tsx - 线性消息渲染
{messages.map((msg, idx) => (
  <div key={idx} className={`flex gap-3 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}>
    ...
  </div>
))}
```

**问题**:
- 消息线性排列，无法区分来自哪个工作流节点
- 无"阶段分节"标题
- 无法识别并行消息

### 1.2 目标效果

```
┌────────────────────────────────────────────────────┐
│ 📍 阶段 1: 系统分析 (node-analyst-1)               │
│ ═══════════════════════════════════════════════════│
│ [Architect Avatar]                                 │
│ ┌─────────────────────────────────────────────────┐│
│ │ 根据系统架构分析，我认为...                     ││
│ │ - 关键点 1                                      ││
│ │ - 关键点 2                                      ││
│ └─────────────────────────────────────────────────┘│
│                                                    │
│ 📍 阶段 2: 并行审查 (node-parallel-1)              │
│ ═══════════════════════════════════════════════════│
│ ┌──────────────────┐  ┌──────────────────┐        │
│ │ [Security]       │  │ [Performance]    │        │
│ │ 安全审查结论...  │  │ 性能审查结论...  │        │
│ └──────────────────┘  └──────────────────┘        │
└────────────────────────────────────────────────────┘
```

---

## 2. 接口规格

### 2.1 Props 定义

```typescript
// components/chat/ChatPanel.types.ts

export interface ChatPanelProps {
  /**
   * 全屏模式
   */
  fullscreen?: boolean;
  
  /**
   * 退出全屏回调
   */
  onExitFullscreen?: () => void;
  
  /**
   * 只读模式 (禁用输入框)
   */
  readOnly?: boolean;
  
  /**
   * Session ID (用于发送用户消息)
   */
  sessionId?: string;
}
```

### 2.2 消息组渲染

```typescript
// 使用 SPEC-001 中定义的 MessageGroup 类型
import type { MessageGroup, Message } from '@/types/session';
```

---

## 3. 组件结构

### 3.1 组件层级

```
ChatPanel
├── ChatHeader
├── MessageGroupList
│   └── MessageGroupCard (per group)
│       ├── GroupHeader (阶段标题)
│       ├── ParallelMessageRow (if parallel)
│       │   └── MessageBubble[]
│       └── SequentialMessage (if sequential)
│           └── MessageBubble
├── TypingIndicator (if streaming)
└── ChatInput
```

### 3.2 核心组件

#### MessageGroupCard

```tsx
// components/chat/MessageGroupCard.tsx

interface MessageGroupCardProps {
  group: MessageGroup;
  isActive: boolean;  // 当前是否正在执行
}

export const MessageGroupCard: FC<MessageGroupCardProps> = ({ group, isActive }) => {
  const { t } = useTranslation('chat');
  
  return (
    <div 
      className={cn(
        "mb-6 transition-all duration-300",
        isActive && "ring-2 ring-blue-500/20 bg-blue-50/30 rounded-lg p-3"
      )}
    >
      {/* 阶段标题 */}
      <GroupHeader 
        nodeName={group.nodeName}
        nodeType={group.nodeType}
        status={group.status}
      />
      
      {/* 消息内容 */}
      <div className="mt-3 pl-4 border-l-2 border-gray-200">
        {group.isParallel ? (
          <ParallelMessageRow messages={group.messages} />
        ) : (
          group.messages.map(msg => (
            <SequentialMessage key={msg.id} message={msg} />
          ))
        )}
      </div>
    </div>
  );
};
```

#### GroupHeader

```tsx
// components/chat/GroupHeader.tsx

interface GroupHeaderProps {
  nodeName: string;
  nodeType: string;
  status: NodeStatus;
}

const nodeTypeIcons: Record<string, string> = {
  start: '🚀',
  agent: '🤖',
  parallel: '⚡',
  sequence: '📝',
  vote: '🗳️',
  loop: '🔄',
  fact_check: '🔍',
  human_review: '👤',
  end: '🏁',
};

const statusColors: Record<NodeStatus, string> = {
  pending: 'text-gray-400',
  running: 'text-blue-500',
  completed: 'text-green-500',
  failed: 'text-red-500',
};

export const GroupHeader: FC<GroupHeaderProps> = ({ nodeName, nodeType, status }) => {
  const icon = nodeTypeIcons[nodeType] || '📍';
  
  return (
    <div className="flex items-center gap-2 text-sm font-medium text-gray-600">
      <span>{icon}</span>
      <span>{nodeName}</span>
      
      {/* 状态指示器 */}
      <span className={cn("ml-auto", statusColors[status])}>
        {status === 'running' && (
          <span className="inline-flex items-center gap-1">
            <LoadingSpinner size={12} />
            进行中
          </span>
        )}
        {status === 'completed' && '✓ 已完成'}
        {status === 'failed' && '✕ 失败'}
      </span>
    </div>
  );
};
```

#### SequentialMessage

```tsx
// components/chat/SequentialMessage.tsx

interface SequentialMessageProps {
  message: Message;
}

export const SequentialMessage: FC<SequentialMessageProps> = ({ message }) => {
  return (
    <div className="flex gap-3 mb-4">
      {/* Agent 头像 */}
      <AgentAvatar 
        name={message.agentName}
        avatar={message.agentAvatar}
        isStreaming={message.isStreaming}
      />
      
      {/* 消息内容 */}
      <div className="flex-1 min-w-0">
        {/* Agent 名称 */}
        <div className="text-sm font-medium text-gray-700 mb-1">
          {message.agentName || 'Agent'}
        </div>
        
        {/* 消息气泡 */}
        <MessageBubble 
          content={message.content}
          isStreaming={message.isStreaming}
          role={message.role}
        />
        
        {/* Token 消耗 (如果有) */}
        {message.tokenUsage && (
          <div className="mt-1 text-xs text-gray-400">
            💰 ${message.tokenUsage.estimatedCostUsd.toFixed(4)} 
            ({message.tokenUsage.outputTokens} tokens)
          </div>
        )}
      </div>
    </div>
  );
};
```

#### MessageBubble (增强版)

```tsx
// components/chat/MessageBubble.tsx

interface MessageBubbleProps {
  content: string;
  isStreaming: boolean;
  role: 'user' | 'agent' | 'system';
}

export const MessageBubble: FC<MessageBubbleProps> = ({ content, isStreaming, role }) => {
  return (
    <div 
      className={cn(
        "p-3 rounded-2xl text-sm",
        role === 'user' 
          ? "bg-blue-600 text-white rounded-br-none ml-auto max-w-[80%]"
          : "bg-gray-50 border border-gray-100 text-gray-800 rounded-bl-none",
        isStreaming && "animate-pulse"
      )}
    >
      <div className="prose prose-sm max-w-none">
        <ReactMarkdown
          remarkPlugins={[remarkGfm]}
          rehypePlugins={[rehypeHighlight]}
        >
          {content}
        </ReactMarkdown>
        
        {/* 流式输入光标 */}
        {isStreaming && (
          <span className="inline-block w-2 h-4 bg-gray-400 animate-blink ml-1" />
        )}
      </div>
    </div>
  );
};
```

---

## 4. 主组件实现

```tsx
// components/chat/ChatPanel.tsx

import { useEffect, useRef, useMemo } from 'react';
import { useSessionStore, selectActiveMessageGroups } from '@/stores/useSessionStore';
import { useConnectStore } from '@/stores/useConnectStore';
import { MessageGroupCard } from './MessageGroupCard';
import { ChatInput } from './ChatInput';

export const ChatPanel: FC<ChatPanelProps> = ({ 
  fullscreen, 
  onExitFullscreen, 
  readOnly,
  sessionId 
}) => {
  const messageGroups = useSessionStore(state => state.messageGroups);
  const currentSession = useSessionStore(state => state.currentSession);
  const activeNodeIds = currentSession?.activeNodeIds ?? [];
  
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  // 自动滚动到底部
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messageGroups]);
  
  // 判断当前组是否活跃
  const isGroupActive = (nodeId: string) => activeNodeIds.includes(nodeId);
  
  return (
    <div 
      className={cn(
        "flex flex-col h-full bg-white border-l border-gray-200 shadow-xl z-10 w-full",
        fullscreen && "fixed inset-0 z-50 p-8"
      )}
    >
      {/* Header */}
      <ChatHeader 
        sessionStatus={currentSession?.status}
        onExitFullscreen={fullscreen ? onExitFullscreen : undefined}
      />
      
      {/* Message Groups */}
      <div className="flex-1 overflow-y-auto p-4">
        {messageGroups.length === 0 ? (
          <EmptyState message="等待会议开始..." />
        ) : (
          messageGroups.map(group => (
            <MessageGroupCard 
              key={group.nodeId}
              group={group}
              isActive={isGroupActive(group.nodeId)}
            />
          ))
        )}
        <div ref={messagesEndRef} />
      </div>
      
      {/* Input */}
      {!readOnly && sessionId && (
        <ChatInput sessionId={sessionId} />
      )}
    </div>
  );
};
```

---

## 5. 样式规格

### 5.1 CSS 变量

```css
/* index.css */

:root {
  /* Chat Panel */
  --chat-group-border-color: #E5E7EB;
  --chat-group-active-bg: rgba(59, 130, 246, 0.05);
  --chat-group-active-border: rgba(59, 130, 246, 0.3);
  
  /* Message Bubble */
  --bubble-user-bg: #2563EB;
  --bubble-agent-bg: #F9FAFB;
  --bubble-agent-border: #E5E7EB;
}

.dark {
  --chat-group-border-color: #374151;
  --chat-group-active-bg: rgba(59, 130, 246, 0.1);
  --bubble-agent-bg: #1F2937;
  --bubble-agent-border: #374151;
}
```

### 5.2 动画

```css
/* 光标闪烁动画 */
@keyframes blink {
  0%, 100% { opacity: 1; }
  50% { opacity: 0; }
}

.animate-blink {
  animation: blink 1s step-end infinite;
}

/* 消息组进入动画 */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.message-group-enter {
  animation: slideIn 0.3s ease-out;
}
```

---

## 6. 测试规格

### 6.1 组件测试

```typescript
// components/chat/__tests__/ChatPanel.test.tsx

import { render, screen } from '@testing-library/react';
import { ChatPanel } from '../ChatPanel';
import { useSessionStore } from '@/stores/useSessionStore';

describe('ChatPanel', () => {
  beforeEach(() => {
    useSessionStore.getState().clearSession();
  });

  it('should render empty state when no messages', () => {
    render(<ChatPanel />);
    expect(screen.getByText('等待会议开始...')).toBeInTheDocument();
  });

  it('should render message groups with headers', () => {
    const store = useSessionStore.getState();
    store.initSession({ /* ... */ });
    store.appendMessage({
      nodeId: 'node-1',
      role: 'agent',
      agentName: 'Analyst',
      content: 'Test message',
      isStreaming: false,
    });
    
    render(<ChatPanel />);
    expect(screen.getByText('Analyst')).toBeInTheDocument();
    expect(screen.getByText('Test message')).toBeInTheDocument();
  });

  it('should highlight active group', () => {
    const store = useSessionStore.getState();
    store.initSession({ /* ... */ });
    store.updateSessionStatus('running');
    store.setActiveNodes(['node-1']);
    
    // ... add message to node-1
    
    render(<ChatPanel />);
    const activeGroup = screen.getByTestId('message-group-node-1');
    expect(activeGroup).toHaveClass('ring-2');
  });

  it('should show streaming indicator', () => {
    const store = useSessionStore.getState();
    store.appendMessage({
      nodeId: 'node-1',
      role: 'agent',
      content: 'Streaming...',
      isStreaming: true,
    });
    
    render(<ChatPanel />);
    expect(screen.getByTestId('streaming-cursor')).toBeInTheDocument();
  });
});
```

### 6.2 快照测试

```typescript
it('should match snapshot', () => {
  const { container } = render(<ChatPanel />);
  expect(container).toMatchSnapshot();
});
```

---

## 7. 检查清单

- [ ] 创建 `MessageGroupCard` 组件
- [ ] 创建 `GroupHeader` 组件
- [ ] 创建 `SequentialMessage` 组件
- [ ] 增强 `MessageBubble` 组件
- [ ] 重构 `ChatPanel` 主组件
- [ ] 添加 CSS 变量和动画
- [ ] 编写组件测试
- [ ] 更新 i18n 翻译

---

## 8. 变更日志

| 日期       | 版本 | 作者 | 变更内容     |
| ---------- | ---- | ---- | ------------ |
| 2025-12-16 | v1.0 | -    | 初始规格创建 |
