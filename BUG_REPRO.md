# Bug Reproduction

## Bug

规则首次部署失败后，即使自动重试成功，状态仍停在 `retrying`，规则也不会出现在可运行列表中。

## Trigger

创建一条规则，让首次部署进入失败状态，再触发一次成功的自动重试。随后读取规则状态并查询可运行规则集合。

## Observed Errors

验证时稳定出现以下结果：

```text
retrying state cannot transition to active
stored state = "retrying", want active
retry state = "retrying", want active
runnable rules = [], want [rule-active]
```
