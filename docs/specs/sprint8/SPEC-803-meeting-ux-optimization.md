# SPEC-803: Meeting UX Optimization

## 元信息

| 属性     | 值       |
| :------- | :------- |
| Spec ID  | SPEC-803 |
| 类型     | UX       |
| 优先级   | P1       |
| 预估工时 | 3h       |
| 依赖     | SPEC-802 |

## 背景

当前会议过程中的 UX/UI 体验需改进：
- 消息组显示 "Unknown Node"
- 无会话状态指示器
- 无进度可视化

## 目标

优化会议过程中的用户体验，提供清晰的状态反馈和进度展示。

## 技术方案

### 1. SessionHeader 组件 (新建)

```
frontend/src/features/meeting/SessionHeader.tsx
```

**功能**：
- 显示会话状态徽章 (RUNNING/PAUSED/COMPLETED)
- 显示当前循环轮次 (Loop 2/3)
- 显示累计 Token 用量
- 显示累计成本

```tsx
<div className="flex items-center justify-between p-4 border-b">
    <div className="flex items-center gap-4">
        <h1>Council Meeting</h1>
        <StatusBadge status={session.status} />
    </div>
    <div className="flex items-center gap-6 text-sm text-gray-500">
        <span>Loop: {currentLoop}/{maxLoops}</span>
        <span>Tokens: {totalTokens.toLocaleString()}</span>
        <span>Cost: ${totalCost.toFixed(4)}</span>
    </div>
</div>
```

### 2. MessageGroup 优化

**改进点**：
1. 使用 Agent 头像和名称
2. 消息状态指示器
3. 改进 Markdown 渲染

```tsx
<div className="flex gap-3">
    <AgentAvatar name={agentName} />
    <div>
        <div className="flex items-center gap-2">
            <span className="font-medium">{agentName}</span>
            <StatusIndicator status={messageStatus} />
        </div>
        <MarkdownRenderer content={message.content} />
    </div>
</div>
```

### 3. 进度条组件

**位置**: ChatPanel 顶部

```tsx
<ProgressBar 
    current={completedNodes.length}
    total={totalNodes.length}
    label={`${completedNodes.length}/${totalNodes.length} 节点完成`}
/>
```

### 4. Agent 头像生成

基于 Agent 名称生成唯一颜色和首字母头像：

```tsx
function AgentAvatar({ name }: { name: string }) {
    const color = stringToColor(name);
    const initial = name.charAt(0).toUpperCase();
    
    return (
        <div 
            className="w-8 h-8 rounded-full flex items-center justify-center text-white font-medium"
            style={{ backgroundColor: color }}
        >
            {initial}
        </div>
    );
}
```

## 验收标准

- [ ] 会话顶部显示状态徽章
- [ ] 显示当前循环轮次
- [ ] 显示累计 Token 和成本
- [ ] 消息组显示正确的 Agent 名称
- [ ] Agent 头像正确显示
- [ ] 消息流式输出时显示 typing 指示器
- [ ] 完成的消息显示完成状态

## 文件清单

| 操作     | 文件路径                                          |
| :------- | :------------------------------------------------ |
| [NEW]    | `frontend/src/features/meeting/SessionHeader.tsx` |
| [MODIFY] | `frontend/src/components/chat/MessageGroup.tsx`   |
| [MODIFY] | `frontend/src/components/chat/ChatPanel.tsx`      |
| [NEW]    | `frontend/src/components/common/AgentAvatar.tsx`  |
| [NEW]    | `frontend/src/components/common/ProgressBar.tsx`  |
