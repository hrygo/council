# Docker Compose配置

<cite>
**本文档引用的文件**  
- [docker-compose.yml](file://docker-compose.yml)
- [.env.example](file://.env.example)
- [README.md](file://README.md)
- [postgres.go](file://internal/infrastructure/db/postgres.go)
- [redis.go](file://internal/infrastructure/cache/redis.go)
- [config.go](file://internal/pkg/config/config.go)
- [migrator.go](file://internal/infrastructure/db/migrator.go)
- [001_init_schema.up.sql](file://internal/infrastructure/db/migrations/001_init_schema.up.sql)
- [service.go](file://internal/core/memory/service.go)
- [llm.go](file://internal/infrastructure/llm/llm.go)
- [Makefile](file://Makefile)
</cite>

## 目录
1. [简介](#简介)
2. [项目结构](#项目结构)
3. [核心组件](#核心组件)
4. [架构概述](#架构概述)
5. [详细组件分析](#详细组件分析)
6. [依赖关系分析](#依赖关系分析)
7. [性能考虑](#性能考虑)
8. [故障排除指南](#故障排除指南)
9. [结论](#结论)

## 简介
本项目是一个可视化多智能体协作系统，名为“The Council”，旨在提供一个私人的AI决策支持平台。系统采用Go语言作为后端开发语言，React作为前端框架，并使用PostgreSQL和Redis作为基础设施服务。通过Docker Compose进行服务编排，实现了数据库、缓存和应用服务的统一管理。系统支持多智能体协作、投票决策、循环优化、事实核查和人工审核等多种高级功能，具备强大的AI驱动能力。

## 项目结构
项目采用分层架构设计，主要包括前端、后端和基础设施三个部分。前端使用React 19构建，后端使用Go 1.21+开发，基础设施包括PostgreSQL 16（含pgvector扩展）和Redis缓存。项目通过Docker Compose进行服务编排，确保开发环境的一致性和可移植性。

```mermaid
graph TB
subgraph "前端"
Frontend[React SPA]
end
subgraph "后端"
Backend[Go/Gin]
end
subgraph "基础设施"
DB[PostgreSQL + pgvector]
Cache[Redis]
ExternalAPIs[LLM, Search]
end
Frontend --> Backend
Backend --> DB
Backend --> Cache
Backend --> ExternalAPIs
```

**图表来源**  
- [README.md](file://README.md#L149-L174)

**本节来源**  
- [README.md](file://README.md#L289-L312)

## 核心组件
系统的核心组件包括工作流引擎、内存服务、中间件和基础设施集成。工作流引擎负责协调多智能体的协作流程，内存服务实现三层记忆机制（隔离区、工作记忆、长期知识），中间件提供安全机制，基础设施集成则负责与外部系统的交互。

**本节来源**  
- [README.md](file://README.md#L294-L307)
- [internal/core/workflow/engine.go](file://internal/core/workflow/engine.go)
- [internal/core/memory/service.go](file://internal/core/memory/service.go)

## 架构概述
系统采用微服务架构，前端通过REST API和WebSocket与后端通信，后端处理业务逻辑并与数据库和缓存交互。PostgreSQL用于持久化存储，Redis用于缓存，外部API（如LLM和搜索服务）通过适配器模式集成。

```mermaid
graph TD
A[前端] --> B[后端]
B --> C[PostgreSQL]
B --> D[Redis]
B --> E[外部API]
C --> F[pgvector]
D --> G[缓存]
E --> H[LLM]
E --> I[搜索]
```

**图表来源**  
- [README.md](file://README.md#L149-L174)
- [internal/infrastructure/db/postgres.go](file://internal/infrastructure/db/postgres.go)
- [internal/infrastructure/cache/redis.go](file://internal/infrastructure/cache/redis.go)

## 详细组件分析

### 数据库服务分析
数据库服务使用pgvector/pgvector:pg16镜像，该镜像在PostgreSQL 16基础上集成了pgvector扩展，支持向量相似性搜索。端口5432映射确保了数据库服务的可访问性。环境变量POSTGRES_USER、POSTGRES_PASSWORD和POSTGRES_DB用于配置数据库的初始用户、密码和数据库名称。postgres_data卷实现了数据的持久化存储，防止容器重启导致数据丢失。

#### 服务配置
```yaml
db:
  image: pgvector/pgvector:pg16
  container_name: council_db
  ports:
    - "5432:5432"
  environment:
    POSTGRES_USER: council
    POSTGRES_PASSWORD: council_password
    POSTGRES_DB: council_db
  volumes:
    - postgres_data:/var/lib/postgresql/data
  restart: always
```

**图表来源**  
- [docker-compose.yml](file://docker-compose.yml#L2-L13)

#### 向量存储实现
系统使用pgvector扩展实现向量存储，通过CREATE EXTENSION IF NOT EXISTS vector;语句启用扩展。memories表包含embedding字段，类型为VECTOR(1536)，用于存储1536维的向量数据。索引使用ivfflat算法，提高向量搜索效率。

```sql
-- Enable pgvector extension
CREATE EXTENSION IF NOT EXISTS vector;

-- Memories Table (Vector Store)
CREATE TABLE memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID REFERENCES groups(id) ON DELETE CASCADE,
    agent_id UUID REFERENCES agents(id),
    session_id UUID REFERENCES sessions(id),
    content TEXT NOT NULL,
    embedding VECTOR(1536),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_memories_embedding ON memories USING ivfflat (embedding vector_cosine_ops) WITH (lists = 100);
```

**图表来源**  
- [001_init_schema.up.sql](file://internal/infrastructure/db/migrations/001_init_schema.up.sql#L1-L96)

**本节来源**  
- [docker-compose.yml](file://docker-compose.yml#L2-L13)
- [001_init_schema.up.sql](file://internal/infrastructure/db/migrations/001_init_schema.up.sql#L1-L96)
- [postgres.go](file://internal/infrastructure/db/postgres.go)

### 缓存服务分析
缓存服务采用redis:alpine轻量镜像，减少了资源占用。6379端口暴露确保了缓存服务的可访问性。服务间通过Docker默认网络通信，实现了无缝集成。

#### 服务配置
```yaml
redis:
  image: redis:alpine
  container_name: council_redis
  ports:
    - "6379:6379"
  restart: always
```

**图表来源**  
- [docker-compose.yml](file://docker-compose.yml#L15-L20)

#### 缓存实现
系统使用Redis作为缓存，通过LPush和LTrim操作实现工作记忆的更新和清理。工作记忆存储最近50条记录，TTL设置为24小时，确保数据的时效性。

```go
func (s *Service) UpdateWorkingMemory(ctx context.Context, groupID string, content string, metadata map[string]interface{}) error {
    // ... 其他代码
    key := fmt.Sprintf("wm:%s", groupID)
    if err := s.cache.LPush(ctx, key, content).Err(); err != nil {
        return fmt.Errorf("failed to push to working memory: %w", err)
    }
    s.cache.Expire(ctx, key, 24*time.Hour)
    s.cache.LTrim(ctx, key, 0, 49)
    return nil
}
```

**图表来源**  
- [service.go](file://internal/core/memory/service.go#L60-L97)

**本节来源**  
- [docker-compose.yml](file://docker-compose.yml#L15-L20)
- [redis.go](file://internal/infrastructure/cache/redis.go)
- [service.go](file://internal/core/memory/service.go#L60-L97)

### 服务依赖关系
服务间通过Docker默认网络通信，无需额外配置。Makefile中的start-db、start-backend和start-frontend命令分别启动数据库、后端和前端服务，确保服务启动顺序正确。

#### 启动命令
```makefile
start-db:
	@echo "$(CYAN)🐳 Starting Docker services...$(RESET)"
	@docker compose up -d
	@echo "$(GREEN)✅ Docker services started$(RESET)"
	@docker compose ps

start-backend:
	@echo "$(CYAN)🔧 Starting Backend on :8080...$(RESET)"
	@lsof -ti:8080 >/dev/null 2>&1 && { echo "$(YELLOW)⚠️ Port 8080 already in use. Stopping...$(RESET)"; make stop-backend; sleep 1; } || true
	@env DATABASE_URL="$(DATABASE_URL)" \
		LLM_PROVIDER="$(LLM_PROVIDER)" \
		LLM_MODEL="$(LLM_MODEL)" \
		GEMINI_API_KEY="$(GEMINI_API_KEY)" \
		go run cmd/council/main.go &
	@sleep 3
	@lsof -ti:8080 >/dev/null 2>&1 && echo "$(GREEN)✅ Backend started$(RESET)" || echo "$(RED)❌ Backend failed to start. Check: make logs-backend$(RESET)"

start-frontend:
	@echo "$(CYAN)🎨 Starting Frontend...$(RESET)"
	@cd frontend && npm run dev &
	@sleep 2
	@echo "$(GREEN)✅ Frontend started$(RESET)"
```

**本节来源**  
- [Makefile](file://Makefile#L71-L129)

### 健康检查配置
系统通过Init函数初始化数据库和Redis连接，确保服务启动时能够成功连接到基础设施。数据库连接通过ping验证，Redis连接通过ping命令验证。

#### 数据库连接初始化
```go
func Init(ctx context.Context, databaseURL string) error {
	var err error
	once.Do(func() {
		pool, err = connect(ctx, databaseURL)
	})
	return err
}

func connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	// ... 其他代码
	if pingErr := p.Ping(ctx); pingErr != nil {
		p.Close()
		return nil, fmt.Errorf("failed to ping database: %w", pingErr)
	}
	// ... 其他代码
}
```

#### Redis连接初始化
```go
func Init(addr string, password string, db int) error {
	var err error
	once.Do(func() {
		client, err = connect(addr, password, db)
	})
	return err
}

func connect(addr string, password string, db int) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if pingErr := c.Ping(context.Background()).Err(); pingErr != nil {
		c.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", pingErr)
	}
	return c, nil
}
```

**本节来源**  
- [postgres.go](file://internal/infrastructure/db/postgres.go#L17-L66)
- [redis.go](file://internal/infrastructure/cache/redis.go#L16-L51)

### 资源限制设置
系统通过环境变量和配置文件管理资源限制，如端口、数据库URL和Redis地址。默认配置确保了服务的可移植性和灵活性。

#### 环境变量配置
```env
# Database
DATABASE_URL=postgres://council:council_password@localhost:5432/council_db?sslmode=disable

# Redis (optional)
REDIS_URL=localhost:6379

# Server
GIN_MODE=debug
PORT=8080
```

**本节来源**  
- [.env.example](file://.env.example#L7-L21)
- [config.go](file://internal/pkg/config/config.go)

## 依赖关系分析
系统依赖关系清晰，前端依赖后端API，后端依赖数据库和缓存服务。外部API（如LLM和搜索服务）通过适配器模式集成，降低了耦合度。

```mermaid
graph TD
Frontend --> Backend
Backend --> DB
Backend --> Cache
Backend --> LLM
Backend --> Search
DB --> pgvector
Cache --> Redis
LLM --> OpenAI
LLM --> Gemini
LLM --> DeepSeek
Search --> Tavily
```

**图表来源**  
- [README.md](file://README.md#L149-L174)
- [internal/infrastructure/llm/router.go](file://internal/infrastructure/llm/router.go)
- [internal/infrastructure/search/tavily.go](file://internal/infrastructure/search/tavily.go)

**本节来源**  
- [README.md](file://README.md#L149-L174)
- [internal/infrastructure/llm/router.go](file://internal/infrastructure/llm/router.go)
- [internal/infrastructure/search/tavily.go](file://internal/infrastructure/search/tavily.go)

## 性能考虑
系统通过三层记忆机制和向量搜索优化性能。隔离区存储原始数据，工作记忆存储近期数据，长期知识存储向量化数据。向量搜索通过pgvector扩展实现，提高了检索效率。

### 三层记忆机制
1. **隔离区（Quarantine）**：存储原始数据，防止低质量数据污染系统。
2. **工作记忆（Working Memory）**：存储近期数据，TTL设置为24小时，确保数据时效性。
3. **长期知识（Long-Term Knowledge）**：存储向量化数据，通过向量搜索实现高效检索。

**本节来源**  
- [README.md](file://README.md#L35-L36)
- [service.go](file://internal/core/memory/service.go)

## 故障排除指南
常见问题包括端口冲突、环境变量未设置和依赖服务未启动。通过Makefile中的status命令可以检查服务状态，logs-db命令可以查看Docker日志。

### 常见问题
1. **端口冲突**：确保5432和6379端口未被占用。
2. **环境变量未设置**：检查.env文件是否正确配置。
3. **依赖服务未启动**：使用make start-db启动数据库服务。

**本节来源**  
- [Makefile](file://Makefile#L52-L65)
- [README.md](file://README.md#L76-L84)

## 结论
本项目通过Docker Compose实现了基础设施的统一管理，使用pgvector扩展支持向量搜索，通过三层记忆机制优化性能。系统架构清晰，依赖关系明确，具备良好的可扩展性和可维护性。通过合理的资源配置和健康检查，确保了系统的稳定运行。