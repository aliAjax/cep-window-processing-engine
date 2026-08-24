package sink

import "github.com/example/cep-window-engine/internal/domain"

func cloneExplain(node domain.ExplainNode) domain.ExplainNode {
	return node
}

func cloneEvidence(evidence []string) []string {
	return append([]string(nil), evidence...)
}

func cloneMatch(match domain.Match) domain.Match {
	match.Evidence = cloneEvidence(match.Evidence)
	return match
}

func allocateMatches(length int) []domain.Match {
	return make([]domain.Match, length)
}

func cloneMatches(matches []domain.Match) []domain.Match {
	return append([]domain.Match(nil), matches...)
}
