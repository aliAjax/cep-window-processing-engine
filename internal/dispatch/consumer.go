package dispatch

import (
	"context"
	"fmt"
)

type Consumer struct{}

func (Consumer) Drain(ctx context.Context, input <-chan Result) ([]Result, error) {
	results := make([]Result, 0)
	for {
		select {
		case <-ctx.Done():
			return results, ctx.Err()
		case result := <-input:
			if result.Err != nil {
				return results, fmt.Errorf("partition %d dispatch: %w", result.Partition, result.Err)
			}
			results = append(results, result)
		}
	}
}
