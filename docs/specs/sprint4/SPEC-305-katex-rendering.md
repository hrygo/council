# SPEC-305: KaTeX 公式渲染

> **优先级**: P2 | **预估工时**: 2h  
> **关联 PRD**: F.4.2 公式渲染

---

## 1. 依赖安装

```bash
npm install remark-math rehype-katex katex
```

---

## 2. CSS 引入

```tsx
// main.tsx 或 App.tsx
import 'katex/dist/katex.min.css';
```

---

## 3. Markdown 渲染配置

```tsx
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeKatex from 'rehype-katex';
import rehypeHighlight from 'rehype-highlight';

export const MarkdownRenderer: FC<{ content: string }> = ({ content }) => (
  <ReactMarkdown
    remarkPlugins={[remarkGfm, remarkMath]}
    rehypePlugins={[rehypeKatex, rehypeHighlight]}
    className="prose prose-sm max-w-none"
  >
    {content}
  </ReactMarkdown>
);
```

---

## 4. 公式语法

### 行内公式
```markdown
爱因斯坦质能方程 $E = mc^2$ 是相对论的核心。
```

### 块级公式
```markdown
$$
\int_{-\infty}^{\infty} e^{-x^2} dx = \sqrt{\pi}
$$
```

---

## 5. 错误处理

```tsx
const SafeKatexRenderer: FC<{ content: string }> = ({ content }) => {
  try {
    return <MarkdownRenderer content={content} />;
  } catch (error) {
    // 公式解析错误时降级显示原始文本
    return (
      <div className="prose prose-sm">
        {content}
        <Alert variant="warning" className="mt-2">
          公式渲染失败，显示原始内容
        </Alert>
      </div>
    );
  }
};
```

---

## 6. 测试要点

- [ ] 行内公式渲染
- [ ] 块级公式渲染
- [ ] 复杂公式支持 (矩阵、积分)
- [ ] 错误公式降级处理
