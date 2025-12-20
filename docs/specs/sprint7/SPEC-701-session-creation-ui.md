# SPEC-701: Session Creation UI

## 1. Background
Currently, the `/chat` (Meeting Room) page is a passive receiver of WebSocket events. Users cannot initiate a new session from the UI, requiring manual API calls (cURL) to start a workflow. This creates a broken user experience.

## 2. Goals
- Provide a user-friendly interface to start a new session directly from the frontend.
- Allow selecting a workflow template.
- Allow inputting initial parameters (e.g., `topic`).
- Automatically connect to the session upon creation.

## 3. User Experience (UX)

### 3.1 Entry Point
- **Empty State**: When user visits `/chat` and no session is active/connected, show a "Start New Session" placeholder/button in the center or a modal overlay.
- **Home Page**: The "Start Session" button on Home page should lead to `/chat` with this initialization flow triggered.

### 3.2 Session Starter Modal/Card
- **Template Selection**: Dropdown or Card grid to select from "System Templates" (e.g., Council Debate) or User Templates.
- **Parameter Configuration**:
  - For MVP: A simple text area or input for `topic`.
  - Future: Dynamic form based on Input Schema.
- **Action**: "Launch Council" button.

### 3.3 Flow
1. User selects "Council Debate".
2. User enters Topic: "Should AI be regulated?".
3. User clicks "Launch".
4. Frontend POSTs to `/api/v1/workflows/execute`.
5. Frontend receives `session_id`.
6. Frontend initializes `useSessionStore` with new ID.
7. `useWebSocketRouter` connects and stream begins.

## 4. Technical Implementation

### 4.1 Frontend Components
- **`SessionStarter`**: New component.
  - Uses `useTemplates` to fetch options.
  - Maintains local state for selection and inputs.
- **Update `MeetingRoom.tsx`**:
  - If `!currentSession`, render `SessionStarter` overlay.
- **Update `useSessionStore`**:
  - Ensure `initSession` clears previous state correctly.

### 4.2 Backend API
- No changes required (`POST /api/v1/workflows/execute` exists).

## 5. Acceptance Criteria
- [ ] User can see a "Start Session" UI when entering `/chat`.
- [ ] User can select "Council Debate" template.
- [ ] User can input a topic.
- [ ] Clicking Launch successfully starts the backend process.
- [ ] The chat interface immediately reflects the new session (WebSocket connects).
