# 9. 部署与打包架构 (Deployment Architecture)

### 9.1 部署架构概览

```mermaid
graph TD
    subgraph Dev["Development"]
        Vite["Vite Dev :5173"]
        Go["Go Backend :8080"]
        PG_Dev["PostgreSQL :5432"]
    end
    
    subgraph Prod["Production - Docker"]
        Nginx["Nginx :80"]
        API["Go API Container"]
        PG_Prod[("PostgreSQL + pgvector")]
    end
    
    Vite --> Go
    Go --> PG_Dev
    Nginx --> API
    API --> PG_Prod
```

### 9.2 Docker Compose 配置

```yaml
# docker-compose.yml
version: '3.8'

services:
  postgres:
    image: pgvector/pgvector:pg16
    environment:
      POSTGRES_USER: council
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: council
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U council"]
      interval: 5s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: .
      dockerfile: Dockerfile.backend
    environment:
      DATABASE_URL: postgres://council:${DB_PASSWORD}@postgres:5432/council?sslmode=disable
      COUNCIL_MASTER_KEY: ${COUNCIL_MASTER_KEY}
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "80:80"
    depends_on:
      - backend

volumes:
  pgdata:
```

### 9.3 Dockerfile 配置

**Backend (Go)**：

```dockerfile
# Dockerfile.backend
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o council-server ./cmd/council

FROM alpine:3.19
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/council-server /usr/local/bin/
EXPOSE 8080
CMD ["council-server"]
```

**Frontend (React SPA)**：

```dockerfile
# frontend/Dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

### 9.4 GitHub Actions CI/CD

```yaml
# .github/workflows/release.yml
name: Build and Deploy

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Login to Docker Hub
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_TOKEN }}
      
      - name: Build and Push Backend
        uses: docker/build-push-action@v5
        with:
          context: .
          file: Dockerfile.backend
          push: true
          tags: |
            yourorg/council-api:${{ github.ref_name }}
            yourorg/council-api:latest
      
      - name: Build and Push Frontend
        uses: docker/build-push-action@v5
        with:
          context: ./frontend
          push: true
          tags: |
            yourorg/council-web:${{ github.ref_name }}
            yourorg/council-web:latest
```

### 9.5 本地开发启动

```bash
# 1. 启动 PostgreSQL
docker-compose up -d postgres

# 2. 启动后端 (热重载)
cd cmd/council && go run . --dev

# 3. 启动前端 (热重载)
cd frontend && npm run dev
```

### 9.6 生产部署

```bash
# 一键启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f backend

# 数据库迁移
docker-compose exec backend council-server migrate up
```
