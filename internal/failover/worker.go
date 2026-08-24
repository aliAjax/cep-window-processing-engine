package failover

import "context"

type Worker struct {
	Plan    *Plan
	Execute func(context.Context, Step) error
}

func (w Worker) Run(ctx context.Context) error {
	for {
		step, ok := w.Plan.Next()
		if !ok {
			return nil
		}
		if w.Execute != nil {
			if err := w.Execute(ctx, step); err != nil {
				return err
			}
		}
	}
}
