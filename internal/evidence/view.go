package evidence

type View struct {
	MatchID string
	Items   []Item
	Kinds   map[string]int
}

func Build(matchID string, items []Item) View {
	view := View{MatchID: matchID, Items: items, Kinds: make(map[string]int)}
	for _, item := range view.Items {
		view.Kinds[item.Kind]++
	}
	return view
}
