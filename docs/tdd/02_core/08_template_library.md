# 2.8 工作流模板模块 (Template Library)

对应 PRD F.3.2，实现模板保存与调用。

**核心数据结构：**

```go
type WorkflowTemplate struct {
    ID              uuid.UUID `gorm:"type:uuid;primary_key"`
    Name            string    `gorm:"size:128;not null"`
    Description     string    `gorm:"size:512"`
    GraphDefinition JSONB     `gorm:"not null"` // React Flow JSON 格式
    IsSystem        bool      `gorm:"default:false"` // 系统内置模板不可删除
    Category        string    `gorm:"size:64"` // "code_review", "brainstorm", "decision", "debate"
    UsageCount      int       `gorm:"default:0"` // 使用次数统计
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type TemplateService interface {
    List(ctx context.Context, category string) ([]WorkflowTemplate, error)
    Get(ctx context.Context, id uuid.UUID) (*WorkflowTemplate, error)
    Create(ctx context.Context, req CreateTemplateRequest) (*WorkflowTemplate, error)
    Delete(ctx context.Context, id uuid.UUID) error
    Apply(ctx context.Context, templateID, groupID uuid.UUID) (*Workflow, error)
}
```

**REST API 设计：**

| Method   | Endpoint                      | 描述                                |
| -------- | ----------------------------- | ----------------------------------- |
| `GET`    | `/api/v1/templates`           | 获取模板列表 (支持 ?category= 过滤) |
| `GET`    | `/api/v1/templates/:id`       | 获取模板详情                        |
| `POST`   | `/api/v1/templates`           | 保存当前工作流为模板                |
| `DELETE` | `/api/v1/templates/:id`       | 删除用户模板 (系统模板禁删)         |
| `POST`   | `/api/v1/templates/:id/apply` | 应用模板到指定群组，创建新 Workflow |

**系统预置模板：**

| 模板名称     | 类别        | 描述                                                                |
| ------------ | ----------- | ------------------------------------------------------------------- |
| 代码评审     | code_review | Start → Parallel(安全/性能/可维护性) → Vote → HumanReview → End     |
| 商业计划压测 | brainstorm  | Start → Loop(魔鬼代言人 3轮) → FactCheck → Vote → HumanReview → End |
| 快速决策     | decision    | Start → Parallel(正/反/中立) → FactCheck → Vote → HumanReview → End |
| 深度辩论     | debate      | Start → Loop(正反方 5轮) → FactCheck → End                          |
