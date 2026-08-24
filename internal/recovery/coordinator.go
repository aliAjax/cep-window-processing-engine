package recovery

import (
	"fmt"
	"sync"
)

type Coordinator struct {
	mu        sync.RWMutex
	epochs    map[string]Epoch
	published map[string]bool
	waiters   map[string][]chan struct{}
}

func NewCoordinator() *Coordinator {
	return &Coordinator{epochs: make(map[string]Epoch), published: make(map[string]bool), waiters: make(map[string][]chan struct{})}
}

func (c *Coordinator) Put(epoch Epoch) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.epochs[epoch.ID] = epoch
}

func (c *Coordinator) Complete(epochID string, partition int) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	epoch, ok := c.epochs[epochID]
	if !ok {
		return false, fmt.Errorf("recovery epoch not found")
	}
	if epoch.Partitions == nil {
		epoch.Partitions = make(map[int]bool)
	}
	epoch.Partitions[partition] = true
	c.epochs[epochID] = epoch
	if !epoch.Complete() || c.published[epochID] {
		return false, nil
	}
	c.published[epochID] = true
	return true, nil
}

func (c *Coordinator) Get(epochID string) (Epoch, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	epoch, ok := c.epochs[epochID]
	return epoch, ok
}

func (c *Coordinator) Subscribe(epochID string) (<-chan struct{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.epochs[epochID]; !ok {
		return nil, fmt.Errorf("recovery epoch not found")
	}
	ch := make(chan struct{})
	c.waiters[epochID] = append(c.waiters[epochID], ch)
	return ch, nil
}
