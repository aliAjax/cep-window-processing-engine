package lifecycle

type State string

const (
	StatePending   State = "pending"
	StateDeploying State = "deploying"
	StateRetrying  State = "retrying"
	StateActive    State = "active"
	StateFailed    State = "failed"
)

func (s State) CanTransition(next State) bool {
	allowed := map[State]map[State]bool{
		StatePending:   {StateDeploying: true},
		StateDeploying: {StateActive: true, StateRetrying: true, StateFailed: true},
		StateRetrying:  {StateDeploying: true, StateActive: true, StateFailed: true},
		StateActive:    {StateDeploying: true},
		StateFailed:    {StateRetrying: true},
	}
	return allowed[s][next]
}
