package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Event struct {
	ID         string         `json:"id"`
	StreamID   string         `json:"stream_id"`
	TenantID   string         `json:"tenant_id"`
	Type       string         `json:"type"`
	Partition  int            `json:"partition"`
	Offset     int64          `json:"offset"`
	EventTime  time.Time      `json:"event_time"`
	IngestedAt time.Time      `json:"ingested_at"`
	Payload    map[string]any `json:"payload"`
	Retracted  bool           `json:"retracted"`
}

func (e Event) Fingerprint() string {
	b, _ := json.Marshal(e)
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}
func (e Event) Field(path string) (any, bool) { v, ok := e.Payload[path]; return v, ok }

type EventSchema struct {
	ID, StreamID, TenantID, Name, Version string
	Fields                                map[string]string
	CreatedAt                             time.Time
}
type Stream struct {
	ID, TenantID, Name string
	Partitions         int
	SchemaID           string
	Active             bool
	CreatedAt          time.Time
}
type WindowPolicy struct {
	ID, Name, Kind   string
	Size, Slide, Gap time.Duration
	Count            int
	AllowedLateness  time.Duration
	Mode             string
}
type PatternDefinition struct {
	ID, TenantID, Name, Version string
	Expression                  string
	Compiled                    *CompiledPlan
	Active                      bool
	CreatedAt                   time.Time
}
type RuleVersion struct {
	ID, PatternID, Version, Status string
	PublishedAt                    time.Time
}
type Watermark struct {
	StreamID  string
	Partition int
	Time      time.Time
	UpdatedAt time.Time
}
type Match struct {
	ID, PatternID, TenantID, StreamID string
	RuleVersion                       string
	StartedAt, CompletedAt            time.Time
	Evidence                          []string
	Explain                           ExplainNode
	Suppressed                        bool
	Reason                            string
}
type ExplainNode struct {
	Operator, Predicate, Result string
	Children                    []ExplainNode `json:"children,omitempty"`
	EventID                     string        `json:"event_id,omitempty"`
}
type CompiledPlan struct {
	Steps      []PlanStep
	Complexity int
	Hash       string
}
type PlanStep struct {
	Alias, EventType, Predicate string
	Negated                     bool
	Repeat                      int
	Within                      time.Duration
}

type Store struct {
	mu         sync.RWMutex
	Streams    map[string]Stream
	Schemas    map[string]EventSchema
	Patterns   map[string]PatternDefinition
	Windows    map[string]WindowPolicy
	Events     []Event
	Matches    map[string]Match
	Watermarks map[string]Watermark
}

type StoreView struct {
	Streams  map[string]Stream
	Schemas  map[string]EventSchema
	Patterns map[string]PatternDefinition
	Matches  map[string]Match
}

func NewStore() *Store {
	return &Store{Streams: map[string]Stream{}, Schemas: map[string]EventSchema{}, Patterns: map[string]PatternDefinition{}, Windows: map[string]WindowPolicy{}, Matches: map[string]Match{}, Watermarks: map[string]Watermark{}}
}

func (s *Store) View() StoreView {
	s.mu.RLock()
	defer s.mu.RUnlock()

	view := StoreView{
		Streams:  make(map[string]Stream, len(s.Streams)),
		Schemas:  make(map[string]EventSchema, len(s.Schemas)),
		Patterns: make(map[string]PatternDefinition, len(s.Patterns)),
		Matches:  make(map[string]Match, len(s.Matches)),
	}
	for id, stream := range s.Streams {
		view.Streams[id] = stream
	}
	for id, schema := range s.Schemas {
		view.Schemas[id] = schema
	}
	for id, pattern := range s.Patterns {
		view.Patterns[id] = pattern
	}
	for id, match := range s.Matches {
		view.Matches[id] = match
	}
	return view
}

func (s *Store) AddStream(v Stream) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v.ID == "" {
		return fmt.Errorf("stream id required")
	}
	if _, ok := s.Streams[v.ID]; ok {
		return fmt.Errorf("stream exists")
	}
	s.Streams[v.ID] = v
	return nil
}
func (s *Store) AddSchema(v EventSchema) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v.ID == "" {
		return fmt.Errorf("schema id required")
	}
	s.Schemas[v.ID] = v
	return nil
}
func (s *Store) AddPattern(v PatternDefinition) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Patterns[v.ID] = v
	return nil
}
func (s *Store) Append(e Event) { s.mu.Lock(); defer s.mu.Unlock(); s.Events = append(s.Events, e) }
func (s *Store) ListEvents(stream string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Event{}
	for _, e := range s.Events {
		if e.StreamID == stream {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].EventTime.Before(out[j].EventTime) })
	return out
}
func (s *Store) PutMatch(m Match) { s.mu.Lock(); defer s.mu.Unlock(); s.Matches[m.ID] = m }
func (s *Store) GetMatch(id string) (Match, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.Matches[id]
	return m, ok
}
func (s *Store) SetWatermark(w Watermark) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Watermarks[fmt.Sprintf("%s/%d", w.StreamID, w.Partition)] = w
}
func (s *Store) Snapshot() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return json.Marshal(s)
}
func (s *Store) Restore(b []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(b, s)
}
