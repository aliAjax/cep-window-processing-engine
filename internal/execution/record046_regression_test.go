package execution

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestExecutionMiddlewareKeepsParentCancellation(t *testing.T) {
	parent, cancelParent := context.WithCancel(context.Background())
	request, cancelRequest := (Middleware{Timeout: time.Hour}).Bind(parent, "pattern-a", nil)
	defer cancelRequest()

	cancelParent()
	select {
	case <-request.Context().Done():
		if !errors.Is(request.Context().Err(), context.Canceled) {
			t.Fatalf("request context error = %v, want context canceled", request.Context().Err())
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("request context did not observe parent cancellation")
	}
}

func TestExecutionRequestKeepsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	request := NewRequest(ctx, "pattern-b", nil)
	if !errors.Is(request.Context().Err(), context.Canceled) {
		t.Fatalf("request context error = %v, want context canceled", request.Context().Err())
	}
}

func TestExecutionStoreUsesCallerContext(t *testing.T) {
	store := NewStore()
	store.Register("pattern-c", func(ctx context.Context, _ Request) (bool, error) {
		return errors.Is(ctx.Err(), context.Canceled), nil
	})
	request := NewRequest(context.Background(), "pattern-c", nil)
	caller, cancel := context.WithCancel(context.Background())
	cancel()

	observed, err := store.Evaluate(caller, request)
	if err != nil {
		t.Fatalf("evaluate returned error: %v", err)
	}
	if !observed {
		t.Fatal("evaluator did not receive the canceled caller context")
	}
}

func TestExecutionRetryStopsAfterCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	started := time.Now()
	err := (ExecutionWorker{Attempts: 2, Delay: 250 * time.Millisecond}).Retry(ctx, func(context.Context) error {
		calls++
		if calls == 1 {
			cancel()
			return errors.New("retryable")
		}
		return nil
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("retry error = %v, want context canceled", err)
	}
	if calls != 1 {
		t.Fatalf("operation calls = %d, want 1", calls)
	}
	if elapsed := time.Since(started); elapsed >= 100*time.Millisecond {
		t.Fatalf("retry returned after %s, want prompt cancellation", elapsed)
	}
}
