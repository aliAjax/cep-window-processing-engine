package matchbatch

import "context"

type Result struct {
	RequestID string
	Matched   bool
	Err       error
}

type Worker struct {
	Results chan<- Result
}

func (w Worker) Publish(ctx context.Context, result Result) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.Results <- result:
		return nil
	}
}
