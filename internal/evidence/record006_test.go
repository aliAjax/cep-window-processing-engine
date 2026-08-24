package evidence

import (
	"reflect"
	"testing"
)

func sampleItems() []Item {
	return []Item{{EventID: "e1", Kind: "input"}, {EventID: "e1", Kind: "input"}, {EventID: "e2", Kind: "match"}}
}

func TestCompactDoesNotMutateInput(t *testing.T) {
	input := sampleItems()
	want := append([]Item(nil), input...)
	_ = Compact(input)
	if !reflect.DeepEqual(input, want) {
		t.Fatalf("compact mutated input: got %#v want %#v", input, want)
	}
}

func TestHistoryKeepsPriorEvidence(t *testing.T) {
	items := sampleItems()
	history := NewHistory()
	history.Append("m1", items)
	items[0].EventID = "changed"
	versions := history.Versions("m1")
	if versions[0][0].EventID != "e1" {
		t.Fatalf("history changed through caller alias: %#v", versions)
	}
}

func TestCacheOwnsSnapshot(t *testing.T) {
	items := sampleItems()
	cache := NewCache()
	cache.Store("m1", items)
	items[0].Value = "mutated"
	first, _ := cache.Load("m1")
	first[0].Value = "reader-mutated"
	second, _ := cache.Load("m1")
	if second[0].Value != "" {
		t.Fatalf("cache exposed aliased snapshot: %#v", second)
	}
}

func TestViewSurvivesLateAppend(t *testing.T) {
	items := make([]Item, 2, 4)
	items[0] = Item{EventID: "e1", Kind: "input"}
	items[1] = Item{EventID: "e2", Kind: "match"}
	view := Build("m1", items)
	items[0].EventID = "late-change"
	items = append(items, Item{EventID: "e3", Kind: "late"})
	if view.Items[0].EventID != "e1" || len(view.Items) != 2 {
		t.Fatalf("view changed after late append: %#v", view.Items)
	}
}
