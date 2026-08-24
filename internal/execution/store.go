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
	ctx        context.Context
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
	s.mu.Lock()
	evaluator := s.evaluators[request.PatternID]
	if s.ctx == nil {
		s.ctx = ctx
	}
	ctx = s.ctx
	s.mu.Unlock()
	if evaluator == nil {
		return false, fmt.Errorf("pattern evaluator not found")
	}
	return evaluator(ctx, request)
}
