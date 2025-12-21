# Sprint 8: Meeting Room Fix

> **目标**: 修复会议室功能，提供良好的人机交互体验，还原 Example 辩论流程

## 概述

基于 `example/skill.md` 调研，当前会议室实现存在以下差距：
- 缺少文档上传入口
- 启动后直接自动运行，无用户确认
- 工作流画布空白
- 消息展示体验差

## 规格文档索引

| Spec ID  | 文档                                                             | 类型    | 优先级 |
| :------- | :--------------------------------------------------------------- | :------ | :----: |
| SPEC-801 | [Session Startup Flow](./SPEC-801-session-startup-flow.md)       | Feature |   P0   |
| SPEC-802 | [Workflow Live Monitor](./SPEC-802-workflow-live-monitor.md)     | Feature |   P1   |
| SPEC-803 | [Meeting UX Optimization](./SPEC-803-meeting-ux-optimization.md) | UX      |   P1   |
| SPEC-804 | [Debate Flow Restoration](./SPEC-804-debate-flow-restoration.md) | Feature |   P0   |

## 执行阶段

| Phase | 名称                 | 工时  | Specs    |
| :---: | :------------------- | :---: | :------- |
|   1   | Session 启动流程重构 |  4h   | SPEC-801 |
|   2   | 工作流实时监控       |  3h   | SPEC-802 |
|   3   | UX/UI 优化           |  3h   | SPEC-803 |
|   4   | 辩论流程还原         |  2h   | SPEC-804 |
|       | **总计**             |  14h  |          |

## 依赖关系

```
SPEC-801 ─► SPEC-802 ─► SPEC-803
    │                       │
    └───────────────────────┴───► SPEC-804
```
