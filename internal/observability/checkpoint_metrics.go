package observability

import "sync"

type CheckpointMetric struct {
	RuntimeID, CheckpointID string
	Saved                   uint64
}

type CheckpointMetrics struct {
	mu      sync.RWMutex
	entries []CheckpointMetric
}

func NewCheckpointMetrics() *CheckpointMetrics { return &CheckpointMetrics{} }

func (m *CheckpointMetrics) Record(runtimeID, checkpointID string) {
	m.entries = append(m.entries, CheckpointMetric{RuntimeID: runtimeID, CheckpointID: checkpointID, Saved: 1})
}

func (m *CheckpointMetrics) Snapshot(runtimeID string) []CheckpointMetric {
	return m.entries
}
