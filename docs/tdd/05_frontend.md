# 5. 前端架构关键点 (Frontend Specifics)

### 5.1 弹性布局实现

使用 `react-resizable-panels`。

```tsx
<PanelGroup direction="horizontal">
  <Panel defaultSize={20} minSize={5} collapsible>
    <WorkflowCanvas readOnly={isRunning} />
  </Panel>
  <PanelResizeHandle />
  <Panel defaultSize={50} minSize={30}>
    <ChatStreamWindow />
  </Panel>
  <PanelResizeHandle />
  <Panel defaultSize={30} minSize={0} collapsible>
    <DocumentReader />
  </Panel>
</PanelGroup>
```

### 5.2 状态管理 (Zustand)

需要管理极其复杂的运行时状态：

```typescript
interface SessionState {
  nodes: Node[];
  edges: Edge[];
  activeNodeIds: string[]; // 当前高亮的节点 (可能多个)
  messages: {
    [nodeId: string]: string; // 增量存储每个节点的输出内容
  };
  layout: {
    leftPanelCollapsed: boolean;
    rightPanelCollapsed: boolean;
  };
}
```

### 5.3 国际化 (i18n)

采用 `react-i18next` 实现中英双语切换，支持运行时语言切换和类型安全。

**初始化配置 (`frontend/src/i18n/index.ts`)：**

```typescript
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// 导入翻译资源
import zhCN from './locales/zh-CN';
import enUS from './locales/en-US';

i18n
  .use(LanguageDetector)      // 自动检测浏览器语言
  .use(initReactI18next)
  .init({
    resources: {
      'zh-CN': zhCN,
      'en-US': enUS,
    },
    fallbackLng: 'zh-CN',     // 默认中文
    interpolation: {
      escapeValue: false,     // React 已处理 XSS
    },
    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage'],
    },
  });

export default i18n;
```

**类型安全定义 (`frontend/src/i18n/types.d.ts`)：**

```typescript
import 'react-i18next';
import type zhCN from './locales/zh-CN';

declare module 'react-i18next' {
  interface CustomTypeOptions {
    defaultNS: 'common';
    resources: typeof zhCN;  // 以中文为基准类型
  }
}
```

**Zustand Store 集成语言状态：**

```typescript
interface ConfigState {
  language: 'zh-CN' | 'en-US';
  godMode: boolean; // 🆕 上帝模式开关
  setLanguage: (lang: 'zh-CN' | 'en-US') => void;
  toggleGodMode: () => void;
}

export const useConfigStore = create<ConfigState>((set) => ({
  language: (localStorage.getItem('i18nextLng') as 'zh-CN' | 'en-US') || 'zh-CN',
  setLanguage: (lang) => {
    i18n.changeLanguage(lang);
    localStorage.setItem('i18nextLng', lang);
  },
  godMode: false,
  toggleGodMode: () => set((state) => ({ godMode: !state.godMode })),
}));
```

### 5.4 并行消息渲染 (Parallel Message UI)

对应 PRD F.4.2，当处于并行节点时多个 Agent 消息并排显示。

**Markdown 渲染技术栈**：

```typescript
// 使用 react-markdown + 插件组合，支持代码块、表格、公式
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';           // GitHub 风格 Markdown (表格、任务列表)
import remarkMath from 'remark-math';         // 数学公式语法
import rehypeKatex from 'rehype-katex';       // LaTeX 公式渲染
import rehypeHighlight from 'rehype-highlight'; // 代码高亮

const MarkdownRenderer: FC<{ content: string }> = ({ content }) => (
  <ReactMarkdown
    remarkPlugins={[remarkGfm, remarkMath]}
    rehypePlugins={[rehypeKatex, rehypeHighlight]}
  >
    {content}
  </ReactMarkdown>
);
```

```tsx
interface ParallelMessageGroup {
  nodeId: string;
  nodeType: 'parallel';
  messages: AgentMessage[];
}

const ParallelMessageRow: FC<{ group: ParallelMessageGroup }> = ({ group }) => {
  return (
    <div className="flex gap-4 w-full">
      {group.messages.map((msg) => (
        <div key={msg.agentId} className="flex-1 min-w-0">
          <AgentAvatar agent={msg.agent} />
          <MessageBubble 
            content={msg.content} 
            isStreaming={msg.isStreaming}
          />
        </div>
      ))}
    </div>
  );
};

// 消息流渲染逻辑
const ChatStreamWindow: FC = () => {
  const { messageGroups } = useSessionStore();
  
  return (
    <div className="flex flex-col gap-6">
      {messageGroups.map((group) => 
        group.nodeType === 'parallel' 
          ? <ParallelMessageRow key={group.nodeId} group={group} />
          : <SequentialMessage key={group.nodeId} message={group.messages[0]} />
      )}
    </div>
  );
};
```

