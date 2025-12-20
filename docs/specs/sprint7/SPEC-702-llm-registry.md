# SPEC-702: Dynamic LLM Registry

## 1. Background
Currently, the system initializes a single `LLMProvider` at startup based on `LLM_PROVIDER` env var. This `LLMProvider` is injected into all agents. If an agent is configured to use a different provider (e.g., DeepSeek) but the system is running on Gemini, the request fails or is routed incorrectly.
The goal is to match `example/llm` flexibility where multiple providers can coexist, and agents select their provider dynamically.

## 2. Architecture Changes

### 2.1 LLM Registry (`internal/infrastructure/llm`)
Replace/Extend `Router` with a `Registry` struct that:
- Holds a map of initialized providers `map[string]LLMProvider`.
- `Get(providerName string)` method:
  - Checks if provider exists.
  - If not, initializes it using global config/env (Lazy Loading).
  - Returns the provider interface.

### 2.2 Node Factory & Processor
- Update `NodeDependencies` to hold `*llm.Registry` instead of `LLMProvider`.
- Update `AgentProcessor` to use `Registry.Get(activeAgent.ModelConfig.Provider)` at runtime.
- Update `WorkflowHandler` to use Registry.

## 3. Implementation Steps

### 3.1 Refactor LLM Package
- Modify `internal/infrastructure/llm/router.go` to implement `Registry`.
- Ensure it supports "openai", "gemini", "deepseek", etc. using `config`.

### 3.2 Update Core Components
- `internal/core/workflow/nodes/factory.go`: Change `LLM` to `Registry`.
- `internal/core/workflow/nodes/agent.go`: Update `Process` method to resolve provider.

### 3.3 Update Wiring (`main.go`)
- Initialize `Registry` instead of single `LLMProvider`.
- Pass registry to `NewWorkflowHandler` and `nodes.NewNodeFactory`.

## 4. verification Plan

### 4.1 Unit Tests
- Update `router_test.go` to verify Registry behavior (caching, lazy load).
- Update `agent_test.go` (if exists) or create new test with mocked Registry.

### 4.2 Manual Verification
- Configure `.env` with `GEMINI_API_KEY` and `DEEPSEEK_API_KEY`.
- Seed an agent with `provider: deepseek`.
- Seed an agent with `provider: gemini`.
- Run a workflow using both agents.
- Verify both succeed.
