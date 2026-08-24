package replayload

import "context"

type ReplayRequest struct {
	ctx     context.Context
	Segment string
}

func NewRequest(ctx context.Context, segment string) ReplayRequest {
	if ctx == nil {
		ctx = context.Background()
	}
	return ReplayRequest{ctx: ctx, Segment: segment}
}

func (r ReplayRequest) Context() context.Context {
	return context.Background()
}
