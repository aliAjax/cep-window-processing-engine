package subscription

import "sync"

type Registry struct {
	mu     sync.RWMutex
	groups map[string][]string
}

func (r *Registry) Add(group, subscriber string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.groups == nil {
		r.groups = map[string][]string{}
	}
	r.groups[group] = append(r.groups[group], subscriber)
}

func (r *Registry) Members(group string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]string(nil), r.groups[group]...)
}