### 5.5 双向文档索引 (Bidirectional Document Reference)

对应 PRD F.4.3，实现 AI 发言与源文档的双向跳转。

**引用格式约定：**

```typescript
// AI 输出中的引用标记格式
// [Ref: P3] 表示引用第3页
// [Ref: L45-50] 表示引用第45-50行

interface DocumentReference {
  type: 'page' | 'line';
  start: number;
  end?: number;
}

// 解析 AI 输出中的引用
const parseReferences = (content: string): DocumentReference[] => {
  const regex = /\[Ref:\s*(P|L)(\d+)(?:-(\d+))?\]/g;
  const refs: DocumentReference[] = [];
  let match;
  
  while ((match = regex.exec(content)) !== null) {
    refs.push({
      type: match[1] === 'P' ? 'page' : 'line',
      start: parseInt(match[2]),
      end: match[3] ? parseInt(match[3]) : undefined,
    });
  }
  return refs;
};
```

**引用点击处理：**

```tsx
const ReferenceLink: FC<{ ref: DocumentReference }> = ({ ref }) => {
  const { scrollToReference, highlightRange } = useDocumentStore();
  
  const handleClick = () => {
    scrollToReference(ref);
    highlightRange(ref.start, ref.end);
  };
  
  return (
    <button 
      onClick={handleClick}
      className="text-blue-500 hover:underline cursor-pointer"
    >
      [Ref: {ref.type === 'page' ? 'P' : 'L'}{ref.start}{ref.end ? `-${ref.end}` : ''}]
    </button>
  );
};
```

### 5.6 布局状态持久化 (Layout Persistence)

对应 PRD F.4.0 状态记忆需求。

```typescript
interface LayoutState {
  panelSizes: [number, number, number]; // 左中右三栏比例
  leftCollapsed: boolean;
  rightCollapsed: boolean;
  maximizedPanel: 'left' | 'center' | 'right' | null;
}

// 使用 Zustand persist 中间件
export const useLayoutStore = create<LayoutState>()(
  persist(
    (set) => ({
      panelSizes: [20, 50, 30],
      leftCollapsed: false,
      rightCollapsed: false,
      maximizedPanel: null,
      
      setPanelSizes: (sizes: [number, number, number]) => set({ panelSizes: sizes }),
      toggleLeftPanel: () => set((s) => ({ leftCollapsed: !s.leftCollapsed })),
      toggleRightPanel: () => set((s) => ({ rightCollapsed: !s.rightCollapsed })),
      maximizePanel: (panel: 'left' | 'center' | 'right' | null) => set({ maximizedPanel: panel }),
    }),
    {
      name: 'council-layout',
      storage: createJSONStorage(() => localStorage),
    }
  )
);
```

### 5.7 全屏专注模式 (Fullscreen Focus Mode)

对应 PRD F.4.0 全屏专注需求，任意栏位可进入沉浸模式。

**最大化按钮组件：**

```tsx
import { Maximize2, Minimize2 } from 'lucide-react';

const PanelMaximizeButton: FC<{ panel: 'left' | 'center' | 'right' }> = ({ panel }) => {
  const { maximizedPanel, maximizePanel } = useLayoutStore();
  const isMaximized = maximizedPanel === panel;
  
  return (
    <button
      onClick={() => maximizePanel(isMaximized ? null : panel)}
      className="p-1.5 rounded hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
      title={isMaximized ? "退出全屏" : "全屏专注"}
    >
      {isMaximized ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
    </button>
  );
};
```

**会议室主布局处理：**

