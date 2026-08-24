package evidence

import "sync"

type Cache struct {
	mu      sync.RWMutex
	entries map[string][]Item
}

func NewCache() *Cache {
	return &Cache{entries: make(map[string][]Item)}
}

func (c *Cache) Store(matchID string, items []Item) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Store a private copy so later mutations to the caller's slice cannot
	// reach into the cached entry.
	c.entries[matchID] = clone(items)
}

func (c *Cache) Load(matchID string) ([]Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items, ok := c.entries[matchID]
	if !ok {
		return nil, false
	}
	// Return a private copy so the caller cannot mutate the cached entry
	// (which would also rewrite previously-loaded versions).
	return clone(items), true
}
