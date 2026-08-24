# Evidence Slice Aliasing

## Bug

The evidence cache, history, compaction, and view boundaries retain shared slice storage. Mutating an input or a returned slice can therefore change cached data, historical versions, compact inputs, and previously built views.

## Trigger

Store evidence and then mutate the caller-owned input or a value returned by `Load` or `Versions`. The same behavior is visible when `Compact` processes duplicate events or when the caller changes the items used to build a view.

## Error

The regression checks fail with messages including:

```text
cache followed input mutation: "mutated-input"
Compact overwrote its input
history followed input mutation: "mutated-input"
view followed input mutation: "mutated"
```
