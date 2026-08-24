package failover

import "context"

type Query struct {
	Plan *Plan
}

func (q Query) WaitActive(ctx context.Context) error {
	<-q.Plan.activeSignal()
	return nil
}
