package runtime

import (
	"fmt"
	"github.com/example/cep-window-engine/internal/dedup"
	"github.com/example/cep-window-engine/internal/domain"
	"github.com/example/cep-window-engine/internal/windowing"
	"sync"
	"time"
)

type Engine struct {
	mu      sync.Mutex
	store   *domain.Store
	buffers map[string]*windowing.Buffer
	dedup   *dedup.Cache
	paused  bool
}

func New(s *domain.Store) *Engine {
	return &Engine{store: s, buffers: map[string]*windowing.Buffer{}, dedup: dedup.New(10 * time.Minute)}
}
func (e *Engine) Pause()       { e.mu.Lock(); e.paused = true; e.mu.Unlock() }
func (e *Engine) Resume()      { e.mu.Lock(); e.paused = false; e.mu.Unlock() }
func (e *Engine) Paused() bool { e.mu.Lock(); defer e.mu.Unlock(); return e.paused }
func (e *Engine) Ingest(ev domain.Event) ([]domain.Match, error) {
	if e.Paused() {
		return nil, fmt.Errorf("runtime paused")
	}
	if !e.dedup.Check(ev.Fingerprint()) {
		return nil, nil
	}
	// Isolate the event from the caller's payload map before storing it in the
	// history log and window buffers; otherwise later caller mutations would
	// rewrite historical events and window snapshots in place.
	ev = ev.Clone()
	e.store.Append(ev)
	e.mu.Lock()
	defer e.mu.Unlock()
	out := []domain.Match{}
	for id, p := range e.store.Patterns {
		if p.Compiled == nil || len(p.Compiled.Steps) == 0 {
			continue
		}
		key := ev.StreamID + ":" + id
		buf := e.buffers[key]
		if buf == nil {
			buf = windowing.NewBuffer(domain.WindowPolicy{Size: 5 * time.Minute})
			e.buffers[key] = buf
		}
		buf.Add(ev)
		if ev.Type != p.Compiled.Steps[0].EventType {
			continue
		}
		match := domain.Match{ID: fmt.Sprintf("m-%d", time.Now().UnixNano()), PatternID: id, TenantID: ev.TenantID, StreamID: ev.StreamID, RuleVersion: p.Version, StartedAt: ev.EventTime, CompletedAt: time.Now(), Evidence: []string{ev.ID}, Explain: domain.ExplainNode{Operator: "sequence", Result: "true", EventID: ev.ID}}
		e.store.PutMatch(match)
		out = append(out, match.Clone())
	}
	return out, nil
}
func (e *Engine) Replay(stream string) (int, error) {
	events := e.store.ListEvents(stream)
	count := 0
	for _, ev := range events {
		if _, err := e.Ingest(ev); err != nil {
			return count, err
		}
		count++
	}
	return count, nil
}
