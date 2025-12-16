# Sprint 3 Specifications: Builder 增强

> **Sprint 周期**: Week 3  
> **目标**: 完善工作流编辑器，支持更多节点类型和模版库  
> **里程碑**: M3 - Builder 完整

---

## 📋 Sprint 3 Specs 索引

| Spec ID  | 文档                                                               | 类型      | 优先级 | 预估 |
| -------- | ------------------------------------------------------------------ | --------- | ------ | ---- |
| SPEC-201 | [PropertyPanel](./SPEC-201-property-panel.md)                      | Component | P0     | 4h   |
| SPEC-202 | [Vote/Loop 节点](./SPEC-202-vote-loop-nodes.md)                    | Component | P1     | 3h   |
| SPEC-203 | [FactCheck/HumanReview](./SPEC-203-factcheck-humanreview-nodes.md) | Component | P1     | 3h   |
| SPEC-204 | [模版库侧边栏](./SPEC-204-template-sidebar.md)                     | Component | P2     | 3h   |
| SPEC-205 | [保存为模版](./SPEC-205-save-template.md)                          | Feature   | P2     | 2h   |
| SPEC-206 | [向导模式](./SPEC-206-wizard-mode.md)                              | Feature   | P0     | 4h   |

---

## 🎯 验收标准

- [ ] 节点可配置属性
- [ ] 支持所有 PRD 定义的节点类型 UI
- [ ] 模版保存和加载

---

## 🔗 API 依赖 (需新增)

| 端点                                | 方法       | 说明             |
| ----------------------------------- | ---------- | ---------------- |
| `/api/v1/templates`                 | GET        | 获取模版列表     |
| `/api/v1/templates`                 | POST       | 创建模版         |
| `/api/v1/templates/:id`             | GET/DELETE | 获取/删除模版    |
| `/api/v1/templates/:id/instantiate` | POST       | 从模版创建工作流 |
