# SPEC-1103: Logic Loop Processor (Delta Logic)

| Metadata     | Value                                   |
| :----------- | :-------------------------------------- |
| **Title**    | Logic Loop Processor (Delta & Rollback) |
| **Author**   | Antigravity                             |
| **Status**   | Draft                                   |
| **Priority** | P0                                      |
| **Sprint**   | Sprint 11                               |

## 1. Goal
Upgrade the `LoopProcessor` node to be "Logic Aware". It should track the **Velocity (Score Delta)** of the optimization process to detect regression (Rollback) or stagnation (Stasis).

## 2. Logic Features

### 2.1 State Management (Context)
The Workflow Engine must persiste the `score_history` in the session context.
`context.scores = [85, 88, 72]`

### 2.2 Delta Calculation
`Delta = CurrentScore - PreviousScore`

-   **Regressions (`Delta < -10`)**:
    -   Action: **ROLLBACK**.
    -   Signal: Output `rollback_signal` event.
    -   **VFS Action**: The `LoopProcessor` (or a dedicated handler) will revert the VFS `CurrentVersion` to `CurrentVersion - 1`.
-   **Stagnation (`Delta between -2 and +2` for 2 rounds)**:
    -   Action: **COOLING**.
    -   Signal: Inject "Try a different angle" instruction to next round.

## 3. JSON Parsing

The Adjudicator's output is now a mixed Markdown + JSON. The Loop Node must robustly extract the JSON block.

```regex
```json\s*(\{.*\})\s*```
```

## 4. TDD Strategy

### 4.1 Test Cases (`nodes/loop_logic_test.go`)
-   [ ] `TestLoop_Extraction`: Feed text with embedded JSON, assert Struct extraction.
-   [ ] `TestLoop_Delta_Positive`: Input: 80 -> 90. Output: `continue`.
-   [ ] `TestLoop_Delta_Negative`: Input: 90 -> 70. Output: `rollback_signal`.
-   [ ] `TestLoop_Delta_Stasis`: Input: 90 -> 91 -> 91. Output: `stasis_signal`.

## 5. Implementation Steps
1.  Enhance `LoopProcessor` struct with `ScoreHistory []int`.
2.  Implement `extractJSON` helper.
3.  Implement Delta logic in `Process` method.
