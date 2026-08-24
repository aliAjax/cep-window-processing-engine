# Execution Cancellation Propagation

## Bug

The execution request chain drops the caller's context at the timeout middleware, request wrapper, evaluator call, and retry loop. An already-cancelled request is therefore treated as active work.

## Trigger

Create an execution request with a cancelled parent context, then pass it through the middleware, request, evaluator, and retry paths. The evaluator still runs and the retry loop performs all configured attempts.

## Error

The reproduction observed `nil` for the context error at each downstream boundary even though the parent was cancelled. The evaluator completed after its wait and retry performed three operations instead of stopping.
