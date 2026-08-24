package repository

import (
	"context"
	"errors"
	"github.com/example/cep-window-engine/internal/domain"
	"sync"
	"time"
)

var ErrNotFound = errors.New("not found")

type StreamRepository interface {
	Create(context.Context, domain.Stream) error
	Get(context.Context, string) (domain.Stream, error)
	List(context.Context) []domain.Stream
}
type PatternRepository interface {
	Create(context.Context, domain.PatternDefinition) error
	Get(context.Context, string) (domain.PatternDefinition, error)
	List(context.Context) []domain.PatternDefinition
}
type Memory struct {
	mu       sync.RWMutex
	streams  map[string]domain.Stream
	patterns map[string]domain.PatternDefinition
}

func New() *Memory {
	return &Memory{streams: map[string]domain.Stream{}, patterns: map[string]domain.PatternDefinition{}}
}
func (m *Memory) Create(ctx context.Context, v domain.Stream) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.streams[v.ID]; ok {
		return errors.New("exists")
	}
	m.streams[v.ID] = v
	return nil
}
func (m *Memory) Get(ctx context.Context, id string) (domain.Stream, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.streams[id]
	if !ok {
		return domain.Stream{}, ErrNotFound
	}
	return v, nil
}
func (m *Memory) List(ctx context.Context) []domain.Stream {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.Stream, 0, len(m.streams))
	for _, v := range m.streams {
		out = append(out, v)
	}
	return out
}
func (m *Memory) CreatePattern(ctx context.Context, v domain.PatternDefinition) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.patterns[v.ID] = v
	return nil
}
func (m *Memory) GetPattern(ctx context.Context, id string) (domain.PatternDefinition, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.patterns[id]
	if !ok {
		return domain.PatternDefinition{}, ErrNotFound
	}
	return v, nil
}
func (m *Memory) ListPattern(ctx context.Context) []domain.PatternDefinition {
	m.mu.RLock()
	defer m.mu.RUnlock()
	o := []domain.PatternDefinition{}
	for _, v := range m.patterns {
		o = append(o, v)
	}
	return o
}

type AuditEntry struct {
	ID, Action, Actor, Resource string
	At                          time.Time
	Metadata                    map[string]string
}
type AuditLog struct {
	mu      sync.Mutex
	entries []AuditEntry
}

func (a *AuditLog) Append(e AuditEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if e.At.IsZero() {
		e.At = time.Now()
	}
	a.entries = append(a.entries, e)
}
func (a *AuditLog) List() []AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]AuditEntry(nil), a.entries...)
}
