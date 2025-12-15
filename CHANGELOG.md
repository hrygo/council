# Changelog

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
