package checkpointrollout

type Query struct{}

func (Query) Runnable(plans []Plan) []Plan {
	out := make([]Plan, 0, len(plans))
	for _, plan := range plans {
		if plan.State == StatePreparing {
			out = append(out, plan)
		}
	}
	return out
}
