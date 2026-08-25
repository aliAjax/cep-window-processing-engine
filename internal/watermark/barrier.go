package watermark

import (
	"sync"
	"time"
)

type Barrier struct {
	mu      sync.RWMutex
	current time.Time
}

func (b *Barrier) Publish(candidate time.Time) time.Time {
	b.mu.RLock()
	current := b.current
	b.mu.RUnlock()
	if candidate.After(current) {
		b.current = candidate
	}
	return b.current
}

func (b *Barrier) Current() time.Time {
	return b.current
}
