# LLM Provider 配置指南

> 本系统支持多个 LLM Provider，可在 Agent 配置中灵活选择。

---

## 1. 可用 Provider 列表

| Provider        | 默认模型         | 特点                 | 推荐场景             |
| :-------------- | :--------------- | :------------------- | :------------------- |
| **gemini**      | gemini-2.0-flash | 超长上下文、多模态   | 文档分析、跨学科推理 |
| **deepseek**    | deepseek-chat    | 逻辑严密、代码能力强 | Bug 修复、数学推导   |
| **siliconflow** | GLM-4.6          | 慢思考、Agent 编排   | 复杂决策、多步推理   |
| **openai**      | gpt-4o           | 速度快、通用性强     | 日常对话、大批量处理 |
| **dashscope**   | qwen-plus        | 中文语义深、文化理解 | 公文写作、RAG 问答   |
| **ollama**      | llama3.2         | 本地部署、隐私保护   | 敏感数据、离线场景   |

---

## 2. 环境变量配置

在项目根目录的 `.env` 文件中配置 API Key：

```bash
# =============================================
# 必选 (默认 Agent 使用的 Provider)
# =============================================

# Google Gemini (Value Defender)
GEMINI_API_KEY=your_gemini_api_key

# DeepSeek (Risk Auditor)
DEEPSEEK_API_KEY=your_deepseek_api_key

# SiliconFlow (Chief Justice)
SILICONFLOW_API_KEY=your_siliconflow_api_key

# =============================================
# 可选 (自定义 Agent 可用)
# =============================================

# OpenAI
OPENAI_API_KEY=your_openai_api_key

# 阿里云 DashScope
DASHSCOPE_API_KEY=your_dashscope_api_key

# 本地 Ollama (无需 API Key)
OLLAMA_BASE_URL=http://localhost:11434
```

---

## 3. 在 Agent 配置中使用

通过 UI 创建或编辑 Agent 时，选择对应的 Provider 和 Model：

### 配置参数说明

| 参数            | 说明          | 取值范围                    | 示例             |
| :-------------- | :------------ | :-------------------------- | :--------------- |
| **provider**    | LLM 服务商    | gemini, deepseek, openai... | gemini           |
| **model**       | 具体模型名    | Provider 支持的模型         | gemini-2.0-flash |
| **temperature** | 创造性/随机性 | 0.0 - 1.0                   | 0.7              |
| **max_tokens**  | 最大输出长度  | 1 - 32768                   | 8192             |
| **top_p**       | 采样概率阈值  | 0.0 - 1.0                   | 0.95             |

### 各 Provider 推荐模型

#### Gemini (Google)

```yaml
provider: gemini
model: gemini-2.0-flash      # 快速响应
model: gemini-2.0-flash-exp  # 实验版
model: gemini-exp-1206       # 最新实验
```

#### DeepSeek

```yaml
provider: deepseek
model: deepseek-chat   # 通用对话
model: deepseek-coder  # 代码专精
```

#### SiliconFlow

```yaml
provider: siliconflow
model: zai-org/GLM-4.6  # 智谱 GLM
model: Qwen/Qwen2.5-72B-Instruct  # 通义千问
```

#### OpenAI

```yaml
provider: openai
model: gpt-4o          # 最新多模态
model: gpt-4o-mini     # 轻量版
model: o1-preview      # 推理增强
```

#### DashScope (阿里云)

```yaml
provider: dashscope
model: qwen-plus       # 通用版
model: qwen-max        # 增强版
model: qwen-turbo      # 快速版
```

#### Ollama (本地)

```yaml
provider: ollama
model: llama3.2        # Meta Llama
model: qwen2.5:7b      # 本地千问
model: deepseek-v2     # 本地 DeepSeek
```

---

## 4. The Council 默认配置

系统内置的 The Council 使用以下模型配置：

| Agent            | Provider    | Model                | Temperature | 设计考量               |
| :--------------- | :---------- | :------------------- | :---------- | :--------------------- |
| 🛡️ Value Defender | gemini      | gemini-3-pro-preview | 0.9         | 高创造力，善于发现价值 |
| 🔍 Risk Auditor   | deepseek    | deepseek-chat        | 0.6         | 逻辑严密，直击痛点     |
| ⚖️ Chief Justice  | siliconflow | GLM-4.6              | 0.2         | 稳定中正，综合能力强   |

### 思维异构策略

> **核心原则**: 使用不同模型家族的 Agent 协作，避免单一模型的思维盲区。

- **Gemini** 家族: 联想力强，适合开拓性思考
- **DeepSeek** 家族: 推理严密，适合逻辑验证
- **GLM** 家族: 中英文兼优，适合综合裁决

---

## 5. 故障排查

### API Key 无效

```
Error: 401 Unauthorized
```

检查：
1. `.env` 文件中的 Key 是否正确
2. Key 是否有余额/配额
3. 服务是否在当前地区可用

### 模型不存在

```
Error: Model 'xxx' not found
```

检查：
1. 模型名称拼写是否正确
2. 该 Key 是否有该模型的访问权限
3. 模型是否已下线

### 连接超时

```
Error: context deadline exceeded
```

检查：
1. 网络是否畅通
2. 是否需要代理配置
3. 对于 Ollama，检查服务是否启动

---

## 6. 下一步

- 阅读 [模型选型策略](./model-selection-strategy.md)
- 查看 [自定义工作流指南](./custom-workflow.md)
- 返回 [The Council 使用指南](./council-debate.md)
