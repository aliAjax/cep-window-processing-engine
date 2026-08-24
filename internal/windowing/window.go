package windowing

import (
	"github.com/example/cep-window-engine/internal/domain"
	"sort"
	"sync"
	"time"
)

type Buffer struct {
	mu        sync.Mutex
	policy    domain.WindowPolicy
	events    []domain.Event
	watermark time.Time
}

func NewBuffer(p domain.WindowPolicy) *Buffer { return &Buffer{policy: p, events: []domain.Event{}} }
func (b *Buffer) Add(e domain.Event) []domain.Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.events = append(b.events, e)
	sort.SliceStable(b.events, func(i, j int) bool { return b.events[i].EventTime.Before(b.events[j].EventTime) })
	return b.expireLocked(e.EventTime)
}
func (b *Buffer) Advance(w time.Time) []domain.Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	if w.Before(b.watermark) {
		return nil
	}
	b.watermark = w
	return b.expireLocked(w)
}
func (b *Buffer) expireLocked(now time.Time) []domain.Event {
	if b.policy.Size <= 0 {
		return nil
	}
	cutoff := now.Add(-b.policy.Size)
	var out []domain.Event
	keep := b.events[:0]
	for _, e := range b.events {
		if e.EventTime.Before(cutoff) {
			out = append(out, e)
		} else {
			keep = append(keep, e)
		}
	}
	b.events = keep
	return out
}
func (b *Buffer) Snapshot() []domain.Event {
	b.mu.Lock()
	events := b.events
	b.mu.Unlock()
	return events
}
func MergeSessions(events []domain.Event, gap time.Duration) [][]domain.Event {
	sort.Slice(events, func(i, j int) bool { return events[i].EventTime.Before(events[j].EventTime) })
	groups := [][]domain.Event{}
	for _, e := range events {
		if len(groups) == 0 || e.EventTime.Sub(groups[len(groups)-1][len(groups[len(groups)-1])-1].EventTime) > gap {
			groups = append(groups, []domain.Event{e})
		} else {
			groups[len(groups)-1] = append(groups[len(groups)-1], e)
		}
	}
	return groups
}
