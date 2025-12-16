# SPEC-303: DocumentReference 双向索引

> **优先级**: P1 | **预估工时**: 4h  
> **关联 PRD**: F.4.3 双向索引

---

## 1. 引用格式

```
[Ref: P3]       - 引用第 3 页
[Ref: L45-50]   - 引用第 45-50 行
```

---

## 2. 解析逻辑

```typescript
interface DocumentReference {
  type: 'page' | 'line';
  start: number;
  end?: number;
}

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

---

## 3. 引用链接渲染

```tsx
// 在 Markdown 渲染中替换引用为可点击链接
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

---

## 4. 文档阅读器

```tsx
// DocumentReader.tsx
const DocumentReader: FC = () => {
  const { document, highlightedRange, scrollTarget } = useDocumentStore();
  const contentRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    if (scrollTarget && contentRef.current) {
      const element = contentRef.current.querySelector(`[data-line="${scrollTarget}"]`);
      element?.scrollIntoView({ behavior: 'smooth', block: 'center' });
    }
  }, [scrollTarget]);
  
  return (
    <div ref={contentRef} className="p-4 overflow-auto">
      {document.lines.map((line, idx) => (
        <div
          key={idx}
          data-line={idx + 1}
          className={cn(
            highlightedRange && 
            idx + 1 >= highlightedRange.start && 
            idx + 1 <= (highlightedRange.end || highlightedRange.start) &&
            "bg-yellow-100"
          )}
        >
          {line}
        </div>
      ))}
    </div>
  );
};
```

---

## 5. Store 集成

```typescript
interface DocumentStore {
  document: { lines: string[] } | null;
  highlightedRange: { start: number; end?: number } | null;
  scrollTarget: number | null;
  
  scrollToReference: (ref: DocumentReference) => void;
  highlightRange: (start: number, end?: number) => void;
  clearHighlight: () => void;
}
```

---

## 6. 测试要点

- [ ] 引用解析正确
- [ ] 点击跳转
- [ ] 高亮显示
- [ ] 自动滚动
