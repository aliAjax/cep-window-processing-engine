# Bug Reproduction

## Bug

事件 payload、窗口快照和匹配结果之间共享嵌套 map 或 slice。调用方修改已经提交或读取的数据时，历史状态与内存 sink 中的结果会一起变化。

## Trigger

向引擎摄入带嵌套 payload 的事件，随后修改调用方保留的原始值；再读取事件历史和窗口快照。对匹配结果执行同样的写后修改，并再次从 store 或内存 sink 读取。

## Observed Errors

```text
stored payload followed caller mutation: caller-mutated
buffer retained caller payload alias: 99
stored match followed returned value mutation: Evidence:["caller-overwrite"]
sink retained input aliases: Evidence:["input-mutated"]
```
