# Sprint 4 Specifications: 高级功能

> **Sprint 周期**: Week 4  
> **目标**: 实现人类裁决、成本预估、文档引用等高级功能  
> **里程碑**: M4 - MVP 完成

---

## 📋 Sprint 4 Specs 索引

| Spec ID  | 文档                                                  | 类型      | 优先级 | 预估 |
| -------- | ----------------------------------------------------- | --------- | ------ | ---- |
| SPEC-301 | [HumanReviewModal](./SPEC-301-human-review-modal.md)  | Component | P0     | 3h   |
| SPEC-302 | [CostEstimator](./SPEC-302-cost-estimator.md)         | Component | P1     | 3h   |
| SPEC-303 | [DocumentReference](./SPEC-303-document-reference.md) | Feature   | P1     | 4h   |
| SPEC-304 | [快捷键](./SPEC-304-fullscreen-shortcuts.md)          | Hook      | P2     | 2h   |
| SPEC-305 | [KaTeX 渲染](./SPEC-305-katex-rendering.md)           | Feature   | P2     | 2h   |

---

## 🎯 验收标准

- [ ] 人类裁决节点可暂停并等待用户输入
- [ ] 启动前显示成本预估
- [ ] 文档引用可点击跳转

---

## 🔗 API 依赖 (需新增)

| 端点                             | 方法 | 说明             |
| -------------------------------- | ---- | ---------------- |
| `/api/v1/workflows/:id/estimate` | POST | 预估成本         |
| `/api/v1/sessions/:id/review`    | POST | 提交人类裁决结果 |
