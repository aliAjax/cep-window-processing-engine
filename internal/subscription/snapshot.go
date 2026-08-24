package subscription

import "sync"

type Snapshot struct {
	Version int64
	Members []string
}

type Snapshotter struct {
	mu        sync.RWMutex
	snapshots map[string]Snapshot
}

func (s *Snapshotter) Capture(group string, snapshot Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.snapshots[group] = snapshot
}

func (s *Snapshotter) Latest(group string) (Snapshot, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	snapshot, ok := s.snapshots[group]
	return snapshot, ok
}
