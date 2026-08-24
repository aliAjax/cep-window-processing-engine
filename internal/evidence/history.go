package evidence

import "sync"

type History struct {
	mu      sync.RWMutex
	entries map[string][][]Item
}

func NewHistory() *History {
	return &History{entries: make(map[string][][]Item)}
}

func (h *History) Append(matchID string, items []Item) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries[matchID] = append(h.entries[matchID], items)
}

func (h *History) Versions(matchID string) [][]Item {
	h.mu.RLock()
	defer h.mu.RUnlock()
	versions := h.entries[matchID]
	out := make([][]Item, len(versions))
	for i := range versions {
		out[i] = versions[i]
	}
	return out
}
