package backpressure

import "sync"

type Event struct {
	ID string
}

type Queue struct {
	mu     sync.Mutex
	state  State
	items  []Event
	closed chan struct{}
	once   sync.Once
}

func NewQueue(events ...Event) *Queue {
	return &Queue{
		state:  StateOpen,
		items:  append([]Event(nil), events...),
		closed: make(chan struct{}),
	}
}

func (q *Queue) RequestClose() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.state = StateDraining
}

func (q *Queue) Pop() (Event, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return Event{}, false
	}
	event := q.items[0]
	q.items = q.items[1:]
	return event, true
}

func (q *Queue) State() State {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.state
}

func (q *Queue) markClosed() {
	q.mu.Lock()
	q.state = StateClosed
	q.mu.Unlock()
	q.once.Do(func() { close(q.closed) })
}

func (q *Queue) closedSignal() <-chan struct{} {
	return q.closed
}
