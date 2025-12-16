# Templates API 设计文档

> **状态**: 待实现  
> **优先级**: Sprint 3  
> **前端依赖**: 模版库侧边栏、保存为模版功能

---

## 概述

模版库 API 用于管理可复用的工作流模版。用户可以将当前工作流保存为模版，也可以从模版创建新工作流。

---

## 数据模型

### Template

```typescript
interface Template {
  id: string;                    // UUID
  name: string;                  // 模版名称
  description: string;           // 模版描述
  category: TemplateCategory;    // 分类
  graph: GraphDefinition;        // 工作流 DAG 定义
  is_system: boolean;            // 是否为系统预置模版
  created_at: string;            // ISO 8601
  updated_at: string;            // ISO 8601
}

type TemplateCategory = 
  | 'code_review'      // 代码评审
  | 'business_plan'    // 商业计划
  | 'quick_decision'   // 快速决策
  | 'custom';          // 用户自定义
```

---

## API 端点

### 1. 获取模版列表

```http
GET /api/v1/templates
```

**Query Parameters:**

| 参数           | 类型    | 必填 | 说明                        |
| -------------- | ------- | ---- | --------------------------- |
| category       | string  | 否   | 按分类筛选                  |
| include_system | boolean | 否   | 是否包含系统模版，默认 true |

**Response 200:**

```json
{
  "templates": [
    {
      "id": "tpl-001",
      "name": "严格代码评审",
      "description": "Start → Parallel(安全/性能/可维护性) → Vote → HumanReview → End",
      "category": "code_review",
      "is_system": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 1
}
```

---

### 2. 获取模版详情

```http
GET /api/v1/templates/:id
```

**Response 200:**

```json
{
  "id": "tpl-001",
  "name": "严格代码评审",
  "description": "...",
  "category": "code_review",
  "is_system": true,
  "graph": {
    "id": "graph-001",
    "name": "Code Review",
    "start_node_id": "node-start",
    "nodes": {
      "node-start": {
        "id": "node-start",
        "type": "start",
        "name": "Start",
        "next_ids": ["node-parallel"]
      }
      // ... more nodes
    }
  },
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

---

### 3. 创建模版

```http
POST /api/v1/templates
```

**Request Body:**

```json
{
  "name": "我的自定义模版",
  "description": "用于快速头脑风暴",
  "category": "custom",
  "graph": {
    "start_node_id": "node-start",
    "nodes": { /* ... */ }
  }
}
```

**Response 201:**

```json
{
  "id": "tpl-uuid-new",
  "name": "我的自定义模版",
  "description": "用于快速头脑风暴",
  "category": "custom",
  "is_system": false,
  "graph": { /* ... */ },
  "created_at": "2024-12-16T10:00:00Z",
  "updated_at": "2024-12-16T10:00:00Z"
}
```

---

### 4. 更新模版

```http
PUT /api/v1/templates/:id
```

> 注意：系统预置模版 (`is_system: true`) 不可修改，返回 403。

**Request Body:**

```json
{
  "name": "更新后的名称",
  "description": "更新后的描述",
  "graph": { /* ... */ }
}
```

**Response 200:** 返回更新后的模版对象

**Response 403:**

```json
{
  "error": "Cannot modify system template"
}
```

---

### 5. 删除模版

```http
DELETE /api/v1/templates/:id
```

> 注意：系统预置模版不可删除。

**Response 204:** 无内容

**Response 403:**

```json
{
  "error": "Cannot delete system template"
}
```

---

### 6. 从模版创建工作流

```http
POST /api/v1/templates/:id/instantiate
```

**Request Body:**

```json
{
  "workflow_name": "基于模版的新工作流",
  "group_id": "group-uuid"  // 可选，关联到群组
}
```

**Response 201:**

```json
{
  "workflow": {
    "id": "wf-new-uuid",
    "name": "基于模版的新工作流",
    "template_id": "tpl-001",
    "graph": { /* 复制自模版 */ },
    "created_at": "2024-12-16T10:00:00Z"
  }
}
```

---

## 系统预置模版

PRD 要求的三个系统预置模版：

| ID                   | 名称         | 描述                                                                | 强制节点               |
| -------------------- | ------------ | ------------------------------------------------------------------- | ---------------------- |
| `sys-code-review`    | 代码评审     | Start → Parallel(安全/性能/可维护性) → Vote → HumanReview → End     | FactCheck, HumanReview |
| `sys-business-plan`  | 商业计划压测 | Start → Loop(魔鬼代言人 3轮) → FactCheck → Vote → HumanReview → End | FactCheck, HumanReview |
| `sys-quick-decision` | 快速决策     | Start → Parallel(正/反/中立) → FactCheck → Vote → HumanReview → End | FactCheck, HumanReview |

---

## 前端 Mock 示例

```typescript
// frontend/src/mocks/templates.ts
export const mockTemplates: Template[] = [
  {
    id: 'sys-code-review',
    name: '严格代码评审',
    description: 'Start → Parallel(安全/性能/可维护性) → Vote → HumanReview → End',
    category: 'code_review',
    is_system: true,
    graph: { /* ... */ },
    created_at: '2024-01-01T00:00:00Z',
    updated_at: '2024-01-01T00:00:00Z',
  },
  // ...
];
```
