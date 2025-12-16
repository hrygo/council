# SPEC-205: 保存为模版功能

> **优先级**: P2 | **预估工时**: 2h  
> **关联 PRD**: F.3.2 保存模版

---

## 1. 触发方式

工作流编辑器工具栏添加 "保存为模版" 按钮：

```tsx
<Toolbar>
  // ...
  <Button variant="outline" onClick={() => setSaveModalOpen(true)}>
    <Save size={16} className="mr-1" /> 保存为模版
  </Button>
</Toolbar>
```

---

## 2. 弹窗设计

```tsx
export const SaveTemplateModal: FC<{ 
  open: boolean; 
  onClose: () => void; 
  currentGraph: GraphDefinition 
}> = ({ open, onClose, currentGraph }) => {
  const { mutate: createTemplate, isLoading } = useCreateTemplate();
  const [form, setForm] = useState({
    name: '',
    description: '',
    category: 'custom' as TemplateCategory,
  });

  const handleSave = () => {
    createTemplate({
      ...form,
      graph: currentGraph,
    }, { 
      onSuccess: () => {
        toast.success('模版保存成功');
        onClose();
      } 
    });
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>保存为模版</DialogTitle>
          <DialogDescription>
            将当前工作流保存为可复用的模版
          </DialogDescription>
        </DialogHeader>
        
        <div className="space-y-4">
          <Input
            label="模版名称"
            placeholder="如: 快速代码评审"
            value={form.name}
            onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
          />
          <Textarea
            label="描述"
            placeholder="描述此模版的用途..."
            value={form.description}
            onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
          />
          <Select
            label="分类"
            value={form.category}
            onChange={category => setForm(f => ({ ...f, category }))}
          >
            <SelectItem value="code_review">代码评审</SelectItem>
            <SelectItem value="business_plan">商业计划</SelectItem>
            <SelectItem value="quick_decision">快速决策</SelectItem>
            <SelectItem value="custom">自定义</SelectItem>
          </Select>
          
          {/* 预览 */}
          <div className="bg-gray-50 p-3 rounded text-sm">
            <strong>包含节点:</strong> {Object.keys(currentGraph.nodes).length} 个
          </div>
        </div>
        
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>取消</Button>
          <Button onClick={handleSave} disabled={!form.name || isLoading}>
            {isLoading ? <Spinner /> : '保存'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
```

---

## 3. API 调用

```typescript
const useCreateTemplate = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTemplateInput) =>
      fetch('/api/v1/templates', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      }),
    onSuccess: () => queryClient.invalidateQueries(['templates']),
  });
};
```

---

## 4. 测试要点

- [ ] 弹窗正确打开
- [ ] 表单验证 (名称必填)
- [ ] 保存成功提示
- [ ] 模版列表刷新
