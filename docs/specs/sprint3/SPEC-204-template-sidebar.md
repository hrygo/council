# SPEC-204: 模版库侧边栏

> **优先级**: P2 | **预估工时**: 3h  
> **关联 PRD**: F.3.2 模版库

---

## 1. 布局

```
┌─────────────────────┐
│ 模版库          [×] │
├─────────────────────┤
│ 筛选: [全部 ▾]      │
├─────────────────────┤
│ 📦 系统模版         │
│ ┌─────────────────┐ │
│ │ 🔍 代码评审     │ │
│ │ 3 Agent · 5节点 │ │
│ └─────────────────┘ │
│ ┌─────────────────┐ │
│ │ 📊 商业计划压测 │ │
│ │ 3 Agent · 7节点 │ │
│ └─────────────────┘ │
├─────────────────────┤
│ 👤 我的模版         │
│ ┌─────────────────┐ │
│ │ 快速决策 v2     │ │
│ │ 2 Agent · 4节点 │ │
│ └─────────────────┘ │
└─────────────────────┘
```

---

## 2. 接口

```typescript
interface TemplateSidebarProps {
  open: boolean;
  onClose: () => void;
  onApply: (template: Template) => void;
}
```

---

## 3. 实现

```tsx
export const TemplateSidebar: FC<TemplateSidebarProps> = ({ open, onClose, onApply }) => {
  const { data: templates, isLoading } = useTemplates();
  const [filter, setFilter] = useState<'all' | 'system' | 'custom'>('all');

  const systemTemplates = templates?.filter(t => t.is_system) || [];
  const customTemplates = templates?.filter(t => !t.is_system) || [];

  return (
    <Sheet open={open} onOpenChange={onClose} side="left">
      <SheetContent className="w-[300px]">
        <SheetHeader>
          <SheetTitle>模版库</SheetTitle>
        </SheetHeader>
        
        <Select value={filter} onChange={setFilter}>
          <SelectItem value="all">全部</SelectItem>
          <SelectItem value="system">系统模版</SelectItem>
          <SelectItem value="custom">我的模版</SelectItem>
        </Select>
        
        {isLoading && <Skeleton count={3} />}
        
        {(filter === 'all' || filter === 'system') && (
          <Section title="📦 系统模版">
            {systemTemplates.map(t => (
              <TemplateCard key={t.id} template={t} onApply={() => onApply(t)} />
            ))}
          </Section>
        )}
        
        {(filter === 'all' || filter === 'custom') && (
          <Section title="👤 我的模版">
            {customTemplates.length === 0 
              ? <EmptyState message="暂无自定义模版" />
              : customTemplates.map(t => (
                  <TemplateCard key={t.id} template={t} onApply={() => onApply(t)} showDelete />
                ))
            }
          </Section>
        )}
      </SheetContent>
    </Sheet>
  );
};
```

---

## 4. TemplateCard

```tsx
const TemplateCard: FC<{ template: Template; onApply: () => void; showDelete?: boolean }> = 
  ({ template, onApply, showDelete }) => (
    <Card className="hover:bg-gray-50 cursor-pointer" onClick={onApply}>
      <CardContent className="p-3">
        <div className="flex items-center gap-2">
          <span className="text-lg">{categoryIcons[template.category]}</span>
          <div className="flex-1 min-w-0">
            <h4 className="font-medium truncate">{template.name}</h4>
            <p className="text-xs text-gray-500">{template.description}</p>
          </div>
          {showDelete && (
            <Button size="icon" variant="ghost" onClick={e => { e.stopPropagation(); handleDelete(template.id); }}>
              <Trash2 size={14} />
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
```

---

## 5. 测试要点

- [ ] 模版列表加载
- [ ] 筛选功能
- [ ] 点击应用模版
- [ ] 删除自定义模版
