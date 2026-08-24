package lifecycle

import "context"

type LifecycleWorker struct {
	Service *Service
}

func (w LifecycleWorker) Retry(ctx context.Context, ruleID string, deploy func(context.Context) error) error {
	if err := w.Service.Transition(ruleID, StateRetrying); err != nil {
		return err
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := deploy(ctx); err != nil {
		_ = w.Service.Transition(ruleID, StateFailed)
		return err
	}
	if w.Service.State(ruleID) == StateRetrying {
		return nil
	}
	return nil
}
