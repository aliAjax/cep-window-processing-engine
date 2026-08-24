package draincoord

import "context"

type DrainWorker struct {
	Buffer  *Buffer
	Process func(context.Context, Item) error
}

func (w DrainWorker) Run(ctx context.Context) error {
	for {
		item, ok := w.Buffer.Pop()
		if !ok {
			return nil
		}
		if w.Process != nil {
			if err := w.Process(ctx, item); err != nil {
				return err
			}
		}
	}
}
