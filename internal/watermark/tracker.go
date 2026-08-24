package watermark

import (
	"fmt"
	"time"
)

type Tracker struct {
	partitions map[int]time.Time
}

func NewTracker() *Tracker {
	return &Tracker{partitions: make(map[int]time.Time)}
}

func (t *Tracker) Advance(partition int, next time.Time) error {
	if current := t.partitions[partition]; !current.IsZero() && next.Before(current) {
		return fmt.Errorf("partition %d watermark regressed", partition)
	}
	t.partitions[partition] = next
	return nil
}

func (t *Tracker) Snapshot() map[int]time.Time {
	return t.partitions
}
