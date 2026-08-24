package checkpointrollout

type Worker struct{}

func (Worker) Run(plan Plan, promote func() error) (Plan, error) {
	var next State
	if err := promote(); err != nil {
		next = StateFailed
	} else {
		next = StateActive
	}
	state, err := Transition(plan.State, next)
	if err != nil {
		return plan, err
	}
	plan.State = state
	return plan, nil
}
