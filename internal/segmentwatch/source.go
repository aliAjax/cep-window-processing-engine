package segmentwatch

import (
	"context"
	"sync/atomic"
)

type Event struct {
	ID string
}

type Source struct {
	gate  <-chan struct{}
	data  map[string][]Event
	reads atomic.Int64
}

func NewSource(gate <-chan struct{}, data map[string][]Event) *Source {
	return &Source{gate: gate, data: data}
}

func (s *Source) Read(ctx context.Context, segment string) ([]Event, error) {
	s.reads.Add(1)
	if s.gate != nil {
		<-s.gate
	}
	return append([]Event(nil), s.data[segment]...), nil
}

func (s *Source) Reads() int64 {
	return s.reads.Load()
}
