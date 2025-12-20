# SPEC-605: Versioning Middleware

> **优先级**: P0  
> **类型**: Feature  
> **预估工时**: 4h

## 1. 概述

实现版本控制中间件，在工作流执行修改操作前自动创建目标文件的备份快照。

## 2. 目标

- 在 HumanReview 节点执行前，自动备份目标文件
- 备份文件存储在 `docs/backup/{session_id}/` 目录
- 支持回滚到任意历史版本

## 3. 技术实现

### 3.1 Middleware 定义

```go
// internal/core/middleware/versioning.go
package middleware

import (
    "context"
    "io"
    "os"
    "path/filepath"
    "time"
    
    "github.com/hrygo/council/internal/core/workflow"
)

type VersioningMiddleware struct {
    BackupDir string
}

func NewVersioningMiddleware(backupDir string) *VersioningMiddleware {
    return &VersioningMiddleware{BackupDir: backupDir}
}

func (v *VersioningMiddleware) Name() string {
    return "versioning"
}

func (v *VersioningMiddleware) BeforeNodeExecution(
    ctx context.Context, 
    session *workflow.Session, 
    node *workflow.Node,
) error {
    // Only trigger for HumanReview nodes with target file
    if node.Type != workflow.NodeTypeHumanReview {
        return nil
    }
    
    targetPath, ok := session.Inputs["target_file"].(string)
    if !ok || targetPath == "" {
        return nil
    }
    
    return v.createBackup(session.ID.String(), targetPath)
}

func (v *VersioningMiddleware) AfterNodeExecution(
    ctx context.Context, 
    session *workflow.Session, 
    node *workflow.Node, 
    output map[string]interface{},
) (map[string]interface{}, error) {
    // No post-processing needed
    return output, nil
}

func (v *VersioningMiddleware) createBackup(sessionID, targetPath string) error {
    // Create backup directory
    backupDir := filepath.Join(v.BackupDir, sessionID)
    os.MkdirAll(backupDir, 0755)
    
    // Generate backup filename
    filename := filepath.Base(targetPath)
    timestamp := time.Now().Format("20060102_150405")
    backupPath := filepath.Join(backupDir, filename+"_"+timestamp+".bak")
    
    // Copy file
    src, err := os.Open(targetPath)
    if err != nil {
        return err
    }
    defer src.Close()
    
    dst, err := os.Create(backupPath)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    _, err = io.Copy(dst, src)
    return err
}
```

### 3.2 注册到 Engine

```go
// cmd/council/main.go
func setupEngine(session *workflow.Session) *workflow.Engine {
    engine := workflow.NewEngine(session)
    
    // Register middleware
    versioningMW := middleware.NewVersioningMiddleware("docs/backup")
    engine.Middlewares = append(engine.Middlewares, versioningMW)
    
    return engine
}
```

## 4. 文件结构

```
internal/
  core/
    middleware/
      versioning.go       # 中间件实现
      versioning_test.go  # 测试
```

备份目录结构：
```
docs/
  backup/
    {session_id}/
      my_doc.md_20241220_142637.bak
      my_doc.md_20241220_143012.bak
```

## 5. 验收标准

- [ ] `internal/core/middleware/versioning.go` 文件存在
- [ ] 执行 HumanReview 节点前，目标文件被自动备份
- [ ] 备份文件命名格式正确 `{filename}_{timestamp}.bak`
- [ ] 备份目录按 session_id 隔离
- [ ] Engine 正确注册并调用 Middleware

## 6. 测试

### 6.1 单元测试

```go
func TestVersioningMiddleware_CreateBackup(t *testing.T) {
    // Setup temp file
    tmpFile, _ := os.CreateTemp("", "test_*.md")
    tmpFile.WriteString("Original content")
    tmpFile.Close()
    
    // Create middleware
    mw := NewVersioningMiddleware(t.TempDir())
    
    // Simulate backup
    err := mw.createBackup("session123", tmpFile.Name())
    assert.NoError(t, err)
    
    // Verify backup exists
    files, _ := filepath.Glob(filepath.Join(t.TempDir(), "session123", "*.bak"))
    assert.Len(t, files, 1)
}
```

### 6.2 集成测试

```bash
# 手动验证
# 1. 启动服务
# 2. 创建会议，选择 Optimize 流程
# 3. 上传一个 .md 文件
# 4. 运行到 HumanReview 步骤
# 5. 检查 docs/backup/ 目录是否有备份
ls docs/backup/
```

## 7. 风险

