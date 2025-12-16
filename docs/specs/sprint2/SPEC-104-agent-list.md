# SPEC-104: AgentList 组件 (CRUD)

> **优先级**: P0 | **预估工时**: 4h  
> **关联 PRD**: F.2.1 角色定义, F.2.2 模型配置

---

## 1. AgentCard 组件

```tsx
interface AgentCardProps {
  agent: Agent;
  onClick: () => void;
}

export const AgentCard: FC<AgentCardProps> = ({ agent, onClick }) => (
  <Card 
    className="hover:shadow-lg transition-all cursor-pointer group"
    onClick={onClick}
  >
    <CardContent className="p-4">
      <div className="flex items-center gap-3">
        <Avatar src={agent.avatar} fallback={agent.name[0]} size={48} />
        <div className="flex-1 min-w-0">
          <h3 className="font-medium truncate">{agent.name}</h3>
          <p className="text-sm text-gray-500 truncate">{agent.persona}</p>
        </div>
      </div>
      
      <div className="mt-3 flex items-center gap-2 text-xs text-gray-400">
        <ProviderBadge provider={agent.model_config.provider} />
        <span>{agent.model_config.model}</span>
      </div>
      
      {/* 能力标签 */}
      <div className="mt-2 flex gap-1">
        {agent.capabilities.web_search && (
          <Badge variant="outline" size="sm">🔍 联网</Badge>
        )}
        {agent.capabilities.code_execution && (
          <Badge variant="outline" size="sm">💻 代码</Badge>
        )}
      </div>
    </CardContent>
  </Card>
);
```

---

## 2. AgentEditDrawer

```tsx
export const AgentEditDrawer: FC<{ 
  agent: Agent | null; 
  open: boolean; 
  onClose: () => void 
}> = ({ agent, open, onClose }) => {
  const isNew = !agent;
  const { mutate: saveAgent, isLoading } = isNew ? useCreateAgent() : useUpdateAgent();
  
  const [form, setForm] = useState<AgentFormData>(
    agent || defaultAgentForm
  );

  return (
    <Sheet open={open} onOpenChange={onClose}>
      <SheetContent className="w-[500px]">
        <SheetHeader>
          <SheetTitle>{isNew ? '创建 Agent' : '编辑 Agent'}</SheetTitle>
        </SheetHeader>
        
        <div className="space-y-6 py-6">
          {/* 基本信息 */}
          <Section title="基本信息">
            <Input label="名称" value={form.name} onChange={...} />
            <AvatarUpload value={form.avatar} onChange={...} />
            <Textarea 
              label="人设提示词 (Persona)" 
              rows={5}
              placeholder="定义角色的性格、语气、思维框架..."
              value={form.persona} 
              onChange={...} 
            />
          </Section>
          
          {/* 模型配置 */}
          <Section title="模型配置">
            <ModelSelector 
              value={form.model_config}
              onChange={config => setForm(f => ({ ...f, model_config: config }))}
            />
          </Section>
          
          {/* 能力开关 */}
          <Section title="能力配置">
            <Switch 
              label="联网搜索" 
              description="启用 Tavily/Serper 进行事实核查"
              checked={form.capabilities.web_search}
              onChange={...}
            />
            <Switch 
              label="代码执行" 
              description="允许执行代码 (Phase 2)"
              checked={form.capabilities.code_execution}
              disabled
            />
          </Section>
        </div>
        
        <SheetFooter>
          <Button variant="outline" onClick={onClose}>取消</Button>
          <Button onClick={() => saveAgent(form, { onSuccess: onClose })}>
            {isLoading ? <Spinner /> : '保存'}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
};
```

---

## 3. AgentGrid vs AgentTable

```tsx
// Grid View
export const AgentGrid: FC<{ agents: Agent[] }> = ({ agents }) => (
  <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
    {agents.map(agent => (
      <AgentCard key={agent.id} agent={agent} />
    ))}
  </div>
);

// Table View
export const AgentTable: FC<{ agents: Agent[] }> = ({ agents }) => (
  <Table>
    <TableHeader>
      <TableRow>
        <TableHead>名称</TableHead>
        <TableHead>人设</TableHead>
        <TableHead>模型</TableHead>
        <TableHead>能力</TableHead>
        <TableHead>操作</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      {agents.map(agent => (
        <TableRow key={agent.id}>
          <TableCell>{agent.name}</TableCell>
          <TableCell className="max-w-[200px] truncate">{agent.persona}</TableCell>
          <TableCell>{agent.model_config.model}</TableCell>
          <TableCell>...</TableCell>
          <TableCell><ActionMenu /></TableCell>
        </TableRow>
      ))}
    </TableBody>
  </Table>
);
```

---

## 4. 测试要点

- [ ] 创建 Agent 成功
- [ ] 编辑 Agent 成功
- [ ] 删除确认弹窗
- [ ] 表单验证
- [ ] 模型配置保存正确
