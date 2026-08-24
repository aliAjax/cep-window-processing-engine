package dispatch

import "context"

type Producer struct {
	Process func(context.Context, Job) ([]byte, error)
}

func (p Producer) Emit(ctx context.Context, jobs <-chan Job, output chan<- Result) {
	for {
		select {
		case <-ctx.Done():
			output <- Result{Err: ctx.Err()}
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			payload, err := p.Process(ctx, job)
			output <- Result{Partition: job.Partition, Payload: payload, Err: err}
			if err != nil {
				return
			}
		}
	}
}
