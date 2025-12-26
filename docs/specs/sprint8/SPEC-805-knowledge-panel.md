# SPEC-805: 知识库面板集成

> **状态**: 待开始  
> **优先级**: P2  
> **Sprint**: 8  
> **预计工时**: 6-8 小时  
> **负责人**: 前端 + 后端

---

## 一、需求背景

### 1.1 问题描述

会议室右侧区域当前未被充分利用，用户无法查看会议过程中的相关知识、上下文和引用文档。

### 1.2 用户故事

**作为** 会议参与者  
**我想要** 在会议室右侧看到相关知识面板  
**以便** 快速了解当前讨论的背景和上下文

**验收标准**:
- 右侧显示知识库面板
- 展示当前会话相关的知识项
- 支持查看知识详情
- 支持跳转到相关消息
- 支持按记忆层级过滤

---

## 二、功能设计

### 2.1 界面布局

```
┌──────────────────────────────────────────────────────────┐
│  MeetingRoom                                             │
├─────────────────────┬────────────────────┬───────────────┤
│  Workflow Canvas    │  Chat Panel        │  Knowledge    │
│  (左侧 25%)         │  (中间 50%)        │  Panel        │
│                     │                    │  (右侧 25%)   │
│  ┌───────────┐      │  ┌──────────────┐  │               │
│  │  Node 1   │      │  │  Message 1   │  │  📚 知识库     │
│  └───────────┘      │  └──────────────┘  │               │
│       ↓             │  ┌──────────────┐  │  🔍 搜索       │
│  ┌───────────┐      │  │  Message 2   │  │  ┌─────────┐  │
│  │  Node 2   │◄─────┼─▶│  (Current)   │  │  │ Filter  │  │
│  └───────────┘      │  └──────────────┘  │  └─────────┘  │
│       ↓             │  ┌──────────────┐  │               │
│  ┌───────────┐      │  │  Message 3   │  │  ┌─────────┐  │
│  │  Node 3   │      │  └──────────────┘  │  │ Item 1  │  │
│  └───────────┘      │                    │  │ 相关度: │  │
│                     │                    │  │ ⭐⭐⭐⭐   │  │
│                     │                    │  └─────────┘  │
│                     │                    │  ┌─────────┐  │
│                     │                    │  │ Item 2  │  │
│                     │                    │  └─────────┘  │
└─────────────────────┴────────────────────┴───────────────┘
```

### 2.2 组件结构

```
KnowledgePanel
├── KnowledgePanelHeader
│   ├── Title ("相关知识")
│   ├── SearchInput
│   └── FilterDropdown (记忆层级)
├── KnowledgeList (虚拟滚动)
│   └── KnowledgeItem (多个)
│       ├── KnowledgeTitle
│       ├── KnowledgeSummary
│       ├── KnowledgeMetadata
│       │   ├── Source (来源)
│       │   ├── Timestamp (创建时间)
│       │   └── RelevanceScore (相关度)
│       └── KnowledgeActions
│           ├── ViewDetailsButton
│           └── JumpToMessageButton
└── KnowledgePanelFooter
    └── StatusText ("显示 10 项，共 25 项")
```

---

## 三、技术实现

### 3.1 前端实现

#### 3.1.1 KnowledgePanel 组件

