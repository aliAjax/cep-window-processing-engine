package auditview

type SessionIndex struct {
	entries map[string][]string
}

func NewSessionIndex() *SessionIndex {
	return &SessionIndex{entries: make(map[string][]string)}
}

func (s *SessionIndex) Add(sessionID, eventID string) {
	s.entries[sessionID] = append(s.entries[sessionID], eventID)
}

func (s *SessionIndex) Snapshot() map[string][]string {
	return s.entries
}
