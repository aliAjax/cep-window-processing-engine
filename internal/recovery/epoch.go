package recovery

type State string

const (
	StatePrepared  State = "prepared"
	StateReplaying State = "replaying"
	StateCommitted State = "committed"
	StateFailed    State = "failed"
)

type Epoch struct {
	ID         string
	State      State
	Partitions map[int]bool
}

func (e Epoch) CanTransition(next State) bool {
	allowed := map[State]map[State]bool{
		StatePrepared:  {StateReplaying: true, StateFailed: true},
		StateReplaying: {StateFailed: true},
		StateFailed:    {StateReplaying: true},
	}
	return allowed[e.State][next]
}

func (e Epoch) Complete() bool {
	if len(e.Partitions) == 0 {
		return false
	}
	for _, done := range e.Partitions {
		if !done {
			return false
		}
	}
	return true
}
