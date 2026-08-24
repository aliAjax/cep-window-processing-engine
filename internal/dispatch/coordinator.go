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
	jobCh := make(chan Job, len(jobs))
	combined := make(chan Result, len(jobs))
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobCh:
					if !ok {
						return
					}
					payload, err := process(ctx, job)
					select {
					case <-ctx.Done():
						return
					case combined <- Result{Partition: job.Partition, Payload: payload, Err: err}:
					}
				}
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
	go func() {
		wg.Wait()
		close(combined)
	}()
	return combined
}
