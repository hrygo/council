# 自定义工作流指南

> 本文档介绍如何创建自定义 Agent、Group 和 Workflow。

---

## 1. 创建自定义 Agent

### 通过 UI 创建

1. 导航到 **Agents** (`/agents`)
2. 点击 **Create Agent** 按钮
3. 填写配置：

| 字段               | 说明                   | 示例                 |
| :----------------- | :--------------------- | :------------------- |
| **Name**           | Agent 显示名称         | "技术专家"           |
| **Description**    | 简短描述               | "负责技术可行性评估" |
| **Persona Prompt** | 系统提示词（定义角色） | 见下方               |
| **Provider**       | LLM 服务商             | gemini               |
| **Model**          | 具体模型               | gemini-2.0-flash     |
| **Temperature**    | 创造性 (0-1)           | 0.7                  |
| **Max Tokens**     | 最大输出长度           | 8192                 |

### Persona Prompt 编写技巧

一个好的 Persona Prompt 应包含：

```markdown
### Role
你是一位 [具体角色]，专注于 [专业领域]。

### Goal
你的任务是 [核心目标]。

### Constraints
1. [约束条件 1]
2. [约束条件 2]

### Workflow
1. [步骤 1]
2. [步骤 2]

### Output Format
**你必须严格按照以下格式输出：**

## 标题
(内容)

## 小结
(结论)
```

### 通过 Prompt 文件创建（开发者）

在 `internal/resources/prompts/` 目录创建 `.md` 文件：

```markdown
---
name: 技术专家
provider: gemini
model: gemini-2.0-flash
temperature: 0.7
max_tokens: 8192
top_p: 0.9
---

### Role
你是一位资深技术架构师...
```

服务启动时会自动加载并注册到数据库。

---

## 2. 创建自定义 Group

### 通过 UI 创建

1. 导航到 **Groups** (`/groups`)
2. 点击 **Create Group** 按钮
3. 填写配置：

| 字段               | 说明           | 示例              |
| :----------------- | :------------- | :---------------- |
| **Name**           | 群组名称       | "技术评审委员会"  |
| **System Prompt**  | 群组共享上下文 | 见下方            |
| **Default Agents** | 默认成员列表   | 选择 2-5 个 Agent |

### System Prompt 示例

```markdown
# 技术评审委员会

你们是公司技术委员会的成员，负责评审技术方案的可行性。

## 评审维度
1. **技术可行性**: 能否实现？
2. **成本效益**: 投入产出比？
3. **风险评估**: 有哪些技术风险？
4. **时间线**: 能按期交付吗？

## 协作规则
- 各司其职，发挥专长
- 基于事实论证，避免主观臆断
- 最终形成统一建议
```

---

## 3. 设计自定义 Workflow

### 使用 Workflow Canvas

1. 导航到 **Builder** (`/builder`)
2. 从左侧 **Node Palette** 拖拽节点到画布
3. 连接节点定义执行流程
4. 点击节点配置属性

### 节点类型说明

| 类型                | 图标 | 用途          | 配置项                       |
| :------------------ | :--- | :------------ | :--------------------------- |
| **Start**           | ▶️    | 流程入口      | 无                           |
| **End**             | ⏹️    | 流程出口      | summary_prompt               |
| **Agent**           | 🤖    | 调用 AI Agent | agent_id                     |
| **Parallel**        | ⚡    | 并行分支      | 无                           |
| **Vote**            | 🗳️    | 投票决策      | threshold, vote_type         |
| **Loop**            | 🔄    | 循环逻辑      | max_rounds, exit_on_score    |
| **FactCheck**       | ✅    | 事实核查      | verify_threshold             |
| **HumanReview**     | 👤    | 人工审核      | timeout_minutes, allow_skip  |
| **MemoryRetrieval** | 📚    | 历史检索      | max_results, time_range_days |

### 工作流设计模式

#### 模式 1: 简单链式

```
Start → Agent A → Agent B → End
```

适用于: 顺序处理任务

#### 模式 2: 并行汇聚

```
Start → Parallel → [Agent A, Agent B] → Agent C → End
```

适用于: 多角度分析后综合

#### 模式 3: 迭代优化

```
Start → Loop → [工作节点] → HumanReview → [Loop/End]
```

适用于: 需要多轮改进的任务

#### 模式 4: 人工网关

```
Start → Agent → HumanReview → [Approved: End] / [Rejected: Loop]
```

适用于: 需要人工确认的关键节点

---

## 4. 最佳实践

### 如何选择模型

| 任务类型   | 推荐模型         | 原因           |
| :--------- | :--------------- | :------------- |
| 创意发散   | Gemini, GPT-4    | 联想力强       |
| 逻辑推理   | DeepSeek         | 推理能力强     |
| 中文写作   | Qwen, GLM        | 中文语义理解好 |
| 代码生成   | DeepSeek, Claude | 结构化输出强   |
| 长文档处理 | Gemini           | 超长上下文     |

### 如何设置 Temperature

| Temperature | 效果                 | 适用场景           |
| :---------- | :------------------- | :----------------- |
| 0.0 - 0.2   | 极稳定，输出确定性高 | 裁判、事实核查     |
| 0.3 - 0.5   | 较稳定，有限变化     | 分析、总结         |
| 0.6 - 0.8   | 平衡，适度创造力     | 一般对话           |
| 0.9 - 1.0   | 高创造力，较发散     | 头脑风暴、正方辩护 |

### 节点连接规则

- **Start** 只能作为源节点
- **End** 只能作为目标节点
- **Parallel** 后应有 2+ 分支
- **Loop** 需要有回边和出口

---

## 5. 保存与复用

### 保存为模板

1. 完成工作流设计
2. 点击 **Save as Template**
3. 输入模板名称和描述
4. 模板将出现在 **Template Library** 中

### 从模板创建

1. 打开 **Template Library** 侧边栏
2. 选择模板，点击 **Use Template**
3. 根据需要修改节点配置
4. 保存为新工作流

---

## 6. 下一步

- 了解 [LLM Provider 配置](./llm-providers.md)
- 阅读 [模型选型策略](./model-selection-strategy.md)
- 查看 [The Council 使用指南](./council-debate.md)
