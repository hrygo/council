# SPEC-202: Vote/Loop 节点 UI

> **优先级**: P1 | **预估工时**: 3h  
> **关联 PRD**: F.3.1 Logic Node

---

## 1. Vote 节点

### 1.1 节点图标和样式

```tsx
// 节点外观
const VoteNodeIcon = () => <Vote className="text-amber-500" />;

const voteNodeStyle = {
  background: 'linear-gradient(135deg, #FEF3C7, #FDE68A)',
  border: '2px solid #F59E0B',
  borderRadius: '8px',
};
```

### 1.2 属性配置

```typescript
interface VoteNodeData {
  label: string;
  threshold: number;      // 0.5-1.0, 默认 0.67
  vote_type: 'yes_no' | 'score_1_10';
  agent_ids: string[];    // 参与投票的 Agent
}
```

### 1.3 属性面板

```tsx
const VoteNodeProperties: FC<{ data: VoteNodeData; onChange: (d: Partial<VoteNodeData>) => void }> = 
  ({ data, onChange }) => (
    <>
      <Slider
        label={`通过阈值: ${Math.round(data.threshold * 100)}%`}
        min={50} max={100} step={5}
        value={data.threshold * 100}
        onChange={v => onChange({ threshold: v / 100 })}
      />
      <RadioGroup
        label="投票类型"
        value={data.vote_type}
        onChange={vote_type => onChange({ vote_type })}
        options={[
          { value: 'yes_no', label: '是/否投票' },
          { value: 'score_1_10', label: '1-10 评分' },
        ]}
      />
      <AgentMultiSelect
        label="参与投票的 Agent"
        value={data.agent_ids}
        onChange={agent_ids => onChange({ agent_ids })}
      />
    </>
  );
```

---

## 2. Loop 节点

### 2.1 节点图标和样式

```tsx
const LoopNodeIcon = () => <RefreshCw className="text-purple-500" />;

const loopNodeStyle = {
  background: 'linear-gradient(135deg, #EDE9FE, #DDD6FE)',
  border: '2px solid #8B5CF6',
  borderRadius: '8px',
};
```

### 2.2 属性配置

```typescript
interface LoopNodeData {
  label: string;
  max_rounds: number;          // 1-10, 默认 3
  exit_condition: 'max_rounds' | 'consensus';
  agent_pairs: [string, string][];  // 辩论配对
}
```

### 2.3 属性面板

```tsx
const LoopNodeProperties: FC<{ data: LoopNodeData; onChange: (d: Partial<LoopNodeData>) => void }> = 
  ({ data, onChange }) => (
    <>
      <NumberInput
        label="最大轮数"
        min={1} max={10}
        value={data.max_rounds}
        onChange={max_rounds => onChange({ max_rounds })}
      />
      <Select
        label="退出条件"
        value={data.exit_condition}
        onChange={exit_condition => onChange({ exit_condition })}
      >
        <SelectItem value="max_rounds">达到最大轮数</SelectItem>
        <SelectItem value="consensus">达成共识</SelectItem>
      </Select>
      <AgentPairEditor
        label="辩论配对"
        value={data.agent_pairs}
        onChange={agent_pairs => onChange({ agent_pairs })}
      />
    </>
  );
```

---

## 3. 节点注册

```typescript
// nodes/index.ts
export const customNodes = {
  vote: VoteNode,
  loop: LoopNode,
  // ...
};

export const nodeDefaults: Record<NodeType, () => NodeData> = {
  vote: () => ({ label: '表决', threshold: 0.67, vote_type: 'yes_no', agent_ids: [] }),
  loop: () => ({ label: '循环辩论', max_rounds: 3, exit_condition: 'max_rounds', agent_pairs: [] }),
};
```

---

## 4. 测试要点

- [ ] 节点正确渲染
- [ ] 属性修改同步
- [ ] 默认值正确
- [ ] 拖拽添加节点
