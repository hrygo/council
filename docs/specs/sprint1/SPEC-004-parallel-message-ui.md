# SPEC-004: 并行消息并排显示

> **优先级**: P1 | **预估工时**: 2h  
> **关联 PRD**: F.4.2 并行 UI | **关联 TDD**: 05_frontend.md §5.4

---

## 1. 设计目标

实现并行节点执行时，多个 Agent 的消息在同一行并排显示。

**目标效果**:
```
┌──────────────────┐  ┌──────────────────┐  ┌──────────────────┐
│ 🛡️ Security      │  │ ⚡ Performance   │  │ 🔧 Maintainability│
├──────────────────┤  ├──────────────────┤  ├──────────────────┤
│ 安全审查结论...  │  │ 性能分析结论...  │  │ 可维护性评估...  │
│ 💰 $0.0032       │  │ 💰 $0.0028       │  │ 💰 $0.0025       │
└──────────────────┘  └──────────────────┘  └──────────────────┘
```

---

## 2. 接口规格

```typescript
// components/chat/ParallelMessageRow.tsx
interface ParallelMessageRowProps {
  messages: Message[];
  maxColumns?: number;  // 默认 3
}

// components/chat/ParallelMessageCard.tsx
interface ParallelMessageCardProps {
  message: Message;
  index: number;
  accentColor: string;
}
```

---

## 3. 核心实现

### 3.1 ParallelMessageRow

```tsx
const accentColors = [
  'border-t-blue-500', 'border-t-green-500', 
  'border-t-purple-500', 'border-t-orange-500'
];

export const ParallelMessageRow: FC<ParallelMessageRowProps> = ({
  messages, maxColumns = 3,
}) => (
  <div 
    className="grid gap-4"
    style={{ gridTemplateColumns: `repeat(${Math.min(messages.length, maxColumns)}, 1fr)` }}
  >
    {messages.map((msg, idx) => (
      <ParallelMessageCard
        key={msg.id}
        message={msg}
        index={idx}
        accentColor={accentColors[idx % accentColors.length]}
      />
    ))}
  </div>
);
```

### 3.2 ParallelMessageCard

```tsx
export const ParallelMessageCard: FC<ParallelMessageCardProps> = ({
  message, index, accentColor,
}) => (
  <div className={cn("border rounded-lg border-t-4", accentColor)}>
    <div className="p-3 border-b flex items-center gap-2">
      <AgentAvatar name={message.agentName} size="sm" />
      <span className="font-medium">{message.agentName}</span>
      {message.isStreaming && <LoadingSpinner size={12} />}
    </div>
    <div className="p-3 prose prose-sm max-h-[400px] overflow-y-auto">
      <ReactMarkdown>{message.content}</ReactMarkdown>
    </div>
    <div className="px-3 py-2 bg-gray-50 text-xs text-gray-500">
      {message.tokenUsage && `💰 $${message.tokenUsage.estimatedCostUsd.toFixed(4)}`}
    </div>
  </div>
);
```

---

## 4. Store 集成

```typescript
// useSessionStore.ts - 处理并行消息
handleParallelNodeStart: (nodeId: string, branches: string[]) => {
  set(state => {
    state.messageGroups.push({
      nodeId,
      nodeName: 'Parallel Review',
      nodeType: 'parallel',
      isParallel: true,
      messages: [],
      status: 'running',
    });
  });
};
```

---

## 5. 测试要点

- [ ] 多卡片网格正确渲染
- [ ] 流式消息时显示加载指示
- [ ] 响应式布局 (小屏幕单列)
- [ ] 不同颜色标识不同 Agent

---

## 6. 变更日志

| 日期       | 版本 | 变更内容 |
| ---------- | ---- | -------- |
| 2025-12-16 | v1.0 | 初始规格 |
