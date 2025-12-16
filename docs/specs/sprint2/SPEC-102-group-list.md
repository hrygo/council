# SPEC-102: GroupList 组件 (CRUD)

> **优先级**: P0 | **预估工时**: 4h  
> **关联 PRD**: F.1.1 创建/管理群组

---

## 1. 组件结构

```
GroupList
├── GroupCard (per group)
│   ├── GroupIcon
│   ├── GroupInfo (name, member count)
│   └── GroupActions (edit, delete)
├── CreateGroupModal
└── EditGroupModal
```

---

## 2. GroupCard 组件

```tsx
interface GroupCardProps {
  group: Group;
  onEdit: () => void;
  onDelete: () => void;
}

export const GroupCard: FC<GroupCardProps> = ({ group, onEdit, onDelete }) => (
  <Card className="hover:shadow-md transition-shadow cursor-pointer">
    <CardHeader className="flex items-center gap-3">
      <GroupIcon icon={group.icon} size={40} />
      <div>
        <h3 className="font-medium">{group.name}</h3>
        <p className="text-sm text-gray-500">
          {group.default_members.length} 位成员
        </p>
      </div>
    </CardHeader>
    <CardFooter className="justify-end gap-2">
      <Button variant="ghost" size="sm" onClick={onEdit}>
        <Pencil size={14} />
      </Button>
      <Button variant="ghost" size="sm" onClick={onDelete}>
        <Trash2 size={14} />
      </Button>
    </CardFooter>
  </Card>
);
```

---

## 3. CreateGroupModal

```tsx
export const CreateGroupModal: FC<{ open: boolean; onClose: () => void }> = ({ open, onClose }) => {
  const { mutate: createGroup, isLoading } = useCreateGroup();
  const [form, setForm] = useState({ name: '', icon: '', system_prompt: '' });

  const handleSubmit = () => {
    createGroup(form, { onSuccess: onClose });
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>创建群组</DialogTitle>
        </DialogHeader>
        
        <div className="space-y-4">
          <Input 
            label="群组名称" 
            value={form.name} 
            onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
          />
          <IconPicker 
            value={form.icon} 
            onChange={icon => setForm(f => ({ ...f, icon }))}
          />
          <Textarea 
            label="群定位 (System Prompt)" 
            placeholder="定义该群的底层逻辑与价值观..."
            value={form.system_prompt}
            onChange={e => setForm(f => ({ ...f, system_prompt: e.target.value }))}
          />
          <AgentMultiSelect
            label="默认成员"
            value={form.default_members}
            onChange={members => setForm(f => ({ ...f, default_members: members }))}
          />
        </div>
        
        <DialogFooter>
          <Button variant="outline" onClick={onClose}>取消</Button>
          <Button onClick={handleSubmit} disabled={isLoading}>
            {isLoading ? <Spinner /> : '创建'}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
```

---

## 4. 删除确认

```tsx
const handleDelete = (group: Group) => {
  confirm({
    title: '确认删除',
    description: `确定要删除群组 "${group.name}" 吗？此操作不可恢复。`,
    confirmText: '删除',
    variant: 'destructive',
    onConfirm: () => deleteGroup(group.id),
  });
};
```

---

## 5. 图标选择器

```tsx
const groupIcons = ['🏢', '🏠', '💼', '🎯', '⚙️', '📊', '🧪', '🎨'];

export const IconPicker: FC<{ value: string; onChange: (v: string) => void }> = 
  ({ value, onChange }) => (
    <div className="grid grid-cols-8 gap-2">
      {groupIcons.map(icon => (
        <button
          key={icon}
          className={cn(
            "p-2 rounded hover:bg-gray-100",
            value === icon && "ring-2 ring-blue-500"
          )}
          onClick={() => onChange(icon)}
        >
          {icon}
        </button>
      ))}
    </div>
  );
```

---

## 6. 测试要点

- [ ] 创建群组成功
- [ ] 编辑群组成功
- [ ] 删除确认弹窗
- [ ] 表单验证 (名称必填)
- [ ] Agent 多选正确
