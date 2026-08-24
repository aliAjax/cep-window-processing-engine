package sink

import (
	"context"
	"testing"

	"github.com/example/cep-window-engine/internal/domain"
)

func TestMemoryListOwnsMatchGraph(t *testing.T) {
	memory := &Memory{}
	match := domain.Match{ID: "m-1", Evidence: []string{"e-1"}, Explain: domain.ExplainNode{Result: "true", Children: []domain.ExplainNode{{EventID: "e-1", Result: "true"}}}}
	if err := memory.Send(context.Background(), match); err != nil {
		t.Fatal(err)
	}
	match.Evidence[0] = "input-mutated"
	match.Explain.Children[0].Result = "false"

	first := memory.List()
	if first[0].Evidence[0] != "e-1" || first[0].Explain.Children[0].Result != "true" {
		t.Fatalf("sink retained input aliases: %#v", first[0])
	}
	first[0].Evidence[0] = "output-mutated"
	first[0].Explain.Children[0].Result = "false"
	second := memory.List()
	if second[0].Evidence[0] != "e-1" || second[0].Explain.Children[0].Result != "true" {
		t.Fatalf("sink snapshot mutation changed stored match: %#v", second[0])
	}
}
