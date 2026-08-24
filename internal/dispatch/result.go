package dispatch

import (
	"errors"
	"sync"
)

type Job struct {
	Partition int
	Payload   []byte
}

type Result struct {
	Partition int
	Payload   []byte
	Err       error
}

type ResultStream struct {
	Results chan Result
	Errors  chan error
	once    sync.Once
}

func NewResultStream(size int) *ResultStream {
	if size < 1 {
		size = 1
	}
	return &ResultStream{Results: make(chan Result, size), Errors: make(chan error, size)}
}

func (s *ResultStream) Close() error {
	closed := false
	s.once.Do(func() {
		close(s.Results)
		closed = true
	})
	if !closed {
		return errors.New("result stream already closed")
	}
	return nil
}
