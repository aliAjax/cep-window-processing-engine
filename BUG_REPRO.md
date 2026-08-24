# Bug Reproduction

## Bug

Cancellation is lost across the segment watch request, coordinator, source, and retry path. A timed-out watch can leave the source blocked and retries running after the caller has returned.

## Trigger

Start a segment watch with a cancelable context, block the source read, and cancel the request while the retry operation is still eligible to run.

## Error

```text
request lost cancellation: <nil>
coordinator did not forward context: <nil>
source stayed blocked after cancellation
retry continued after cancellation: err=temporary calls=4
```
