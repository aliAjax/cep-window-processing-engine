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
	w.Results <- result
	return nil
}
