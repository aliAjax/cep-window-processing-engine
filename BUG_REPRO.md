# Concurrent Watermark Advancement Failure

## Bug

Concurrent partition updates race while reading and publishing shared watermark state. Snapshot readers can also share the tracker's internal map with writers. Separately, the coordinator can publish a stale global value, and a configured maximum step rejects the first watermark report.

## Trigger

Advance the same watermark barrier and tracker from concurrent goroutines while taking snapshots. Then advance a partition and inspect the global watermark, and submit the first watermark through a policy with a maximum step configured.

## Error

The red verification observed a race between reads and writes in `barrier.go`, followed by `race detected during execution of test`. Concurrent tracker access also produced `fatal error: concurrent map read and map write`. The coordinator returned `2023-11-15 06:13:20 +0800 CST` instead of `2023-11-15 06:13:25 +0800 CST`, and the first policy advance failed with `first watermark was rejected`.
