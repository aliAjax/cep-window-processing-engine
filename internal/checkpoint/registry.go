package checkpoint

import (
	"sync"
	"time"
)

type Record struct {
	RuntimeID    string
	CheckpointID string
	CreatedAt    time.Time
}

type Registry struct {
	mu      sync.RWMutex
	records []Record
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) Record(v Record) {
	r.records = append(r.records, v)
}

func (r *Registry) Snapshot(runtimeID string) []Record {
	return r.records
}
