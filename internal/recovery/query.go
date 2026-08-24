package recovery

type Query struct {
	Coordinator *Coordinator
}

func (q Query) Current(epochID string) (State, bool) {
	if q.Coordinator == nil {
		return "", false
	}
	epoch, ok := q.Coordinator.Get(epochID)
	if !ok {
		return "", false
	}
	return epoch.State, true
}
