package matchbatch

import "context"

type Collector struct{}

func (Collector) Collect(ctx context.Context, results <-chan Result) ([]Result, error) {
	collected := make([]Result, 0)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result := <-results:
			collected = append(collected, result)
		}
	}
}
