# SPEC-105: ModelSelector 组件

> **优先级**: P1 | **预估工时**: 3h  
> **关联 PRD**: F.2.2 模型配置 (Model Agnostic)

---

## 1. 组件接口

```typescript
interface ModelSelectorProps {
  value: ModelConfig;
  onChange: (config: ModelConfig) => void;
  showAdvanced?: boolean;  // 是否显示高级参数 (godMode)
}

interface ModelConfig {
  provider: ModelProvider;
  model: string;
  temperature: number;
  top_p: number;
  max_tokens: number;
}

type ModelProvider = 'openai' | 'anthropic' | 'google' | 'deepseek' | 'dashscope';
```

---

## 2. 组件布局

```
┌─────────────────────────────────────────────┐
│ 模型配置                                    │
├─────────────────────────────────────────────┤
│ 供应商:  [OpenAI ▾]                         │
│ 模型:    [gpt-4-turbo ▾]                    │
├─────────────────────────────────────────────┤
│ ▸ 高级参数 (仅 God Mode)                    │
│   Temperature: [0.7] ──●────────            │
│   Top P:       [1.0] ───────────●           │
│   Max Tokens:  [4096]                       │
└─────────────────────────────────────────────┘
```

---

## 3. 实现

```tsx
const providers: Record<ModelProvider, { name: string; icon: string; models: string[] }> = {
  openai: {
    name: 'OpenAI',
    icon: '🟢',
    models: ['gpt-4-turbo', 'gpt-4o', 'gpt-4o-mini', 'o1-preview', 'o1-mini'],
  },
  anthropic: {
    name: 'Anthropic',
    icon: '🟠',
    models: ['claude-3.5-sonnet', 'claude-3-opus', 'claude-3-haiku'],
  },
  google: {
    name: 'Google',
    icon: '🔵',
    models: ['gemini-1.5-pro', 'gemini-1.5-flash', 'gemini-2.0-flash-exp'],
  },
  deepseek: {
    name: 'DeepSeek',
    icon: '🟣',
    models: ['deepseek-chat', 'deepseek-reasoner'],
  },
  dashscope: {
    name: 'DashScope',
    icon: '🟡',
    models: ['qwen-max', 'qwen-plus', 'qwen-turbo'],
  },
};

export const ModelSelector: FC<ModelSelectorProps> = ({ value, onChange, showAdvanced }) => {
  const { godMode } = useConfigStore();
  const showParams = showAdvanced || godMode;

  return (
    <div className="space-y-4">
      {/* Provider Select */}
      <Select
        label="供应商"
        value={value.provider}
        onChange={provider => onChange({ 
          ...value, 
          provider, 
          model: providers[provider].models[0] 
        })}
      >
        {Object.entries(providers).map(([key, p]) => (
          <SelectItem key={key} value={key}>
            {p.icon} {p.name}
          </SelectItem>
        ))}
      </Select>

      {/* Model Select */}
      <Select
        label="模型"
        value={value.model}
        onChange={model => onChange({ ...value, model })}
      >
        {providers[value.provider].models.map(m => (
          <SelectItem key={m} value={m}>{m}</SelectItem>
        ))}
      </Select>

      {/* Advanced Parameters */}
      {showParams && (
        <Collapsible title="高级参数">
          <Slider
            label="Temperature (创造力)"
            min={0} max={2} step={0.1}
            value={value.temperature}
            onChange={temperature => onChange({ ...value, temperature })}
          />
          <Slider
            label="Top P"
            min={0} max={1} step={0.05}
            value={value.top_p}
            onChange={top_p => onChange({ ...value, top_p })}
          />
          <NumberInput
            label="Max Tokens"
            min={100} max={128000}
            value={value.max_tokens}
            onChange={max_tokens => onChange({ ...value, max_tokens })}
          />
        </Collapsible>
      )}
    </div>
  );
};
```

---

## 4. 默认值

```typescript
const defaultModelConfig: ModelConfig = {
  provider: 'openai',
  model: 'gpt-4o',
  temperature: 0.7,
  top_p: 1.0,
  max_tokens: 4096,
};
```

---

## 5. 测试要点

- [ ] 供应商切换时模型自动更新
- [ ] 高级参数仅在 God Mode 显示
- [ ] Slider 值变化正确
- [ ] 配置正确传递给父组件
