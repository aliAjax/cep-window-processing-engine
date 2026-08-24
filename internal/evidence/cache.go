package evidence

import "sync"

type Cache struct {
	mu      sync.RWMutex
	entries map[string][]Item
}

func NewCache() *Cache {
	return &Cache{entries: make(map[string][]Item)}
}

// Store keeps a private copy of items. Callers retain and mutate their own
// slice after this returns (compaction, late-arriving events, ...); storing a
// reference to the caller's backing array would let those later mutations
// leak into the cache and corrupt the cached evidence.
func (c *Cache) Store(matchID string, items []Item) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[matchID] = clone(items)
}

// Load returns a private copy of the cached items so the caller cannot mutate
// the cache's internal slice by mutating the returned slice.
func (c *Cache) Load(matchID string) ([]Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items, ok := c.entries[matchID]
	if !ok {
		return nil, false
	}
	return clone(items), true
}
