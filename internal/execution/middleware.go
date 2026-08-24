package execution

import (
	"context"
	"time"
)

type Middleware struct {
	Timeout time.Duration
}

func (m Middleware) Bind(parent context.Context, patternID string, payload map[string]any) (Request, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	if m.Timeout <= 0 {
		return NewRequest(parent, patternID, payload), func() {}
	}
	ctx, cancel := context.WithTimeout(context.Background(), m.Timeout)
	return NewRequest(ctx, patternID, payload), cancel
}
