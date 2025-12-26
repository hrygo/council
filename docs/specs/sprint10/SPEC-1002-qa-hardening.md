# SPEC-1002: 质量防线硬化 (QA Hardening)

## 1. 需求背景
目前核心层（Core）的测试覆盖率为 73.3%，其中 `nodes` 包仅为 66.1%。考虑到 `nodes` 承载了所有 AI 编排逻辑，其不稳定性将直接导致工作流执行失败。同时，WebSocket 契约缺乏端到端的运行时校验。

## 2. 目标定义
- `internal/core/workflow/nodes` 包覆盖率提升至 **80%** 以上。
- 实现 E2E 级别的 WebSocket 字段值一致性校验。
- 建立基于契约的防退化测试。

## 3. 实现路线

### 3.1 覆盖率补全
1.  **Agent 节点异常路径**: 测试 Provider 返回错误、Context 取消、Token 超时等场景。
2.  **并行与循环逻辑**: 模拟 `Parallel` 节点中部分支路失败的容错逻辑。
3.  **HumanReview**: 测试人工审核过程中的状态持久化与并发竞争。

### 3.2 E2E 契约验证 (`e2e/tests/websocket-contract.spec.ts`)
- **现状**: 仅校验了字段是否存在。
- **增强**: 
    - 针对 `token_stream` 消息，验证 `chunk` 不为空。
    - 针对 `node_state_change`，验证 `status` 必须为枚举值（running/completed/failed）。
    - 验证 `timestamp` 的严格 ISO 8601 格式。

### 3.3 CI 集成
- 在 `Makefile` 的 `test` 目标中增加覆盖率阈值检查。
- 如果 Core 包覆盖率低于 80%，则构建失败。

## 4. 验收标准
- [ ] `go test ./internal/core/workflow/nodes -cover` 显示结果 $\ge 80\%$。
- [ ] E2E 测试包含对异常消息格式的拦截校验。
- [ ] 所有 `workflow.generated.ts` 中的类型在 E2E 测试中都有对应的运行时采样。
