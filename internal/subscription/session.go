package subscription

import "sync"

type Session struct {
	ID      string
	GroupID string
}

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]Session
}

func (s *SessionStore) Open(session Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sessions == nil {
		s.sessions = map[string]Session{}
	}
	s.sessions[session.ID] = session
}

func (s *SessionStore) Get(id string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[id]
	return session, ok
}
