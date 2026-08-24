package lifecycle

import (
	"fmt"
	"sync"
)

type Service struct {
	mu     sync.RWMutex
	states map[string]State
}

func NewService() *Service {
	return &Service{states: make(map[string]State)}
}

func (s *Service) Transition(ruleID string, next State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	current := s.states[ruleID]
	if current == "" {
		current = StatePending
	}
	if !current.CanTransition(next) {
		return fmt.Errorf("invalid rule state transition %s -> %s", current, next)
	}
	s.states[ruleID] = next
	return nil
}

func (s *Service) State(ruleID string) State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.states[ruleID]
}

func (s *Service) Snapshot() map[string]State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]State, len(s.states))
	for id, state := range s.states {
		out[id] = state
	}
	return out
}
