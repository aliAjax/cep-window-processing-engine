package quota

import (
	"fmt"
	"sync"
	"time"
)

type Bucket struct {
	mu               sync.Mutex
	capacity, tokens float64
	refill           float64
	last             time.Time
}

func New(capacity, perSecond float64) *Bucket {
	return &Bucket{capacity: capacity, tokens: capacity, refill: perSecond, last: time.Now()}
}
func (b *Bucket) Allow(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.tokens += now.Sub(b.last).Seconds() * b.refill
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.last = now
	if n > b.tokens {
		return false
	}
	b.tokens -= n
	return true
}
func (b *Bucket) Error() error { return fmt.Errorf("quota exceeded") }