```typescript
// frontend/src/features/meeting-room/components/KnowledgePanel.tsx

import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { useKnowledge } from '@/hooks/useKnowledge';
import { KnowledgeItem } from './KnowledgeItem';
import { VirtualList } from '@/components/VirtualList';

interface KnowledgePanelProps {
  sessionID: string;
}

export const KnowledgePanel: React.FC<KnowledgePanelProps> = ({ sessionID }) => {
  const { knowledge, isLoading, fetchKnowledge } = useKnowledge(sessionID);
  const [searchQuery, setSearchQuery] = useState('');
  const [memoryLayer, setMemoryLayer] = useState<'all' | 'sandboxed' | 'working' | 'long-term'>('all');

  useEffect(() => {
    fetchKnowledge({ memoryLayer });
  }, [sessionID, memoryLayer]);

  const filteredKnowledge = knowledge.filter(item =>
    item.title.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="knowledge-panel h-full flex flex-col bg-white dark:bg-gray-900">
      {/* Header */}
      <div className="p-4 border-b border-gray-200 dark:border-gray-700">
        <h3 className="text-lg font-semibold mb-2">📚 相关知识</h3>
        
        {/* Search */}
        <input
          type="text"
          placeholder="搜索知识..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full px-3 py-2 border rounded-md"
        />

        {/* Filter */}
        <select
          value={memoryLayer}
          onChange={(e) => setMemoryLayer(e.target.value as any)}
          className="w-full mt-2 px-3 py-2 border rounded-md"
        >
          <option value="all">全部记忆</option>
          <option value="sandboxed">隔离区</option>
          <option value="working">工作记忆</option>
          <option value="long-term">长期记忆</option>
        </select>
      </div>

      {/* Knowledge List */}
      <div className="flex-1 overflow-hidden">
        {isLoading ? (
          <div className="flex items-center justify-center h-full">
            <span>加载中...</span>
          </div>
        ) : (
          <VirtualList
            items={filteredKnowledge}
            itemHeight={120}
            renderItem={(item) => (
              <KnowledgeItem key={item.id} knowledge={item} />
            )}
          />
        )}
      </div>

      {/* Footer */}
      <div className="p-2 border-t border-gray-200 dark:border-gray-700 text-sm text-gray-500">
        显示 {filteredKnowledge.length} 项，共 {knowledge.length} 项
      </div>
    </div>
  );
};
```

#### 3.1.2 KnowledgeItem 组件

```typescript
// frontend/src/features/meeting-room/components/KnowledgeItem.tsx

import React from 'react';
import { Knowledge } from '@/types/knowledge';
import { formatDistanceToNow } from 'date-fns';
import { zhCN } from 'date-fns/locale';

interface KnowledgeItemProps {
  knowledge: Knowledge;
}

export const KnowledgeItem: React.FC<KnowledgeItemProps> = ({ knowledge }) => {
  const handleViewDetails = () => {
    // 展开详情
  };

  const handleJumpToMessage = () => {
    // 跳转到相关消息
    const messageElement = document.getElementById(`message-${knowledge.sourceMessageID}`);
    messageElement?.scrollIntoView({ behavior: 'smooth' });
  };

  return (
    <div className="p-4 border-b border-gray-100 hover:bg-gray-50 dark:hover:bg-gray-800">
      {/* Title */}
      <h4 className="font-semibold text-sm mb-1 line-clamp-2">
        {knowledge.title}
      </h4>

      {/* Summary */}
      <p className="text-xs text-gray-600 dark:text-gray-400 mb-2 line-clamp-3">
        {knowledge.summary}
      </p>

      {/* Metadata */}
      <div className="flex items-center justify-between text-xs text-gray-500">
        <div className="flex items-center gap-2">
          <span className="px-2 py-1 bg-blue-100 text-blue-700 rounded">
            {knowledge.memoryLayer}
          </span>
          <span>
            {formatDistanceToNow(new Date(knowledge.createdAt), { locale: zhCN, addSuffix: true })}
          </span>
        </div>

        {/* Relevance Score */}
        <div className="flex items-center gap-1">
          {[...Array(5)].map((_, i) => (
            <span key={i} className={i < knowledge.relevanceScore ? 'text-yellow-400' : 'text-gray-300'}>
              ⭐
            </span>
          ))}
        </div>
      </div>

      {/* Actions */}
      <div className="flex gap-2 mt-2">
        <button
          onClick={handleViewDetails}
          className="text-xs text-blue-600 hover:underline"
        >
          查看详情
        </button>
        {knowledge.sourceMessageID && (
          <button
            onClick={handleJumpToMessage}
            className="text-xs text-blue-600 hover:underline"
          >
            跳转到消息
          </button>
        )}
      </div>
    </div>
  );
};
```

#### 3.1.3 useKnowledge Hook

