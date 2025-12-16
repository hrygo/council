# SPEC-406: Templates CRUD API

> **优先级**: P1 | **详细设计**: docs/api/templates.md

---

## 1. 端点概览

| 方法   | 路径                                | 说明   |
| ------ | ----------------------------------- | ------ |
| GET    | `/api/v1/templates`                 | 列表   |
| POST   | `/api/v1/templates`                 | 创建   |
| GET    | `/api/v1/templates/:id`             | 详情   |
| PUT    | `/api/v1/templates/:id`             | 更新   |
| DELETE | `/api/v1/templates/:id`             | 删除   |
| POST   | `/api/v1/templates/:id/instantiate` | 实例化 |

---

## 2. Handler 实现

```go
func (h *TemplateHandler) Register(r *gin.RouterGroup) {
    r.GET("/templates", h.List)
    r.POST("/templates", h.Create)
    r.GET("/templates/:id", h.Get)
    r.PUT("/templates/:id", h.Update)
    r.DELETE("/templates/:id", h.Delete)
    r.POST("/templates/:id/instantiate", h.Instantiate)
}
```

---

## 3. 系统预置模版

```go
var SystemTemplates = []Template{
    {
        ID:          "sys-code-review",
        Name:        "严格代码评审",
        Category:    "code_review",
        IsSystem:    true,
        Graph:       codeReviewGraph,
    },
    {
        ID:          "sys-business-plan",
        Name:        "商业计划压测",
        Category:    "business_plan",
        IsSystem:    true,
        Graph:       businessPlanGraph,
    },
    {
        ID:          "sys-quick-decision",
        Name:        "快速决策",
        Category:    "quick_decision",
        IsSystem:    true,
        Graph:       quickDecisionGraph,
    },
}
```

---

## 4. 测试用例

- CRUD 操作
- 系统模版不可修改/删除
- 实例化创建工作流
