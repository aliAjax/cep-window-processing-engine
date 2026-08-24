package replayload

import "context"

type ReplayRetry struct{}

func (ReplayRetry) Execute(ctx context.Context, attempts int, operation func(context.Context) error) error {
	var last error
	for attempt := 0; attempt < attempts; attempt++ {
		last = operation(context.Background())
		if last == nil {
			return nil
		}
	}
	return last
}