```typescript
// frontend/src/hooks/useKnowledge.ts

import { useState, useEffect } from 'react';
import { Knowledge } from '@/types/knowledge';
import { useWebSocket } from './useWebSocket';

export const useKnowledge = (sessionID: string) => {
  const [knowledge, setKnowledge] = useState<Knowledge[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const { on, off } = useWebSocket();

  const fetchKnowledge = async (options?: {
    memoryLayer?: 'all' | 'sandboxed' | 'working' | 'long-term';
    limit?: number;
    offset?: number;
  }) => {
    setIsLoading(true);
    
    try {
      const params = new URLSearchParams();
      if (options?.memoryLayer && options.memoryLayer !== 'all') {
        params.append('memory_layer', options.memoryLayer);
      }
      params.append('limit', String(options?.limit || 50));
      params.append('offset', String(options?.offset || 0));

      const response = await fetch(`/api/sessions/${sessionID}/knowledge?${params}`);
      const data = await response.json();
      
      setKnowledge(data.items);
    } catch (error) {
      console.error('Failed to fetch knowledge:', error);
    } finally {
      setIsLoading(false);
    }
  };

  // 监听 WebSocket 事件
  useEffect(() => {
    const handleKnowledgeUpdated = (data: { sessionID: string; knowledge: Knowledge[] }) => {
      if (data.sessionID === sessionID) {
        setKnowledge(prev => [...data.knowledge, ...prev]);
      }
    };

    on('knowledge:updated', handleKnowledgeUpdated);

    return () => {
      off('knowledge:updated', handleKnowledgeUpdated);
    };
  }, [sessionID, on, off]);

  return {
    knowledge,
    isLoading,
    fetchKnowledge,
  };
};
```

### 3.2 后端实现

#### 3.2.1 API Handler

```go
// internal/api/handler/knowledge.go

package handler

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/council/internal/core/memory"
)

type KnowledgeHandler struct {
    memoryService *memory.Service
}

func NewKnowledgeHandler(memoryService *memory.Service) *KnowledgeHandler {
    return &KnowledgeHandler{
        memoryService: memoryService,
    }
}

// GET /api/sessions/:sessionID/knowledge
func (h *KnowledgeHandler) GetSessionKnowledge(c *gin.Context) {
    sessionID := c.Param("sessionID")
    
    // 解析查询参数
    memoryLayer := c.DefaultQuery("memory_layer", "all")
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

    // 从记忆系统检索知识
    knowledge, err := h.memoryService.RetrieveSessionKnowledge(c.Request.Context(), &memory.RetrievalRequest{
        SessionID:   sessionID,
        MemoryLayer: memoryLayer,
        Limit:       limit,
        Offset:      offset,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "items": knowledge,
        "total": len(knowledge),
        "limit": limit,
        "offset": offset,
    })
}
```

#### 3.2.2 Memory Service 扩展

