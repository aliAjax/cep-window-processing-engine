# CEP Window Engine

纯 Go 声明式复杂事件处理与流窗口计算引擎参考实现。它接收有序或乱序事件，使用受限模式 DSL 编译为有限状态执行计划，配合事件时间窗口、水位线、去重和 checkpoint/replay 生成带证据与解释树的匹配结果。

## 快速运行

```bash
go run ./cmd/cep-engine
curl http://127.0.0.1:8080/healthz
```

The embedded operations console is available at `http://127.0.0.1:8080/console/`.

核心 API：`POST /api/v1/streams`、`POST /api/v1/schemas`、`POST /api/v1/patterns/validate`、`POST /api/v1/patterns/publish`、`POST /api/v1/events`、`GET /api/v1/matches`、`GET /api/v1/matches/{id}/explain`、`POST /api/v1/runtimes/checkpoint` 和 `POST /api/v1/runtimes/replay`。

模式示例：`login -> purchase`，谓词可写为 `login[user_id=42]`。解析器只允许受限 AST（sequence、not、within、predicate），禁止脚本执行；编译器提供复杂度上限。

## 模块

`internal/domain` 定义 Event/Pattern/Window/Match 聚合；`pattern/parser` 与 `pattern/compiler` 完成 DSL 校验和计划编译；`runtime` 管理去重、证据和解释；`windowing` 提供排序窗口、过期和 session 合并；`watermark`、`dispatch`、`admission`、`provider` 与 `execution` 负责分区推进、任务分发、规则准入和隔离执行；`checkpoint` 使用 SHA-256 内容寻址快照；`sink`、`quota`、`repository`、`observability`、`enterprise` 为可替换端口和运维适配层。

## 验证

```bash
go test ./...
go test -race ./...
go vet ./...
go build ./...
scripts/smoke.sh
```

生产部署可使用 `Dockerfile` 和 `deployments/kubernetes.yaml`。checkpoint 目录通过 `CEP_CHECKPOINT_DIR` 配置，监听地址通过 `CEP_ADDR` 配置。
