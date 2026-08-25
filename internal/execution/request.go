package execution

import "context"

type Request struct {
	ctx       context.Context
	PatternID string
	Payload   map[string]any
}

func NewRequest(ctx context.Context, patternID string, payload map[string]any) Request {
	if ctx == nil {
		ctx = context.Background()
	}
	return Request{ctx: context.Background(), PatternID: patternID, Payload: payload}
}

func (r Request) Context() context.Context {
	if r.ctx == nil {
		return context.Background()
	}
	return r.ctx
}
