# SPEC-504: 安全强化 (Security Hardening)

> **优先级**: P1  
> **类型**: Security  
> **预估工时**: 8h

## 1. 概述

强化系统安全性，实现细粒度权限控制和 API 速率限制，为生产环境部署做准备。

## 2. 目标

- 实现基于角色的访问控制 (RBAC)
- API 速率限制防止滥用
- 敏感数据脱敏处理
- 安全审计日志

## 3. 技术方案

### 3.1 RBAC 权限模型

**角色定义:**

| 角色     | 权限                           |
| :------- | :----------------------------- |
| `viewer` | 只读访问工作流、查看执行结果   |
| `editor` | 创建/编辑工作流、执行工作流    |
| `admin`  | 管理用户、管理 Agent、系统配置 |
| `owner`  | 完全控制，包括删除和转让       |

**权限资源:**

```go
// internal/core/auth/permissions.go
type Permission string

const (
    WorkflowRead    Permission = "workflow:read"
    WorkflowWrite   Permission = "workflow:write"
    WorkflowExecute Permission = "workflow:execute"
    WorkflowDelete  Permission = "workflow:delete"
    
    AgentRead       Permission = "agent:read"
    AgentWrite      Permission = "agent:write"
    AgentDelete     Permission = "agent:delete"
    
    GroupManage     Permission = "group:manage"
    UserManage      Permission = "user:manage"
)

var RolePermissions = map[string][]Permission{
    "viewer": {WorkflowRead, AgentRead},
    "editor": {WorkflowRead, WorkflowWrite, WorkflowExecute, AgentRead, AgentWrite},
    "admin":  {WorkflowRead, WorkflowWrite, WorkflowExecute, WorkflowDelete, 
               AgentRead, AgentWrite, AgentDelete, GroupManage},
    "owner":  {"*"}, // All permissions
}
```

**中间件实现:**

```go
// internal/api/middleware/auth.go
func RequirePermission(permission Permission) gin.HandlerFunc {
    return func(c *gin.Context) {
        user := GetCurrentUser(c)
        if user == nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return
        }
        
        if !HasPermission(user.Role, permission) {
            c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden"})
            return
        }
        
        c.Next()
    }
}

// 使用示例
router.DELETE("/workflows/:id", 
    RequirePermission(WorkflowDelete), 
    handler.DeleteWorkflow,
)
```

### 3.2 API 速率限制

**配置:**

```go
// internal/api/middleware/ratelimit.go
type RateLimitConfig struct {
    RequestsPerMinute int
    BurstSize         int
    KeyFunc           func(*gin.Context) string // IP or UserID
}

var DefaultLimits = map[string]RateLimitConfig{
    "api:general":    {RequestsPerMinute: 60, BurstSize: 10},
    "api:execute":    {RequestsPerMinute: 10, BurstSize: 2},
    "api:generate":   {RequestsPerMinute: 5, BurstSize: 1},
    "ws:connect":     {RequestsPerMinute: 30, BurstSize: 5},
}
```

**实现 (使用 Redis + 滑动窗口):**

```go
func RateLimiter(config RateLimitConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := config.KeyFunc(c)
        
        // 使用 Redis INCR + EXPIRE 实现滑动窗口
        count, err := redis.Incr(ctx, "ratelimit:"+key).Result()
        if err != nil {
            c.Next()
            return
        }
        
        if count == 1 {
            redis.Expire(ctx, "ratelimit:"+key, time.Minute)
        }
        
        if count > int64(config.RequestsPerMinute) {
            c.Header("Retry-After", "60")
            c.AbortWithStatusJSON(429, gin.H{
                "error": "Rate limit exceeded",
                "retry_after": 60,
            })
            return
        }
        
        c.Header("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
        c.Header("X-RateLimit-Remaining", strconv.Itoa(config.RequestsPerMinute - int(count)))
        c.Next()
    }
}
```

### 3.3 敏感数据处理

**API Key 脱敏:**

