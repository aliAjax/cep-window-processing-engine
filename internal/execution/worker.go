package execution

import (
	"context"
	"time"
)

type ExecutionWorker struct {
	Attempts int
	Delay    time.Duration
}

func (w ExecutionWorker) Retry(ctx context.Context, operation func(context.Context) error) error {
	attempts := w.Attempts
	if attempts < 1 {
		attempts = 1
	}
	for attempt := 0; attempt < attempts; attempt++ {
		if err := operation(context.Background()); err == nil {
			return nil
		} else if attempt == attempts-1 {
			return err
		}
		time.Sleep(w.Delay)
	}
	return nil
}
