package windowing

import (
	"testing"
	"time"

	"github.com/example/cep-window-engine/internal/domain"
)

func TestBufferSnapshotOwnsEventPayload(t *testing.T) {
	buffer := NewBuffer(domain.WindowPolicy{Size: time.Hour})
	event := domain.Event{ID: "e-1", EventTime: time.Unix(10, 0), Payload: map[string]any{"temperature": 4}}
	buffer.Add(event)
	event.Payload["temperature"] = 99

	first := buffer.Snapshot()
	if got := first[0].Payload["temperature"]; got != 4 {
		t.Fatalf("buffer retained caller payload alias: %v", got)
	}
	first[0].Payload["temperature"] = -10
	if got := buffer.Snapshot()[0].Payload["temperature"]; got != 4 {
		t.Fatalf("snapshot mutation changed buffered event: %v", got)
	}
}
