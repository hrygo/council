---
description: Run an iterative design optimization loop using The Council (Dialecta) with human-in-the-loop application.
---

# Design Optimization Loop
// turbo-all

This workflow guides you through the process of refining the design draft using AI debates and consistency checks, with a focus on observability and stability.

## Parameters
*   `max_loops`: (Optional) Maximum number of iteration cycles. Default: 3.

## Process Overview

1.  **Extract History**: Agent intelligently summarizes past debate verdicts.
2.  **Convene**: The AI Council debates the current draft.
3.  **Monitor**: Real-time status with throttled updates (avoiding noise).
4.  **Verify**: Flip-Flop detection.
5.  **Apply & Loop**: Apply changes and loop if conditions met.

## Step-by-Step Instructions

### Step 1: Compress History Context (Agent Task)

**Before starting**: Check `docs/reports/` for recent debate reports.

*   **Task**: Read the most recent 2-3 debate reports.
*   **Action**: Create or Update `docs/reports/history_summary.md`.
    *   *Format*: A chronological list of "Verdicts" and "Key Decisions".
    *   *Goal*: This file serves as the "Common Law" (判例法).

### Step 2: Convene The Council

Run the Python script.
*   Note: Output is now throttled (updates every ~2s) and transient status lines are excluded from the saved report.

```bash
python3 scripts/dialecta_debate.py docs/design_draft.md docs/PRD.md --history docs/reports/history_summary.md
```

### Step 3: Verify Consistency (Agent Logic)

Read the new report (path printed at end of Step 2).
Compare "Verdict" against `docs/reports/history_summary.md`.

*   **Flip-Flop Check**:
    *   **If FLIP-FLOP DETECTED**: STOP. Inform user.
    *   **If CONSISTENT**: Proceed.

### Step 4: Apply Changes & Check Loop Condition

If report is valid:
1.  Read **Verdict** and **Recommendations**.
2.  Edit `docs/design_draft.md`.
3.  **Commit**: Update `version` in header.
4.  **Loop Check**:
    *   Check current loop count.
    *   **IF** `Current_Loop < max_loops` **AND** `Verdict != "Approved"`:
        *   Go to **Step 1** for next iteration.
    *   **ELSE**:
        *   Terminate workflow. Report final status.
