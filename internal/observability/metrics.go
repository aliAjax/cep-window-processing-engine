package observability

import (
	"fmt"
	"sync"
	"time"
)

type Counter struct {
	mu    sync.Mutex
	Value uint64
}

func (c *Counter) Inc()         { c.mu.Lock(); c.Value++; c.mu.Unlock() }
func (c *Counter) Add(n uint64) { c.mu.Lock(); c.Value += n; c.mu.Unlock() }
func (c *Counter) Get() uint64  { c.mu.Lock(); defer c.mu.Unlock(); return c.Value }

type Metrics struct {
	Events, Matches, Duplicates, Late, Errors Counter
	Started                                   time.Time
}

func New() *Metrics { return &Metrics{Started: time.Now()} }
func (m *Metrics) Snapshot() map[string]any {
	return map[string]any{"events": m.Events.Get(), "matches": m.Matches.Get(), "duplicates": m.Duplicates.Get(), "late": m.Late.Get(), "errors": m.Errors.Get(), "uptime_seconds": time.Since(m.Started).Seconds()}
}

type Logger struct {
	mu     sync.Mutex
	Fields map[string]string
}

func (l *Logger) Log(level, msg string, fields map[string]any) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Printf(`{"level":%q,"msg":%q,"ts":%q}`, level, msg, time.Now().UTC().Format(time.RFC3339))
	for k, v := range fields {
		fmt.Printf(`,"%s":%v`, k, v)
	}
	fmt.Println()
}
