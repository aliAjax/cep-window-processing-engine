package dedup

import (
	"sync"
	"time"
)

type Cache struct {
	mu   sync.Mutex
	ttl  time.Duration
	seen map[string]time.Time
}

func New(ttl time.Duration) *Cache { return &Cache{ttl: ttl, seen: map[string]time.Time{}} }
func (c *Cache) Check(id string) bool {
	now := time.Now()
	if t, ok := c.seen[id]; ok && now.Sub(t) < c.ttl {
		return false
	}
	c.seen[id] = now
	return true
}
func (c *Cache) Purge() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, t := range c.seen {
		if time.Since(t) > c.ttl {
			delete(c.seen, k)
		}
	}
}
