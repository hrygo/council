# 2.4 群组管理模块 (Group Management)

对应 PRD F.1.x，实现场景隔离和上下文管理。

**核心数据结构：**

```go
type Group struct {
    ID              uuid.UUID       `gorm:"type:uuid;primary_key"`
    Name            string          `gorm:"size:128;not null"`
    Icon            string          `gorm:"size:256"`                    // 图标 URL 或 emoji
    SystemPrompt    string          `gorm:"type:text"`                   // 群定位 (宪法)
    DefaultAgentIDs []uuid.UUID     `gorm:"type:uuid[];serializer:json"` // 默认班底 (3个)
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// GroupService 群组业务逻辑
type GroupService interface {
    Create(ctx context.Context, req CreateGroupRequest) (*Group, error)
    Update(ctx context.Context, id uuid.UUID, req UpdateGroupRequest) (*Group, error)
    Delete(ctx context.Context, id uuid.UUID) error
    List(ctx context.Context) ([]Group, error)
    GetWithAgents(ctx context.Context, id uuid.UUID) (*GroupWithAgents, error)
}

// 创建群组时的校验逻辑 (PRD F.1.1: 默认班底必须为 3 位)
func (s *groupService) Create(ctx context.Context, req CreateGroupRequest) (*Group, error) {
    // 校验默认班底数量
    if len(req.DefaultAgentIDs) != 3 {
        return nil, fmt.Errorf("default agents must be exactly 3, got %d", len(req.DefaultAgentIDs))
    }
    // ... 其他创建逻辑
}
```

**REST API 设计：**

| Method   | Endpoint             | 描述                       |
| -------- | -------------------- | -------------------------- |
| `GET`    | `/api/v1/groups`     | 获取群组列表               |
| `POST`   | `/api/v1/groups`     | 创建群组                   |
| `GET`    | `/api/v1/groups/:id` | 获取群组详情（含默认班底） |
| `PUT`    | `/api/v1/groups/:id` | 更新群组配置               |
| `DELETE` | `/api/v1/groups/:id` | 删除群组（级联删除记忆）   |
