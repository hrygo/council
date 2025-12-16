# API Specification v1.5 - UI Support

> **Version**: v1.0 | **Status**: Draft | **Base**: /api/v1

This specification defines the backend endpoints required to support the new UI/UX designs defined in `ui_design_v1.5.md`.

## 1. Dashboard (Council Hall)

### 1.1 Global Stats
**GET** `/dashboard/stats`
Returns aggregated metrics for the "Council Hall" overview.

**Response:**
```json
{
  "total_sessions": 42,
  "active_sessions": 3,
  "total_token_usage": 1542000,
  "total_cost_usd": 4.52,
  "system_status": "operational" // or "maintenance"
}
```

### 1.2 Group Cards (Enhanced)
**GET** `/groups`
*Extension of existing endpoint to include metrics.*

**Response (List Item):**
```json
{
  "id": "uuid",
  "name": "SaaS Architecture",
  "status": "active", // active, idle, archived
  "stats": {
    "token_usage_7d": 120000,
    "cost_7d": 0.35,
    "session_count": 5
  },
  "icon": "server"
}
```

---

## 2. Run Mode (Meeting Room)

### 2.1 Timeline Events
**GET** `/sessions/:id/timeline`
Retrieves the structured event log for the Timeline Sidebar.

**Response:**
```json
[
  {
    "id": "step-1",
    "node_id": "node-analyst-1",
    "name": "System Analysis",
    "type": "agent", // agent, logic, human_review
    "status": "completed", // pending, running, completed, failed
    "started_at": "2025-12-16T10:00:00Z",
    "ended_at": "2025-12-16T10:02:00Z",
    "agents": ["Architect"] // For parallel visibility
  }
]
```

### 2.2 Intervention Controls
**POST** `/sessions/:id/control`
Allocates user control commands to the execution engine.

**Request:**
```json
{
  "action": "pause" // pause, resume, stop
}
```

### 2.3 Human Review (Decision)
**POST** `/sessions/:id/review/:node_id`
Submits the user's decision for a `HumanReview` node.

**Request:**
```json
{
  "decision": "approve", // approve, reject, revise
  "feedback": "Please reconsider the security implications...", // Optional revision comments
  "modified_content": "..." // If user edited the proposal directly (Diff View)
}
```

---

## 3. Build Mode (Workflow Editor)

### 3.1 Templates
**GET** `/templates`
List available system and user templates.

**Response:**
```json
[
  {
    "id": "tpl-code-review",
    "name": "Strict Code Review",
    "description": "3-round debate focused on security and perf.",
    "complexity": "moderate",
    "estimated_cost": 0.15,
    "tags": ["dev", "security"]
  }
]
```

### 3.2 Validate Graph
**POST** `/workflows/validate`
Checks DAG validity before execution.

**Request:**
```json
{
  "graph": { ... } // React Flow graph format or internal adjacency list
}
```

**Response:**
```json
{
  "valid": false,
  "errors": [
    { "node_id": "node-3", "message": "Isolated node detected" },
    { "node_id": "node-5", "message": "Missing required input from 'Analyst'" }
  ]
}
```