| 风险                 | 缓解措施                   |
| :------------------- | :------------------------- |
| 备份文件过多占用空间 | 添加定期清理逻辑 (Phase 2) |
| 备份失败阻塞流程     | 捕获错误并 log，不中断流程 |

## 8. Rollback 交互逻辑 (Issue 6 Remediation)

### 8.1 用户交互流程

```
┌─────────────────────────────────────────────────────────────────┐
│                    HumanReview 节点界面                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  📊 裁决评分: 78/100                                            │
│                                                                 │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ Adjudicator 的修改建议:                                  │   │
│  │ 1. 第 45 行：增加错误处理逻辑                            │   │
│  │ 2. 第 78 行：补充边界条件说明                            │   │
│  └─────────────────────────────────────────────────────────┘   │
│                                                                 │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │ 可用操作:                                                 │  │
│  │                                                           │  │
│  │  [继续优化]  [应用修改]  [退出]  [↩️ 回滚]               │  │
│  │                                                           │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                 │
│  💾 备份状态: 已创建 (my_doc.md_20241220_142637.bak)           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 8.2 Rollback 按钮行为

| 用户操作      | 系统行为                                                                                   |
| :------------ | :----------------------------------------------------------------------------------------- |
| 点击 [↩️ 回滚] | 弹出确认对话框                                                                             |
| 确认回滚      | 1. 读取最近备份<br>2. 覆盖目标文件<br>3. 回退 history_summary.md<br>4. 重新进入上一轮 Loop |
| 取消回滚      | 关闭对话框，保持当前状态                                                                   |

### 8.3 后端 Rollback API

```go
// internal/api/handlers/rollback.go
type RollbackRequest struct {
    SessionID string `json:"session_id"`
    TargetFile string `json:"target_file"`
}

type RollbackResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    RestoredFrom string `json:"restored_from"`
}

func (h *WorkflowHandler) Rollback(c *gin.Context) {
    var req RollbackRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }
    
    // 1. 查找最近备份
    backupPath := h.versioning.FindLatestBackup(req.SessionID, req.TargetFile)
    if backupPath == "" {
        c.JSON(404, gin.H{"error": "No backup found"})
        return
    }
    
    // 2. 恢复文件
    if err := h.versioning.RestoreFromBackup(backupPath, req.TargetFile); err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(200, RollbackResponse{
        Success: true,
        Message: "File restored successfully",
        RestoredFrom: backupPath,
    })
}
```

### 8.4 Versioning Middleware 扩展

```go
// internal/core/middleware/versioning.go

// FindLatestBackup returns the path to the most recent backup file.
func (v *VersioningMiddleware) FindLatestBackup(sessionID, targetPath string) string {
    filename := filepath.Base(targetPath)
    pattern := filepath.Join(v.BackupDir, sessionID, filename+"_*.bak")
    
    matches, _ := filepath.Glob(pattern)
    if len(matches) == 0 {
        return ""
    }
    
    // Sort by timestamp (descending) and return latest
    sort.Sort(sort.Reverse(sort.StringSlice(matches)))
    return matches[0]
}

// RestoreFromBackup copies backup content to original file.
func (v *VersioningMiddleware) RestoreFromBackup(backupPath, targetPath string) error {
    src, err := os.Open(backupPath)
    if err != nil {
        return fmt.Errorf("failed to open backup: %w", err)
    }
    defer src.Close()
    
    dst, err := os.Create(targetPath)
    if err != nil {
        return fmt.Errorf("failed to create target: %w", err)
    }
    defer dst.Close()
    
    _, err = io.Copy(dst, src)
    return err
}
```

### 8.5 前端 Rollback 按钮

```typescript
// frontend/src/components/workflow/HumanReviewPanel.tsx
const handleRollback = async () => {
  const confirmed = await confirm("确定要回滚到上一个版本吗？当前修改将丢失。");
  if (!confirmed) return;
  
  try {
    const response = await api.post('/api/v1/workflow/rollback', {
      session_id: sessionId,
      target_file: targetFile,
    });
    
    toast.success(`已回滚到: ${response.data.restored_from}`);
    onRollbackComplete();
  } catch (error) {
    toast.error('回滚失败: ' + error.message);
  }
};
```

### 8.6 验收标准补充

- [ ] HumanReview 界面显示 [↩️ 回滚] 按钮
- [ ] 点击回滚弹出确认对话框
- [ ] 确认后成功恢复文件内容
- [ ] 回滚后 UI 正确刷新
- [ ] 备份状态显示当前备份文件名

