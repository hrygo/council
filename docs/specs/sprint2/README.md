# Sprint 2 Specifications: 管理页面

> **Sprint 周期**: Week 2  
> **目标**: 完成群组管理和 Agent 管理的 CRUD 页面  
> **里程碑**: M2 - 管理页面完整

---

## 📋 Sprint 2 Specs 索引

| Spec ID  | 文档                                          | 类型      | 优先级 | 预估 |
| -------- | --------------------------------------------- | --------- | ------ | ---- |
| SPEC-101 | [Groups 页面](./SPEC-101-groups-page.md)      | Page      | P0     | 3h   |
| SPEC-102 | [GroupList 组件](./SPEC-102-group-list.md)    | Component | P0     | 4h   |
| SPEC-103 | [Agents 页面](./SPEC-103-agents-page.md)      | Page      | P0     | 3h   |
| SPEC-104 | [AgentList 组件](./SPEC-104-agent-list.md)    | Component | P0     | 4h   |
| SPEC-105 | [ModelSelector](./SPEC-105-model-selector.md) | Component | P1     | 3h   |

---

## 🎯 验收标准

- [ ] 可创建/编辑/删除群组
- [ ] 群组可配置名称、图标、默认成员
- [ ] 可创建/编辑/删除 Agent
- [ ] Agent 可配置模型供应商和参数

---

## 🔗 API 依赖 (已实现)

| 端点                 | 方法           | 说明                 |
| -------------------- | -------------- | -------------------- |
| `/api/v1/groups`     | GET/POST       | 群组列表/创建        |
| `/api/v1/groups/:id` | GET/PUT/DELETE | 群组详情/更新/删除   |
| `/api/v1/agents`     | GET/POST       | Agent 列表/创建      |
| `/api/v1/agents/:id` | GET/PUT/DELETE | Agent 详情/更新/删除 |
