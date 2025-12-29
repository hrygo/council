# SPEC-1203: Loop Analytics (Score Chart)

| Metadata     | Value                       |
| :----------- | :-------------------------- |
| **Title**    | Optimization Loop Analytics |
| **Type**     | Frontend Feature            |
| **Sprint**   | Sprint 12                   |
| **Priority** | P1                          |

## 1. Goal
Visualize the progress of the `council_optimize` workflow. Users should see how the score improves (or regresses) over iterations.

## 2. Requirements
- **Live Chart**: A line chart showing `Score` vs `Round Number`.
- **Indicators**:
    - **Baseline**: Initial score.
    - **Target**: The exit threshold (e.g., 90).
    - **Current**: Latest score.
- **Events**: Markers on the chart for significant events (e.g., "Surgeon applied patch", "Rollback").

## 3. Data Source
- `Session Context`: The backend now provides `session.context_data["score_history"]`.
- API: We need to expose this via the `GET /api/v1/sessions/:id` or WebSocket updates.
    - *Plan*: The `StreamEvent` from `LoopProcessor` already broadcasts the score. The frontend should aggregate these events or fetch the Session context on load.

## 4. UI Component
- **Library**: Recharts or Chart.js.
- **Location**: "Meeting Header" or a dedicated "Analytics" widget in the left/right panel.

## 5. Acceptance Criteria
- [ ] Line chart updates in real-time as `LoopProcessor` emits events.
- [ ] Y-axis: 0-100 Score. X-axis: Rounds.
- [ ] User can clearly see if the optimization is converging.