```go
// 存储时加密
func EncryptAPIKey(plainKey string) (string, error) {
    // 使用 AES-256-GCM 加密
    block, _ := aes.NewCipher(secretKey)
    gcm, _ := cipher.NewGCM(block)
    nonce := make([]byte, gcm.NonceSize())
    io.ReadFull(rand.Reader, nonce)
    return base64.StdEncoding.EncodeToString(
        gcm.Seal(nonce, nonce, []byte(plainKey), nil),
    ), nil
}

// 返回时脱敏
func MaskAPIKey(key string) string {
    if len(key) < 8 {
        return "****"
    }
    return key[:4] + "****" + key[len(key)-4:]
}
```

**日志脱敏:**

```go
// 自动过滤敏感字段
type SanitizedLogger struct {
    sensitiveFields []string
}

func (l *SanitizedLogger) Log(data map[string]interface{}) {
    for _, field := range l.sensitiveFields {
        if _, exists := data[field]; exists {
            data[field] = "[REDACTED]"
        }
    }
    log.Printf("%+v", data)
}

var logger = SanitizedLogger{
    sensitiveFields: []string{"password", "api_key", "token", "secret"},
}
```

### 3.4 安全审计日志

```go
// internal/infrastructure/audit/logger.go
type AuditEvent struct {
    Timestamp  time.Time         `json:"timestamp"`
    UserID     string            `json:"user_id"`
    Action     string            `json:"action"`
    Resource   string            `json:"resource"`
    ResourceID string            `json:"resource_id"`
    Result     string            `json:"result"` // success, failure, denied
    IP         string            `json:"ip"`
    UserAgent  string            `json:"user_agent"`
    Metadata   map[string]string `json:"metadata,omitempty"`
}

func LogAudit(c *gin.Context, action string, resource string, resourceID string, result string) {
    event := AuditEvent{
        Timestamp:  time.Now(),
        UserID:     GetUserID(c),
        Action:     action,
        Resource:   resource,
        ResourceID: resourceID,
        Result:     result,
        IP:         c.ClientIP(),
        UserAgent:  c.GetHeader("User-Agent"),
    }
    
    // 写入审计日志表或发送到日志服务
    auditLogger.Log(event)
}
```

**审计事件示例:**

```json
{
  "timestamp": "2025-12-20T10:30:00Z",
  "user_id": "user-123",
  "action": "workflow:delete",
  "resource": "workflow",
  "resource_id": "wf-456",
  "result": "success",
  "ip": "192.168.1.100",
  "user_agent": "Mozilla/5.0..."
}
```

## 4. 数据库迁移

```sql
-- migrations/003_add_audit_logs.up.sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id UUID REFERENCES users(id),
    action VARCHAR(64) NOT NULL,
    resource VARCHAR(64) NOT NULL,
    resource_id VARCHAR(128),
    result VARCHAR(16) NOT NULL,
    ip INET,
    user_agent TEXT,
    metadata JSONB
);

CREATE INDEX idx_audit_logs_user ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
```

## 5. 验收标准

- [ ] RBAC 权限检查在所有受保护端点生效
- [ ] 超过速率限制返回 429 状态码
- [ ] API Key 等敏感信息不在日志中明文出现
- [ ] 审计日志记录关键操作
- [ ] 安全测试通过 (OWASP Top 10 检查)

## 6. 安全检查清单

| 项目                                   | 状态 |
| :------------------------------------- | :--- |
| SQL 注入防护 (参数化查询)              | [ ]  |
| XSS 防护 (模板转义)                    | [ ]  |
| CSRF 保护 (Token 验证)                 | [ ]  |
| 认证 Token 安全 (HttpOnly, Secure)     | [ ]  |
| 敏感端点 HTTPS Only                    | [ ]  |
| 依赖漏洞扫描 (npm audit / govulncheck) | [ ]  |

## 7. Makefile 集成

```makefile
# Security
security-scan: ## 🔒 Run security scans
	@echo "Scanning Go dependencies..."
	@govulncheck ./...
	@echo "Scanning Node dependencies..."
	@cd frontend && npm audit

security-audit: ## 📋 View audit logs
	@psql $(DATABASE_URL) -c "SELECT * FROM audit_logs ORDER BY timestamp DESC LIMIT 50;"
```
