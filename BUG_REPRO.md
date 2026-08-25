# Checkpoint Commit Error Identity Loss

## Bug

When a checkpoint operation and its cleanup both fail, the commit chain loses one or more underlying error identities. Save and write failures can be replaced by later cleanup failures, rollback failures are discarded, and a lease release wrapper prevents `errors.Is` from finding the backend error.

## Trigger

Exercise the lease release, service commit, transaction finish, and partition writer paths with paired operation and cleanup failures. Inspect each returned error with `errors.Is` for both the primary operation error and the later release, rollback, or close error.

## Error

The red verification reported:

- `release identity lost: release checkpoint lease "lease-1": lease backend unavailable`
- `save failure was overwritten: release checkpoint guard: guard release failed`
- `joined failures were not preserved: checkpoint operation failed`
- `write failure was overwritten: partition close failed`
