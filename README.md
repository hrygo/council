# The Council (理事会)

**Version**: 0.5.0

The Council is a visualized Multi-Agent Collaboration System & Personal Private Think Tank.

## Architecture
- **Backend**: Go (Gin)
- **Frontend**: React (Vite + atomic CSS / Tailwind)
- **Database**: PostgreSQL (pgvector)
- **Infrastructure**: Docker

## Getting Started

### Prerequisites
- Docker Engine >= 20.10
- Docker Compose v2.x (Plugin, command: `docker compose`)
- Go 1.21+
- Node.js 20+

### Installation

1. Start database:
   ```bash
   docker compose up -d
   ```

2. Start Backend:
   ```bash
   go run cmd/council/main.go
   ```

3. Start Frontend:
   ```bash
   cd frontend
   npm install
   npm run dev
   ```
