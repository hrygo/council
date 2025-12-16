# SPEC-304: 全屏快捷键

> **优先级**: P2 | **预估工时**: 2h  
> **关联 PRD**: F.4.0 全屏专注

---

## 1. 快捷键定义

| 快捷键         | 功能             |
| -------------- | ---------------- |
| `Cmd/Ctrl + 1` | 切换左侧面板全屏 |
| `Cmd/Ctrl + 2` | 切换中间面板全屏 |
| `Cmd/Ctrl + 3` | 切换右侧面板全屏 |
| `Escape`       | 退出全屏         |

---

## 2. Hook 实现

```typescript
// hooks/useFullscreenShortcuts.ts
import { useEffect } from 'react';
import { useLayoutStore } from '@/stores/useLayoutStore';

export const useFullscreenShortcuts = () => {
  const { maximizedPanel, maximizePanel } = useLayoutStore();
  
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Escape 退出全屏
      if (e.key === 'Escape' && maximizedPanel) {
        maximizePanel(null);
        return;
      }
      
      // Cmd/Ctrl + 数字键
      if ((e.metaKey || e.ctrlKey) && !e.shiftKey && !e.altKey) {
        const panelMap: Record<string, 'left' | 'center' | 'right'> = {
          '1': 'left',
          '2': 'center',
          '3': 'right',
        };
        
        const panel = panelMap[e.key];
        if (panel) {
          e.preventDefault();
          maximizePanel(maximizedPanel === panel ? null : panel);
        }
      }
    };
    
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [maximizedPanel, maximizePanel]);
};
```

---

## 3. 使用

```tsx
// MeetingRoom.tsx
export const MeetingRoom: FC = () => {
  useFullscreenShortcuts();
  // ...
};
```

---

## 4. 快捷键提示

```tsx
const ShortcutHint: FC = () => {
  const isMac = navigator.platform.includes('Mac');
  const mod = isMac ? '⌘' : 'Ctrl';
  
  return (
    <div className="text-xs text-gray-400 space-x-4">
      <span>{mod}+1 工作流</span>
      <span>{mod}+2 对话</span>
      <span>{mod}+3 文档</span>
      <span>Esc 退出</span>
    </div>
  );
};
```

---

## 5. 测试要点

- [ ] 各快捷键生效
- [ ] 切换正确
- [ ] Escape 退出
- [ ] 不与浏览器快捷键冲突
