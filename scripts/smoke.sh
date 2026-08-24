#!/usr/bin/env bash
set -euo pipefail
base="${CEP_SMOKE_BASE:-http://127.0.0.1:8080}"
curl -fsS "$base/healthz" >/dev/null
curl -fsS -X POST "$base/api/v1/streams" -H 'content-type: application/json' -d '{"id":"orders","tenant_id":"t1","name":"orders","partitions":1}' >/dev/null
curl -fsS -X POST "$base/api/v1/schemas" -H 'content-type: application/json' -d '{"id":"order-v1","stream_id":"orders","version":"1","fields":{"user":"string"}}' >/dev/null
curl -fsS -X POST "$base/api/v1/patterns/validate" -H 'content-type: application/json' -d '{"expression":"purchase"}' | grep -q 'true'
curl -fsS -X POST "$base/api/v1/patterns/publish" -H 'content-type: application/json' -d '{"id":"purchase-rule","tenant_id":"t1","version":"1","expression":"purchase"}' >/dev/null
curl -fsS -X POST "$base/api/v1/events" -H 'content-type: application/json' -d '{"id":"e1","stream_id":"orders","tenant_id":"t1","type":"purchase","payload":{"user":"u1"}}' | grep -q accepted
curl -fsS "$base/api/v1/matches" | grep -q purchase-rule
curl -fsS -X POST "$base/api/v1/runtimes/checkpoint" | grep -q checkpoint_id
curl -fsS -X POST "$base/api/v1/runtimes/replay" -H 'content-type: application/json' -d '{"stream_id":"orders"}' | grep -q replayed
echo "smoke ok"
