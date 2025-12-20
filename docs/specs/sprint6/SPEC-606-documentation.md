# SPEC-606: Documentation Updates

> **优先级**: P1  
> **类型**: Documentation  
> **预估工时**: 3h

## 1. 概述

更新项目文档，将 "The Council" 作为系统的旗舰示例进行介绍和说明。

## 2. 目标

- 更新 README.md，突出 "The Council" 开箱即用体验
- 创建用户指南，解释如何使用 Debate 和 Optimize 流程
- 为开发者提供自定义 Agent/Workflow 的教程

## 3. 文档清单

### 3.1 README.md 更新

**新增章节**:

```markdown
## 🚀 Quick Start: The Council

The Council is the built-in "AI Governance Board" that comes pre-configured out-of-the-box.

### Default Agents
- **Value Defender** 🛡️: Advocates for the strategic value of your proposal
- **Risk Auditor** 🔍: Identifies potential risks and logical gaps
- **Chief Justice** ⚖️: Synthesizes arguments and delivers the final verdict

### Available Workflows
1. **Debate**: A single round of three-way debate
2. **Optimize**: An iterative loop with human-in-the-loop review

### Try It Now
1. Start the server: `go run cmd/council/main.go`
2. Open browser: `http://localhost:8080`
3. Select "The Council" group
4. Create a new meeting with the "Debate" workflow
5. Upload your document and watch the AI council deliberate!
```

### 3.2 docs/guide/council-debate.md (新建)

**内容大纲**:

1. **什么是 The Council**
   - 设计哲学：对抗性协作
   - 三个角色的分工

2. **如何使用 Debate 流程**
   - 步骤截图说明
   - 输入/输出说明

3. **如何使用 Optimize 流程**
   - 循环逻辑解释
   - HumanReview 节点说明

4. **自定义你的 Council**
   - 如何修改 Agent Persona
   - 如何调整模型配置

### 3.3 docs/guide/custom-workflow.md (新建)

**内容大纲**:

1. **创建自定义 Agent**
   - 使用 Agent Factory UI
   - Persona Prompt 编写技巧

2. **创建自定义 Group**
   - 设置 System Prompt
   - 选择默认成员

3. **设计自定义 Workflow**
   - 使用 Workflow Canvas
   - 节点类型说明

4. **最佳实践**
   - 如何选择模型
   - 如何设置 Temperature

### 3.4 docs/guide/llm-providers.md (新建) - Gap 1 Remediation

**内容**: 完整的 LLM Provider 配置指南

```markdown
# LLM Provider 配置指南

本系统支持 6 个 LLM Provider，可在 Agent 配置中灵活选择。

## 可用 Provider 列表

| Provider        | 默认模型             | 特点                  | 推荐场景             |
| --------------- | -------------------- | --------------------- | -------------------- |
| **gemini**      | gemini-3-pro-preview | 超长上下文、多模态    | 文档分析、跨学科推理 |
| **deepseek**    | deepseek-chat        | 逻辑严密、代码能力强  | Bug 修复、数学推导   |
| **siliconflow** | GLM-4.6              | 慢思考、Agent 编排    | 复杂决策、多步推理   |
| **openai**      | gpt-5-mini           | 速度快、成本低        | 日常对话、大批量处理 |
| **dashscope**   | qwen-plus            | 中文语义深、文化理解  | 公文写作、RAG 问答   |
| **openrouter**  | grok-4               | 256k 上下文、风格犀利 | 创意发散、反直觉观点 |

## 环境变量配置

在 `.env` 文件中配置 API Key:

\`\`\`bash
# 必选 (默认 Agent 使用)
GEMINI_API_KEY=your_key
DEEPSEEK_API_KEY=your_key
SILICONFLOW_API_KEY=your_key

# 可选 (自定义 Agent 可用)
OPENAI_API_KEY=your_key
DASHSCOPE_API_KEY=your_key
OPENROUTER_API_KEY=your_key
\`\`\`

## 在 Agent 配置中使用

通过 UI 创建 Agent 时，选择 Provider 和 Model:

| 参数        | 说明         | 示例                 |
| ----------- | ------------ | -------------------- |
| provider    | LLM 服务商   | gemini               |
| model       | 具体模型名   | gemini-3-pro-preview |
| temperature | 创造性 (0-1) | 0.7                  |
| max_tokens  | 最大输出长度 | 8192                 |
| top_p       | 采样范围     | 0.95                 |
```

### 3.5 docs/guide/model-selection-strategy.md (新建) - Gap 2 Remediation

**来源**: 迁移自 `example/prompts/model_selection_strategy.md`

**内容摘要**:

```markdown
# 🧠 SparkForge Council: Model Orchestration Strategy

## 核心设计哲学：思维异构 (Think Heterogeneously)

为避免单一模型家族的思维同质化盲区，采用 **\"模型联邦\" (Model Federation)** 策略。

## 角色阵营与模型配置

### ✊ 正方 (Value Defender)
- **模型**: Google Gemini 3.0 Pro
- **选型逻辑**: 联想力强、创造性高
- **Temperature**: 0.9 (鼓励发散)

### 👊 反方 (Risk Auditor)
- **模型**: DeepSeek-V3
- **选型逻辑**: 逻辑严密、直击痛点
- **Temperature**: 0.6 (抑制幻觉)

### ⚖️ 裁判 (Chief Justice)
- **模型**: Zhipu GLM-4.6
- **选型逻辑**: 慢思考、中正平和
- **Temperature**: 0.2 (保证稳定)

## 调整指南

1. **正方不够兴奋**: 提高 temperature 至 1.0
2. **反方攻击性不足**: 确认使用 DeepSeek 原生接口
3. **裁判逻辑混乱**: 降低 temperature 至 0.1

> **SparkForge Principle**: 正确的模型放在正确的位置，才能涌现群体智慧。
```

## 4. 验收标准

- [ ] README.md 包含 "Quick Start: The Council" 章节
- [ ] `docs/guide/council-debate.md` 文件存在且内容完整
- [ ] `docs/guide/custom-workflow.md` 文件存在且内容完整
- [ ] `docs/guide/llm-providers.md` 文件存在且包含 6 个 Provider 配置 **(新增)**
- [ ] `docs/guide/model-selection-strategy.md` 文件存在 **(新增)**
- [ ] 所有文档无 Broken Links
- [ ] 文档语言：中英双语或仅中文 (根据项目语言策略)

## 5. 依赖

- **SPEC-601/602/603**: 所有 Seeder 完成后，可进行实际操作截图

