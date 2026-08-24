package checkpointrollout

type Worker struct{}

func (Worker) Run(plan Plan, promote func() error) (Plan, error) {
	if err := promote(); err != nil {
		plan.State = StateActive
		return plan, nil
	}

	state, err := Transition(plan.State, StateActive)
	if err != nil {
		return plan, err
	}
	plan.State = state
	return plan, nil
}
