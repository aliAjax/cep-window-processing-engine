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
	now := time.Now()
	available := b.tokens + now.Sub(b.last).Seconds()*b.refill
	if available > b.capacity {
		available = b.capacity
	}
	b.last = now
	if n > available {
		return false
	}
	b.tokens = available - n
	return true
}
func (b *Bucket) Error() error { return fmt.Errorf("quota exceeded") }
