package runtime

import (
	"testing"
	"time"

	"github.com/example/cep-window-engine/internal/domain"
)

func TestEngineIngestOwnsCallerEvent(t *testing.T) {
	store := domain.NewStore()
	store.Patterns["p-1"] = domain.PatternDefinition{ID: "p-1", Version: "v1", Compiled: &domain.CompiledPlan{Steps: []domain.PlanStep{{EventType: "reading"}}}}
	engine := New(store)
	event := domain.Event{ID: "e-1", StreamID: "s-1", TenantID: "t-1", Type: "reading", EventTime: time.Unix(10, 0), Payload: map[string]any{"sensor": "a"}}

	matches, err := engine.Ingest(event)
	if err != nil || len(matches) != 1 {
		t.Fatalf("ingest failed: matches=%d err=%v", len(matches), err)
	}
	matches[0].Evidence[0] = "caller-overwrite"
	matches[0].Explain.Result = "false"
	stored, ok := store.GetMatch(matches[0].ID)
	if !ok || stored.Evidence[0] != "e-1" || stored.Explain.Result != "true" {
		t.Fatalf("stored match followed returned value mutation: %#v", stored)
	}
}
