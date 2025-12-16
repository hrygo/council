# SPEC-302: CostEstimator 成本预估

> **优先级**: P1 | **预估工时**: 3h  
> **关联 PRD**: F.4.4 成本预估模块 | **API**: cost_estimation.md

---

## 1. 触发位置

在点击 "开始会议" 前显示预估面板：

```tsx
// MeetingRoom.tsx
const [showEstimate, setShowEstimate] = useState(false);

const handleStartMeeting = async () => {
  setShowEstimate(true);  // 先显示预估
};
```

---

## 2. 面板布局

```
┌─────────────────────────────────────────────┐
│ 💰 成本预估                                 │
├─────────────────────────────────────────────┤
│ 总预估成本: $0.35          耗时: ~2 分钟   │
├─────────────────────────────────────────────┤
│ 分项明细:                                   │
│ ├─ CEO (gpt-4-turbo)      $0.15    5k tokens│
│ ├─ CFO (claude-3.5)       $0.12    4.5k     │
│ └─ 事实核查 (tavily)       $0.02    1k      │
├─────────────────────────────────────────────┤
│ ⚠️ 本次会议预估成本超过 $0.30              │
│    建议: 切换为更经济的模型                 │
├─────────────────────────────────────────────┤
│            [取消] [调整配置] [确认启动]     │
└─────────────────────────────────────────────┘
```

---

## 3. 组件实现

```tsx
interface CostEstimatorProps {
  workflowId: string;
  proposal?: string;
  onConfirm: () => void;
  onCancel: () => void;
}

export const CostEstimator: FC<CostEstimatorProps> = ({ 
  workflowId, proposal, onConfirm, onCancel 
}) => {
  const { data, isLoading, error } = useQuery({
    queryKey: ['estimate', workflowId],
    queryFn: () => fetch(`/api/v1/workflows/${workflowId}/estimate`, {
      method: 'POST',
      body: JSON.stringify({ proposal }),
    }).then(r => r.json()),
  });

  if (isLoading) return <Skeleton />;
  if (error) return <ErrorState />;

  return (
    <Dialog open>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>💰 成本预估</DialogTitle>
        </DialogHeader>
        
        <div className="grid grid-cols-2 gap-4 py-4">
          <div className="text-center">
            <div className="text-2xl font-bold">${data.estimated_cost.amount.toFixed(2)}</div>
            <div className="text-sm text-gray-500">预估成本</div>
          </div>
          <div className="text-center">
            <div className="text-2xl font-bold">{data.estimated_duration.formatted}</div>
            <div className="text-sm text-gray-500">预估耗时</div>
          </div>
        </div>
        
        {/* 分项明细 */}
        <div className="space-y-2">
          {data.breakdown.map((item: CostBreakdownItem) => (
            <div key={item.node_id} className="flex justify-between text-sm">
              <span>{item.agent_name || item.node_name}</span>
              <span className="text-gray-500">{item.model}</span>
              <span>${item.estimated_cost.toFixed(4)}</span>
            </div>
          ))}
        </div>
        
        {/* 警告 */}
        {data.warnings.map((w: CostWarning, i: number) => (
          <Alert key={i} variant="warning">
            {w.message}
            {w.suggestion && <p className="text-xs mt-1">{w.suggestion}</p>}
          </Alert>
        ))}
        
        <DialogFooter>
          <Button variant="outline" onClick={onCancel}>取消</Button>
          <Button variant="outline">调整配置</Button>
          <Button onClick={onConfirm}>确认启动</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
```

---

## 4. 测试要点

- [ ] API 调用正确
- [ ] 分项明细显示
- [ ] 警告正确触发
- [ ] 确认后启动会议
