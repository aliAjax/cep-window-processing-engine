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
	b.current = candidate
	return b.current
}

func (b *Barrier) Current() time.Time {
	return b.current
}
