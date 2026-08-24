package subscription

import "sync"

type Notification struct {
	ID      string
	Payload []byte
}

type Fanout struct {
	mu     sync.RWMutex
	queues map[string][]Notification
}

func (f *Fanout) Enqueue(group string, notification Notification) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queues[group] = append(f.queues[group], notification)
}

func (f *Fanout) Pending(group string) []Notification {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return append([]Notification(nil), f.queues[group]...)
}
