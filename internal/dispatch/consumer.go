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
		case result, ok := <-input:
			if !ok {
				return results, nil
			}
			if result.Err != nil {
				return nil, fmt.Errorf("partition %d dispatch failed", result.Partition)
			}
			results = append(results, result)
		}
	}
}
