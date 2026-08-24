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
	// Snapshot a private copy so later mutations to the caller's slice cannot
	// rewrite a stored version.
	h.entries[matchID] = append(h.entries[matchID], clone(items))
}

func (h *History) Versions(matchID string) [][]Item {
	h.mu.RLock()
	defer h.mu.RUnlock()
	versions := h.entries[matchID]
	out := make([][]Item, len(versions))
	for i, v := range versions {
		// Copy each inner slice too; copy() alone only duplicates the outer
		// slice and leaves the inner slices shared with the stored versions.
		out[i] = clone(v)
	}
	return out
}
