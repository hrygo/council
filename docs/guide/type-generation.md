# 类型生成指南

## 概述

本项目使用 [tygo](https://github.com/gzuidhof/tygo) 工具自动从 Go 结构体生成 TypeScript 类型定义，确保前后端类型一致性。

## 快速开始

### 生成类型

```bash
make generate-types
```

### 验证类型一致性

```bash
make verify-types
```

## 添加新的 WebSocket 消息类型

### 1. 在 Go 代码中定义结构体

在 `internal/core/workflow/` 目录下定义你的消息结构体：

```go
// internal/core/workflow/context.go

// StreamEvent represents a real-time event sent to the client
type StreamEvent struct {
    Type      string                 `json:"event"`
    Timestamp time.Time              `json:"timestamp"`
    NodeID    string                 `json:"node_id,omitempty"`
    Data      map[string]interface{} `json:"data,omitempty"`
}
```

### 2. 更新 tygo.yaml 配置

如果需要导出新的类型，在 `tygo.yaml` 中添加：

```yaml
packages:
  - path: "github.com/hrygo/council/internal/core/workflow"
    type_mappings:
      time.Time: "string"
      interface{}: "any"
    output_path: "frontend/src/types/workflow.generated.ts"
    include:
      - StreamEvent
      - YourNewType  # 添加新类型
```

### 3. 生成 TypeScript 类型

```bash
make generate-types
```

生成的 TypeScript 类型文件位于：
- `frontend/src/types/workflow.generated.ts`

### 4. 在前端使用生成的类型

```typescript
import { StreamEvent } from '@/types/workflow.generated';

// WebSocket 消息处理
ws.onmessage = (event) => {
  const message: StreamEvent = JSON.parse(event.data);
  console.log(message.event, message.timestamp);
};
```

## 类型映射规则

tygo 会自动将 Go 类型映射为对应的 TypeScript 类型：

| Go 类型 | TypeScript 类型 |
|---------|----------------|
| `string` | `string` |
| `int`, `int64`, `float64` | `number` |
| `bool` | `boolean` |
| `time.Time` | `string` (配置) |
| `interface{}` | `any` (配置) |
| `map[string]T` | `{ [key: string]: T}` |
| `[]T` | `T[]` |
| `*T` (指针) | `T \| undefined` |

## JSON 标签处理

tygo 会读取 Go 结构体的 JSON 标签：

```go
type Example struct {
    ID       string `json:"id"`           // 必填字段
    Name     string `json:"name,omitempty"` // 可选字段
    Internal string `json:"-"`            // 忽略字段
}
```

生成的 TypeScript 类型：

```typescript
export interface Example {
  id: string;
  name?: string;
  // Internal 字段被忽略
}
```

## CI/CD 集成

在 CI 流程中，类型一致性会被自动验证。如果 Go 代码更新但 TypeScript 类型未重新生成，CI 会失败：

```yaml
# .github/workflows/ci.yml
- name: Verify type consistency
  run: make verify-types
```

## 常见问题

### Q: 为什么需要 `make generate-types`？

**A**: tygo 不会自动监听文件变化。每次修改 Go 结构体后，需要手动运行 `make generate-types` 来同步 TypeScript 类型。

### Q: 生成的文件可以手动编辑吗？

**A**: 不可以。所有 `*.generated.ts` 文件都由 tygo 自动生成，手动修改会在下次生成时被覆盖。

### Q: 如何排除某些类型不导出？

**A**: 在 `tygo.yaml` 中使用 `include` 字段明确指定要导出的类型，未列出的类型不会被导出。

### Q: 可以自定义类型映射吗？

**A**: 可以。在 `tygo.yaml` 的 `type_mappings` 中添加自定义映射：

```yaml
type_mappings:
  time.Time: "string"
  uuid.UUID: "string"
  your.CustomType: "YourTSType"
```

## 最佳实践

1. **每次修改 Go 结构体后立即运行 `make generate-types`**
2. **将 `*.generated.ts` 文件提交到版本控制**
3. **在 PR 中检查生成的类型变更**
4. **使用生成的类型替代手写类型定义**
5. **定期运行 `make verify-types` 确保类型同步**

## 相关资源

- [tygo GitHub 仓库](https://github.com/gzuidhof/tygo)
- [项目配置文件](../../tygo.yaml)
- [Makefile 命令](../../Makefile)
