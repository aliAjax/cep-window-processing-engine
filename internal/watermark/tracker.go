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
	return t.partitions
}
