# Bug Reproduction

## Bug

Evidence compaction, history, cache, and view construction retain shared slice backing arrays. Compaction or a late append can silently rewrite evidence that callers already received, while cached and historical snapshots drift apart.

## Trigger

Return a full evidence slice to a caller, store it in history and cache, build a view, then compact the original slice and append a late event. Mutating either the caller-owned input or a loaded snapshot exposes the same ownership leak.

## Observed Errors

```text
compact mutated input
history changed through caller alias
cache exposed aliased snapshot
view changed after late append
```
