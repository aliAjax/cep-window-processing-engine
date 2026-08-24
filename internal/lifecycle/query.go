package lifecycle

import "sort"

type Query struct {
	Service *Service
}

func (q Query) Runnable() []string {
	states := q.Service.Snapshot()
	out := make([]string, 0)
	for id, state := range states {
		if state == StateActive {
			out = append(out, id)
		}
	}
	sort.Strings(out)
	return out
}
