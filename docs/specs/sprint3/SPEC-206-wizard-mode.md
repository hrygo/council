# SPEC-206: 向导模式 (Wizard Mode)

> **优先级**: P0 | **预估工时**: 4h  
> **关联 PRD**: F.3.3 向导模式 | **关联 TDD**: 02_core/13_nl2workflow.md

---

## 1. 功能概述

用户输入自然语言描述，系统自动生成 DAG 工作流。

> PRD 示例：用户输入"帮我设计一个严格的代码评审流程"，系统自动生成 DAG 图。

---

## 2. 组件结构

```
WizardMode
├── IntentInput           # 自然语言输入框
├── TemplateRecommender   # 相似模版推荐
├── GeneratedPreview      # 生成的 DAG 预览
└── ActionButtons         # 立即运行 / 进入编辑
```

---

## 3. 接口定义

```typescript
interface WizardModeProps {
  onComplete: (workflow: GraphDefinition) => void;
  onOpenBuilder: (workflow: GraphDefinition) => void;
}

interface GenerateWorkflowRequest {
  intent: string;
  context?: {
    group_id?: string;
    preferred_agents?: string[];
  };
}

interface GenerateWorkflowResponse {
  workflow: GraphDefinition;
  similar_templates: Template[];
  confidence: number;
}
```

---

## 4. 三步向导流程

```tsx
const WizardMode: FC<WizardModeProps> = ({ onComplete, onOpenBuilder }) => {
  const [step, setStep] = useState<1 | 2 | 3>(1);
  const [intent, setIntent] = useState('');
  const [generated, setGenerated] = useState<GenerateWorkflowResponse | null>(null);

  // Step 1: 意图输入
  const Step1Intent = () => (
    <div className="space-y-4">
      <h2>描述你的会议需求</h2>
      <Textarea
        placeholder="例如：帮我设计一个严格的代码评审流程，需要从安全、性能、可维护性三个角度分析"
        value={intent}
        onChange={e => setIntent(e.target.value)}
        rows={4}
      />
      <Button onClick={handleGenerate} disabled={!intent.trim()}>
        生成流程
      </Button>
    </div>
  );

  // Step 2: 流程选择
  const Step2FlowSelect = () => (
    <div className="space-y-4">
      <h2>推荐流程</h2>
      {generated?.similar_templates.map(tpl => (
        <TemplateCard key={tpl.id} template={tpl} onSelect={...} />
      ))}
      <Divider>或使用 AI 生成的流程</Divider>
      <WorkflowPreview graph={generated?.workflow} />
    </div>
  );

  // Step 3: 成本预估
  const Step3CostPreview = () => (
    <CostEstimator workflowId={generated?.workflow.id} />
  );

  return (
    <Dialog>
      <StepIndicator current={step} steps={['意图', '流程', '预估']} />
      {step === 1 && <Step1Intent />}
      {step === 2 && <Step2FlowSelect />}
      {step === 3 && <Step3CostPreview />}
      <DialogFooter>
        <Button variant="outline" onClick={() => onOpenBuilder(generated.workflow)}>
          进入构建态微调
        </Button>
        <Button onClick={() => onComplete(generated.workflow)}>
          立即开始会议
        </Button>
      </DialogFooter>
    </Dialog>
  );
};
```

---

## 5. API 端点

```http
POST /api/v1/workflows/generate
```

**Request:**
```json
{
  "intent": "帮我设计一个严格的代码评审流程",
  "context": {
    "group_id": "group-123"
  }
}
```

**Response:**
```json
{
  "workflow": {
    "id": "wf-generated",
    "nodes": { /* ... */ }
  },
  "similar_templates": [
    { "id": "sys-code-review", "name": "代码评审", "similarity": 0.85 }
  ],
  "confidence": 0.78
}
```

---

## 6. 测试要点

- [ ] 自然语言生成 DAG
- [ ] 模版推荐匹配
- [ ] 低置信度警告
- [ ] 进入编辑器二次调整
