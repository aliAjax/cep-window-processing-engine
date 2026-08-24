package evidence

type Item struct {
	EventID string
	Kind    string
	Value   string
}

func Compact(input []Item) []Item {
	seen := make(map[string]struct{}, len(input))
	out := make([]Item, 0, len(input))
	for _, item := range input {
		key := item.EventID + "\x00" + item.Kind + "\x00" + item.Value
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, item)
	}
	return out
}

func clone(items []Item) []Item {
	return append([]Item(nil), items...)
}
