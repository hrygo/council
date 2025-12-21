# SPEC-804: Debate Flow Restoration

## 元信息

| 属性     | 值                 |
| :------- | :----------------- |
| Spec ID  | SPEC-804           |
| 类型     | Feature            |
| 优先级   | P0                 |
| 预估工时 | 2h                 |
| 依赖     | SPEC-801, SPEC-803 |

## 背景

`example/skill.md` 定义了完整的辩论优化流程，当前系统未完全还原该流程。

## 目标

验证并完善 Council Debate 辩论流程，确保 Affirmative → Negative → Adjudicator 循环正确执行。

## Example 流程对照

### skill.md 定义的 6 步流程

| Step  | 名称               | 对应实现              | 状态  |
| :---: | :----------------- | :-------------------- | :---: |
|   1   | 压缩历史上下文     | Memory Retrieval Node |   ✅   |
|   2   | 召集理事会辩论     | Parallel + Agent      |   ⚠️   |
|   3   | 验证一致性/收敛性  | Vote Node             |   ✅   |
|   4   | 快照备份           | Versioning Middleware |   ✅   |
|   5   | 执行修改 (Surgeon) | Human Review Node     |   ⚠️   |
|   6   | 状态更新/循环决策  | Loop Node             |   ⚠️   |

### 待验证/修复项

#### 1. Agent Prompts 同步

确保 `seeder.go` 中的 Agent prompts 与 `example/prompts/` 内容一致：

| Agent       | Prompt 文件    | 模型配置            |
| :---------- | :------------- | :------------------ |
| Affirmative | affirmative.md | gemini-3-pro        |
| Negative    | negative.md    | deepseek-chat       |
| Adjudicator | adjudicator.md | siliconflow/GLM-4.6 |

#### 2. 辩论输入结构

```go
// 传入 Agent 的 input 应包含:
input := map[string]interface{}{
    "target_material":         documentContent,  // 待审议文档
    "instructions":            objective,        // 优化目标
    "history_summary":         historyContext,   // 历史摘要 (来自 Memory Node)
    "current_loop":            loopIndex,        // 当前轮次
}
```

#### 3. Loop Node 配置

```json
{
    "type": "loop",
    "properties": {
        "max_rounds": 3,
        "exit_condition": "score >= 90 OR verdict == 'approved'"
    }
}
```

## 验证计划

### 1. Agent Prompts 验证

```bash
# 检查 seeder.go 中的 prompt 是否与 example/prompts/ 一致
diff <(sqlite3 council.db "SELECT system FROM agents WHERE name='Affirmative'") example/prompts/affirmative.md
```

### 2. 辩论流程测试

1. 启动会话并上传测试文档
2. 观察 Affirmative 输出格式是否符合 prompt 定义
3. 观察 Negative 输出格式是否符合 prompt 定义
4. 观察 Adjudicator 评分是否正确解析
5. 验证循环是否按 max_rounds 或 exit_condition 终止

### 3. 人工审核节点测试

1. 在 Adjudicator 完成后触发 HumanReview
2. 验证 Modal 正确弹出
3. 验证 Approve/Reject 操作正确恢复/终止流程

## 修改清单

### 后端

| 操作     | 文件路径                                       |
| :------- | :--------------------------------------------- |
| [VERIFY] | `internal/resources/seeder.go`                 |
| [VERIFY] | `internal/core/workflow/nodes/loop.go`         |
| [VERIFY] | `internal/core/workflow/nodes/human_review.go` |

### 前端

| 操作     | 文件路径                                                          |
| :------- | :---------------------------------------------------------------- |
| [VERIFY] | `frontend/src/hooks/useWebSocketRouter.ts`                        |
| [VERIFY] | `frontend/src/features/execution/components/HumanReviewModal.tsx` |

## 验收标准

- [ ] Affirmative Agent 输出符合 prompt 格式要求
- [ ] Negative Agent 输出符合 prompt 格式要求
- [ ] Adjudicator 输出包含评分和判决结论
- [ ] Loop 节点正确执行多轮迭代
- [ ] 达到 exit_condition 时正确终止
- [ ] HumanReview 节点在预期位置触发
- [ ] 用户可通过 Modal 批准/拒绝
