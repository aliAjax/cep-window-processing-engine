package dispatch

import (
	"context"
	"sync"
)

type Coordinator struct {
	Workers int
}

func (c Coordinator) Run(ctx context.Context, jobs []Job, process func(context.Context, Job) ([]byte, error)) <-chan Result {
	workers := c.Workers
	if workers < 1 {
		workers = 1
	}
	jobCh := make(chan Job)
	combined := make(chan Result, len(jobs))
	var wg sync.WaitGroup
	go func() {
		wg.Wait()
		close(combined)
	}()
	for i := 0; i < workers; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for job := range jobCh {
				payload, err := process(ctx, job)
				combined <- Result{Partition: job.Partition, Payload: payload, Err: err}
			}
		}()
	}
	go func() {
		defer close(jobCh)
		for _, job := range jobs {
			select {
			case <-ctx.Done():
				return
			case jobCh <- job:
			}
		}
	}()
	return combined
}
