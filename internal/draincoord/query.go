package draincoord

import "context"

type DrainQuery struct {
	Buffer *Buffer
}

func (q DrainQuery) Wait(ctx context.Context) error {
	<-q.Buffer.closedSignal()
	return nil
}
