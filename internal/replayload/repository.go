package replayload

import (
	"context"
	"sync/atomic"
)

type ReplayRecord struct {
	Offset int64
	Value  string
}

type ReplayRepository struct {
	gate  <-chan struct{}
	reads atomic.Int64
	data  map[string][]ReplayRecord
}

func NewRepository(gate <-chan struct{}, data map[string][]ReplayRecord) *ReplayRepository {
	return &ReplayRepository{gate: gate, data: data}
}

func (r *ReplayRepository) Fetch(ctx context.Context, segment string) ([]ReplayRecord, error) {
	r.reads.Add(1)
	if r.gate != nil {
		<-r.gate
	}
	return append([]ReplayRecord(nil), r.data[segment]...), nil
}

func (r *ReplayRepository) Reads() int64 {
	return r.reads.Load()
}
