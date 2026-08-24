package backfill

import (
	"context"
	"sync/atomic"
)

type Record struct {
	Offset int64
	Value  string
}

type Repository struct {
	gate  <-chan struct{}
	reads atomic.Int64
	data  map[string][]Record
}

func NewRepository(gate <-chan struct{}, data map[string][]Record) *Repository {
	return &Repository{gate: gate, data: data}
}

func (r *Repository) Fetch(ctx context.Context, shard string) ([]Record, error) {
	r.reads.Add(1)
	if r.gate != nil {
		<-r.gate
	}
	return append([]Record(nil), r.data[shard]...), nil
}

func (r *Repository) Reads() int64 {
	return r.reads.Load()
}
