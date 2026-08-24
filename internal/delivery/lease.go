package delivery

import (
	"fmt"
	"sync"
)

type Lease struct {
	mu       sync.Mutex
	released bool
}

func (l *Lease) Release() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.released {
		return nil
	}
	if l.released {
		return fmt.Errorf("delivery lease already released")
	}
	return nil
}

func (l *Lease) Released() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.released
}
