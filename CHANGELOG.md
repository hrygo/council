# Changelog

## [0.14.0] - 2025-12-21

### Added
- **Meeting Terminology Unification**: Complete migration from "Chat" to "Meeting" terminology across the entire application (UI, Routes, i18n).
- **Dynamic LLM Registry (Sprint 7)**:
  - Database-backed LLM provider and model management.
  - New `llm_providers` and `llm_models` tables with automatic migrations.
  - Dynamic Model Selector in Agent settings fetching options from API.
- **Robust Session Initialization**:
  - `SessionStarter` now gracefully handles missing or malformed template graph data.
  - Added comprehensive unit tests for session startup flow.

### Changed
- **Database Driver**: Migrated from `lib/pq` to `pgx/v5` for enhanced performance and modern PostgreSQL features.
- **Backend Routing**: Improved LLM provider name normalization (e.g., case-insensitive "Gemini").
- **I18n**: Restructured translation keys for better "Meeting" context.

### Fixed
- **SessionStarter Crash**: Fixed a critical runtime error when accessing `nodes` on empty graph definitions.
- **Lint & Types**: Resolved several TypeScript and Go linting warnings.

## [0.13.0] - 2025-12-20

### Added
- **High-End Landing Page**: New dynamic HomePage with animations, glassmorphism UI, feature grid, and statistics.
- **Routing Overhaul**:
  - `/` is now the HomePage (Landing).
  - `/chat` is now the Meeting Room (Session).
  - Navigation updated to reflect new structure.
- **Save Notifications**: Modern `Toast` notification system for workflow save actions (replacing `alert()`).
- **UI Enhancements**:
  - Unified "Exit Fullscreen" button with consistent dark mode support.
  - Clickable Logo returning to Home.
  - Prevent accidental text selection on navigation buttons.

### Fixed
- **Workflow Save Logic**: Resolved 404/500 errors by stripping IDs from Templates/Wizard results to force reliable creation (POST).
- **Dark Mode Visibility**: Fixed invisible text in workflow nodes by enforcing custom node mapping over React Flow defaults.
- **Layout**: Fixed bottom styling/overflow issues in Landing Page and sidebar interactions.

### Changed
- **Navigation**: Removed redundant `/meeting` route entry.

## [0.12.0] - 2025-12-20

### Added
- **QA Enhancement**: Comprehensive test coverage improvements across backend and frontend.
  - Handler coverage: 81.4%
  - Frontend stores coverage: 82%
  - Infrastructure tests for DB, Cache, Search
- **PDF Tooling**: Universal Markdown to PDF converter with theme support.
- **Auth**: User sign-out with `next-auth` integration.
- **UI**: Regenerate button and improved chat input auto-resize.

### Changed
- **GEMINI.md**: Added "Pragmatism" (务求实效) as first principle.
- **Refactor**: `workflow.Repository` interface for better testability.
- **.gitignore**: Optimized to exclude test artifacts.

---

## [0.11.0] - 2025-12-17

### Added
- **Docs Overhaul**: Comprehensive README update with Architecture diagrams and Usage Guide.
- **Development Plan**: Added **Phase 2 (Post-MVP)** roadmap including i18n (P0) and E2E Testing.

### Changed
- **Cleanup**: Removed manual TOC from README for cleaner layout.
- **Fix**: Resolved SVG badge rendering issues in README.
- **Refactor**: Eliminated all legacy TODO items in codebase.

---

## [0.10.0] - 2025-12-17

### Added
- **Sprint 4 Complete**: Advanced Features for MVP Readiness.
- **Human-in-the-Loop (SPEC-301, SPEC-405)**:
  - **HumanReviewModal**: Frontend modal for pausing workflow execution and awaiting human decisions.
  - **Backend**: `StatusSuspended` node state, `ErrSuspended` signal, and `POST /sessions/:id/review` API for Approve/Reject/Modify actions.
  - **Engine**: `ResumeNode()` method to resume suspended workflows with injected output.
- **Cost Estimation (SPEC-302, SPEC-407)**:
  - **CostEstimator Widget**: Real-time cost/token/agent breakdown displayed in `WorkflowEditor`.
  - **Backend**: `EstimateWorkflowCost()` logic with model-based pricing and `POST /workflows/estimate` API.
