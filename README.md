# The Council (理事会)

<div align="center">

<img src="https://img.shields.io/badge/version-0.14.0-blue" alt="Version" />
<img src="https://img.shields.io/badge/Go-1.21%2B-00ADD8?logo=go" alt="Go" />
<img src="https://img.shields.io/badge/React-19-61DAFB?logo=react" alt="React" />
<img src="https://img.shields.io/badge/license-MIT-green" alt="License" />

**A Visualized Multi-Agent Collaboration System & Personal Private Think Tank**

</div>

---

## ✨ Features

### 🤖 Multi-Agent Orchestration
- **Visual Workflow Builder**: Drag-and-drop DAG editor powered by React Flow
- **AI-Driven Generation**: Describe your workflow in natural language, let AI design it
- **Template Library**: Save and reuse workflow patterns

### 🧠 Intelligent Nodes
| Node Type       | Description                                             |
| :-------------- | :------------------------------------------------------ |
| **Agent**       | LLM-powered AI agents with customizable personas        |
| **Vote**        | Democratic decision-making with configurable thresholds |
| **Loop**        | Iterative refinement until consensus or max rounds      |
| **FactCheck**   | Web search integration for claim verification           |
| **HumanReview** | Pause execution for human oversight                     |

### 💡 Advanced Capabilities
- **Human-in-the-Loop**: Pause workflows for human decisions
- **Cost Estimation**: Real-time token & cost tracking before execution
- **Three-Tier Memory**: Quarantine → Working Memory → Long-Term Knowledge
- **Anti-Hallucination**: Built-in fact verification and circuit breakers

### 🎨 Modern UI/UX
- **Resizable Panels**: Workflow | Chat | Documents
- **Keyboard Shortcuts**: `Cmd+1/2/3` for panel focus, `Esc` to exit
- **KaTeX Rendering**: Math formula support in chat
- **Document References**: Clickable `[Ref: ID]` links

---

## 🚀 Quick Start

### Prerequisites

| Requirement    | Version |
| :------------- | :------ |
| Docker         | ≥ 20.10 |
| Docker Compose | v2.x    |
| Go             | ≥ 1.21  |
| Node.js        | ≥ 20    |

### One-Command Start

```bash
# Clone the repository
git clone https://github.com/hrygo/council.git
cd council

# Start all services (Docker + Backend + Frontend)
make start
```

**Access the application:**
- 🌐 Frontend: http://localhost:5173
- 🔌 Backend API: http://localhost:8080
- 📊 WebSocket: ws://localhost:8080/ws

### Manual Start

```bash
# 1. Start infrastructure (PostgreSQL + Redis)
make start-db

# 2. Start Backend
make start-backend

# 3. Start Frontend
make start-frontend
```

### Stop Services

```bash
make stop
```

---

## 🏛️ The Council: Out-of-Box Experience

The Council is the **built-in AI Governance Board** that comes pre-configured, allowing you to experience multi-agent collaboration immediately.

### Default Agents

| Agent                | Role                             | Model        | Strategy             |
| :------------------- | :------------------------------- | :----------- | :------------------- |
| 🛡️ **Value Defender** | Advocates for strategic value    | Gemini 3 Pro | Creative (temp: 0.9) |
| 🔍 **Risk Auditor**   | Identifies risks and gaps        | DeepSeek     | Logical (temp: 0.6)  |
| ⚖️ **Chief Justice**  | Synthesizes and delivers verdict | GLM-4.6      | Balanced (temp: 0.2) |

### Available Workflows

1. **Council Debate** - Single round three-way debate with verdict
2. **Council Optimize** - Iterative refinement loop with human-in-the-loop review

### Try It Now

```bash
# 1. Start the server
make start

# 2. Open browser
open http://localhost:5173

# 3. Select "The Council" group from sidebar
# 4. Create a new meeting with "Council Debate" workflow
# 5. Upload your document and watch the AI council deliberate!
```

### How the Optimize Loop Works

```
📄 Your Document
       ↓
🧠 Memory Retrieval (历史上下文)
       ↓
   ┌───┴───┐
   ↓       ↓
🛡️ Value  🔍 Risk
Defender  Auditor
   ↓       ↓
   └───┬───┘
       ↓
⚖️ Chief Justice (评分裁决)
       ↓
👤 Human Review (继续/应用/回滚)
       ↓
   [Score < 90?] → 🔄 Loop Back
   [Score ≥ 90?] → ✅ Complete
```

---

## 🏗 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (React SPA)                   │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Workflow   │  │    Chat     │  │     Document        │  │
│  │   Editor    │  │   Panel     │  │     Reader          │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└───────────────────────────┬─────────────────────────────────┘
                            │ REST + WebSocket
┌───────────────────────────┴─────────────────────────────────┐
│                      Backend (Go/Gin)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Workflow   │  │   Memory    │  │       LLM           │  │
│  │   Engine    │  │   Service   │  │    Providers        │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└───────────────────────────┬─────────────────────────────────┘
                            │
┌───────────────────────────┴─────────────────────────────────┐
│                     Infrastructure                          │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │ PostgreSQL  │  │    Redis    │  │   External APIs     │  │
│  │ (pgvector)  │  │   (Cache)   │  │  (LLM, Search)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Tech Stack

