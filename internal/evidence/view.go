package evidence

type View struct {
	MatchID string
	Items   []Item
	Kinds   map[string]int
}

func Build(matchID string, items []Item) View {
	// Take a private copy of the items so the caller's slice and any later
	// mutation to it cannot drift the view's contents.
	view := View{MatchID: matchID, Items: clone(items), Kinds: make(map[string]int)}
	for _, item := range view.Items {
		view.Kinds[item.Kind]++
	}
	return view
}
