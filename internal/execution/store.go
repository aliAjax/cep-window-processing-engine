package execution

import (
	"context"
	"fmt"
	"sync"
)

type Evaluator func(context.Context, Request) (bool, error)

type Store struct {
	mu         sync.RWMutex
	evaluators map[string]Evaluator
}

func NewStore() *Store {
	return &Store{evaluators: make(map[string]Evaluator)}
}

func (s *Store) Register(patternID string, evaluator Evaluator) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evaluators[patternID] = evaluator
}

func (s *Store) Evaluate(ctx context.Context, request Request) (bool, error) {
	s.mu.RLock()
	evaluator := s.evaluators[request.PatternID]
	s.mu.RUnlock()
	if evaluator == nil {
		return false, fmt.Errorf("pattern evaluator not found")
	}
	if ctx == nil {
		ctx = request.Context()
	}
	return evaluator(request.Context(), request)
}
