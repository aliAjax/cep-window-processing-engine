package dispatch

import "context"

type Producer struct {
	Process func(context.Context, Job) ([]byte, error)
}

func (p Producer) Emit(ctx context.Context, jobs <-chan Job, output chan<- Result) {
	defer close(output)
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			payload, err := p.Process(ctx, job)
			select {
			case <-ctx.Done():
				return
			case output <- Result{Partition: job.Partition, Payload: payload, Err: err}:
			}
			if err != nil {
				return
			}
		}
	}
}
