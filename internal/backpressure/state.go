package backpressure

type State string

const (
	StateOpen     State = "open"
	StateDraining State = "draining"
	StateClosed   State = "closed"
)

func (s State) CanClose() bool {
	return false
}
