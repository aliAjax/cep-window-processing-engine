# Bug Reproduction

## Bug

分区投递中任一 worker 失败后，请求可能挂起或因向已关闭 channel 发送而崩溃；正常分片结果也会丢失，错误流消费方无法退出。

## Trigger

并发投递多个分区任务，让其中一个分区返回错误，同时消费结果与错误 channel。重复运行可观察 WaitGroup 时序竞争、失败后的 feeder 阻塞、成功结果被丢弃以及 Errors channel 不关闭。

## Observed Errors

```text
panic: send on closed channel
fatal error: all goroutines are asleep - deadlock!
Drain returned 0 results (expected 5 successful out of 6 jobs)
Errors range blocks FOREVER: Close() never closes the Errors channel
```
