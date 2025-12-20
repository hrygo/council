# 🔴 CRITICAL AUDIT REPORT: Default Experience Implementation Plan

**审计员**: 独立架构专家  
**日期**: 2024-12-20  
**审计对象**: Sprint 6 Specs (SPEC-601 ~ SPEC-606)  
**参照物**: `example/skill.md`, `example/scripts/dialecta_debate.py`

---

## 审计结论: ⚠️ 部分通过，需重大修正

**当前覆盖率: ~70% (非 100%)**

---

## 🔴 严重问题 (Critical Issues)

### Issue 1: `council_optimize` Workflow 严重简化

| `skill.md` 原始步骤                             | SPEC-603 Workflow        | 状态  |
| ----------------------------------------------- | ------------------------ | ----- |
| Step 1: Compress History (读取历史报告)         | ❌ 缺失                   | 🔴 GAP |
| Step 2: Convene The Council (运行辩论)          | ✅ `debate_step` node     | ✅     |
| Step 3: Verify Consistency (多维评分+Delta分析) | ❌ 缺失                   | 🔴 GAP |
| Step 4: Snapshot Backup (备份)                  | ❌ 仅 Middleware (P1)     | 🟡     |
| Step 5: The Surgeon (智能编辑)                  | ❌ 仅 `human_review` 暂停 | 🔴 GAP |
| Step 6: State Update (更新历史)                 | ❌ 缺失                   | 🔴 GAP |

**问题描述**:
- `SPEC-603` 中的 `council_optimize` workflow 仅包含 5 个节点：`Start → Loop → Agent → HumanReview → End`
- 原始 `skill.md` 有 **6 个复杂步骤**，包括文件读写、评分计算、Delta 分析、回滚逻辑
- 当前 workflow 定义无法复刻 `skill.md` 的完整逻辑

**影响**: **示例的 "Optimize" 流程无法真正运行**，用户只会看到一个简化版的循环，缺少关键的智能判断能力。

---

### Issue 2: SQL 中嵌入长文本 Prompt 的可维护性风险

**当前方案** (SPEC-601):
```sql
INSERT INTO agents (persona_prompt) VALUES
(E'### Role\n\n你是一位极具前瞻性的【SparkForge 价值辩护人】...\n(3000+ 字符)');
```

**问题**:
1. **特殊字符转义**: Prompt 中包含 `'`, `"`, `\n`, `|`, `[`, `]` 等，SQL 转义极易出错
2. **版本 Diff 困难**: Prompt 变更时，SQL 文件 diff 几乎不可读
3. **跨团队协作障碍**: 产品/内容团队无法直接编辑 SQL 文件
4. **无校验**: 缺少 JSON Schema 或类型检查

**建议**: 考虑混合方案 - Prompt 存储在独立 `.md` 文件，Migration 通过 Go 代码读取并插入。

---

### Issue 3: Workflow 的 JSONB 嵌套在 SQL 中难以验证

**当前方案** (SPEC-603):
```sql
INSERT INTO workflow_templates (graph_definition) VALUES
('{ "start_node_id": "start", "nodes": { ... 50+ 行 JSON ... } }'::jsonb);
```

**问题**:
1. **无 Schema 验证**: 如果 JSON 结构错误，只有运行时才会发现
2. **IDE 支持差**: 嵌入 SQL 的 JSON 没有语法高亮或自动补全
3. **测试困难**: 无法对 JSON 结构进行单元测试

---

## 🟡 中等问题 (Medium Issues)

### Issue 4: Versioning Middleware 是 P1，但 skill.md 的 Backup 是必需步骤

- `skill.md` Step 4 明确要求在编辑前创建备份
- `SPEC-605` 将其标记为 P1 (可选)
- 如果 P0 完成后不实现 P1，"Optimize" 流程将失去回滚能力

**建议**: 将 `SPEC-605` 提升为 **P0**，或明确记录为 MVP 限制。

---

### Issue 5: 3-Tier Memory 系统是否已实现？

**假设**: 方案假设 `skill.md` 的 "Historian Read/Write" 由现有 3-Tier Memory 覆盖。

**实际情况**: 需要验证:
1. 会话结束后，报告是否自动写入 Tier 1 (Quarantine)?
2. 新会话开始时，是否自动检索相关历史?
3. "Smart Digest" CronJob 是否已实现?

如果 Memory 系统未完成，`skill.md` 的历史上下文管理将无法工作。

---

### Issue 6: Agent Model Config 硬编码了特定供应商

**当前方案**:
```json
{"provider": "gemini", "model": "gemini-3-pro-preview"}
```

**问题**:
- 如果用户没有 Gemini API Key，Affirmative Agent 将无法运行
- 没有 Fallback 机制
- 没有 "使用默认模型" 选项

**建议**: 增加配置校验，若 API Key 缺失，提示用户或使用通用 Fallback。

---

## ✅ 通过的部分 (Passed Items)

| 检查项                          | 状态                                    |
| ------------------------------- | --------------------------------------- |
| 骨架代码无 Council 耦合         | ✅ `make verify-decoupling` 通过         |
| Parallel Node 支持并行执行      | ✅ `engine.go:handleParallel()` 已验证   |
| LLM Router 支持多供应商         | ✅ `router.go` 支持 6 个供应商           |
| HumanReview Node 可暂停流程     | ✅ `human_review.go` 返回 `ErrSuspended` |
| Agent Node 支持 Token Streaming | ✅ `agent.go` 发送 `token_stream` 事件   |

---

## 📋 修正建议清单

| 优先级 | 修正项  | 建议                                                              |
| ------ | ------- | ----------------------------------------------------------------- |
| **P0** | Issue 1 | 重新设计 `council_optimize` workflow，增加更多节点或承认 MVP 限制 |
| **P0** | Issue 2 | 考虑混合方案：Prompt 存 `.md`，Migration 读取插入                 |
| **P1** | Issue 4 | 将 SPEC-605 提升为 P0                                             |
| **P1** | Issue 5 | 添加 Memory 系统验证测试用例                                      |
| **P2** | Issue 6 | 增加 API Key 校验和 Fallback 逻辑                                 |

---

## 最终结论

**如果目标是 100% 复刻 `skill.md`**: ❌ 当前方案不满足，需要补充更多节点或 Agent 逻辑。

**如果目标是提供一个简化版 "Debate" 示例作为 out-of-box**: ✅ `council_debate` workflow 可以工作。

**建议**:
1. 明确 MVP 范围：是否仅覆盖 "Debate" 而将 "Optimize" 留给 Phase 2？
2. 如果坚持 100% 覆盖，需要投入更多工时实现 `skill.md` 的完整逻辑。
