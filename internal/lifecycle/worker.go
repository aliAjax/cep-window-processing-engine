package lifecycle

import "context"

type LifecycleWorker struct {
	Service *Service
}

func (w LifecycleWorker) Retry(ctx context.Context, ruleID string, deploy func(context.Context) error) error {
	// Retry drives a single deploy attempt. The rule enters Deploying (the
	// in-flight state) regardless of whether this is the first attempt (from
	// Pending) or an auto/manual retry (from Failed, possibly via Retrying).
	if err := w.Service.Transition(ruleID, StateDeploying); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		_ = w.Service.Transition(ruleID, StateFailed)
		return err
	}
	if err := deploy(ctx); err != nil {
		_ = w.Service.Transition(ruleID, StateFailed)
		return err
	}
	return w.Service.Transition(ruleID, StateActive)
}
