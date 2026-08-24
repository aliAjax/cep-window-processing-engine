package domain

import (
	"testing"
	"time"
)

func TestStoreListEventsOwnsPayload(t *testing.T) {
	store := NewStore()
	event := Event{ID: "e-1", StreamID: "s-1", EventTime: time.Unix(10, 0), Payload: map[string]any{"zone": "north"}}
	store.Append(event)
	event.Payload["zone"] = "caller-mutated"

	first := store.ListEvents("s-1")
	if got := first[0].Payload["zone"]; got != "north" {
		t.Fatalf("stored payload followed caller mutation: %v", got)
	}
	first[0].Payload["zone"] = "snapshot-mutated"
	second := store.ListEvents("s-1")
	if got := second[0].Payload["zone"]; got != "north" {
		t.Fatalf("later snapshot followed earlier result mutation: %v", got)
	}
}
