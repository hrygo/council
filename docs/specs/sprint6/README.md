# Sprint 6: Default Experience (Council Out-of-the-Box)

本 Sprint 专注于将 `example/` 目录中的 "Council Debate" 示例迁移为系统的默认开箱即用体验。

## 目标

1. 用户启动系统后，立即看到 "The Council" 作为预配置的群组。
2. 用户可以直接使用 "Debate" 或 "Optimize" 流程模板。
3. 三个默认 Agent（正方、反方、裁决官）开箱可用。
4. "Optimize" 流程完整覆盖 `skill.md` 的 6 个步骤。

## 实现策略

**核心方法**: Go Seeder + Prompt Embed
- Prompt 存储在 `.md` 文件，使用 `//go:embed` 嵌入
- 服务启动时通过 Seeder 注入数据
- 避免 SQL 内嵌长文本的可维护性问题

**分层设计**:
- **确定性层 (Workflow)**: DAG 定义流程结构
- **智能层 (Prompt)**: Agent Prompt 内含评分矩阵和判断逻辑

## SPEC 列表

| SPEC ID  | 名称                  | 类型           | 优先级 | 预估工时 |
| :------- | :-------------------- | :------------- | :----- | :------- |
| SPEC-609 | Architecture Fixes    | Bug Fix        | **P0** | 6h       |
| SPEC-607 | Memory Retrieval Node | Go Node        | **P0** | 4h       |
| SPEC-608 | Prompt Embed 机制     | Infrastructure | **P0** | 4h       |
| SPEC-601 | Default Agents        | Go Seeder      | **P0** | 4h       |
| SPEC-602 | Default Group         | Go Seeder      | **P0** | 2h       |
| SPEC-603 | Default Workflows     | Go Seeder      | **P0** | 6h       |
| SPEC-605 | Versioning Middleware | Go Middleware  | **P0** | 4h       |
| SPEC-606 | Documentation Updates | Docs           | P1     | 3h       |

**总工时**: 33h

## 依赖关系

```
SPEC-608 (Prompt Embed) ─► SPEC-601 (Agents) ─┐
                                              ├─► SPEC-602 (Group) ─► SPEC-603 (Workflows)
                       SPEC-607 (Memory Node) ─┘
                       SPEC-605 (Versioning) ─► [Parallel]
```

## 验收标准

### 功能验收
- [ ] 服务启动后，数据库中存在 3 个系统 Agent (Prompt 完整)
- [ ] 数据库中存在 "The Council" 群组
- [ ] 模板库中存在 "Debate" 和 "Optimize" 流程
- [ ] Optimize 流程包含 `memory_retrieval` 节点
- [ ] HumanReview 前自动创建备份
- [ ] 用户可通过 UI 直接运行完整的 Optimize 循环

### 解耦验证
- [ ] `make verify-decoupling` 通过
- [ ] `internal/resources/prompts/*.md` 存在，非 SQL 内嵌

### skill.md 覆盖
- [ ] Step 1 (History): Memory Retrieval Node ✅
- [ ] Step 2 (Convene): Parallel + Agent Nodes ✅
- [ ] Step 3 (Verify): Adjudicator Prompt 含评分 ✅
- [ ] Step 4 (Backup): Versioning Middleware ✅
- [ ] Step 5 (Surgeon): HumanReview + UI ✅
- [ ] Step 6 (Loop): Loop Node ✅
