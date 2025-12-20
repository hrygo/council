# SPEC-602: Default Group (Go Seeder)

> **优先级**: P0  
> **类型**: Data Seed  
> **预估工时**: 2h  
> **依赖**: SPEC-601 (Default Agents)

## 1. 概述

通过 Go Seeder 机制将 "The Council" 默认群组注入数据库，与 SPEC-601 使用相同的注入方式。

## 2. 目标

- 服务启动时自动创建 "The Council" 群组
- 在 `SeedAgents()` 之后调用 `SeedGroups()`
- 幂等性：重复启动不会创建重复数据

## 3. Group 定义

| 属性                  | 值                                                          |
| :-------------------- | :---------------------------------------------------------- |
| **ID**                | `system_council`                                            |
| **Name**              | The Council                                                 |
| **System Prompt**     | 见下方                                                      |
| **Default Agent IDs** | `[system_affirmative, system_negative, system_adjudicator]` |

### System Prompt

```markdown
# The Council - 多智能体协作治理体

你是 "理事会" (The Council) 的一部分，这是一个由多个 AI 专家组成的治理机构。

## 核心原则

1. **对抗性协作**: 通过正反辩论锻造卓越决策。
2. **工业级标准**: 所有输出必须具备可执行性。
3. **全局统筹**: 始终以用户的"初始目标"为最高准则。
```

## 4. 技术实现

### 4.1 Seeder 方法

```go
// internal/resources/seeder.go

const councilSystemPrompt = `# The Council - 多智能体协作治理体

你是 "理事会" (The Council) 的一部分，这是一个由多个 AI 专家组成的治理机构。

## 核心原则

1. **对抗性协作**: 通过正反辩论锻造卓越决策。
2. **工业级标准**: 所有输出必须具备可执行性。
3. **全局统筹**: 始终以用户的"初始目标"为最高准则。
`

func (s *Seeder) SeedGroups(ctx context.Context) error {
    defaultAgentIDs := []string{
        "system_affirmative",
        "system_negative", 
        "system_adjudicator",
    }
    
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO groups (id, name, system_prompt, default_agent_ids, created_at, updated_at)
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        ON CONFLICT (id) DO NOTHING
    `, "system_council", "The Council", councilSystemPrompt, pq.Array(defaultAgentIDs))
    
    return err
}
```

### 4.2 集成到启动流程

```go
// cmd/council/main.go
func main() {
    // ... 初始化 DB ...
    
    seeder := resources.NewSeeder(db)
    
    // 顺序执行：先 Agent，再 Group
    if err := seeder.SeedAgents(ctx); err != nil {
        log.Fatal(err)
    }
    if err := seeder.SeedGroups(ctx); err != nil {
        log.Fatal(err)
    }
    if err := seeder.SeedWorkflows(ctx); err != nil {
        log.Fatal(err)
    }
    
    // ... 启动服务 ...
}
```

### 4.3 统一调用入口 (推荐)

```go
// internal/resources/seeder.go

func (s *Seeder) SeedAll(ctx context.Context) error {
    if err := s.SeedAgents(ctx); err != nil {
        return fmt.Errorf("seed agents: %w", err)
    }
    if err := s.SeedGroups(ctx); err != nil {
        return fmt.Errorf("seed groups: %w", err)
    }
    if err := s.SeedWorkflows(ctx); err != nil {
        return fmt.Errorf("seed workflows: %w", err)
    }
    return nil
}
```

```go
// cmd/council/main.go
seeder := resources.NewSeeder(db)
if err := seeder.SeedAll(ctx); err != nil {
    log.Fatalf("Failed to seed data: %v", err)
}
```

## 5. 验收标准

- [ ] `SeedGroups()` 方法实现完成
- [ ] 服务启动后，数据库存在 `system_council` 群组
- [ ] `default_agent_ids` 正确包含 3 个 Agent ID
- [ ] 重复启动不会创建重复数据
- [ ] Agent Seeder 在 Group Seeder 之前执行

## 6. 测试

```bash
# 启动服务
make start-backend

# 验证数据
psql -c "SELECT id, name, default_agent_ids FROM groups WHERE id = 'system_council'"

# 预期:
#       id       |    name     |                    default_agent_ids
# ---------------+-------------+----------------------------------------------------------
#  system_council | The Council | {system_affirmative,system_negative,system_adjudicator}
```

## 7. 与 SPEC-601 关系

**执行顺序**: SPEC-601 (SeedAgents) → SPEC-602 (SeedGroups) → SPEC-603 (SeedWorkflows)

所有数据注入在服务启动时统一执行，避免 SQL Migration 时序问题。
