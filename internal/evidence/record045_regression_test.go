package evidence

import (
	"reflect"
	"testing"
)

func evidenceItems() []Item {
	return []Item{
		{EventID: "event-1", Kind: "source", Value: "alpha"},
		{EventID: "event-2", Kind: "result", Value: "beta"},
	}
}

func TestEvidenceCacheDetachesInputAndOutput(t *testing.T) {
	cache := NewCache()
	input := evidenceItems()
	cache.Store("match-1", input)
	input[0].Value = "mutated-input"
	first, ok := cache.Load("match-1")
	if !ok {
		t.Fatal("cached evidence was not found")
	}
	if got := first[0].Value; got != "alpha" {
		t.Fatalf("cache followed input mutation: %q", got)
	}
	first[0].Value = "mutated-output"
	second, _ := cache.Load("match-1")
	if got := second[0].Value; got != "alpha" {
		t.Fatalf("cache exposed stored evidence: %q", got)
	}
}

func TestEvidenceCompactPreservesInput(t *testing.T) {
	a := Item{EventID: "event-1", Kind: "source", Value: "alpha"}
	b := Item{EventID: "event-2", Kind: "result", Value: "beta"}
	input := []Item{a, a, b}
	wantInput := append([]Item(nil), input...)
	got := Compact(input)
	if !reflect.DeepEqual(input, wantInput) {
		t.Fatalf("Compact overwrote its input: got %#v want %#v", input, wantInput)
	}
	if want := []Item{a, b}; !reflect.DeepEqual(got, want) {
		t.Fatalf("unexpected compacted evidence: got %#v want %#v", got, want)
	}
}

func TestEvidenceHistoryDetachesVersions(t *testing.T) {
	history := NewHistory()
	input := evidenceItems()
	history.Append("match-1", input)
	input[0].Value = "mutated-input"
	first := history.Versions("match-1")
	if got := first[0][0].Value; got != "alpha" {
		t.Fatalf("history followed input mutation: %q", got)
	}
	first[0][0].Value = "mutated-output"
	second := history.Versions("match-1")
	if got := second[0][0].Value; got != "alpha" {
		t.Fatalf("history exposed a stored version: %q", got)
	}
}

func TestEvidenceViewDetachesItems(t *testing.T) {
	input := evidenceItems()
	view := Build("match-1", input)
	input[0].Value = "mutated"
	if got := view.Items[0].Value; got != "alpha" {
		t.Fatalf("view followed input mutation: %q", got)
	}
	if view.Kinds["source"] != 1 || view.Kinds["result"] != 1 {
		t.Fatalf("view kind counts changed: %#v", view.Kinds)
	}
}
