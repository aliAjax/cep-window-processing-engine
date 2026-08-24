package checkpointrollout

import "fmt"

type State string

const (
	StateQueued    State = "queued"
	StatePreparing State = "preparing"
	StateActive    State = "active"
	StateFailed    State = "failed"
)

func CanTransition(current, next State) bool {
	allowed := map[State]map[State]bool{
		StateQueued:    {StatePreparing: true},
		StatePreparing: {StateActive: true, StateFailed: true},
	}
	return allowed[current][next]
}

func Transition(current, next State) (State, error) {
	if !CanTransition(current, next) {
		return current, fmt.Errorf("invalid checkpoint rollout transition %s -> %s", current, next)
	}
	return next, nil
}
