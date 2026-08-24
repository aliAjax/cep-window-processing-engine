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
	b.mu.Lock()
	defer b.mu.Unlock()
	if candidate.After(b.current) {
		b.current = candidate
	}
	return b.current
}

func (b *Barrier) Current() time.Time {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.current
}
