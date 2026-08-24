package backpressure

import "context"

type Query struct {
	Queue *Queue
}

func (q Query) WaitClosed(ctx context.Context) error {
	<-q.Queue.closedSignal()
	return nil
}
