package partition

import (
	"sync"
)

type Registry struct {
	mu          sync.RWMutex
	assignments map[int]string
	generation  uint64
}

func NewRegistry() *Registry {
	return &Registry{assignments: make(map[int]string)}
}

func (r *Registry) Restore(snapshot Snapshot) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.assignments = make(map[int]string, len(snapshot.Assignments))
	for partition, owner := range snapshot.Assignments {
		r.assignments[partition] = owner
	}
	r.generation = snapshot.Generation
}

func (r *Registry) Assign(partition int, owner string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.assignments == nil {
		r.assignments = make(map[int]string)
	}
	r.assignments[partition] = owner
}

func (r *Registry) Snapshot() Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := Snapshot{Assignments: make(map[int]string, len(r.assignments)), Generation: r.generation}
	for partition, owner := range r.assignments {
		out.Assignments[partition] = owner
	}
	return out
}
