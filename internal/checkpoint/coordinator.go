package checkpoint

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/example/cep-window-engine/internal/domain"
	"github.com/example/cep-window-engine/internal/observability"
	"github.com/example/cep-window-engine/internal/repository"
)

type BatchResult struct {
	RuntimeID     string
	CheckpointIDs []string
}

type Coordinator struct {
	manager  *Manager
	registry *Registry
	jobs     *repository.CheckpointJobs
	metrics  *observability.CheckpointMetrics
	sequence atomic.Uint64
	mu       sync.RWMutex
	results  []BatchResult
}

func NewCoordinator(manager *Manager, registry *Registry, jobs *repository.CheckpointJobs, metrics *observability.CheckpointMetrics) *Coordinator {
	return &Coordinator{manager: manager, registry: registry, jobs: jobs, metrics: metrics}
}

func (c *Coordinator) SaveBatch(ctx context.Context, runtimeID string, stores []*domain.Store) ([]string, error) {
	jobID := fmt.Sprintf("%s-%d", runtimeID, c.sequence.Add(1))
	c.jobs.Begin(jobID, runtimeID)
	ids := make([]string, 0, len(stores))
	for _, store := range stores {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		id, err := c.manager.Save(store)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
		c.registry.Record(Record{RuntimeID: runtimeID, CheckpointID: id, CreatedAt: time.Now()})
		c.metrics.Record(runtimeID, id)
	}
	c.results = append(c.results, BatchResult{RuntimeID: runtimeID, CheckpointIDs: ids})
	c.jobs.Complete(jobID)
	return ids, nil
}

func (c *Coordinator) Results(runtimeID string) []BatchResult {
	return c.results
}
