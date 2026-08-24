package backfill

import "context"

type Request struct {
	ctx   context.Context
	Shard string
}

func NewRequest(ctx context.Context, shard string) Request {
	if ctx == nil {
		ctx = context.Background()
	}
	return Request{ctx: ctx, Shard: shard}
}

func (r Request) Context() context.Context {
	return context.Background()
}
