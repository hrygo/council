# SPEC-503: 性能优化 (Performance Optimization)

> **优先级**: P2  
> **类型**: Refactor  
> **预估工时**: 6h

## 1. 概述

优化前端应用性能，提升首屏加载速度和运行时流畅度。

## 2. 目标

- 首屏加载时间 (FCP) < 1.5s
- 首次可交互时间 (TTI) < 2.5s
- Lighthouse Performance 分数 > 90
- 聊天列表支持虚拟滚动 (10000+ 消息无卡顿)

## 3. 优化策略

### 3.1 代码分割 (Route-based Code Splitting)

**Before:**
```typescript
import { WorkflowBuilder } from './features/WorkflowBuilder';
import { MeetingRoom } from './features/MeetingRoom';
import { AgentsPage } from './pages/AgentsPage';
```

**After:**
```typescript
import { lazy, Suspense } from 'react';

const WorkflowBuilder = lazy(() => import('./features/WorkflowBuilder'));
const MeetingRoom = lazy(() => import('./features/MeetingRoom'));
const AgentsPage = lazy(() => import('./pages/AgentsPage'));

// 在 Router 中
<Suspense fallback={<LoadingSpinner />}>
  <Routes>
    <Route path="/builder/*" element={<WorkflowBuilder />} />
    <Route path="/meeting/*" element={<MeetingRoom />} />
    <Route path="/agents" element={<AgentsPage />} />
  </Routes>
</Suspense>
```

### 3.2 Bundle 分析与优化

```bash
# 分析 bundle
npm run build -- --analyze

# 预期优化目标
Initial Bundle: < 150KB (gzipped)
Route Chunks: < 50KB each
```

**大型依赖外置:**
```typescript
// vite.config.ts
build: {
  rollupOptions: {
    output: {
      manualChunks: {
        'react-vendor': ['react', 'react-dom', 'react-router-dom'],
        'flow-vendor': ['@xyflow/react'],
        'markdown': ['react-markdown', 'rehype-katex', 'remark-math'],
      },
    },
  },
},
```

### 3.3 虚拟列表优化

**聊天面板使用虚拟滚动:**

```typescript
import { useVirtualizer } from '@tanstack/react-virtual';

function ChatMessageList({ messages }: { messages: Message[] }) {
  const parentRef = useRef<HTMLDivElement>(null);
  
  const virtualizer = useVirtualizer({
    count: messages.length,
    getScrollElement: () => parentRef.current,
    estimateSize: () => 80, // 估计每条消息高度
    overscan: 5,
  });

  return (
    <div ref={parentRef} className="h-full overflow-auto">
      <div style={{ height: virtualizer.getTotalSize() }}>
        {virtualizer.getVirtualItems().map((item) => (
          <MessageBubble 
            key={item.key}
            message={messages[item.index]}
            style={{ transform: `translateY(${item.start}px)` }}
          />
        ))}
      </div>
    </div>
  );
}
```

### 3.4 图片与资源优化

```typescript
// 使用 WebP 格式
import agentAvatar from './assets/agent-avatar.webp';

// 懒加载图片
<img loading="lazy" src={avatarUrl} alt="Agent" />
```

### 3.5 React 渲染优化

```typescript
// 使用 memo 避免不必要的重渲染
const MessageBubble = memo(({ message }: Props) => {
  return <div>{message.content}</div>;
});

// Store 选择器细化
const nodeStatus = useWorkflowRunStore(
  (state) => state.nodes.find(n => n.id === nodeId)?.data.status
);
```

## 4. 性能预算 (Performance Budget)

| 指标       | 当前   | 目标    |
| :--------- | :----- | :------ |
| Initial JS | ~300KB | < 150KB |
| LCP        | ~2.5s  | < 1.5s  |
| FCP        | ~2.0s  | < 1.0s  |
| TTI        | ~3.5s  | < 2.5s  |
| CLS        | 0.15   | < 0.1   |

## 5. 监控与度量

```typescript
// 使用 Web Vitals 采集
import { onLCP, onFCP, onCLS, onTTFB } from 'web-vitals';

function sendToAnalytics(metric: Metric) {
  console.log(metric.name, metric.value);
  // 发送到监控服务
}

onLCP(sendToAnalytics);
onFCP(sendToAnalytics);
onCLS(sendToAnalytics);
onTTFB(sendToAnalytics);
```

## 6. 验收标准

- [ ] Lighthouse Performance > 90
- [ ] 首屏加载 < 2s (3G 网络模拟)
- [ ] 10000 条消息列表滚动流畅 (60fps)
- [ ] Bundle 大小符合预算

## 7. Makefile 集成

```makefile
# Performance
perf-analyze: ## 📊 Analyze bundle size
	@cd frontend && npm run build -- --analyze

perf-lighthouse: ## 🔦 Run Lighthouse audit
	@npx lighthouse http://localhost:5173 --output=html --output-path=./lighthouse-report.html
```