```tsx
const MeetingRoom: FC = () => {
  const { maximizedPanel, panelSizes, leftCollapsed, rightCollapsed } = useLayoutStore();
  
  // 全屏模式：只渲染单个面板
  if (maximizedPanel) {
    const panelMap = {
      left: <WorkflowCanvas fullscreen onExitFullscreen={() => maximizePanel(null)} />,
      center: <ChatStreamWindow fullscreen onExitFullscreen={() => maximizePanel(null)} />,
      right: <DocumentReader fullscreen onExitFullscreen={() => maximizePanel(null)} />,
    };
    
    return (
      <div className="h-screen w-screen">
        {panelMap[maximizedPanel]}
      </div>
    );
  }
  
  // 正常三栏布局
  return (
    <PanelGroup direction="horizontal" onLayout={setPanelSizes}>
      <Panel 
        defaultSize={panelSizes[0]} 
        minSize={5} 
        collapsible 
        collapsed={leftCollapsed}
      >
        <div className="relative h-full">
          <PanelMaximizeButton panel="left" />
          <WorkflowCanvas readOnly={isRunning} />
        </div>
      </Panel>
      
      <PanelResizeHandle className="w-1 bg-gray-200 hover:bg-blue-400 transition-colors" />
      
      <Panel defaultSize={panelSizes[1]} minSize={30}>
        <div className="relative h-full">
          <PanelMaximizeButton panel="center" />
          <ChatStreamWindow />
        </div>
      </Panel>
      
      <PanelResizeHandle className="w-1 bg-gray-200 hover:bg-blue-400 transition-colors" />
      
      <Panel 
        defaultSize={panelSizes[2]} 
        minSize={0} 
        collapsible 
        collapsed={rightCollapsed}
      >
        <div className="relative h-full">
          <PanelMaximizeButton panel="right" />
          <DocumentReader />
        </div>
      </Panel>
    </PanelGroup>
  );
};
```

**键盘快捷键支持：**

```typescript
// useFullscreenShortcuts.ts
import { useEffect } from 'react';
import { useLayoutStore } from '@/stores/layoutStore';

export const useFullscreenShortcuts = () => {
  const { maximizedPanel, maximizePanel } = useLayoutStore();
  
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Escape 退出全屏
      if (e.key === 'Escape' && maximizedPanel) {
        maximizePanel(null);
        return;
      }
      
      // Cmd/Ctrl + 1/2/3 切换全屏
      if ((e.metaKey || e.ctrlKey) && !e.shiftKey) {
        switch (e.key) {
          case '1':
            maximizePanel(maximizedPanel === 'left' ? null : 'left');
            break;
          case '2':
            maximizePanel(maximizedPanel === 'center' ? null : 'center');
            break;
          case '3':
            maximizePanel(maximizedPanel === 'right' ? null : 'right');
            break;
        }
      }
    };
    
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [maximizedPanel, maximizePanel]);
};

### 5.8 人类裁决交互 (Human Review UI)

对应 PRD F.3.1 HumanReview 节点。

```tsx
const HumanReviewModal: FC<{ draft: string; onApprove: (content: string) => void; onReject: (reason: string) => void }> = ({ draft, onApprove, onReject }) => {
  const [content, setContent] = useState(draft);
  const [rejectReason, setRejectReason] = useState("");
  const [isRejecting, setIsRejecting] = useState(false);

  return (
    <Dialog open={true}>
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>🛡️ 需要人类裁决 (Human Review Required)</DialogTitle>
          <DialogDescription>
             AI 已生成决策草案，请仔细审查。您具有最终决定权。
          </DialogDescription>
        </DialogHeader>
        
        {isRejecting ? (
           <div className="space-y-4">
              <Textarea 
                placeholder="请输入驳回理由..." 
                value={rejectReason}
                onChange={e => setRejectReason(e.target.value)}
              />
              <div className="flex justify-end gap-2">
                <Button variant="ghost" onClick={() => setIsRejecting(false)}>返回</Button>
                <Button variant="destructive" onClick={() => onReject(rejectReason)}>确认驳回</Button>
              </div>
           </div>
        ) : (
           <div className="space-y-4">
              <Textarea 
                className="min-h-[300px] font-mono"
                value={content}
                onChange={e => setContent(e.target.value)}
              />
              <div className="flex justify-end gap-2">
                <Button variant="outline" onClick={() => setIsRejecting(true)}>驳回</Button>
                <Button onClick={() => onApprove(content)}>签署并通过</Button>
              </div>
           </div>
        )}
      </DialogContent>
    </Dialog>
  );
};
```
```
