package draincoord

import "sync"

type Item struct {
	ID string
}

type Buffer struct {
	mu     sync.Mutex
	state  DrainState
	items  []Item
	closed chan struct{}
	once   sync.Once
}

func NewBuffer(items ...Item) *Buffer {
	return &Buffer{state: StateOpen, items: append([]Item(nil), items...), closed: make(chan struct{})}
}

func (b *Buffer) RequestStop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state = StateDraining
}

func (b *Buffer) Pop() (Item, bool) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.items) == 0 {
		return Item{}, false
	}
	item := b.items[0]
	b.items = b.items[1:]
	return item, true
}

func (b *Buffer) State() DrainState {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

func (b *Buffer) markClosed() {
	b.mu.Lock()
	b.state = StateClosed
	b.mu.Unlock()
	b.once.Do(func() { close(b.closed) })
}

func (b *Buffer) closedSignal() <-chan struct{} {
	return b.closed
}
