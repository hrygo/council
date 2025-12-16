# API 设计文档

> **状态**: 设计完成，待实现  
> **版本**: v1.0

---

## 📋 API 文档索引

| 文档                                       | 描述                     | Sprint   | 实现规格                                                        |
| ------------------------------------------ | ------------------------ | -------- | --------------------------------------------------------------- |
| [templates.md](./templates.md)             | 模版库 CRUD API          | Sprint 3 | [SPEC-406](../specs/backend/SPEC-406-templates-api.md)          |
| [human_review.md](./human_review.md)       | 人类裁决 API + WebSocket | Sprint 4 | [SPEC-405](../specs/backend/SPEC-405-human-review-processor.md) |
| [cost_estimation.md](./cost_estimation.md) | 成本预估 API             | Sprint 4 | [SPEC-407](../specs/backend/SPEC-407-cost-estimation-api.md)    |

---

## 📌 文档层级说明

```
docs/
├── api/                  # API 设计文档 (契约定义)
│   ├── templates.md      # 端点、请求/响应格式
│   ├── human_review.md
│   └── cost_estimation.md
└── specs/backend/        # 实现规格 (如何实现)
    ├── SPEC-405-*.md     # 处理器代码结构
    ├── SPEC-406-*.md
    └── SPEC-407-*.md
```

- **API 文档**: 定义接口契约（前后端协作依据）
- **Specs**: 定义实现细节（开发者编码依据）

---

## 🔗 其他 API 参考

已实现的 API 端点请参考：

- `POST /api/v1/workflows/execute` - 启动工作流
- `POST /api/v1/sessions/:id/control` - 暂停/恢复/停止
- `GET/POST /api/v1/groups` - 群组 CRUD
- `GET/POST /api/v1/agents` - Agent CRUD
- `WebSocket /ws` - 实时事件流

详见 [api_spec_v1.5.md](../api_spec_v1.5.md)
