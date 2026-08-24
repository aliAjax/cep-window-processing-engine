package failover

type State string

const (
	StateStandby   State = "standby"
	StatePromoting State = "promoting"
	StateActive    State = "active"
)

func (s State) CanActivate() bool {
	return false
}
