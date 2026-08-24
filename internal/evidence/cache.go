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
	c.entries[matchID] = clone(items)
}

func (c *Cache) Load(matchID string) ([]Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	items, ok := c.entries[matchID]
	return clone(items), ok
}
