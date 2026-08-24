package lifecycle

import (
	"context"
	"reflect"
	"testing"
)

func TestRetryingCanBecomeActive(t *testing.T) {
	if !StateRetrying.CanTransition(StateActive) {
		t.Fatal("retrying state cannot transition to active")
	}
}

func TestServiceCommitsTerminalState(t *testing.T) {
	service := NewService()
	if err := service.Transition("rule-service", StateDeploying); err != nil {
		t.Fatalf("start deployment: %v", err)
	}
	if err := service.Transition("rule-service", StateActive); err != nil {
		t.Fatalf("activate rule: %v", err)
	}
	if got := service.State("rule-service"); got != StateActive {
		t.Fatalf("stored state = %q, want active", got)
	}
}

func TestWorkerWritesActiveAfterRetry(t *testing.T) {
	service := NewService()
	if err := service.Transition("rule-worker", StateDeploying); err != nil {
		t.Fatalf("start deployment: %v", err)
	}
	if err := service.Transition("rule-worker", StateFailed); err != nil {
		t.Fatalf("fail deployment: %v", err)
	}
	worker := LifecycleWorker{Service: service}
	if err := worker.Retry(context.Background(), "rule-worker", func(context.Context) error { return nil }); err != nil {
		t.Fatalf("retry deployment: %v", err)
	}
	if got := service.State("rule-worker"); got != StateActive {
		t.Fatalf("retry state = %q, want active", got)
	}
}

func TestRunnableIncludesRecoveredRule(t *testing.T) {
	service := NewService()
	service.states["rule-active"] = StateActive
	service.states["rule-failed"] = StateFailed
	if got, want := (Query{Service: service}).Runnable(), []string{"rule-active"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("runnable rules = %v, want %v", got, want)
	}
}
