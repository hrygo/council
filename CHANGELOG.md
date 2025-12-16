# Changelog

## [0.6.1] - 2025-12-16

### Fixed
- **CI**: Resolved `errcheck` and `staticcheck` lint errors in `agent_test.go` and `hub.go`.
- **Frontend**: Removed unused `useTranslation` import in `App.tsx` causing build failure.

## [0.6.0] - 2025-12-16

### Added
- **Three-Tier Memory Protocol**:
  - **Tier 1 (Quarantine)**: `quarantine_logs` table (PostgreSQL) and logging middleware.
  - **Tier 2 (Working Memory)**: Redis-backed `WorkingMemoryBuffer` with Ingress Filter.
  - **Tier 3 (LTM)**: Stub for `KnowledgePromoter` and Hybrid Retrieval (PGVector + Redis).
- **Defense Mechanisms (Safety)**:
  - `LogicCircuitBreaker`: Middleware to prevent infinite loops (recursion depth check).
  - `AntiHallucination`: Regex-based trigger to flag unverified claims (`[Specific Metric]`).
  - `MemoryMiddleware`: Auto-injection of node outputs into Memory Protocol.
- **Frontend Dual-Mode Architecture**:
  - **Stores**: `useLayoutStore` (Panel Persistence), `useConfigStore` (God Mode).
  - **MeetingRoom**: Resizable 3-pane layout (Canvas, Chat, Docs) with Fullscreen Focus.
  - **UI Features**: `CostEstimator` and `ParallelMessageRow` visualization.
- **Infrastructure**: Added `redis` service to `docker-compose.yml` and automated DB migration runner.

## [0.5.0] - 2025-12-16

### Added
- **Execution Runtime**: Implemented Session State Machine (`Pending`, `Running`, `Completed`, `Failed`).
- **WebSocket Infrastructure**: `Hub` and `ServeWs` for real-time event broadcasting (`/ws` endpoint).
- **API Integration**:
  - `POST /workflows/execute`: Integrated with Engine and Hub.
  - Runtime events (`StreamEvent`) pushed to connected clients via WebSocket.

## [0.4.0] - 2025-12-16

### Added
- **Workflow Engine**: Implemented `DAG` based execution engine.
  - **Core**: `GraphDefinition`, `Node`, `engine.go` (Concurrent execution with Factory pattern).
  - **Validation**: Cycle detection and connectivity checks.
  - **Processors**: `StartNode` (Inputs), `AgentNode` (LLM Integration), `EndNode` (Summarization).
  - **TDD**: 100% unit test coverage for engine components.

## [0.3.0] - 2025-12-15

### Added
- **Group & Agent Management**: Core domain logic, Entities, and Repositories.
- **API**: RESTful Handlers for Groups and Agents.
- **CI/CD**: GitHub Actions for Backend (Go) and Frontend (React/Vite).
- **Testing**: TDD adoption with high coverage (>92%).

## [0.2.0] - 2025-12-15

### Service Infrastructure
- **LLM & Embedding Separation**: Refactored `config` and `llm` package to decouple Chat Models (`LLMConfig`) from Vector Models (`EmbeddingConfig`).
- **Provider Support**: Added **SiliconFlow** (`siliconflow`), **DashScope** (`dashscope`), **Gemini** (`google`), and **Ollama** (`ollama`) support.
- **Embedding Standardization**: Standardized `Embed` interface across all providers. Activated `SiliconFlow` (`Qwen/Qwen3-Embedding-8B`, 1536 dim) by default.
- **Gemini Streaming**: Implemented native streaming support for Google Gemini.

### Added
- **Configuration**: Added `config/config.go` structural refactor.
- **Environment**: Added `.env.example` with detailed configuration guide.

## [0.1.0] - 2025-12-15

### Added
- **Project Structure**: Initialized root directories (`cmd`, `internal`, `pkg`, `prompts`) and Go module `github.com/hrygo/council`.
- **Backend**: Basic Go server structure with health check endpoint.
- **Frontend**: React + Vite + TailwindCSS setup with Zustand stores and i18n.
- **Infrastructure**: `docker-compose.yml` for PostgreSQL with pgvector support.
