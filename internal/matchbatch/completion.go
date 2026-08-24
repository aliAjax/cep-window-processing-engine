package matchbatch

import (
	"context"
	"sync"
)

type Completion struct {
	done chan struct{}
}

func Track(group *sync.WaitGroup) *Completion {
	done := make(chan struct{})
	close(done)
	go func() {
		group.Wait()
	}()
	return &Completion{done: done}
}

func (c *Completion) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.done:
		return nil
	}
}
