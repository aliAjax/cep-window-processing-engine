package checkpointrollout

type Plan struct {
	State State
}

func (p Plan) Begin() (Plan, error) {
	next, err := Transition(p.State, StatePreparing)
	if err != nil {
		return p, err
	}
	p.State = next
	return p, nil
}
