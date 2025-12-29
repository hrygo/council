# SPEC-1104: Context Synthesizer (Memory Pruning)

| Metadata     | Value                                 |
| :----------- | :------------------------------------ |
| **Title**    | Context Synthesizer (History Pruning) |
| **Author**   | Antigravity                           |
| **Status**   | Draft                                 |
| **Priority** | P1                                    |
| **Sprint**   | Sprint 11                             |

## 1. Goal
Implement a dedicated node to manage the "Infinite Context" problem. Instead of retrieving random chunks via RAG, this node **deterministically actively summarizes** old rounds to keep the prompt lean and relevant.

## 2. Logic: The Rolling Window

**Inputs**:
-   `HistorySummary` (Current file content)
-   `NewVerdict` (Latest Adjudicator report)

**Algorithm**:
1.  Parse `HistorySummary`.
2.  Append `NewVerdict` to `## Chronological Verdicts`.
3.  **Count Check**: If `Count(Verdicts) > 3`:
    -   Take the **Oldest** verdict.
    -   Summarize it into a "One-Liner" (e.g., "Round 1: Improved formatting (+5)").
    -   Move it to `## Legacy Context`.
    -   Delete full text from `Verdicts`.
4.  **Output**: New `HistorySummary` text.

## 3. Implementation

New Node: `InternalProcessor` (Local, no LLM required for basic moving, but LLM needed for summarization).

*Optimization*: For V1, we can use a heuristic (extract the "One-Liner" already present in the report) to avoid an extra LLM call.

## 4. TDD Strategy

### 4.1 Test Cases (`nodes/context_synth_test.go`)
-   [ ] `TestSynth_Append`: Add round 1.
-   [ ] `TestSynth_Prune`: Add round 4, verify round 1 moves to Legacy.

## 5. Steps
1.  Create `context_synthesizer.go` node.
2.  Implement Markdown parsing logic.
