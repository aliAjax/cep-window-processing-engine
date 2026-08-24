package backpressure

import "context"

type Worker struct {
	Queue   *Queue
	Process func(context.Context, Event) error
}

func (w Worker) Drain(ctx context.Context) error {
	for {
		event, ok := w.Queue.Pop()
		if !ok {
			return nil
		}
		if w.Process != nil {
			if err := w.Process(ctx, event); err != nil {
				return err
			}
		}
	}
}
