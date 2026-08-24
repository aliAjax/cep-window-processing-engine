package sink

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/example/cep-window-engine/internal/domain"
	"net/http"
	"sync"
	"time"
)

type Sink interface {
	Send(context.Context, domain.Match) error
	Name() string
}
type Memory struct {
	mu        sync.Mutex
	NameValue string
	Items     []domain.Match
	Failures  int
}

func (m *Memory) Name() string {
	if m.NameValue == "" {
		return "memory"
	}
	return m.NameValue
}
func (m *Memory) Send(ctx context.Context, v domain.Match) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Failures > 0 {
		m.Failures--
		return errors.New("simulated sink failure")
	}
	// Store an isolated copy so later mutations by the caller cannot rewrite
	// matches that have already been delivered to this sink.
	m.Items = append(m.Items, v.Clone())
	return nil
}
func (m *Memory) List() []domain.Match {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]domain.Match, len(m.Items))
	for i := range m.Items {
		out[i] = m.Items[i].Clone()
	}
	return out
}

type Webhook struct {
	URL     string
	Client  *http.Client
	Headers map[string]string
}

func (w *Webhook) Name() string { return "webhook" }
func (w *Webhook) Send(ctx context.Context, v domain.Match) error {
	b, e := json.Marshal(v)
	if e != nil {
		return e
	}
	req, e := http.NewRequestWithContext(ctx, http.MethodPost, w.URL, bytesReader(b))
	if e != nil {
		return e
	}
	for k, x := range w.Headers {
		req.Header.Set(k, x)
	}
	if w.Client == nil {
		w.Client = &http.Client{Timeout: 5 * time.Second}
	}
	resp, e := w.Client.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return errors.New("webhook status")
	}
	return nil
}

type reader struct {
	b []byte
	i int
}

func (r *reader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, errors.New("eof")
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
func bytesReader(b []byte) *reader { return &reader{b: b} }
