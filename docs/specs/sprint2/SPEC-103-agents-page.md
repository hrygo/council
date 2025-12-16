# SPEC-103: Agents 页面路由和布局

> **优先级**: P0 | **预估工时**: 3h  
> **关联 PRD**: F.2 AI 理事工厂

---

## 1. 路由配置

```typescript
{ path: '/agents', element: <AgentsPage /> }
{ path: '/agents/:id', element: <AgentDetailPage /> }
```

---

## 2. 页面布局

```
┌─────────────────────────────────────────────────────────────┐
│ Header: Agent 工厂                   [Grid] [List] [+ 新建] │
├─────────────────────────────────────────────────────────────┤
│ 筛选: [所有] [CEO] [CFO] [CTO] ...        搜索: [______]   │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │  👔 CEO      │  │  💰 CFO      │  │  💻 CTO      │      │
│  │  战略决策    │  │  财务分析    │  │  技术架构    │      │
│  │  GPT-4      │  │  Claude-3.5  │  │  Gemini-Pro  │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. 组件结构

```tsx
export const AgentsPage: FC = () => {
  const { agents, isLoading } = useAgents();
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [filter, setFilter] = useState('');
  const [search, setSearch] = useState('');

  const filteredAgents = useMemo(() => 
    agents?.filter(a => 
      a.name.includes(search) || a.persona.includes(search)
    ), [agents, search]
  );

  return (
    <PageContainer>
      <PageHeader title="Agent 工厂">
        <ViewToggle value={viewMode} onChange={setViewMode} />
        <Button>+ 新建 Agent</Button>
      </PageHeader>
      
      <FilterBar>
        <SearchInput value={search} onChange={setSearch} />
      </FilterBar>
      
      {viewMode === 'grid' 
        ? <AgentGrid agents={filteredAgents} />
        : <AgentTable agents={filteredAgents} />
      }
    </PageContainer>
  );
};
```

---

## 4. 数据 Hook

```typescript
export function useAgents() {
  return useQuery({
    queryKey: ['agents'],
    queryFn: () => fetch('/api/v1/agents').then(r => r.json()),
  });
}

export function useCreateAgent() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateAgentInput) =>
      fetch('/api/v1/agents', { method: 'POST', body: JSON.stringify(data) }),
    onSuccess: () => queryClient.invalidateQueries(['agents']),
  });
}
```

---

## 5. 类型定义

```typescript
interface Agent {
  id: string;
  name: string;
  avatar: string;
  persona: string;          // 人设提示词
  model_config: ModelConfig;
  capabilities: {
    web_search: boolean;
    code_execution: boolean;
  };
  created_at: string;
  updated_at: string;
}

interface ModelConfig {
  provider: 'openai' | 'anthropic' | 'google' | 'deepseek';
  model: string;
  temperature: number;
  top_p: number;
  max_tokens: number;
}
```

---

## 6. 测试要点

- [ ] 视图切换正常
- [ ] 搜索过滤生效
- [ ] Agent 卡片点击进入编辑
- [ ] 创建按钮功能
