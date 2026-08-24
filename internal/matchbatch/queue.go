package matchbatch

import (
	"context"
	"errors"
)

var ErrQueueClosed = errors.New("match batch queue closed")

type Request struct {
	ID      string
	Payload []byte
}

type Queue struct {
	items  chan Request
	closed chan struct{}
}

func NewQueue(capacity int) *Queue {
	return &Queue{
		items:  make(chan Request, capacity),
		closed: make(chan struct{}),
	}
}

func (q *Queue) Items() <-chan Request {
	return q.items
}

func (q *Queue) Submit(ctx context.Context, request Request) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-q.closed:
		return ErrQueueClosed
	case q.items <- request:
		return nil
	}
}

func (q *Queue) Close() {
	select {
	case <-q.closed:
	default:
		close(q.closed)
	}
}
