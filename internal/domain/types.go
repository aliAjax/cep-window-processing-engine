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

// clonePayload returns a deep copy of a payload map so callers cannot mutate
// events that have already been admitted into the engine. Values that are
// themselves maps or slices are copied recursively to break nested aliasing.
func clonePayload(p map[string]any) map[string]any {
	if p == nil {
		return nil
	}
	out := make(map[string]any, len(p))
	for k, v := range p {
		out[k] = cloneValue(v)
	}
	return out
}

func cloneValue(v any) any {
	switch val := v.(type) {
	case map[string]any:
		return clonePayload(val)
	case []any:
		out := make([]any, len(val))
		for i := range val {
			out[i] = cloneValue(val[i])
		}
		return out
	default:
		return v
	}
}

// Clone returns an isolated copy of the event. The Payload map is copied so
// that later mutations by the caller (or any shared reference) cannot alter
// the historical record or window snapshots held by the engine.
func (e Event) Clone() Event {
	cp := e
	cp.Payload = clonePayload(e.Payload)
	return cp
}

// Clone returns an isolated copy of the match. The Evidence slice and the
// Explain tree (which holds nested child slices) are copied so that mutations
// to a returned match cannot propagate back into the store or to sinks.
func (m Match) Clone() Match {
	cp := m
	cp.Evidence = append([]string(nil), m.Evidence...)
	cp.Explain = cloneExplain(m.Explain)
	return cp
}

func cloneExplain(n ExplainNode) ExplainNode {
	cp := n
	if len(n.Children) > 0 {
		cp.Children = make([]ExplainNode, len(n.Children))
		for i := range n.Children {
			cp.Children[i] = cloneExplain(n.Children[i])
		}
	}
	return cp
}

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

func NewStore() *Store {
	return &Store{Streams: map[string]Stream{}, Schemas: map[string]EventSchema{}, Patterns: map[string]PatternDefinition{}, Windows: map[string]WindowPolicy{}, Matches: map[string]Match{}, Watermarks: map[string]Watermark{}}
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
func (s *Store) Append(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Events = append(s.Events, e.Clone())
}
func (s *Store) ListEvents(stream string) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []Event{}
	for _, e := range s.Events {
		if e.StreamID == stream {
			out = append(out, e.Clone())
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].EventTime.Before(out[j].EventTime) })
	return out
}
func (s *Store) PutMatch(m Match) { s.mu.Lock(); defer s.mu.Unlock(); s.Matches[m.ID] = m.Clone() }
func (s *Store) GetMatch(id string) (Match, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.Matches[id]
	if !ok {
		return Match{}, false
	}
	return m.Clone(), true
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
