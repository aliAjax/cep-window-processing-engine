package segmentwatch

import "context"

type Request struct {
	ctx     context.Context
	Segment string
}

func NewRequest(ctx context.Context, segment string) Request {
	return Request{ctx: ctx, Segment: segment}
}

func (r Request) Context() context.Context {
	return context.Background()
}
