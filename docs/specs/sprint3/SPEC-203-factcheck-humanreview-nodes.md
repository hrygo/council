# SPEC-203: FactCheck/HumanReview 节点 UI

> **优先级**: P1 | **预估工时**: 3h  
> **关联 PRD**: F.3.1 FactCheck, HumanReview

---

## 1. FactCheck 节点

### 1.1 节点样式

```tsx
const FactCheckNodeIcon = () => <Search className="text-teal-500" />;

const factCheckNodeStyle = {
  background: 'linear-gradient(135deg, #CCFBF1, #99F6E4)',
  border: '2px solid #14B8A6',
};
```

### 1.2 属性配置

```typescript
interface FactCheckNodeData {
  label: string;
  search_sources: ('tavily' | 'serper' | 'local_kb')[];
  max_queries: number;        // 最大搜索次数
  verify_threshold: number;   // 置信度阈值
}
```

### 1.3 属性面板

```tsx
const FactCheckNodeProperties: FC<{ data: FactCheckNodeData; onChange: ... }> = 
  ({ data, onChange }) => (
    <>
      <CheckboxGroup
        label="搜索源"
        value={data.search_sources}
        onChange={search_sources => onChange({ search_sources })}
        options={[
          { value: 'tavily', label: '🌐 Tavily (联网)' },
          { value: 'serper', label: '🔍 Serper (联网)' },
          { value: 'local_kb', label: '📚 本地知识库' },
        ]}
      />
      <NumberInput
        label="最大搜索次数"
        min={1} max={10}
        value={data.max_queries}
        onChange={max_queries => onChange({ max_queries })}
      />
      <Slider
        label={`置信度阈值: ${Math.round(data.verify_threshold * 100)}%`}
        min={50} max={100}
        value={data.verify_threshold * 100}
        onChange={v => onChange({ verify_threshold: v / 100 })}
      />
    </>
  );
```

---

## 2. HumanReview 节点

### 2.1 节点样式

```tsx
const HumanReviewNodeIcon = () => <UserCheck className="text-rose-500" />;

const humanReviewNodeStyle = {
  background: 'linear-gradient(135deg, #FFE4E6, #FECDD3)',
  border: '2px solid #F43F5E',
};
```

### 2.2 属性配置

```typescript
interface HumanReviewNodeData {
  label: string;
  review_type: 'approve_reject' | 'edit_content';
  timeout_minutes: number;    // 超时时间
  allow_skip: boolean;        // 是否允许跳过
}
```

### 2.3 属性面板

```tsx
const HumanReviewNodeProperties: FC<{ data: HumanReviewNodeData; onChange: ... }> = 
  ({ data, onChange }) => (
    <>
      <Alert variant="warning" className="mb-4">
        ⚠️ 此节点将暂停工作流，等待人类审核
      </Alert>
      <Select
        label="审核类型"
        value={data.review_type}
        onChange={review_type => onChange({ review_type })}
      >
        <SelectItem value="approve_reject">通过/驳回</SelectItem>
        <SelectItem value="edit_content">编辑内容</SelectItem>
      </Select>
      <NumberInput
        label="超时时间 (分钟)"
        min={5} max={60}
        value={data.timeout_minutes}
        onChange={timeout_minutes => onChange({ timeout_minutes })}
      />
      <Switch
        label="允许跳过"
        description="若超时，自动通过"
        checked={data.allow_skip}
        onChange={allow_skip => onChange({ allow_skip })}
      />
    </>
  );
```

---

## 3. 默认值

```typescript
export const nodeDefaults = {
  fact_check: () => ({
    label: '事实核查',
    search_sources: ['tavily'],
    max_queries: 3,
    verify_threshold: 0.7,
  }),
  human_review: () => ({
    label: '人类裁决',
    review_type: 'approve_reject',
    timeout_minutes: 30,
    allow_skip: false,
  }),
};
```

---

## 4. 测试要点

- [ ] 节点样式正确
- [ ] 搜索源多选生效
- [ ] 超时设置保存
- [ ] 强制安全节点提示
