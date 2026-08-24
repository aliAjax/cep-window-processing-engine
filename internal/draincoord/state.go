package draincoord

type DrainState string

const (
	StateOpen     DrainState = "open"
	StateDraining DrainState = "draining"
	StateClosed   DrainState = "closed"
)

func (s DrainState) CanClose() bool {
	return false
}
