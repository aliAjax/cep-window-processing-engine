# Bug Reproduction

## Bug

The zero-value subscription registry, session store, fanout queue, and snapshotter write to uninitialized maps. The first operation on any of these components panics.

## Trigger

Run any first-write path on a zero-value component: add the first subscriber, open the first session, enqueue the first notification, or capture the first snapshot.

## Error

```text
panic: assignment to entry in nil map
```

The four reproduction tests fail in `Registry.Add`, `SessionStore.Open`, `Fanout.Enqueue`, and `Snapshotter.Capture`, respectively.
