---
description: Run an iterative design optimization loop using The Council (Dialecta) with human-in-the-loop application.
---

# Design Optimization Loop
// turbo-all

This workflow guides you through the process of refining the design draft using AI debates and consistency checks, with a focus on observability and stability.

## Process Overview

1.  **Extract History**: Agent intelligently summarizes past debate verdicts into a compact context file to prevent amnesia.
2.  **Convene**: The AI Council debates the current draft against the PRD.
3.  **Monitor**: Observe the real-time "Thinking" output to ensure the debate is progressing.
4.  **Verify**: Agent runs a Flip-Flop detection check against the history summary.
5.  **Apply**: You (the Human or Agent) apply changes.

## Step-by-Step Instructions

### Step 1: Compress History Context (Agent Task)

**Before starting**: Check `docs/reports/` for recent debate reports.

*   **Task**: Read the most recent 2-3 debate reports.
*   **Action**: Create or Update `docs/reports/history_summary.md`.
    *   *Format*: A chronological list of "Verdicts" and "Key Decisions".
    *   *Goal*: This file serves as the "Common Law" (判例法) for the system to reference.

### Step 2: Convene The Council

Run the Python script to generate the debate report. The script has been optimized to stream output in real-time.

```bash
python3 scripts/dialecta_debate.py docs/design_draft.md docs/PRD.md --history docs/reports/history_summary.md
```

**Observability Tip**:
*   You will see lines like `⏳ Status: 🔵 Pro [Thinking...]` streaming in the terminal.
*   If the output halts for >60 seconds, check the `command_status`.

### Step 3: Verify Consistency (Agent Logic)

Read the new report generated (the path is printed at the end of Step 2).
Compare its "Verdict" against `docs/reports/history_summary.md`.

*   **Flip-Flop Check**: Does this new verdict contradict a decision made in the *immediate* previous session?
    *   *Example*: Session N-1 said "Remove Feature X", Session N says "Bring back Feature X" without new context.
    *   **If FLIP-FLOP DETECTED**: 
        *   STOP. Do not edit the design draft.
        *   Inform the user: "Detected policy flip-flop. Manual review required."
    *   **If CONSISTENT**: Proceed to Step 4.

### Step 4: Apply Changes

If the report is valid:
1.  Read the **Verdict** and **Actionable Recommendations** sections carefully.
2.  Edit `docs/design_draft.md`.
3.  **Commit**: Update the `version` in the draft's header.
4.  **Loop**: Go back to Step 1 for the next iteration if needed.
