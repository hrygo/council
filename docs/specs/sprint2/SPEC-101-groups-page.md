# SPEC-101: Groups 页面路由和布局

> **优先级**: P0 | **预估工时**: 3h  
> **关联 PRD**: F.1 群组管理体系

---

## 1. 路由配置

```typescript
// router.tsx
{ path: '/groups', element: <GroupsPage /> }
{ path: '/groups/:id', element: <GroupDetailPage /> }
```

---

## 2. 页面布局

```
┌─────────────────────────────────────────────────────┐
│ Header: 群组管理                      [+ 新建群组]  │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐             │
│  │ Group 1 │  │ Group 2 │  │ Group 3 │  ...        │
│  │ 🏢 SaaS │  │ 🏠 家庭 │  │ 💼 投资 │             │
│  └─────────┘  └─────────┘  └─────────┘             │
│                                                     │
└─────────────────────────────────────────────────────┘
```

---

## 3. 组件结构

```tsx
// pages/GroupsPage.tsx
export const GroupsPage: FC = () => {
  const { groups, isLoading, error } = useGroups();
  const [isCreateOpen, setIsCreateOpen] = useState(false);

  return (
    <PageContainer>
      <PageHeader 
        title={t('groups.title')}
        action={<Button onClick={() => setIsCreateOpen(true)}>+ 新建群组</Button>}
      />
      
      {isLoading && <Skeleton />}
      {error && <ErrorState error={error} />}
      {groups && <GroupList groups={groups} />}
      
      <CreateGroupModal 
        open={isCreateOpen} 
        onClose={() => setIsCreateOpen(false)} 
      />
    </PageContainer>
  );
};
```

---

## 4. 数据 Hook

```typescript
// hooks/useGroups.ts
export function useGroups() {
  return useQuery({
    queryKey: ['groups'],
    queryFn: () => fetch('/api/v1/groups').then(r => r.json()),
  });
}

export function useGroup(id: string) {
  return useQuery({
    queryKey: ['groups', id],
    queryFn: () => fetch(`/api/v1/groups/${id}`).then(r => r.json()),
  });
}

export function useCreateGroup() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateGroupInput) => 
      fetch('/api/v1/groups', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => queryClient.invalidateQueries(['groups']),
  });
}
```

---

## 5. 类型定义

```typescript
// types/group.ts
interface Group {
  id: string;
  name: string;
  icon: string;
  system_prompt: string;
  default_members: string[];  // Agent IDs
  created_at: string;
  updated_at: string;
}

interface CreateGroupInput {
  name: string;
  icon?: string;
  system_prompt?: string;
  default_members?: string[];
}
```

---

## 6. 测试要点

- [ ] 列表正确渲染
- [ ] 加载状态显示
- [ ] 错误状态处理
- [ ] 创建按钮可用
