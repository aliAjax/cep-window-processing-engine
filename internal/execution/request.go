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
	return Request{ctx: ctx, PatternID: patternID, Payload: payload}
}

func (r Request) Context() context.Context {
	return context.Background()
}
