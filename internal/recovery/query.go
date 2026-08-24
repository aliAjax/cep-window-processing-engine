package recovery

import (
	"context"
	"fmt"
)

type Query struct {
	Coordinator *Coordinator
}

func (q Query) Current(epochID string) (State, bool) {
	if q.Coordinator == nil {
		return "", false
	}
	epoch, ok := q.Coordinator.Get(epochID)
	if !ok {
		return "", false
	}
	return epoch.State, true
}

func (q Query) Wait(ctx context.Context, epochID string) (State, error) {
	if q.Coordinator == nil {
		return "", fmt.Errorf("recovery coordinator is unavailable")
	}
	done, err := q.Coordinator.Subscribe(epochID)
	if err != nil {
		return "", err
	}
	<-done
	state, ok := q.Current(epochID)
	if !ok {
		return "", fmt.Errorf("recovery epoch not found")
	}
	return state, nil
}