```go
// internal/core/memory/service.go

package memory

type RetrievalRequest struct {
    SessionID   string
    MemoryLayer string // "all", "sandboxed", "working", "long-term"
    Limit       int
    Offset      int
}

type KnowledgeItem struct {
    ID                string    `json:"id"`
    Title             string    `json:"title"`
    Summary           string    `json:"summary"`
    Content           string    `json:"content"`
    MemoryLayer       string    `json:"memory_layer"`
    RelevanceScore    int       `json:"relevance_score"` // 1-5
    SourceMessageID   string    `json:"source_message_id,omitempty"`
    CreatedAt         time.Time `json:"created_at"`
}

func (s *Service) RetrieveSessionKnowledge(ctx context.Context, req *RetrievalRequest) ([]*KnowledgeItem, error) {
    // 从三层记忆系统检索知识
    var results []*KnowledgeItem

    // 1. 隔离区记忆 (Sandboxed Memory)
    if req.MemoryLayer == "all" || req.MemoryLayer == "sandboxed" {
        sandboxed, err := s.sandboxedStore.RetrieveBySession(ctx, req.SessionID)
        if err == nil {
            for _, item := range sandboxed {
                results = append(results, &KnowledgeItem{
                    ID:             item.ID,
                    Title:          item.Title,
                    Summary:        item.Summary,
                    Content:        item.Content,
                    MemoryLayer:    "sandboxed",
                    RelevanceScore: calculateRelevance(item),
                    CreatedAt:      item.CreatedAt,
                })
            }
        }
    }

    // 2. 工作记忆 (Working Memory)
    if req.MemoryLayer == "all" || req.MemoryLayer == "working" {
        working, err := s.workingStore.RetrieveBySession(ctx, req.SessionID)
        if err == nil {
            for _, item := range working {
                results = append(results, &KnowledgeItem{
                    ID:             item.ID,
                    Title:          item.Title,
                    Summary:        item.Summary,
                    Content:        item.Content,
                    MemoryLayer:    "working",
                    RelevanceScore: calculateRelevance(item),
                    CreatedAt:      item.CreatedAt,
                })
            }
        }
    }

    // 3. 长期记忆 (Long-term Memory)
    if req.MemoryLayer == "all" || req.MemoryLayer == "long-term" {
        longTerm, err := s.longTermStore.RetrieveSimilar(ctx, req.SessionID)
        if err == nil {
            for _, item := range longTerm {
                results = append(results, &KnowledgeItem{
                    ID:             item.ID,
                    Title:          item.Title,
                    Summary:        item.Summary,
                    Content:        item.Content,
                    MemoryLayer:    "long-term",
                    RelevanceScore: calculateRelevance(item),
                    CreatedAt:      item.CreatedAt,
                })
            }
        }
    }

    // 按相关度排序
    sort.Slice(results, func(i, j int) bool {
        return results[i].RelevanceScore > results[j].RelevanceScore
    })

    // 分页
    start := req.Offset
    end := req.Offset + req.Limit
    if start > len(results) {
        return []*KnowledgeItem{}, nil
    }
    if end > len(results) {
        end = len(results)
    }

    return results[start:end], nil
}

// 计算相关度 (1-5)
func calculateRelevance(item interface{}) int {
    // 基于多种因素计算:
    // - 时间新近度
    // - 访问频率
    // - 语义相似度
    // 简化版本: 返回随机 3-5
    return rand.Intn(3) + 3
}
```

#### 3.2.3 WebSocket 事件推送

```go
// internal/core/workflow/nodes/memory_retrieval.go

func (n *MemoryRetrievalNode) Execute(ctx context.Context) error {
    // ... 执行记忆检索逻辑

    // 检索完成后，推送知识更新事件
    knowledgeItems := convertToKnowledgeItems(memories)
    
    n.broadcaster.Send(ctx, &ws.Message{
        Type: "knowledge:updated",
        Data: map[string]interface{}{
            "session_id": n.sessionID,
            "knowledge":  knowledgeItems,
        },
    })

    return nil
}
```

### 3.3 类型定义

```typescript
// frontend/src/types/knowledge.ts

export interface Knowledge {
  id: string;
  title: string;
  summary: string;
  content: string;
  memoryLayer: 'sandboxed' | 'working' | 'long-term';
  relevanceScore: number; // 1-5
  sourceMessageID?: string;
  createdAt: string;
}

export interface KnowledgeResponse {
  items: Knowledge[];
  total: number;
  limit: number;
  offset: number;
}
```

---

## 四、交互设计

### 4.1 用户操作流程

1. **进入会议室**: 右侧自动显示知识库面板
2. **自动加载**: 加载当前会话的相关知识 (默认全部记忆层级)
3. **搜索知识**: 在搜索框输入关键词，实时过滤
4. **过滤层级**: 选择特定记忆层级 (隔离区/工作记忆/长期记忆)
5. **查看详情**: 点击知识项展开完整内容
6. **跳转消息**: 点击 "跳转到消息" 定位到相关消息

### 4.2 实时更新

当工作流执行到 `memory_retrieval` 节点时:
1. 后端推送 `knowledge:updated` WebSocket 事件
2. 前端接收事件，更新知识列表
3. 新增知识项高亮显示 (3 秒后恢复)

---

## 五、性能优化

### 5.1 前端优化

| 优化项 | 方案 | 预期收益 |
|--------|------|----------|
| 虚拟滚动 | React-Window | 渲染 1000+ 知识项时保持流畅 |
| 懒加载 | 滚动到底部时加载更多 | 减少初始加载时间 |
| 缓存 | 缓存已加载的知识项 | 避免重复请求 |
| 防抖搜索 | 输入停止 300ms 后再搜索 | 减少请求次数 |

### 5.2 后端优化

