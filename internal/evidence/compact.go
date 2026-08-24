package evidence

type Item struct {
	EventID string
	Kind    string
	Value   string
}

// Compact returns a deduplicated copy of input.
//
// It MUST NOT mutate the caller's slice or alias its backing array: the same
// evidence slice is often handed to a caller (e.g. the explain tree response),
// cached, and recorded in history before compaction runs. Writing back into
// input's array (input[:0] + append) would silently corrupt every other
// holder of that array. Always allocate a fresh slice.
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

// clone returns a deep copy of items so callers can mutate the result without
// affecting any other slice that shares the same backing array.
func clone(items []Item) []Item {
	if items == nil {
		return nil
	}
	out := make([]Item, len(items))
	copy(out, items)
	return out
}
