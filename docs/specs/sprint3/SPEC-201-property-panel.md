# SPEC-201: PropertyPanel 节点属性面板

> **优先级**: P0 | **预估工时**: 4h  
> **关联 PRD**: F.3.1 节点编辑器

---

## 1. 布局

```
┌─────────────────────────────┐
│ 节点属性                    │
│ ─────────────────────────── │
│ 类型: [Agent Node]          │
│ 名称: [________________]    │
│                             │
│ ▼ 基础配置                  │
│   Agent: [CEO ▾]            │
│   输入源: [上一节点输出 ▾]  │
│                             │
│ ▼ 高级配置                  │
│   超时: [30] 秒             │
│   重试: [3] 次              │
│                             │
│ [删除节点]                  │
└─────────────────────────────┘
```

---

## 2. 接口

```typescript
interface PropertyPanelProps {
  node: Node | null;
  onUpdate: (nodeId: string, data: Partial<NodeData>) => void;
  onDelete: (nodeId: string) => void;
}
```

---

## 3. 动态表单

```tsx
const nodeFormConfigs: Record<NodeType, FormFieldConfig[]> = {
  agent: [
    { key: 'agent_id', type: 'select', label: 'Agent', options: 'agents' },
    { key: 'input_source', type: 'select', label: '输入源' },
  ],
  vote: [
    { key: 'threshold', type: 'slider', label: '通过阈值', min: 0.5, max: 1.0 },
    { key: 'vote_type', type: 'radio', label: '投票类型', options: ['yes_no', 'score'] },
  ],
  loop: [
    { key: 'max_rounds', type: 'number', label: '最大轮数', min: 1, max: 10 },
    { key: 'exit_condition', type: 'select', label: '退出条件' },
  ],
  // ...
};

export const PropertyPanel: FC<PropertyPanelProps> = ({ node, onUpdate, onDelete }) => {
  if (!node) return <EmptyState message="选择节点以编辑属性" />;
  
  const formConfig = nodeFormConfigs[node.type];
  
  return (
    <div className="p-4 space-y-4">
      <h3 className="font-medium">节点属性</h3>
      <Badge>{nodeTypeLabels[node.type]}</Badge>
      
      <Input 
        label="名称"
        value={node.data.label}
        onChange={e => onUpdate(node.id, { label: e.target.value })}
      />
      
      {formConfig.map(field => (
        <DynamicFormField 
          key={field.key}
          config={field}
          value={node.data[field.key]}
          onChange={v => onUpdate(node.id, { [field.key]: v })}
        />
      ))}
      
      <Button variant="destructive" onClick={() => onDelete(node.id)}>
        删除节点
      </Button>
    </div>
  );
};
```

---

## 4. 测试要点

- [ ] 选中节点时面板显示
- [ ] 属性修改实时同步到画布
- [ ] 删除节点后面板清空
- [ ] 不同节点类型显示不同表单
