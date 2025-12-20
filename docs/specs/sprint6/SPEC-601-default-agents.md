# SPEC-601: Default Agents (Go Seeder)

> **优先级**: P0  
> **类型**: Data Seed  
> **预估工时**: 4h  
> **依赖**: SPEC-608 (Prompt Embed)

## 1. 概述

通过 Go Seeder 机制将默认 Agent 数据注入数据库，Prompt 内容从嵌入的 `.md` 文件读取。

## 2. 目标

- 服务启动时自动创建 3 个系统级 Agent
- Prompt 使用 Go Embed 机制，避免 SQL 转义问题
- 幂等性：重复启动不会创建重复数据

## 3. Agent 定义

| ID                   | Name           | Provider    | Model                | Prompt Source                   |
| :------------------- | :------------- | :---------- | :------------------- | :------------------------------ |
| `system_affirmative` | Value Defender | gemini      | gemini-3-pro-preview | `prompts/system_affirmative.md` |
| `system_negative`    | Risk Auditor   | deepseek    | deepseek-chat        | `prompts/system_negative.md`    |
| `system_adjudicator` | Chief Justice  | siliconflow | GLM-4.6              | `prompts/system_adjudicator.md` |

## 4. 实现方式

### 4.1 Prompt 文件

见 **SPEC-608**，Prompt 存储在 `internal/resources/prompts/*.md`。

### 4.2 Seeder 调用

```go
// cmd/council/main.go
seeder := resources.NewSeeder(db)
if err := seeder.SeedAgents(ctx); err != nil {
    log.Fatal(err)
}
```

### 4.3 幂等性

使用 `ON CONFLICT (id) DO NOTHING` 保证重复执行不会报错。

## 5. 验收标准

- [ ] `internal/resources/prompts/` 包含 3 个 Prompt 文件
- [ ] 服务首次启动后，数据库存在 3 个 `system_` 前缀的 Agent
- [ ] `persona_prompt` 字段内容完整（非 SQL 转义的乱码）
- [ ] `model_config` JSONB 字段正确
- [ ] 重复启动不会创建重复数据

## 6. 测试

```bash
# 启动服务
make start-backend

# 验证数据
psql -c "SELECT id, name, length(persona_prompt) as prompt_len FROM agents WHERE id LIKE 'system_%'"

# 预期:
#         id         |      name       | prompt_len
# -------------------+-----------------+------------
#  system_affirmative| Value Defender  |       2500+
#  system_negative   | Risk Auditor    |       2500+
#  system_adjudicator| Chief Justice   |       3000+
```

## 7. 与 SQL Migration 的关系

**原方案**: 使用 SQL Migration 内嵌 Prompt  
**新方案**: SQL Migration 仅创建表结构，Seeder 负责数据注入

Migration 文件 (`seed_default_agents.sql`) 已移除或简化。
