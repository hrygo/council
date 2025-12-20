# Sprint 6 审计修正方案 (Remediation Plan)

**基于**: [Sprint 6 审计报告](./sprint6_audit_report.md)  
**日期**: 2024-12-20  
**目标**: 解决审计发现的问题，实现 100% 示例覆盖

---

## 核心策略调整

### 方案选择: 分层实现

将 `skill.md` 的复杂逻辑分为两层：

1. **确定性层 (Deterministic)**: Workflow DAG 定义流程结构
2. **智能层 (Agentic)**: Agent Prompt 内嵌智能判断逻辑

**关键洞察**: `skill.md` 中大量逻辑（评分、Delta分析、回滚判断）是 **Adjudicator Prompt 的职责**，而非 Workflow Node 的职责。

---

## 修正方案详情

### 🔴 Issue 1 修正: 重新设计 `council_optimize` Workflow

#### 问题
当前 Workflow 节点过于简化，缺少 `skill.md` 的关键步骤。

#### 方案: 增加 Workflow 节点 + 增强 Prompt

**新 Workflow 结构**:
```
Start
  │
  ▼
[Memory Retrieval] ──► 读取历史上下文 (Tier 2 Memory)
  │
  ▼
[Parallel: Aff + Neg]
  │
  ▼
[Adjudicator] ──► Prompt 内含评分矩阵 + 退出条件判断
  │
  ▼
[Human Review] ──► 显示评分/建议，用户决定 继续/应用/退出
  │
  ├─► (继续) → Loop 回到 Memory Retrieval
  └─► (退出) → End
```

**关键变化**:

| 原方案                 | 新方案                                     |
| ---------------------- | ------------------------------------------ |
| 5 节点                 | 7 节点                                     |
| 无 Memory 节点         | 新增 `memory_retrieval` 节点               |
| 简单 Loop              | Loop 内含完整辩论子流程                    |
| Adjudicator 仅输出观点 | Adjudicator 输出 **结构化评分 + 行动建议** |

**新增节点类型**: `NodeTypeMemoryRetrieval`
- 从 Tier 2 Working Memory 检索相关历史
- 注入到后续 Agent 的 Context 中

#### 新 SPEC 需求
- `SPEC-607`: Memory Retrieval Node 实现

---

### 🔴 Issue 2 修正: Prompt 存储策略

#### 问题
SQL 中嵌入长文本 Prompt 难以维护。

#### 方案: 混合存储

**目录结构**:
```
internal/resources/
  prompts/
    system_affirmative.md   # 完整 Prompt (从 example 复制)
    system_negative.md
    system_adjudicator.md
  migrations/
    embed.go                # 使用 //go:embed 读取 .md 文件
```

**Migration 改为**:
```go
//go:embed prompts/*.md
var promptFiles embed.FS

func SeedAgents(db *sql.DB) error {
    affirmativePrompt, _ := promptFiles.ReadFile("prompts/system_affirmative.md")
    
    _, err := db.Exec(`
        INSERT INTO agents (id, name, persona_prompt, ...) 
        VALUES ($1, $2, $3, ...)
        ON CONFLICT (id) DO NOTHING
    `, "system_affirmative", "Value Defender", string(affirmativePrompt), ...)
    
    return err
}
```

**优势**:
- Prompt 保持 `.md` 格式，易于编辑和 Diff
- 无 SQL 转义问题
- 仍通过 Migration 机制执行

#### 新 SPEC 需求
- `SPEC-608`: Prompt 嵌入机制 (Go Embed)

---

### 🟡 Issue 4 修正: 提升 SPEC-605 优先级

#### 问题
Versioning Middleware 是 P1，但 `skill.md` 明确要求备份。

#### 方案
- 将 `SPEC-605` 从 **P1 → P0**
- 备份在 HumanReview 节点触发前自动执行
- 如果用户在 HumanReview 时选择 "Rollback"，恢复备份

---

### 🟡 Issue 5 修正: Memory 系统验证

#### 问题
假设 3-Tier Memory 已实现，但未验证。

#### 方案
- 添加 **集成测试用例** 验证 Memory 读写
- 如果 Memory 系统不完整，需在 Sprint 6 中补充或将 `council_optimize` 降级

#### 验证清单
- [ ] Session 结束后，对话历史是否自动写入 Tier 1?
- [ ] 新 Session 开始时，`memory_retrieval` 节点能否检索历史?
- [ ] Manual Promotion (用户标记重要) 是否可用?

---

### 🟡 Issue 6 修正: Model Config Fallback

#### 问题
Agent 硬编码特定供应商，无 Fallback。

#### 方案
在 `LLM Router` 中增加 Fallback 逻辑：

```go
func (r *Router) Route(config ModelConfig) (Provider, error) {
    provider, err := r.getProvider(config.Provider)
    if err != nil {
        // Fallback to default provider
        log.Warn("Provider %s unavailable, falling back to default", config.Provider)
        return r.getDefaultProvider()
    }
    return provider, nil
}
```

**也可以**: 在 UI 中提示用户配置 API Key，而非静默 Fallback。

---

## 更新后的 Sprint 6 任务清单

| SPEC ID      | 名称                     | 类型        | 优先级 | 状态   | 备注                 |
| ------------ | ------------------------ | ----------- | ------ | ------ | -------------------- |
| SPEC-601     | Default Agents Migration | SQL + Embed | P0     | 🔄 修改 | 改用 Go Embed        |
| SPEC-602     | Default Group Migration  | SQL         | P0     | ✅ 保留 |                      |
| SPEC-603     | Default Workflows        | SQL         | P0     | 🔄 修改 | 增加节点             |
| SPEC-605     | Versioning Middleware    | Go          | **P0** | 🔄 提升 | 从 P1 → P0           |
| SPEC-606     | Documentation            | Docs        | P1     | ✅ 保留 |                      |
| **SPEC-607** | Memory Retrieval Node    | Go          | **P0** | 🆕 新增 | 支持历史上下文       |
| **SPEC-608** | Prompt Embed 机制        | Go          | **P0** | 🆕 新增 | 解决 SQL Prompt 问题 |

**新增工时**: +8h (SPEC-607: 4h, SPEC-608: 4h)

---

## 新依赖关系

```
SPEC-608 (Prompt Embed) ─► SPEC-601 (Agents) ─┐
                                              ├─► SPEC-602 (Group) ─► SPEC-603 (Workflows)
                       SPEC-607 (Memory Node) ─┘
                       SPEC-605 (Versioning) ─► [Parallel]
```

---

## 验收标准更新

### 功能验收
- [ ] `make migrate` 后，数据库中存在 3 个系统 Agent (Prompt 完整)
- [ ] Workflow Canvas 可正确渲染包含 `memory_retrieval` 节点的 Optimize 流程
- [ ] 运行 Debate 流程，三个 Agent 可正常调用各自的 LLM
- [ ] 运行 Optimize 流程，Memory 节点可检索历史上下文
- [ ] HumanReview 前自动创建备份

### 解耦验证
- [ ] `make verify-decoupling` 通过
- [ ] `internal/resources/prompts/*.md` 存在，非 SQL 内嵌

---

## 结论

通过以下调整，可实现 **100% 示例覆盖**：

1. 新增 `NodeTypeMemoryRetrieval` 节点
2. 改用 Go Embed 存储 Prompt
3. 提升 Versioning Middleware 优先级
4. 增强 Adjudicator Prompt 包含评分矩阵

**总工时调整**: 原 19h → 新 27h (+8h)
