package evidence

type View struct {
	MatchID string
	Items   []Item
	Kinds   map[string]int
}

// Build constructs an immutable View over a private copy of items. Callers
// keep mutating their own slice (compaction, late-arriving events, ...) after
// this returns; storing a reference to the caller's backing array would let
// those mutations change the evidence already exposed through the View.
func Build(matchID string, items []Item) View {
	cp := clone(items)
	view := View{MatchID: matchID, Items: cp, Kinds: make(map[string]int, len(cp))}
	for _, item := range cp {
		view.Kinds[item.Kind]++
	}
	return view
}
