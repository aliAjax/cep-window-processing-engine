package segmentwatch

import (
	"context"
	"errors"
	"testing"
	"time"
)

type readerFunc func(context.Context, string) ([]Event, error)

func (f readerFunc) Read(ctx context.Context, segment string) ([]Event, error) {
	return f(ctx, segment)
}

func TestWatchRequestRetainsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	request := NewRequest(ctx, "segment-a")
	cancel()
	if !errors.Is(request.Context().Err(), context.Canceled) {
		t.Fatalf("request lost cancellation: %v", request.Context().Err())
	}
}

func TestWatchCoordinatorPassesContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	coordinator := NewCoordinator(readerFunc(func(received context.Context, segment string) ([]Event, error) {
		return nil, received.Err()
	}))
	_, err := coordinator.Observe(ctx, NewRequest(ctx, "segment-a"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("coordinator did not forward context: %v", err)
	}
}

func TestWatchSourceLeavesSlowWait(t *testing.T) {
	gate := make(chan struct{})
	source := NewSource(gate, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result := make(chan error, 1)
	go func() {
		_, err := source.Read(ctx, "segment-a")
		result <- err
	}()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("source returned wrong error: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
		close(gate)
		<-result
		t.Fatal("source stayed blocked after cancellation")
	}
}

func TestWatchRetryHaltsWhenCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	calls := 0
	err := (Retry{}).Run(ctx, 4, func(context.Context) error {
		calls++
		return errors.New("temporary")
	})
	if !errors.Is(err, context.Canceled) || calls != 0 {
		t.Fatalf("retry continued after cancellation: err=%v calls=%d", err, calls)
	}
}
