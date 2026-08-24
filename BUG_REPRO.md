# Bug Reproduction

## Bug

Batch delivery holds leases until the whole function returns, individual lease release does not mark the lease released, failed sends skip service-level release, and transaction cleanup discards the original delivery error.

## Trigger

Send a batch through a bounded lease pool and make a sink fail after several matches. Then pass the delivery error through transaction cleanup. Repeating the batch consumes the remaining lease capacity while cleanup either returns only the rollback error or returns nil.

## Error

The first failure is reported as `deliver match m2: send failed: m2`. Subsequent work reaches resource exhaustion because acquired leases are not returned. When rollback also fails, `Finish` returns `rollback exploded` and the original delivery error is absent from the error chain.
