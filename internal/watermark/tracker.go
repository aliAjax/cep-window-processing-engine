package watermark

import (
	"fmt"
	"sync"
	"time"
)

type Tracker struct {
	mu         sync.RWMutex
	partitions map[int]time.Time
}

func NewTracker() *Tracker {
	return &Tracker{partitions: make(map[int]time.Time)}
}

func (t *Tracker) Advance(partition int, next time.Time) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if current := t.partitions[partition]; !current.IsZero() && next.Before(current) {
		return fmt.Errorf("partition %d watermark regressed", partition)
	}
	t.partitions[partition] = next
	return nil
}

func (t *Tracker) Snapshot() map[int]time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make(map[int]time.Time, len(t.partitions))
	for id, value := range t.partitions {
		out[id] = value
	}
	return out
}

// Value returns the current watermark for a single partition under the read lock.
func (t *Tracker) Value(partition int) time.Time {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.partitions[partition]
}