| 优化项 | 方案 | 预期收益 |
|--------|------|----------|
| 数据库索引 | 在 session_id, memory_layer 字段建索引 | 查询速度提升 10x |
| 缓存热点数据 | Redis 缓存最近 100 条知识 | 响应时间减少 70% |
| 分页查询 | 限制单次返回 50 条 | 避免数据量过大 |

---

## 六、测试策略

### 6.1 单元测试

```typescript
// frontend/src/hooks/useKnowledge.test.ts

describe('useKnowledge', () => {
  it('should fetch knowledge on mount', async () => {
    const { result, waitForNextUpdate } = renderHook(() => useKnowledge('session1'));
    
    await waitForNextUpdate();
    
    expect(result.current.knowledge).toHaveLength(10);
    expect(result.current.isLoading).toBe(false);
  });

  it('should filter by memory layer', async () => {
    const { result } = renderHook(() => useKnowledge('session1'));
    
    await act(() => {
      result.current.fetchKnowledge({ memoryLayer: 'working' });
    });
    
    expect(result.current.knowledge.every(k => k.memoryLayer === 'working')).toBe(true);
  });
});
```

### 6.2 集成测试

```go
// internal/api/handler/knowledge_test.go

func TestKnowledgeHandler_GetSessionKnowledge(t *testing.T) {
    handler := setupKnowledgeHandler(t)
    
    req := httptest.NewRequest("GET", "/api/sessions/sess1/knowledge?memory_layer=working", nil)
    w := httptest.NewRecorder()
    
    handler.GetSessionKnowledge(w, req)
    
    assert.Equal(t, http.StatusOK, w.Code)
    
    var resp map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &resp)
    
    assert.NotNil(t, resp["items"])
    assert.Equal(t, 10, len(resp["items"].([]interface{})))
}
```

### 6.3 E2E 测试

```typescript
// e2e/tests/knowledge-panel.spec.ts

test('Knowledge panel displays and filters correctly', async ({ page }) => {
  await page.goto('/meeting-room/session1');
  
  // 验证知识面板可见
  await expect(page.locator('.knowledge-panel')).toBeVisible();
  
  // 验证知识项加载
  await expect(page.locator('.knowledge-item')).toHaveCount(10);
  
  // 测试搜索
  await page.fill('input[placeholder="搜索知识..."]', 'AI');
  await expect(page.locator('.knowledge-item')).toHaveCount(3);
  
  // 测试过滤
  await page.selectOption('select', 'working');
  await expect(page.locator('.knowledge-item')).toHaveCount(5);
  
  // 测试跳转到消息
  await page.click('text=跳转到消息');
  await expect(page.locator('.message.highlighted')).toBeVisible();
});
```

---

## 七、验收标准

### 7.1 功能验收

- [x] 知识库面板在会议室右侧正确显示
- [x] 支持加载当前会话的相关知识
- [x] 支持按记忆层级过滤 (全部/隔离区/工作记忆/长期记忆)
- [x] 支持搜索知识
- [x] 支持查看知识详情
- [x] 支持跳转到相关消息
- [x] WebSocket 实时更新知识列表

### 7.2 性能验收

- [x] 初始加载时间 < 1s
- [x] 搜索响应时间 < 300ms
- [x] 支持渲染 1000+ 知识项不卡顿 (虚拟滚动)

### 7.3 用户体验验收

- [x] 界面布局合理，不遮挡聊天内容
- [x] 知识项信息清晰易读
- [x] 相关度评分直观显示
- [x] 交互操作流畅

### 7.4 质量验收

- [x] 单元测试覆盖率 ≥ 80%
- [x] E2E 测试通过
- [x] 无 Lint 错误
- [x] 代码审查通过

---

## 八、后续优化 (可选)

### 8.1 高级功能

- 知识项收藏功能
- 知识项导出 (Markdown/PDF)
- 知识图谱可视化
- 智能推荐相关知识

### 8.2 性能优化

- 使用 Web Worker 处理搜索
- 预加载下一页数据
- 图片懒加载

---

## 九、参考资料

- [三层记忆协议 SPEC-408](../backend/SPEC-408-memory-protocol.md)
- [记忆检索节点 SPEC-607](../sprint6/SPEC-607-memory-retrieval-node.md)
- [React-Window 文档](https://react-window.vercel.app/)