| Layer        | Technology                                       |
| :----------- | :----------------------------------------------- |
| **Frontend** | React 19, Vite, TailwindCSS, Zustand, React Flow |
| **Backend**  | Go 1.21+, Gin, gorilla/websocket                 |
| **Database** | PostgreSQL 16 + pgvector                         |
| **Cache**    | Redis                                            |
| **LLM**      | OpenAI, Gemini, DeepSeek, DashScope, Ollama      |

---

## 📖 Usage Guide

### 1. Creating a Workflow

#### Option A: Visual Builder
1. Navigate to **Workflow Builder** (`/builder`)
2. Drag nodes from the palette to the canvas
3. Connect nodes to define the execution flow
4. Configure each node via the **Property Panel**
5. Click **Save Workflow**

#### Option B: AI Wizard
1. Click the **Wizard** button
2. Describe your workflow in natural language:
   > "Create a debate between a Pro and Con agent, then have a Judge summarize"
3. Review and edit the generated workflow
4. Save to your library

### 2. Managing Agents

1. Navigate to **Agents** (`/agents`)
2. Click **Create Agent**
3. Configure:
   - **Name & Role**: e.g., "Legal Advisor"
   - **System Prompt**: Define the agent's persona
   - **Model**: Select LLM provider and model
   - **Capabilities**: Enable Web Search, Code Execution

### 3. Running a Workflow

1. Open a saved workflow
2. Click **Run** or navigate to **Meeting Room**
3. Watch the execution in real-time:
   - Left panel: Workflow progress visualization
   - Center panel: Agent conversations
   - Right panel: Document references
4. If a **HumanReview** node is reached:
   - Review the proposal in the modal
   - Click **Approve** or **Reject**

### 4. Keyboard Shortcuts

| Shortcut       | Action                  |
| :------------- | :---------------------- |
| `Cmd/Ctrl + 1` | Maximize Workflow panel |
| `Cmd/Ctrl + 2` | Maximize Chat panel     |
| `Cmd/Ctrl + 3` | Maximize Document panel |
| `Escape`       | Exit fullscreen mode    |

---

## 🔌 API Reference

### Workflows

| Method | Endpoint                     | Description              |
| :----- | :--------------------------- | :----------------------- |
| GET    | `/api/v1/workflows`          | List all workflows       |
| POST   | `/api/v1/workflows`          | Create a workflow        |
| GET    | `/api/v1/workflows/:id`      | Get workflow details     |
| PUT    | `/api/v1/workflows/:id`      | Update a workflow        |
| POST   | `/api/v1/workflows/generate` | AI-generate workflow     |
| POST   | `/api/v1/workflows/estimate` | Estimate execution cost  |
| POST   | `/api/v1/workflows/execute`  | Start workflow execution |

### Sessions

| Method | Endpoint                       | Description                  |
| :----- | :----------------------------- | :--------------------------- |
| POST   | `/api/v1/sessions/:id/control` | Pause/Resume/Stop session    |
| POST   | `/api/v1/sessions/:id/signal`  | Send signal to session       |
| POST   | `/api/v1/sessions/:id/review`  | Submit human review decision |

### Resources

| Method | Endpoint            | Description               |
| :----- | :------------------ | :------------------------ |
| CRUD   | `/api/v1/groups`    | Manage user groups        |
| CRUD   | `/api/v1/agents`    | Manage AI agents          |
| CRUD   | `/api/v1/templates` | Manage workflow templates |

### WebSocket

Connect to `ws://localhost:8080/ws` for real-time events:

```typescript
// Event types
type WSEventType = 
  | 'token_stream'           // Streaming LLM output
  | 'node_state_change'      // Node status updates
  | 'token_usage'            // Token consumption
  | 'execution:paused'       // Session paused
  | 'execution:completed'    // Session completed
  | 'human_interaction_required' // Need human review
  | 'error';                 // Error occurred
```

---

## 🛠 Development

### Project Structure

```
council/
├── cmd/council/          # Application entrypoint
├── internal/
│   ├── api/              # HTTP handlers & WebSocket
│   ├── core/             # Business logic
│   │   ├── workflow/     # Workflow engine
│   │   ├── memory/       # Three-tier memory
│   │   └── middleware/   # Safety mechanisms
│   └── infrastructure/   # External integrations
├── frontend/
│   └── src/
│       ├── components/   # Reusable UI components
│       ├── features/     # Feature modules
│       ├── stores/       # Zustand state management
│       └── hooks/        # Custom React hooks
├── docs/
│   ├── specs/            # Feature specifications
│   └── reports/          # Audit reports
└── prompts/              # LLM prompt templates
```

### Commands

```bash
# Development
make start          # Start all services
make stop           # Stop all services
make restart        # Restart all services

# Quality
make test           # Run all tests
make lint           # Lint backend
npm run lint        # Lint frontend (in frontend/)

# Database
make reset-db       # Reset database
make migrate        # Run migrations
```

---

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.

---

## 🙏 Acknowledgments

- [React Flow](https://reactflow.dev/) - Workflow visualization
- [Gin](https://gin-gonic.com/) - HTTP framework
- [pgvector](https://github.com/pgvector/pgvector) - Vector similarity search
- [Zustand](https://zustand-demo.pmnd.rs/) - State management

---

<div align="center">
  <sub>Built with ❤️ by the Council Team</sub>
</div>
