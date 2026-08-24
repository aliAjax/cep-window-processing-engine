package evidence

import "sync"

type History struct {
	mu      sync.RWMutex
	entries map[string][][]Item
}

func NewHistory() *History {
	return &History{entries: make(map[string][][]Item)}
}

// Append records a snapshot of items as a new version. Callers keep mutating
// their own slice (compaction, late-arriving events, ...) after this returns;
// storing a reference to the caller's backing array would let those mutations
// retroactively rewrite previously-recorded versions.
func (h *History) Append(matchID string, items []Item) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries[matchID] = append(h.entries[matchID], clone(items))
}

// Versions returns private copies of every recorded version so the caller
// cannot mutate the history's internal slices by mutating the returned slices.
func (h *History) Versions(matchID string) [][]Item {
	h.mu.RLock()
	defer h.mu.RUnlock()
	versions := h.entries[matchID]
	out := make([][]Item, len(versions))
	for i := range versions {
		out[i] = clone(versions[i])
	}
	return out
}