- **Knowledge & Experience**:
  - **Document Reference (SPEC-303)**: `[Ref: ID]` in chat messages transformed to clickable links opening Document Reader.
  - **KaTeX Rendering (SPEC-305)**: LaTeX math formula support via `remark-math` + `rehype-katex`.
  - **Fullscreen Shortcuts (SPEC-304)**: `Cmd/Ctrl + 1/2/3` for panel maximize, `Escape` to exit.
- **Three-Tier Memory Protocol (SPEC-408)**:
  - **Ingress Filter**: Confidence and content-length validation before Working Memory insertion.
  - **Cleanup**: `CleanupWorkingMemory()` placeholder for future scheduled purging.
- **Web Search Integration (SPEC-411)**:
  - **TavilyClient**: Implementation for web search API (`internal/infrastructure/search`).
  - **FactCheckProcessor**: Integrated real Tavily search for claim verification.

### Changed
- **Development Plan**: All Sprint 1-4 items marked as Done. Project at MVP milestone.

---

## [0.9.0] - 2025-12-17

### Added
- **Sprint 3 Complete**: Advanced Workflow Builder & Backend Processors.
- **Workflow Builder Enhancements**:
  - **Property Panel (SPEC-201)**: Dynamic configuration panel for all node types with strict type validation.
  - **New Node Types (SPEC-202, SPEC-203)**:
    - **Vote**: Threshold-based approval logic.
    - **Loop**: Iteration control with consensus or max-round exit conditions.
    - **FactCheck**: Source-based verification logic.
    - **HumanReview**: Workflow suspension/resume capability.
  - **Template System (SPEC-204, SPEC-205)**: Full stack implementation (Sidebar, Save Modal, Backend API) for workflow templates.
  - **Wizard Mode (SPEC-206)**: AI-driven workflow generation from natural language descriptions.
- **Backend Core**: 
  - **Processor Registry**: Implemented `VoteProcessor`, `LoopProcessor`, `FactCheckProcessor`, and `HumanReviewProcessor`.
  - **Template API**: Endpoints for Template CRUD operations (`/api/v1/templates`).
- **Frontend Visuals**: 
  - **Custom Nodes**: Rich visual React Flow components with icons and status indicators.

## [0.8.0] - 2025-12-17

### Added
- **Sprint 2 Complete**: Full implementation of Groups and Agents Management.
- **Groups Management (SPEC-101, 102)**:
  - **CRUD Operations**: Complete management of User Groups for multi-tenant isolation.
  - **UI/UX**: Dedicated `/groups` page with Grid view and Create/Edit modals.
- **Agents Management (SPEC-103, 104)**:
  - **Agent Factory**: Interface to design specialized AI agents.
  - **Model Configuration (SPEC-105)**: Advanced Model Selector with support for OpenAI, Anthropic, Google, DeepSeek, and DashScope.
  - **Capabilities**: Toggle for Web Search and Code Execution.
- **Backend**: 
  - **Extended `Agent` entity with `ModelConfig` supporting `top_p` and `max_tokens`.
  - **JSONB schema flexibility leveraged for seamless updates.
- **DevEx**: Added `/agents` and `/groups` routes to primary `App` navigation.

## [0.7.0] - 2025-12-17

### Added
- **Sprint 1 Complete**: Full implementation of Sprint 1 Specifications (SPEC-001 to SPEC-005).
- **Session Store (SPEC-001)**: `useSessionStore` rewrite with strict normalized state, streaming support, and `useAuthStore` separation.
- **Workflow Run Store (SPEC-002)**: Complete runtime stage management (Pending/Running/Completed/Failed) with execution timers and controls.
- **Chat Panel (SPEC-003, SPEC-004)**: 
  - **Grouped Messages**: Visual grouping by Node ID and Agent identity.
  - **Parallel Layout**: Side-by-side rendering for parallel agent execution steps.
- **WebSocket Optimization (SPEC-005)**:
  - **Robustness**: Auto-reconnect with exponential backoff and Heartbeat mechanism.
  - **Type Safety**: New `useWebSocketRouter` hook with strict `WSMessage` typing.
- **Quality Assurance**: Added "Strict Quality Gates" to `GEMINI.md` requiring Acceptance Criteria + CI checks for every Spec.

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
